/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 */
package optimizer

import (
	"fmt"
	"sync"
	"time"

	l8c "github.com/saichler/l8common/go/common"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/vend/route/routes"
)

const (
	DefaultMaxRouteDist      = 50.0
	DefaultMaxDetour         = 3.0
	DefaultReloadMinutes     = 30
	DefaultAvgSpeedMph       = 25.0
	DefaultBreakAfterMinutes = 240 // 4 hours
	DefaultBreakDuration     = 30
	DefaultServiceMinutes = 20
	DefaultFuelPriceGal   = 3.50
)

// Mutex to serialize concurrent route generation requests
var generateMtx sync.Mutex

// GenerateRoutes orchestrates the full route optimization pipeline.
func GenerateRoutes(nic ifs.IVNic, req *vend.VendRouteOptRequest) ([]*vend.VendRoute, error) {
	generateMtx.Lock()
	defer generateMtx.Unlock()

	config := buildConfig(req)
	maxDetour := DefaultMaxDetour
	if req.MaxDetourDistance > 0 {
		maxDetour = req.MaxDetourDistance
	}

	// Create Router — uses OSRM if available, falls back to haversine
	osrmURL := DefaultOSRMUrl // TODO: read from login.json app.osrmUrl
	router := NewRouter(osrmURL)

	// Step 1: Build demand lists + machine address info
	listA, listB, machineInfo, err := BuildDemandLists(nic)
	if err != nil {
		return nil, fmt.Errorf("failed to build demand lists: %v", err)
	}
	nic.Resources().Logger().Info(fmt.Sprintf("Route optimizer: List A=%d, List B=%d, machines=%d", len(listA), len(listB), len(machineInfo)))
	if len(listA) == 0 {
		return nil, nil
	}

	req.ListACount = int32(len(listA))

	// Load resources
	allTrucks := LoadAllTrucks(nic)
	allDrivers := LoadAllDrivers(nic)
	allFacilities := LoadAllFacilities(nic)

	plannedDate := req.PlannedDate
	if plannedDate == 0 {
		plannedDate = time.Now().Add(24 * time.Hour).Unix()
	}
	dayOfWeek := toDayOfWeek(time.Unix(plannedDate, 0))

	startTime := req.StartTime
	if startTime == 0 {
		startTime = plannedDate
	}

	// Step 2: Assign machines to drivers by geographic proximity (uses Router)
	driverRoutes := AssignMachinesToDrivers(listA, listB, allTrucks, allDrivers,
		allFacilities, dayOfWeek, maxDetour, config, router, nic)

	if len(driverRoutes) == 0 {
		dayName := time.Unix(plannedDate, 0).Weekday().String()
		req.Error = fmt.Sprintf("No drivers available on %s with %d urgent machines.", dayName, len(listA))
		return nil, nil
	}

	var generatedRoutes []*vend.VendRoute
	listBAdded := 0

	for i, dr := range driverRoutes {
		// Step 3: Build route with end-location awareness + facility reloads (uses Router)
		built := BuildRouteForDriver(&dr, allFacilities, config, router)

		// Step 4: Apply traffic statistics to leg durations
		ApplyTrafficToLegs(built.Legs, startTime, config.ServiceMinutes, config.ReloadMinutes)
		built.Metrics = ComputeRouteMetrics(built.Legs, startTime, built.TruckMPG,
			config.FuelPriceGal, config.ServiceMinutes, config.ReloadMinutes)

		// Step 5: Optionally refine with Google Maps real-time traffic
		RefineWithTraffic(built, startTime, config, nic)

		// Step 5: Generate VendRoute
		route := toVendRouteFromDriver(built, &dr, allFacilities, machineInfo, plannedDate, i+1)
		l8c.GenerateID(&route.RouteId)
		vendcommon.PostEntity(routes.ServiceName, routes.ServiceArea, route, nic)

		generatedRoutes = append(generatedRoutes, route)
		req.GeneratedRouteIds = append(req.GeneratedRouteIds, route.RouteId)

		for _, m := range dr.Machines {
			if m.Urgency == "low" {
				listBAdded++
			}
		}
	}

	req.GeneratedCount = int32(len(generatedRoutes))
	req.ListBAdded = int32(listBAdded)

	return generatedRoutes, nil
}

