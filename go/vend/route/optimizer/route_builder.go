/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 */
package optimizer

import (
	"github.com/saichler/l8vendingmachine/go/types/vend"
)

// RouteConfig holds tuning parameters for route building.
type RouteConfig struct {
	AvgSpeedMph    float64
	ServiceMinutes     int32
	ReloadMinutes      int32
	FuelPriceGal       float64
	BreakAfterMinutes  int32 // insert break after this many driving minutes (default 240)
	BreakDurationMinutes int32 // break length (default 30)
}

// RouteStop represents a stop in the built route (machine, facility reload, or end-of-day).
type RouteStop struct {
	MachineId  string
	FacilityId string
	Lat        float64
	Lng        float64
	Urgency    string // "high", "low", "reload", "end", "break"
	IsReload   bool
	IsEnd      bool
	IsBreak    bool
	Products   map[string]int32 // sku → qty to restock at this stop
}

// BuiltRoute is the result of BuildRoute — ready for traffic refinement.
type BuiltRoute struct {
	Stops    []RouteStop
	Legs     []RouteLeg
	Metrics  RouteMetrics
	TruckMPG float64
}

// BuildRouteForDriver orders stops with end-location awareness, inserts facility
// reloads (preferring home depot near end), and applies 2-opt including the end leg.
func BuildRouteForDriver(dr *DriverRoute,
	facilities []*vend.VendStockingFacility, config *RouteConfig, router *Router) *BuiltRoute {

	// Step 1: Order machines by nearest-neighbor, anchoring last stop near end location
	ordered := nearestNeighborWithEnd(dr.Machines, dr.StartLat, dr.StartLng, dr.EndLat, dr.EndLng)

	// Step 2: 2-opt including the end-location leg
	improved := twoOptWithEnd(ordered, dr.EndLat, dr.EndLng)

	// Step 3: Insert facility reloads (prefer home depot near end of route)
	truckStock := buildTruckStockMap(dr.Truck)
	stops := insertFacilityReloads(improved, truckStock, facilities, config)

	// Step 4: Prefer home depot for last reload
	preferHomeDepotForLastReload(stops, dr.Truck, facilities)

	// Step 5: Insert break stops after N hours of cumulative driving/service
	stops = insertBreaks(stops, dr.StartLat, dr.StartLng, config, router)

	// Step 6: Add end-of-day stop
	stops = append(stops, RouteStop{
		Lat: dr.EndLat, Lng: dr.EndLng,
		Urgency: "end", IsReload: false, IsEnd: true,
	})

	// Step 6: Build legs using Router (OSRM road distance or haversine fallback)
	legs := buildLegsWithRouter(dr.StartLat, dr.StartLng, stops, router)
	mpg := dr.Truck.MilesPerGallon
	metrics := ComputeRouteMetrics(legs, 0, mpg, config.FuelPriceGal,
		config.ServiceMinutes, config.ReloadMinutes, config.BreakDurationMinutes)

	return &BuiltRoute{Stops: stops, Legs: legs, Metrics: metrics, TruckMPG: mpg}
}

// nearestNeighborWithEnd: save machine nearest to end location as last stop,
// then nearest-neighbor the rest from start, then append the saved last stop.
func nearestNeighborWithEnd(machines []MachineDemand, startLat, startLng, endLat, endLng float64) []MachineDemand {
	if len(machines) <= 1 {
		return machines
	}

	// Find machine nearest to end location — reserve it as last stop
	lastIdx := 0
	bestEndDist := Haversine(machines[0].Lat, machines[0].Lng, endLat, endLng)
	for i := 1; i < len(machines); i++ {
		d := Haversine(machines[i].Lat, machines[i].Lng, endLat, endLng)
		if d < bestEndDist {
			bestEndDist = d
			lastIdx = i
		}
	}
	lastMachine := machines[lastIdx]
	remaining := make([]MachineDemand, 0, len(machines)-1)
	for i, m := range machines {
		if i != lastIdx {
			remaining = append(remaining, m)
		}
	}

	// Nearest-neighbor the rest from start
	ordered := nearestNeighborOrder(remaining, startLat, startLng)

	// Append the reserved last stop
	ordered = append(ordered, lastMachine)
	return ordered
}

