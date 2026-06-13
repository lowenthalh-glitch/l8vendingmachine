/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 */
package optimizer

// Cluster is a group of machines to be served by one truck in one route.
type Cluster struct {
	Machines    []MachineDemand
	CentroidLat float64
	CentroidLng float64
	TotalDemand map[string]int32 // sku → total units needed
}

// ClusterMachines groups List A machines by proximity, then inserts List B
// machines into existing clusters if the detour is small enough.
func ClusterMachines(listA, listB []MachineDemand, maxDistMiles, maxDetourMiles float64) []Cluster {
	if len(listA) == 0 {
		return nil
	}

	assigned := make([]bool, len(listA))
	var clusters []Cluster

	for {
		// Find the first unassigned machine
		seedIdx := -1
		for i, a := range assigned {
			if !a {
				seedIdx = i
				break
			}
		}
		if seedIdx < 0 {
			break
		}

		cluster := Cluster{TotalDemand: make(map[string]int32)}
		addToCluster(&cluster, listA[seedIdx])
		assigned[seedIdx] = true

		// Grow cluster by adding nearest unassigned machines
		for {
			bestIdx := -1
			bestDist := maxDistMiles + 1

			for i, m := range listA {
				if assigned[i] {
					continue
				}
				dist := Haversine(cluster.CentroidLat, cluster.CentroidLng, m.Lat, m.Lng)
				if dist < bestDist {
					bestDist = dist
					bestIdx = i
				}
			}

			if bestIdx < 0 || bestDist > maxDistMiles {
				break
			}

			addToCluster(&cluster, listA[bestIdx])
			assigned[bestIdx] = true
		}

		clusters = append(clusters, cluster)
	}

	// Insert List B machines into existing clusters if detour is small
	insertListB(clusters, listB, maxDetourMiles)

	return clusters
}

func addToCluster(c *Cluster, m MachineDemand) {
	c.Machines = append(c.Machines, m)
	for sku, qty := range m.Products {
		c.TotalDemand[sku] += qty
	}
	// Recalculate centroid
	points := make([][2]float64, len(c.Machines))
	for i, machine := range c.Machines {
		points[i] = [2]float64{machine.Lat, machine.Lng}
	}
	c.CentroidLat, c.CentroidLng = Centroid(points)
}

// insertListB adds List B machines to the nearest cluster if the insertion
// cost (detour) is below the threshold.
func insertListB(clusters []Cluster, listB []MachineDemand, maxDetourMiles float64) {
	for _, m := range listB {
		bestCluster := -1
		bestDetour := maxDetourMiles + 1

		for ci, c := range clusters {
			if len(c.Machines) < 2 {
				dist := Haversine(c.CentroidLat, c.CentroidLng, m.Lat, m.Lng)
				if dist < bestDetour {
					bestDetour = dist
					bestCluster = ci
				}
				continue
			}
			// Check insertion cost between each pair of consecutive stops
			detour := insertionCost(c.Machines, m)
			if detour < bestDetour {
				bestDetour = detour
				bestCluster = ci
			}
		}

		if bestCluster >= 0 {
			m.Urgency = "low"
			addToCluster(&clusters[bestCluster], m)
		}
	}
}

// insertionCost finds the cheapest place to insert m between existing stops.
func insertionCost(stops []MachineDemand, m MachineDemand) float64 {
	best := Haversine(stops[len(stops)-1].Lat, stops[len(stops)-1].Lng, m.Lat, m.Lng)
	for i := 0; i < len(stops)-1; i++ {
		a := stops[i]
		b := stops[i+1]
		direct := Haversine(a.Lat, a.Lng, b.Lat, b.Lng)
		via := Haversine(a.Lat, a.Lng, m.Lat, m.Lng) + Haversine(m.Lat, m.Lng, b.Lat, b.Lng)
		detour := via - direct
		if detour < best {
			best = detour
		}
	}
	return best
}
