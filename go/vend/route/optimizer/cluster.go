/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 */
package optimizer

const (
	MaxRouteDurationSecs = 8 * 3600 // 8 hours max per route
)

// Cluster is a group of machines to be served by one truck in one route.
type Cluster struct {
	Machines    []MachineDemand
	CentroidLat float64
	CentroidLng float64
	TotalDemand map[string]int32 // sku → total units needed
}

// ClusterMachines groups List A machines by proximity, then inserts List B
// machines into existing clusters if the detour is small enough.
// Clusters are capped at 8 hours estimated duration.
func ClusterMachines(listA, listB []MachineDemand, maxDistMiles, maxDetourMiles float64,
	avgSpeedMph float64, serviceMinutes int32) []Cluster {

	if len(listA) == 0 {
		return nil
	}

	assigned := make([]bool, len(listA))
	var clusters []Cluster

	for {
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
		clusterDist := 0.0

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

			// Check if adding this stop would exceed 8 hours
			newDist := clusterDist + bestDist
			nStops := int32(len(cluster.Machines)) + 1
			estSecs := int64(newDist/avgSpeedMph*3600) + int64(nStops)*int64(serviceMinutes)*60
			if estSecs > MaxRouteDurationSecs {
				break
			}

			addToCluster(&cluster, listA[bestIdx])
			assigned[bestIdx] = true
			clusterDist = newDist
		}

		clusters = append(clusters, cluster)
	}

	insertListB(clusters, listB, maxDetourMiles, avgSpeedMph, serviceMinutes)

	return clusters
}

func addToCluster(c *Cluster, m MachineDemand) {
	c.Machines = append(c.Machines, m)
	for sku, qty := range m.Products {
		c.TotalDemand[sku] += qty
	}
	points := make([][2]float64, len(c.Machines))
	for i, machine := range c.Machines {
		points[i] = [2]float64{machine.Lat, machine.Lng}
	}
	c.CentroidLat, c.CentroidLng = Centroid(points)
}

func insertListB(clusters []Cluster, listB []MachineDemand, maxDetourMiles float64,
	avgSpeedMph float64, serviceMinutes int32) {

	for _, m := range listB {
		bestCluster := -1
		bestDetour := maxDetourMiles + 1

		for ci, c := range clusters {
			nStops := int32(len(c.Machines))
			estDist := estimateClusterDistance(c.Machines)
			estSecs := int64(estDist/avgSpeedMph*3600) + int64(nStops)*int64(serviceMinutes)*60
			if estSecs+int64(serviceMinutes)*60 > MaxRouteDurationSecs {
				continue
			}

			if len(c.Machines) < 2 {
				dist := Haversine(c.CentroidLat, c.CentroidLng, m.Lat, m.Lng)
				if dist < bestDetour {
					bestDetour = dist
					bestCluster = ci
				}
				continue
			}
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

func estimateClusterDistance(machines []MachineDemand) float64 {
	if len(machines) < 2 {
		return 0
	}
	total := 0.0
	for i := 0; i < len(machines)-1; i++ {
		total += Haversine(machines[i].Lat, machines[i].Lng, machines[i+1].Lat, machines[i+1].Lng)
	}
	return total
}

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
