/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 *
 * Workload balancing: redistributes machines between drivers after initial
 * geographic assignment to equalize stop counts or estimated durations.
 */
package optimizer

// BalanceWorkload redistributes machines between driver routes to equalize workload.
// mode: "stops" = balance by stop count, "duration" = balance by estimated duration.
// This is a POST-PROCESSING step on the initial geographic assignment.
func BalanceWorkload(driverRoutes []DriverRoute, mode string, config *RouteConfig) {
	if mode == "" || len(driverRoutes) < 2 {
		return
	}

	switch mode {
	case "stops":
		balanceByStops(driverRoutes, config)
	case "duration":
		balanceByDuration(driverRoutes, config)
	}
}

func balanceByStops(routes []DriverRoute, config *RouteConfig) {
	for iterations := 0; iterations < 50; iterations++ {
		maxIdx, minIdx := findMaxMinByStops(routes)
		if maxIdx == minIdx {
			return
		}
		diff := len(routes[maxIdx].Machines) - len(routes[minIdx].Machines)
		if diff <= 1 {
			return // balanced enough
		}

		// Move the machine from max-driver that is nearest to min-driver
		moved := moveBestMachine(&routes[maxIdx], &routes[minIdx], config)
		if !moved {
			return
		}
	}
}

func balanceByDuration(routes []DriverRoute, config *RouteConfig) {
	for iterations := 0; iterations < 50; iterations++ {
		maxIdx, minIdx := findMaxMinByDuration(routes, config)
		if maxIdx == minIdx {
			return
		}
		maxDur := estimateRouteDuration(routes[maxIdx].Machines, config)
		minDur := estimateRouteDuration(routes[minIdx].Machines, config)
		if maxDur-minDur < int64(config.ServiceMinutes)*60*2 {
			return // within 2 stops of balance
		}

		moved := moveBestMachine(&routes[maxIdx], &routes[minIdx], config)
		if !moved {
			return
		}
	}
}

func findMaxMinByStops(routes []DriverRoute) (int, int) {
	maxIdx, minIdx := 0, 0
	for i := range routes {
		if len(routes[i].Machines) > len(routes[maxIdx].Machines) {
			maxIdx = i
		}
		if len(routes[i].Machines) < len(routes[minIdx].Machines) {
			minIdx = i
		}
	}
	return maxIdx, minIdx
}

func findMaxMinByDuration(routes []DriverRoute, config *RouteConfig) (int, int) {
	maxIdx, minIdx := 0, 0
	maxDur := estimateRouteDuration(routes[0].Machines, config)
	minDur := maxDur
	for i := 1; i < len(routes); i++ {
		dur := estimateRouteDuration(routes[i].Machines, config)
		if dur > maxDur {
			maxDur = dur
			maxIdx = i
		}
		if dur < minDur {
			minDur = dur
			minIdx = i
		}
	}
	return maxIdx, minIdx
}

// moveBestMachine moves one machine from the overloaded driver to the underloaded driver.
// Picks the machine from the overloaded driver that is nearest to the underloaded driver's
// start/end centroid. Returns false if no move improves the balance or would exceed shift.
func moveBestMachine(from, to *DriverRoute, config *RouteConfig) bool {
	if len(from.Machines) <= 1 {
		return false
	}

	// Check if 'to' driver has capacity
	toDur := estimateRouteDuration(to.Machines, config)
	if toDur+int64(config.ServiceMinutes)*60 > to.ShiftDurationSecs {
		return false
	}

	// Find machine in 'from' nearest to 'to' driver's centroid
	toCenterLat := (to.StartLat + to.EndLat) / 2
	toCenterLng := (to.StartLng + to.EndLng) / 2

	bestIdx := -1
	bestDist := 999999.0
	for i, m := range from.Machines {
		d := Haversine(m.Lat, m.Lng, toCenterLat, toCenterLng)
		if d < bestDist {
			bestDist = d
			bestIdx = i
		}
	}
	if bestIdx < 0 {
		return false
	}

	// Move the machine
	to.Machines = append(to.Machines, from.Machines[bestIdx])
	from.Machines = append(from.Machines[:bestIdx], from.Machines[bestIdx+1:]...)
	return true
}
