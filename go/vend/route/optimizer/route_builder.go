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
	ServiceMinutes int32
	ReloadMinutes  int32
	FuelPriceGal   float64
}

// RouteStop represents a stop in the built route (machine or facility reload).
type RouteStop struct {
	MachineId  string
	FacilityId string
	Lat        float64
	Lng        float64
	Urgency    string // "high", "low", "reload"
	IsReload   bool
	Products   map[string]int32 // sku → qty to restock at this stop
}

// BuiltRoute is the result of BuildRoute — ready for traffic refinement.
type BuiltRoute struct {
	Stops    []RouteStop
	Legs     []RouteLeg
	Metrics  RouteMetrics
	TruckMPG float64
}

// BuildRoute orders stops, inserts facility reloads when truck stock depletes,
// and applies 2-opt improvement.
func BuildRoute(cluster *Cluster, assignment *Assignment,
	facilities []*vend.VendStockingFacility, config *RouteConfig) *BuiltRoute {

	// Step 1: Order machines by nearest-neighbor from driver start
	ordered := nearestNeighborOrder(cluster.Machines, assignment.StartLat, assignment.StartLng)

	// Step 2: Walk through ordered stops, track stock, insert facility reloads
	truckStock := buildTruckStockMap(assignment.Truck)
	stops := insertFacilityReloads(ordered, truckStock, facilities, config)

	// Step 3: Apply 2-opt on machine stops only, then re-insert facility reloads
	machineStops := extractMachineStops(stops)
	improved := twoOpt(machineStops)
	truckStock = buildTruckStockMap(assignment.Truck) // reset stock
	stops = insertFacilityReloads(improved, truckStock, facilities, config)

	// Step 4: Build legs and compute metrics
	legs := buildLegs(assignment.StartLat, assignment.StartLng, stops, config.AvgSpeedMph)
	mpg := assignment.Truck.MilesPerGallon
	metrics := ComputeRouteMetrics(legs, 0, mpg, config.FuelPriceGal,
		config.ServiceMinutes, config.ReloadMinutes)

	return &BuiltRoute{Stops: stops, Legs: legs, Metrics: metrics, TruckMPG: mpg}
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

	// Phase 1: Walk through machines in order. Serve all that current stock allows.
	// Accumulate the ones we can't serve into a "deferred" list.
	var served []RouteStop
	var deferred []MachineDemand
	curLat, curLng := 0.0, 0.0
	if len(machines) > 0 {
		curLat, curLng = machines[0].Lat, machines[0].Lng
	}

	for _, m := range machines {
		if hasStock(truckStock, m.Products) {
			deductStock(truckStock, m.Products)
			served = append(served, RouteStop{
				MachineId: m.MachineId, Lat: m.Lat, Lng: m.Lng,
				Urgency: m.Urgency, Products: m.Products,
			})
			curLat, curLng = m.Lat, m.Lng
		} else {
			deferred = append(deferred, m)
		}
	}

	if len(deferred) == 0 {
		return served // No reload needed — truck had enough stock
	}

	// Phase 2: Need to reload. Find the optimal facility that minimizes total detour.
	// The reload goes between the last served stop and the first deferred stop.
	// Evaluate each facility: cost = (lastServed → facility) + (facility → firstDeferred)
	// vs direct (lastServed → firstDeferred).
	// Pick the facility with the minimum added cost.
	var stops []RouteStop
	stops = append(stops, served...)

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

func buildLegs(startLat, startLng float64, stops []RouteStop, avgSpeed float64) []RouteLeg {
	legs := make([]RouteLeg, len(stops))
	curLat, curLng := startLat, startLng
	for i, s := range stops {
		legs[i] = HaversineLeg(curLat, curLng, s.Lat, s.Lng, avgSpeed, s.IsReload)
		curLat, curLng = s.Lat, s.Lng
	}
	return legs
}
