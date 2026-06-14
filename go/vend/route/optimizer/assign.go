/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 */
package optimizer

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/vend/route/drivers"
	"github.com/saichler/l8vendingmachine/go/vend/route/trucks"
	"github.com/saichler/l8vendingmachine/go/vend/warehouse/facilities"
)

const DefaultShiftDurationSecs = 8 * 3600 // 8 hours default

// DriverRoute holds a driver's assignment: who they are, their truck, locations, and assigned machines.
type DriverRoute struct {
	Driver           *vend.VendDriver
	Truck            *vend.VendDeliveryTruck
	Schedule         *vend.VendDriverScheduleDay
	StartLat         float64
	StartLng         float64
	EndLat           float64
	EndLng           float64
	ShiftDurationSecs int64
	Machines         []MachineDemand
}

// AssignMachinesToDrivers assigns machines to drivers based on geographic proximity
// to each driver's start and end location. Replaces cluster-then-assign pattern.
func AssignMachinesToDrivers(listA, listB []MachineDemand,
	allTrucks []*vend.VendDeliveryTruck, allDrivers []*vend.VendDriver,
	allFacilities []*vend.VendStockingFacility,
	dayOfWeek vend.VendDayOfWeek, maxDetourMiles float64,
	config *RouteConfig, router *Router, nic ifs.IVNic) []DriverRoute {

	// Build available driver routes (driver+truck pairs that work on this day)
	var driverRoutes []DriverRoute
	usedTrucks := make(map[string]bool)

	for _, driver := range allDrivers {
		if !driver.IsActive {
			continue
		}
		sched := findScheduleDay(driver.Schedule, dayOfWeek)
		if sched == nil {
			continue
		}
		truck := findTruckForDriver(driver, allTrucks, usedTrucks)
		if truck == nil {
			continue
		}
		usedTrucks[truck.TruckId] = true

		startLat, startLng := resolveStartLocation(driver, sched, nic)
		endLat, endLng := resolveEndLocation(driver, sched, truck, allFacilities, startLat, startLng, nic)

		shiftSecs := int64(DefaultShiftDurationSecs)
		if sched.ShiftDurationMinutes > 0 {
			shiftSecs = int64(sched.ShiftDurationMinutes) * 60
		}
		driverRoutes = append(driverRoutes, DriverRoute{
			Driver: driver, Truck: truck, Schedule: sched,
			StartLat: startLat, StartLng: startLng,
			EndLat: endLat, EndLng: endLng,
			ShiftDurationSecs: shiftSecs,
		})
	}

	if len(driverRoutes) == 0 {
		return nil
	}

	// Assign List A machines to nearest driver (60% start proximity, 40% end proximity)
	assignMachines(listA, driverRoutes, config, router)

	// Insert List B machines if within detour threshold
	for _, m := range listB {
		bestDriver := -1
		bestCost := maxDetourMiles + 1
		for di, dr := range driverRoutes {
			if len(dr.Machines) == 0 {
				continue
			}
			cost := insertionCostForDriver(m, dr)
			if cost < bestCost {
				bestCost = cost
				bestDriver = di
			}
		}
		if bestDriver >= 0 {
			m.Urgency = "low"
			driverRoutes[bestDriver].Machines = append(driverRoutes[bestDriver].Machines, m)
		}
	}

	// Filter out drivers with no machines
	var result []DriverRoute
	for _, dr := range driverRoutes {
		if len(dr.Machines) > 0 {
			result = append(result, dr)
		}
	}
	return result
}