func nearestNeighborOrder(machines []MachineDemand, startLat, startLng float64) []MachineDemand {
	remaining := make([]MachineDemand, len(machines))
	copy(remaining, machines)
	var ordered []MachineDemand

	curLat, curLng := startLat, startLng
	for len(remaining) > 0 {
		bestIdx := 0
		bestDist := Haversine(curLat, curLng, remaining[0].Lat, remaining[0].Lng)
		for i := 1; i < len(remaining); i++ {
			d := Haversine(curLat, curLng, remaining[i].Lat, remaining[i].Lng)
			if d < bestDist {
				bestDist = d
				bestIdx = i
			}
		}
		chosen := remaining[bestIdx]
		ordered = append(ordered, chosen)
		curLat, curLng = chosen.Lat, chosen.Lng
		remaining = append(remaining[:bestIdx], remaining[bestIdx+1:]...)
	}
	return ordered
}

func buildTruckStockMap(truck *vend.VendDeliveryTruck) map[string]int32 {
	stock := make(map[string]int32)
	for _, item := range truck.Stock {
		stock[item.Sku] = item.Quantity
	}
	return stock
}

func insertFacilityReloads(machines []MachineDemand, truckStock map[string]int32,
	facilities []*vend.VendStockingFacility, config *RouteConfig) []RouteStop {

	// Walk through machines in 2-opt order. When we can't serve one:
	// - If the machine is NEAR our current position (< nearbyThreshold), reload now
	//   and serve it — don't drive away and come back later.
	// - If it's far, defer it for later.
	const nearbyThresholdMiles = 5.0

	var stops []RouteStop
	var deferred []MachineDemand
	curLat, curLng := 0.0, 0.0
	if len(machines) > 0 {
		curLat, curLng = machines[0].Lat, machines[0].Lng
	}

	for _, m := range machines {
		if hasStock(truckStock, m.Products) {
			deductStock(truckStock, m.Products)
			stops = append(stops, RouteStop{
				MachineId: m.MachineId, Lat: m.Lat, Lng: m.Lng,
				Urgency: m.Urgency, Products: m.Products,
			})
			curLat, curLng = m.Lat, m.Lng
		} else {
			// Check if this machine is nearby — if so, reload now instead of deferring
			dist := Haversine(curLat, curLng, m.Lat, m.Lng)
			if dist <= nearbyThresholdMiles {
				fac := findOptimalFacility(curLat, curLng, []MachineDemand{m}, facilities)
				if fac != nil {
					stops = append(stops, RouteStop{
						FacilityId: fac.FacilityId,
						Lat: fac.Coordinates.Latitude, Lng: fac.Coordinates.Longitude,
						Urgency: "reload", IsReload: true,
					})
					reloadTruckFromFacility(truckStock, fac)
					curLat, curLng = fac.Coordinates.Latitude, fac.Coordinates.Longitude
					// Now serve the machine
					deductStock(truckStock, m.Products)
					stops = append(stops, RouteStop{
						MachineId: m.MachineId, Lat: m.Lat, Lng: m.Lng,
						Urgency: m.Urgency, Products: m.Products,
					})
					curLat, curLng = m.Lat, m.Lng
					continue
				}
			}
			deferred = append(deferred, m)
		}
	}

	if len(deferred) == 0 {
		return stops
	}

	for len(deferred) > 0 {
		// Find optimal facility for reload
		fac := findOptimalFacility(curLat, curLng, deferred, facilities)
		if fac == nil {
			break
		}

		stops = append(stops, RouteStop{
			FacilityId: fac.FacilityId,
			Lat:        fac.Coordinates.Latitude,
			Lng:        fac.Coordinates.Longitude,
			Urgency:    "reload",
			IsReload:   true,
		})
		reloadTruckFromFacility(truckStock, fac)
		curLat, curLng = fac.Coordinates.Latitude, fac.Coordinates.Longitude

		// Serve as many deferred machines as possible (nearest first)
		var stillDeferred []MachineDemand
		for len(deferred) > 0 {
			bestIdx := -1
			bestDist := 999999.0
			for i, m := range deferred {
				if hasStock(truckStock, m.Products) {
					d := Haversine(curLat, curLng, m.Lat, m.Lng)
					if d < bestDist {
						bestDist = d
						bestIdx = i
					}
				}
			}
			if bestIdx == -1 {
				// Remaining deferred can't be served — need another reload
				stillDeferred = deferred
				break
			}
			m := deferred[bestIdx]
			deductStock(truckStock, m.Products)
			stops = append(stops, RouteStop{
				MachineId: m.MachineId, Lat: m.Lat, Lng: m.Lng,
				Urgency: m.Urgency, Products: m.Products,
			})
			curLat, curLng = m.Lat, m.Lng
			deferred = append(deferred[:bestIdx], deferred[bestIdx+1:]...)
		}
		deferred = stillDeferred
	}

	return stops
}