func buildConfig(req *vend.VendRouteOptRequest) *RouteConfig {
	reloadMin := int32(DefaultReloadMinutes)
	if req.ReloadTimeMinutes > 0 {
		reloadMin = req.ReloadTimeMinutes
	}
	breakAfter := int32(DefaultBreakAfterMinutes)
	if req.BreakAfterMinutes > 0 {
		breakAfter = req.BreakAfterMinutes
	}
	breakDur := int32(DefaultBreakDuration)
	if req.BreakDurationMinutes > 0 {
		breakDur = req.BreakDurationMinutes
	}
	return &RouteConfig{
		AvgSpeedMph:          DefaultAvgSpeedMph,
		ServiceMinutes:       DefaultServiceMinutes,
		ReloadMinutes:        reloadMin,
		FuelPriceGal:         DefaultFuelPriceGal,
		BreakAfterMinutes:    breakAfter,
		BreakDurationMinutes: breakDur,
	}
}

func toVendRouteFromDriver(built *BuiltRoute, dr *DriverRoute,
	facilities []*vend.VendStockingFacility, machineInfo map[string]MachineInfo,
	plannedDate int64, seq int) *vend.VendRoute {

	t := time.Unix(plannedDate, 0)
	name := fmt.Sprintf("Route %s-%02d", t.Format("2006-01-02"), seq)

	// Determine primary facility (first reload, or nearest to start)
	facilityId := ""
	for _, s := range built.Stops {
		if s.IsReload {
			facilityId = s.FacilityId
			break
		}
	}
	if facilityId == "" && len(facilities) > 0 {
		best := findNearestOpenFacility(dr.StartLat, dr.StartLng, facilities)
		if best != nil {
			facilityId = best.FacilityId
		}
	}

	stops := make([]*vend.VendRouteStop, len(built.Stops))
	for i, s := range built.Stops {
		arrival := int64(0)
		if i < len(built.Metrics.PlannedArrivals) {
			arrival = built.Metrics.PlannedArrivals[i]
		}
		stopType := "machine"
		urgency := s.Urgency
		machineId := s.MachineId
		stopFacilityId := ""
		if s.IsReload {
			stopType = "reload"
			urgency = "reload"
			machineId = ""
			stopFacilityId = s.FacilityId
		} else if s.IsBreak {
			stopType = "break"
			urgency = "break"
			machineId = ""
		} else if s.IsEnd {
			stopType = "end"
			urgency = "end"
		}
		mName := ""
		mAddr := ""
		mCity := ""
		if s.IsBreak {
			mName = "Driver Break"
		} else if !s.IsReload && !s.IsEnd && !s.IsBreak {
			if mi, ok := machineInfo[machineId]; ok {
				mName = mi.Name
				mAddr = mi.Address
				mCity = mi.City
			}
		} else if s.IsReload {
			for _, f := range facilities {
				if f.FacilityId == stopFacilityId {
					mName = f.Name
					if f.Address != nil {
						mAddr = f.Address.Line1
						mCity = f.Address.City
					}
					break
				}
			}
		} else if s.IsEnd {
			mName = "End of Day"
		}
		stops[i] = &vend.VendRouteStop{
			StopOrder:       int32(i + 1),
			MachineId:       machineId,
			PlannedArrival:  arrival,
			ServiceUrgency:  urgency,
			StopType:        stopType,
			FacilityId:      stopFacilityId,
			MachineName:     mName,
			LocationAddress: mAddr,
			LocationCity:    mCity,
		}
	}

	return &vend.VendRoute{
		Name:              name,
		Status:            vend.VendRouteStatus_VEND_ROUTE_STATUS_PLANNED,
		DriverId:          dr.Driver.DriverId,
		VehicleId:         dr.Truck.TruckId,
		FacilityId:        facilityId,
		PlannedDate:       plannedDate,
		TotalDistance:      built.Metrics.TotalDistanceMiles,
		TotalDuration:     int32(built.Metrics.TotalDurationSecs / 60),
		EstimatedFuelCost: built.Metrics.EstimatedFuelCost,
		Stops:             stops,
	}
}

func toDayOfWeek(t time.Time) vend.VendDayOfWeek {
	switch t.Weekday() {
	case time.Monday:
		return vend.VendDayOfWeek_VEND_DAY_OF_WEEK_MONDAY
	case time.Tuesday:
		return vend.VendDayOfWeek_VEND_DAY_OF_WEEK_TUESDAY
	case time.Wednesday:
		return vend.VendDayOfWeek_VEND_DAY_OF_WEEK_WEDNESDAY
	case time.Thursday:
		return vend.VendDayOfWeek_VEND_DAY_OF_WEEK_THURSDAY
	case time.Friday:
		return vend.VendDayOfWeek_VEND_DAY_OF_WEEK_FRIDAY
	case time.Saturday:
		return vend.VendDayOfWeek_VEND_DAY_OF_WEEK_SATURDAY
	case time.Sunday:
		return vend.VendDayOfWeek_VEND_DAY_OF_WEEK_SUNDAY
	}
	return vend.VendDayOfWeek_VEND_DAY_OF_WEEK_UNSPECIFIED
}
