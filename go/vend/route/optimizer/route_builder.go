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

	var stops []RouteStop
	curLat, curLng := 0.0, 0.0
	if len(machines) > 0 {
		curLat, curLng = machines[0].Lat, machines[0].Lng
	}

	for _, m := range machines {
		// Check if truck has enough stock for this machine
		if !hasStock(truckStock, m.Products) {
			// Insert facility reload
			fac := findNearestOpenFacility(curLat, curLng, facilities)
			if fac != nil {
				stops = append(stops, RouteStop{
					FacilityId: fac.FacilityId,
					Lat:        fac.Coordinates.Latitude,
					Lng:        fac.Coordinates.Longitude,
					Urgency:    "reload",
					IsReload:   true,
				})
				reloadTruckFromFacility(truckStock, fac)
				curLat, curLng = fac.Coordinates.Latitude, fac.Coordinates.Longitude
			}
		}

		// Deduct from truck stock
		for sku, qty := range m.Products {
			truckStock[sku] -= qty
			if truckStock[sku] < 0 {
				truckStock[sku] = 0
			}
		}

		stops = append(stops, RouteStop{
			MachineId: m.MachineId,
			Lat:       m.Lat,
			Lng:       m.Lng,
			Urgency:   m.Urgency,
			Products:  m.Products,
		})
		curLat, curLng = m.Lat, m.Lng
	}
	return stops
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