// findOptimalFacility evaluates all facilities and picks the one that minimizes
// the total detour: (current position → facility → nearest deferred machine).
func findOptimalFacility(curLat, curLng float64, deferred []MachineDemand,
	facilities []*vend.VendStockingFacility) *vend.VendStockingFacility {

	var best *vend.VendStockingFacility
	bestCost := 999999.0

	// Find the centroid of deferred machines (where we want to end up after reload)
	points := make([][2]float64, len(deferred))
	for i, m := range deferred {
		points[i] = [2]float64{m.Lat, m.Lng}
	}
	deferredLat, deferredLng := Centroid(points)

	for _, f := range facilities {
		if f.Status != vend.VendFacilityStatus_VEND_FACILITY_STATUS_ACTIVE {
			continue
		}
		if f.Coordinates == nil {
			continue
		}
		fLat, fLng := f.Coordinates.Latitude, f.Coordinates.Longitude
		// Cost = distance to facility + distance from facility to deferred centroid
		cost := Haversine(curLat, curLng, fLat, fLng) + Haversine(fLat, fLng, deferredLat, deferredLng)
		if cost < bestCost {
			bestCost = cost
			best = f
		}
	}
	return best
}

func deductStock(stock map[string]int32, products map[string]int32) {
	for sku, qty := range products {
		stock[sku] -= qty
		if stock[sku] < 0 {
			stock[sku] = 0
		}
	}
}

func hasStock(truckStock map[string]int32, demand map[string]int32) bool {
	for sku, qty := range demand {
		if truckStock[sku] < qty {
			return false
		}
	}
	return true
}

func findNearestOpenFacility(lat, lng float64, facilities []*vend.VendStockingFacility) *vend.VendStockingFacility {
	var best *vend.VendStockingFacility
	bestDist := 999999.0
	for _, f := range facilities {
		if f.Status != vend.VendFacilityStatus_VEND_FACILITY_STATUS_ACTIVE {
			continue
		}
		if f.Coordinates == nil {
			continue
		}
		d := Haversine(lat, lng, f.Coordinates.Latitude, f.Coordinates.Longitude)
		if d < bestDist {
			bestDist = d
			best = f
		}
	}
	return best
}

func reloadTruckFromFacility(truckStock map[string]int32, fac *vend.VendStockingFacility) {
	for _, item := range fac.Stock {
		if item.Quantity > 0 {
			truckStock[item.Sku] = item.Quantity
		}
	}
}

func extractMachineStops(stops []RouteStop) []MachineDemand {
	var machines []MachineDemand
	for _, s := range stops {
		if !s.IsReload {
			machines = append(machines, MachineDemand{
				MachineId: s.MachineId,
				Lat:       s.Lat,
				Lng:       s.Lng,
				Products:  s.Products,
				Urgency:   s.Urgency,
			})
		}
	}
	return machines
}

func buildLegsWithRouter(startLat, startLng float64, stops []RouteStop, router *Router) []RouteLeg {
	legs := make([]RouteLeg, len(stops))
	curLat, curLng := startLat, startLng
	for i, s := range stops {
		dist, dur := router.Distance(curLat, curLng, s.Lat, s.Lng)
		legs[i] = RouteLeg{
			DistanceMiles:   dist,
			DurationSeconds: dur,
			IsReload:        s.IsReload,
			IsBreak:         s.IsBreak,
			IsEnd:           s.IsEnd,
		}
		curLat, curLng = s.Lat, s.Lng
	}
	return legs
}

// twoOptWithEnd runs 2-opt including the end-location leg in distance calculations.
func twoOptWithEnd(stops []MachineDemand, endLat, endLng float64) []MachineDemand {
	if len(stops) < 4 {
		return stops
	}
	improved := make([]MachineDemand, len(stops))
	copy(improved, stops)

	totalDist := func(s []MachineDemand) float64 {
		d := 0.0
		for i := 0; i < len(s)-1; i++ {
			d += Haversine(s[i].Lat, s[i].Lng, s[i+1].Lat, s[i+1].Lng)
		}
		// Include end-location leg
		if len(s) > 0 {
			last := s[len(s)-1]
			d += Haversine(last.Lat, last.Lng, endLat, endLng)
		}
		return d
	}

	changed := true
	for changed {
		changed = false
		bestDist := totalDist(improved)
		for i := 0; i < len(improved)-2; i++ {
			for j := i + 2; j < len(improved); j++ {
				candidate := make([]MachineDemand, len(improved))
				copy(candidate, improved)
				reverseSegmentDemand(candidate, i+1, j)
				d := totalDist(candidate)
				if d < bestDist-0.01 {
					copy(improved, candidate)
					bestDist = d
					changed = true
				}
			}
		}
	}
	return improved
}