func assignMachines(machines []MachineDemand, driverRoutes []DriverRoute, config *RouteConfig, router *Router) {
	assigned := make([]bool, len(machines))

	for {
		bestMachine := -1
		bestDriver := -1
		bestScore := 999999.0

		for mi, m := range machines {
			if assigned[mi] {
				continue
			}
			for di, dr := range driverRoutes {
				if !driverHasSkills(dr.Driver, m.RequiredSkills) {
					continue
				}
				svcMin := m.ServiceMinutes
				if svcMin <= 0 {
					svcMin = config.ServiceMinutes
				}
				estDuration := estimateRouteDuration(dr.Machines, config)
				if estDuration+int64(svcMin)*60 > dr.ShiftDurationSecs {
					continue
				}
				startDist, _ := router.Distance(m.Lat, m.Lng, dr.StartLat, dr.StartLng)
				endDist, _ := router.Distance(m.Lat, m.Lng, dr.EndLat, dr.EndLng)
				score := 0.6*startDist + 0.4*endDist
				// Priority weighting: urgent machines get a 20% distance bonus
				if m.Urgency == "high" {
					score *= 0.8
				}
				if score < bestScore {
					bestScore = score
					bestMachine = mi
					bestDriver = di
				}
			}
		}

		if bestMachine < 0 {
			break // All assigned or no capacity
		}

		assigned[bestMachine] = true
		driverRoutes[bestDriver].Machines = append(driverRoutes[bestDriver].Machines, machines[bestMachine])
	}
}

func estimateRouteDuration(machines []MachineDemand, config *RouteConfig) int64 {
	if len(machines) < 2 {
		return int64(len(machines)) * int64(config.ServiceMinutes) * 60
	}
	dist := 0.0
	for i := 0; i < len(machines)-1; i++ {
		dist += Haversine(machines[i].Lat, machines[i].Lng, machines[i+1].Lat, machines[i+1].Lng)
	}
	travelSecs := int64(dist / config.AvgSpeedMph * 3600)
	serviceSecs := int64(len(machines)) * int64(config.ServiceMinutes) * 60
	return travelSecs + serviceSecs
}

func insertionCostForDriver(m MachineDemand, dr DriverRoute) float64 {
	if len(dr.Machines) == 0 {
		return Haversine(dr.StartLat, dr.StartLng, m.Lat, m.Lng)
	}
	// Cost of adding m after the last machine
	last := dr.Machines[len(dr.Machines)-1]
	return Haversine(last.Lat, last.Lng, m.Lat, m.Lng)
}

func findTruckForDriver(driver *vend.VendDriver, allTrucks []*vend.VendDeliveryTruck, used map[string]bool) *vend.VendDeliveryTruck {
	for _, truck := range allTrucks {
		if truck.TruckId != driver.TruckId {
			continue
		}
		if truck.Status != vend.VendTruckStatus_VEND_TRUCK_STATUS_ACTIVE {
			continue
		}
		if used[truck.TruckId] {
			continue
		}
		if !licenseCovers(driver.LicenseClass, licenseRequired(truck.Type)) {
			continue
		}
		return truck
	}
	return nil
}

// resolveEndLocation resolves the driver's end-of-day location using the fallback chain:
// 1. schedule.endLocationId → VendLocation coordinates
// 2. driver home address → currentLatitude/currentLongitude
// 3. truck homeDepotId → facility coordinates
// 4. driver start location (round trip)
func resolveEndLocation(driver *vend.VendDriver, sched *vend.VendDriverScheduleDay,
	truck *vend.VendDeliveryTruck, allFacilities []*vend.VendStockingFacility,
	startLat, startLng float64, nic ifs.IVNic) (float64, float64) {

	// 1. Schedule end location
	if sched.EndLocationId != "" {
		lat, lng := resolveLocationCoords(sched.EndLocationId, nic)
		if lat != 0 || lng != 0 {
			return lat, lng
		}
	}
	// 2. Driver home/current position
	if driver.CurrentLatitude != 0 || driver.CurrentLongitude != 0 {
		return driver.CurrentLatitude, driver.CurrentLongitude
	}
	// 3. Truck home depot facility
	if truck.HomeDepotId != "" {
		for _, f := range allFacilities {
			if f.FacilityId == truck.HomeDepotId && f.Coordinates != nil {
				return f.Coordinates.Latitude, f.Coordinates.Longitude
			}
		}
	}
	// 4. Round trip to start
	return startLat, startLng
}

