/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 */
package optimizer

// twoOpt improves stop ordering by swapping edge pairs to reduce total distance.
func twoOpt(stops []MachineDemand) []MachineDemand {
	if len(stops) < 4 {
		return stops
	}

	improved := make([]MachineDemand, len(stops))
	copy(improved, stops)

	changed := true
	for changed {
		changed = false
		for i := 0; i < len(improved)-2; i++ {
			for j := i + 2; j < len(improved); j++ {
				if twoOptGain(improved, i, j) > 0 {
					reverseSegment(improved, i+1, j)
					changed = true
				}
			}
		}
	}
	return improved
}

// twoOptGain returns the distance saved by reversing the segment between i+1 and j.
func twoOptGain(stops []MachineDemand, i, j int) float64 {
	a, b := stops[i], stops[i+1]
	c, d := stops[j], stops[j%len(stops)]
	if j+1 < len(stops) {
		d = stops[j+1]
	} else {
		// Last element — compare with wrap-around cost
		return 0
	}

	oldDist := Haversine(a.Lat, a.Lng, b.Lat, b.Lng) +
		Haversine(c.Lat, c.Lng, d.Lat, d.Lng)
	newDist := Haversine(a.Lat, a.Lng, c.Lat, c.Lng) +
		Haversine(b.Lat, b.Lng, d.Lat, d.Lng)

	return oldDist - newDist
}

func reverseSegment(stops []MachineDemand, i, j int) {
	for i < j {
		stops[i], stops[j] = stops[j], stops[i]
		i++
		j--
	}
}
