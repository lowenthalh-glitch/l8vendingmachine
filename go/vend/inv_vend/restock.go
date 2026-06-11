/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 */

package main

import (
	"fmt"
	"time"

	"github.com/saichler/l8types/go/ifs"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/vend/fleet/machines"
	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

const (
	restockService     = "Restock"
	restockServiceArea = byte(10)
)

func computeRestockRecommendations(nic ifs.IVNic) {
	// Wait for profiles to be populated
	time.Sleep(3 * time.Minute)

	for {
		// Load profiles
		profileResults, err := vendcommon.GetEntities(profileService, profileServiceArea,
			&vend.VendMachineProfile{}, nic)
		if err != nil || len(profileResults) == 0 {
			time.Sleep(30 * time.Minute)
			continue
		}
		profiles := make(map[string]*vend.VendMachineProfile)
		for _, elem := range profileResults {
			p, ok := elem.(*vend.VendMachineProfile)
			if ok && p != nil {
				profiles[p.MachineId] = p
			}
		}

		// Load current fleet machines for live stock data
		machineResults, _ := vendcommon.GetEntities(
			machines.ServiceName, machines.ServiceArea,
			&vend.VendFleetMachine{}, nic)

		// Compute revenue ranks for priority adjustment
		revenueRanks := computeRevenueRanks(profiles)

		// Run scenarios per machine
		var allCandidates []*vend.VendRestockRecommendation
		for _, elem := range machineResults {
			m, ok := elem.(*vend.VendFleetMachine)
			if !ok || m == nil {
				continue
			}
			profile := profiles[m.MachineId]
			if profile == nil {
				continue
			}

			// Scenario 1: Day-of-week demand
			if rec := evaluateDayOfWeekDemand(m, profile); rec != nil {
				allCandidates = append(allCandidates, rec)
			}
			// Scenario 3: Fast movers
			if rec := evaluateFastMovers(m, profile); rec != nil {
				allCandidates = append(allCandidates, rec)
			}
			// Scenario 10: Critical threshold prediction
			if rec := evaluateCriticalPrediction(m, profile); rec != nil {
				allCandidates = append(allCandidates, rec)
			}
		}

		// Scenario 9: Revenue-based priority adjustment
		applyRevenuePriority(allCandidates, revenueRanks)

		// Merge: one recommendation per machine (highest priority wins)
		merged := mergeRecommendations(allCandidates)

		// Write recommendations
		now := time.Now().Unix()
		count := 0
		for _, rec := range merged {
			rec.CreatedAt = now
			rec.ExpiresAt = now + 24*3600 // 24 hour expiry
			existing, _ := vendcommon.GetEntity(restockService, restockServiceArea,
				&vend.VendRestockRecommendation{RecommendationId: rec.RecommendationId}, nic)
			if existing != nil {
				vendcommon.PutEntity(restockService, restockServiceArea, rec, nic)
			} else {
				vendcommon.PostEntity(restockService, restockServiceArea, rec, nic)
			}
			count++
		}

		// Expire old recommendations
		expireOldRecommendations(now, nic)

		if count > 0 {
			fmt.Printf("[RESTOCK] Generated %d recommendations from %d machines\n", count, len(machineResults))
		}

		time.Sleep(30 * time.Minute)
	}
}

func mergeRecommendations(candidates []*vend.VendRestockRecommendation) []*vend.VendRestockRecommendation {
	best := make(map[string]*vend.VendRestockRecommendation)
	for _, rec := range candidates {
		existing := best[rec.MachineId]
		if existing == nil || rec.Priority > existing.Priority {
			best[rec.MachineId] = rec
		} else if rec.Priority == existing.Priority {
			// Same priority — combine reasons
			existing.Reason = existing.Reason + " | " + rec.Reason
		}
	}
	result := make([]*vend.VendRestockRecommendation, 0, len(best))
	for _, rec := range best {
		result = append(result, rec)
	}
	return result
}

func expireOldRecommendations(now int64, nic ifs.IVNic) {
	results, err := vendcommon.GetEntities(restockService, restockServiceArea,
		&vend.VendRestockRecommendation{}, nic)
	if err != nil {
		return
	}
	expired := 0
	for _, elem := range results {
		rec, ok := elem.(*vend.VendRestockRecommendation)
		if !ok || rec == nil {
			continue
		}
		if rec.ExpiresAt > 0 && rec.ExpiresAt < now {
			nic.Unicast("", restockService, restockServiceArea, ifs.DELETE,
				&vend.VendRestockRecommendation{RecommendationId: rec.RecommendationId})
			expired++
		}
	}
	if expired > 0 {
		fmt.Printf("[RESTOCK] Expired %d old recommendations\n", expired)
	}
}

func computeRevenueRanks(profiles map[string]*vend.VendMachineProfile) map[string]int32 {
	type entry struct {
		id  string
		rev int64
	}
	sorted := make([]entry, 0, len(profiles))
	for id, p := range profiles {
		sorted = append(sorted, entry{id, p.TotalRevenue_30D})
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].rev > sorted[i].rev {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	ranks := make(map[string]int32)
	for i, e := range sorted {
		percentile := int32(float64(i+1) / float64(len(sorted)) * 100)
		ranks[e.id] = percentile
	}
	return ranks
}
