/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 */
package optimizer

import "math"

const (
	EarthRadiusMiles = 3958.8
	DegToRad         = math.Pi / 180.0
)

// RouteLeg represents one leg between two consecutive stops.
// Source-agnostic: populated by haversine or Google Maps.
type RouteLeg struct {
	DistanceMiles   float64
	DurationSeconds int64
	IsReload        bool
}

// RouteMetrics holds computed totals for a route.
type RouteMetrics struct {
	TotalDistanceMiles float64
	TotalDurationSecs  int64
	EstimatedFuelCost  float64
	PlannedArrivals    []int64
}

// Haversine returns the great-circle distance in miles between two GPS points.
func Haversine(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := (lat2 - lat1) * DegToRad
	dLng := (lng2 - lng1) * DegToRad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*DegToRad)*math.Cos(lat2*DegToRad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadiusMiles * c
}

// Centroid returns the geographic center of a set of points.
func Centroid(points [][2]float64) (float64, float64) {
	if len(points) == 0 {
		return 0, 0
	}
	var sumLat, sumLng float64
	for _, p := range points {
		sumLat += p[0]
		sumLng += p[1]
	}
	n := float64(len(points))
	return sumLat / n, sumLng / n
}

// ComputeRouteMetrics calculates totals from a list of legs.
// Used by both haversine and Google Maps paths — single implementation.
func ComputeRouteMetrics(legs []RouteLeg, startTime int64, mpg float64,
	fuelPrice float64, serviceMinutes int32, reloadMinutes int32) RouteMetrics {

	m := RouteMetrics{}
	currentTime := startTime
	m.PlannedArrivals = make([]int64, len(legs))

	for i, leg := range legs {
		currentTime += leg.DurationSeconds
		m.PlannedArrivals[i] = currentTime
		m.TotalDistanceMiles += leg.DistanceMiles

		// Add service or reload time after arriving
		if leg.IsReload {
			currentTime += int64(reloadMinutes) * 60
		} else {
			currentTime += int64(serviceMinutes) * 60
		}
	}

	m.TotalDurationSecs = currentTime - startTime
	if mpg > 0 {
		m.EstimatedFuelCost = (m.TotalDistanceMiles / mpg) * fuelPrice
	}
	return m
}

// HaversineLeg creates a RouteLeg from haversine distance with a fixed avg speed.
func HaversineLeg(lat1, lng1, lat2, lng2 float64, avgSpeedMph float64, isReload bool) RouteLeg {
	dist := Haversine(lat1, lng1, lat2, lng2)
	dur := int64(0)
	if avgSpeedMph > 0 {
		dur = int64(dist / avgSpeedMph * 3600)
	}
	return RouteLeg{DistanceMiles: dist, DurationSeconds: dur, IsReload: isReload}
}
