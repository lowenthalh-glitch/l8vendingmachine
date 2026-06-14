/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 *
 * Statistical traffic multipliers based on time-of-day and day-of-week.
 * Applied to OSRM base durations to approximate real traffic conditions.
 */
package optimizer

import "time"

// Default time-of-day multipliers (index = hour 0-23)
var hourMultipliers = [24]float64{
	0.8, 0.8, 0.8, 0.8, 0.8, 0.9, // 00-05: night/early
	1.0, 1.4, 1.4, 1.1, 1.1, 1.1, // 06-11: morning rush then mid-morning
	1.2, 1.2, 1.1, 1.1, 1.4, 1.4, // 12-17: lunch then afternoon rush
	1.3, 0.9, 0.9, 0.8, 0.8, 0.8, // 18-23: evening then night
}

// Day-of-week factors (Sunday=0 .. Saturday=6)
var dayFactors = [7]float64{
	0.6, // Sunday
	1.0, // Monday
	1.0, // Tuesday
	1.0, // Wednesday
	1.0, // Thursday
	1.0, // Friday
	0.7, // Saturday
}

// ApplyTrafficStats adjusts a base duration (from OSRM) using statistical
// traffic multipliers for the given arrival time.
func ApplyTrafficStats(baseDurationSecs int64, arrivalTime int64) int64 {
	t := time.Unix(arrivalTime, 0)
	hour := t.Hour()
	day := int(t.Weekday())

	multiplier := hourMultipliers[hour] * dayFactors[day]
	return int64(float64(baseDurationSecs) * multiplier)
}

// ApplyTrafficToLegs adjusts all leg durations in-place using traffic stats.
// startTime is the route departure time; each leg's arrival is computed
// cumulatively so later legs get the correct time-of-day multiplier.
func ApplyTrafficToLegs(legs []RouteLeg, startTime int64, serviceMinutes int32, reloadMinutes int32) {
	currentTime := startTime
	for i := range legs {
		adjusted := ApplyTrafficStats(legs[i].DurationSeconds, currentTime)
		legs[i].DurationSeconds = adjusted
		currentTime += adjusted
		if legs[i].IsReload {
			currentTime += int64(reloadMinutes) * 60
		} else {
			currentTime += int64(serviceMinutes) * 60
		}
	}
}
