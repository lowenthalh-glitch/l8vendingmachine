/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 *
 * Post-build shift enforcement: after all routes are built, check which
 * exceed the driver's shift duration. Move overflow machines to other
 * drivers with capacity. Defer to next day if no driver can take them.
 */
package optimizer

import (
	"github.com/saichler/l8vendingmachine/go/types/vend"
)

// builtDriverRoute pairs a DriverRoute with its built route for post-processing.
type builtDriverRoute struct {
	DR    *DriverRoute
	Built *BuiltRoute
}

// EnforceShiftLimits checks all built routes against driver shift durations.
// Overflow machines are moved to drivers with capacity, or deferred.
// Returns deferred machines that couldn't fit any route (for next-day planning).
func EnforceShiftLimits(driverRoutes []DriverRoute, builtRoutes []*BuiltRoute,
	facilities []*vend.VendStockingFacility, config *RouteConfig, router *Router) []MachineDemand {

	var deferred []MachineDemand

	// Loop until no route exceeds its shift (or no more moves possible)
	for iteration := 0; iteration < 20; iteration++ {
		overIdx := findOverShiftRoute(driverRoutes, builtRoutes)
		if overIdx < 0 {
			break // All routes fit within shift
		}

		// Extract the overflow machine — pick lowest priority first, then farthest from route center
		machine, removed := extractOverflowMachine(&driverRoutes[overIdx])
		if !removed {
			break // Can't remove anything (only 1 machine left)
		}

		// Try to move to another driver with capacity
		moved := false
		for i := range driverRoutes {
			if i == overIdx {
				continue
			}
			// Check if this driver has capacity
			if builtRoutes[i].Metrics.TotalDurationSecs+int64(config.ServiceMinutes)*60 > driverRoutes[i].ShiftDurationSecs {
				continue
			}
			// Add machine to this driver
			driverRoutes[i].Machines = append(driverRoutes[i].Machines, machine)
			// Rebuild both routes
			builtRoutes[overIdx] = BuildRouteForDriver(&driverRoutes[overIdx], facilities, config, router)
			builtRoutes[i] = BuildRouteForDriver(&driverRoutes[i], facilities, config, router)
			moved = true
			break
		}

		if !moved {
			// No driver can take it — defer to next day
			machine.Urgency = "low" // downgrade to "can wait"
			deferred = append(deferred, machine)
			// Rebuild the over-shift route without the removed machine
			builtRoutes[overIdx] = BuildRouteForDriver(&driverRoutes[overIdx], facilities, config, router)
		}
	}

	return deferred
}

// findOverShiftRoute returns the index of the first route that exceeds its driver's shift.
func findOverShiftRoute(driverRoutes []DriverRoute, builtRoutes []*BuiltRoute) int {
	for i, dr := range driverRoutes {
		if dr.ShiftDurationSecs > 0 && builtRoutes[i].Metrics.TotalDurationSecs > dr.ShiftDurationSecs {
			return i
		}
	}
	return -1
}

// extractOverflowMachine removes and returns the best candidate to drop:
// 1. "low" urgency machines first (can wait)
// 2. If all "high", pick the one farthest from route center
func extractOverflowMachine(dr *DriverRoute) (MachineDemand, bool) {
	if len(dr.Machines) <= 1 {
		return MachineDemand{}, false
	}

	// First try to find a "low" urgency machine
	for i := len(dr.Machines) - 1; i >= 0; i-- {
		if dr.Machines[i].Urgency == "low" {
			m := dr.Machines[i]
			dr.Machines = append(dr.Machines[:i], dr.Machines[i+1:]...)
			return m, true
		}
	}

	// All high priority — pick the one farthest from the route centroid
	points := make([][2]float64, len(dr.Machines))
	for i, m := range dr.Machines {
		points[i] = [2]float64{m.Lat, m.Lng}
	}
	cLat, cLng := Centroid(points)

	farthestIdx := 0
	farthestDist := 0.0
	for i, m := range dr.Machines {
		d := Haversine(m.Lat, m.Lng, cLat, cLng)
		if d > farthestDist {
			farthestDist = d
			farthestIdx = i
		}
	}

	m := dr.Machines[farthestIdx]
	dr.Machines = append(dr.Machines[:farthestIdx], dr.Machines[farthestIdx+1:]...)
	return m, true
}
