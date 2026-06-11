/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 */

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/saichler/l8types/go/ifs"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

const DefaultRetentionDays = 30

func getRetentionDays() int {
	if val := os.Getenv("VEND_SNAPSHOT_RETENTION_DAYS"); val != "" {
		if days, err := strconv.Atoi(val); err == nil && days > 0 {
			return days
		}
	}
	return DefaultRetentionDays
}

func cleanOldSnapshots(nic ifs.IVNic) {
	// Wait for first of next month
	for {
		sleepUntilFirstOfMonth()
		cutoff := time.Now().AddDate(0, 0, -getRetentionDays()).Unix()

		results, err := vendcommon.GetEntities(
			snapshotService, snapshotServiceArea,
			&vend.VendInventorySnapshot{}, nic)
		if err != nil {
			fmt.Printf("[RETENTION] Error fetching snapshots: %v\n", err)
			continue
		}

		deleted := 0
		for _, elem := range results {
			snap, ok := elem.(*vend.VendInventorySnapshot)
			if !ok || snap == nil {
				continue
			}
			if snap.Timestamp < cutoff {
				nic.Unicast("", snapshotService, snapshotServiceArea, ifs.DELETE,
					&vend.VendInventorySnapshot{SnapshotId: snap.SnapshotId})
				deleted++
			}
		}

		if deleted > 0 {
			fmt.Printf("[RETENTION] Deleted %d snapshots older than %d days\n", deleted, getRetentionDays())
		}

		// Decay profile sample counts (halve so old data fades)
		decayProfiles(nic)
	}
}

func decayProfiles(nic ifs.IVNic) {
	results, err := vendcommon.GetEntities(profileService, profileServiceArea,
		&vend.VendMachineProfile{}, nic)
	if err != nil {
		return
	}
	count := 0
	for _, elem := range results {
		p, ok := elem.(*vend.VendMachineProfile)
		if !ok || p == nil {
			continue
		}
		for i := range p.DowSampleCount {
			p.DowSampleCount[i] = p.DowSampleCount[i] / 2
		}
		for i := range p.HodSampleCount {
			p.HodSampleCount[i] = p.HodSampleCount[i] / 2
		}
		p.RestockCount_30D = p.RestockCount_30D / 2
		p.TotalRevenue_30D = p.TotalRevenue_30D / 2
		vendcommon.PutEntity(profileService, profileServiceArea, p, nic)
		count++
	}
	if count > 0 {
		fmt.Printf("[RETENTION] Decayed %d machine profiles\n", count)
	}
}

func sleepUntilFirstOfMonth() {
	now := time.Now()
	next := time.Date(now.Year(), now.Month()+1, 1, 2, 0, 0, 0, now.Location())
	time.Sleep(time.Until(next))
}