func resolveLocationCoords(locationId string, nic ifs.IVNic) (float64, float64) {
	results, err := vendcommon.GetEntities("Location", byte(10), &vend.VendLocation{LocationId: locationId}, nic)
	if err == nil && len(results) > 0 {
		if loc, ok := results[0].(*vend.VendLocation); ok && loc.Coordinates != nil {
			return loc.Coordinates.Latitude, loc.Coordinates.Longitude
		}
	}
	return 0, 0
}

func findScheduleDay(schedule []*vend.VendDriverScheduleDay, day vend.VendDayOfWeek) *vend.VendDriverScheduleDay {
	for _, s := range schedule {
		if s.Day == day {
			return s
		}
	}
	return nil
}

func resolveStartLocation(driver *vend.VendDriver, sched *vend.VendDriverScheduleDay, nic ifs.IVNic) (float64, float64) {
	if sched.StartLocationId != "" {
		lat, lng := resolveLocationCoords(sched.StartLocationId, nic)
		if lat != 0 || lng != 0 {
			return lat, lng
		}
	}
	if driver.CurrentLatitude != 0 || driver.CurrentLongitude != 0 {
		return driver.CurrentLatitude, driver.CurrentLongitude
	}
	return 0, 0
}

func driverHasSkills(driver *vend.VendDriver, required []string) bool {
	if len(required) == 0 {
		return true
	}
	driverSkills := make(map[string]bool)
	for _, s := range driver.Skills {
		driverSkills[s] = true
	}
	for _, r := range required {
		if !driverSkills[r] {
			return false
		}
	}
	return true
}

func licenseRequired(truckType vend.VendTruckType) vend.VendLicenseClass {
	switch truckType {
	case vend.VendTruckType_VEND_TRUCK_TYPE_BOX_TRUCK,
		vend.VendTruckType_VEND_TRUCK_TYPE_REFRIGERATED:
		return vend.VendLicenseClass_VEND_LICENSE_CLASS_B
	default:
		return vend.VendLicenseClass_VEND_LICENSE_CLASS_C
	}
}

func licenseCovers(driverClass, required vend.VendLicenseClass) bool {
	return driverClass >= required
}

// LoadAllTrucks fetches all delivery trucks.
func LoadAllTrucks(nic ifs.IVNic) []*vend.VendDeliveryTruck {
	results, err := vendcommon.GetEntities(trucks.ServiceName, trucks.ServiceArea, &vend.VendDeliveryTruck{}, nic)
	if err != nil {
		return nil
	}
	out := make([]*vend.VendDeliveryTruck, 0, len(results))
	for _, r := range results {
		if t, ok := r.(*vend.VendDeliveryTruck); ok {
			out = append(out, t)
		}
	}
	return out
}

// LoadAllDrivers fetches all drivers.
func LoadAllDrivers(nic ifs.IVNic) []*vend.VendDriver {
	results, err := vendcommon.GetEntities(drivers.ServiceName, drivers.ServiceArea, &vend.VendDriver{}, nic)
	if err != nil {
		return nil
	}
	out := make([]*vend.VendDriver, 0, len(results))
	for _, r := range results {
		if d, ok := r.(*vend.VendDriver); ok {
			out = append(out, d)
		}
	}
	return out
}

// LoadAllFacilities fetches all stocking facilities.
func LoadAllFacilities(nic ifs.IVNic) []*vend.VendStockingFacility {
	results, err := vendcommon.GetEntities(facilities.ServiceName, facilities.ServiceArea, &vend.VendStockingFacility{}, nic)
	if err != nil {
		return nil
	}
	out := make([]*vend.VendStockingFacility, 0, len(results))
	for _, r := range results {
		if f, ok := r.(*vend.VendStockingFacility); ok {
			out = append(out, f)
		}
	}
	return out
}