func reverseSegmentDemand(stops []MachineDemand, i, j int) {
	for i < j {
		stops[i], stops[j] = stops[j], stops[i]
		i++
		j--
	}
}

// preferHomeDepotForLastReload checks the last reload stop — if the truck's home depot
// is a facility and is within 2× the distance of the current reload, swap it.
func preferHomeDepotForLastReload(stops []RouteStop, truck *vend.VendDeliveryTruck,
	facilities []*vend.VendStockingFacility) {

	if truck.HomeDepotId == "" {
		return
	}
	// Find last reload stop
	lastReloadIdx := -1
	for i := len(stops) - 1; i >= 0; i-- {
		if stops[i].IsReload {
			lastReloadIdx = i
			break
		}
	}
	if lastReloadIdx < 0 {
		return
	}
	if stops[lastReloadIdx].FacilityId == truck.HomeDepotId {
		return // Already at home depot
	}

	// Find home depot facility
	var homeDepot *vend.VendStockingFacility
	for _, f := range facilities {
		if f.FacilityId == truck.HomeDepotId && f.Status == vend.VendFacilityStatus_VEND_FACILITY_STATUS_ACTIVE {
			homeDepot = f
			break
		}
	}
	if homeDepot == nil || homeDepot.Coordinates == nil {
		return
	}

	// Compare distances: current reload vs home depot from the previous stop
	prevLat, prevLng := stops[lastReloadIdx].Lat, stops[lastReloadIdx].Lng
	if lastReloadIdx > 0 {
		prevLat, prevLng = stops[lastReloadIdx-1].Lat, stops[lastReloadIdx-1].Lng
	}
	currentDist := Haversine(prevLat, prevLng, stops[lastReloadIdx].Lat, stops[lastReloadIdx].Lng)
	homeDepotDist := Haversine(prevLat, prevLng, homeDepot.Coordinates.Latitude, homeDepot.Coordinates.Longitude)

	// Swap if home depot is within 2× the distance
	if homeDepotDist <= currentDist*2 {
		stops[lastReloadIdx].FacilityId = homeDepot.FacilityId
		stops[lastReloadIdx].Lat = homeDepot.Coordinates.Latitude
		stops[lastReloadIdx].Lng = homeDepot.Coordinates.Longitude
	}
}

// insertBreaks walks through stops, tracks cumulative time, and inserts a break
// stop when the driver has been working longer than breakAfterMinutes.
// If a facility reload is near the break point, the break is combined with it.
func insertBreaks(stops []RouteStop, startLat, startLng float64, config *RouteConfig, router *Router) []RouteStop {
	breakAfter := int64(config.BreakAfterMinutes) * 60
	if breakAfter <= 0 {
		return stops // breaks disabled
	}

	var result []RouteStop
	cumulativeSecs := int64(0)
	lastBreakAt := int64(0)
	curLat, curLng := startLat, startLng

	for _, s := range stops {
		// Estimate travel time to this stop
		_, travelSecs := router.Distance(curLat, curLng, s.Lat, s.Lng)
		cumulativeSecs += travelSecs

		// Check if we need a break before this stop
		timeSinceBreak := cumulativeSecs - lastBreakAt
		if timeSinceBreak >= breakAfter && !s.IsReload {
			// Insert break at current position (driver stops where they are)
			result = append(result, RouteStop{
				Lat: curLat, Lng: curLng,
				Urgency: "break", IsBreak: true,
			})
			lastBreakAt = cumulativeSecs
		}

		// If this is a reload stop, it counts as a break too (driver rests while loading)
		if s.IsReload {
			lastBreakAt = cumulativeSecs
		}

		result = append(result, s)
		curLat, curLng = s.Lat, s.Lng

		// Add service time
		if s.IsReload {
			cumulativeSecs += int64(config.ReloadMinutes) * 60
		} else {
			cumulativeSecs += int64(config.ServiceMinutes) * 60
		}
	}
	return result
}
