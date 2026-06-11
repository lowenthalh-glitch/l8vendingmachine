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
	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

const (
	profileService     = "MachProf"
	profileServiceArea = byte(10)
)

// prevSnapshotFill stores previous fill% per machine for depletion calculation.
var prevSnapshotFill = make(map[string]int32)
var prevSnapshotTime = make(map[string]int64)

// prevSlotStockForProfile stores previous stock per slot for product depletion.
var prevSlotStockForProfile = make(map[string]int32)

var lastClassifyTime int64

func updateMachineProfile(m *vend.VendFleetMachine, fillPct int32, dailyRevenue int64, nic ifs.IVNic) {
	now := time.Now().Unix()
	dow := int(time.Now().Weekday()) // 0=Sun..6=Sat
	hod := time.Now().Hour()         // 0-23

	// Fetch or create profile
	profile := getOrCreateProfile(m, nic)

	// Ensure repeated arrays are sized
	ensureArraySizes(profile)

	// 1. Depletion rate from fill% change
	prevFill, hasPrev := prevSnapshotFill[m.MachineId]
	prevTime := prevSnapshotTime[m.MachineId]
	if hasPrev && fillPct < prevFill && prevTime > 0 {
		hoursElapsed := float64(now-prevTime) / 3600
		if hoursElapsed > 0 {
			depletionRate := float64(prevFill-fillPct) / hoursElapsed
			profile.DowDepletionRate[dow] = rollingAvg(profile.DowDepletionRate[dow], profile.DowSampleCount[dow], depletionRate)
			profile.DowSampleCount[dow]++
			profile.HodDepletionRate[hod] = rollingAvg(profile.HodDepletionRate[hod], profile.HodSampleCount[hod], depletionRate)
			profile.HodSampleCount[hod]++
			profile.AvgHourlyDepletion = rollingAvg(profile.AvgHourlyDepletion,
				profile.DowSampleCount[0]+profile.DowSampleCount[1]+profile.DowSampleCount[2]+
					profile.DowSampleCount[3]+profile.DowSampleCount[4]+profile.DowSampleCount[5]+profile.DowSampleCount[6],
				depletionRate)
		}
	} else if hasPrev && fillPct > prevFill+10 {
		// Restock event detected
		profile.RestockCount_30D++
	}
	prevSnapshotFill[m.MachineId] = fillPct
	prevSnapshotTime[m.MachineId] = now

	// 2. Revenue tracking
	profile.AvgDailyRevenue = rollingAvg(float64(profile.AvgDailyRevenue), 30, float64(dailyRevenue))
	profile.TotalRevenue_30D += dailyRevenue

	// 3. Fill %
	profile.AvgFillPct = fillPct

	// 4. Per-product depletion
	updateProductProfiles(profile, m, now)

	// 5. Location classification (every 6 hours)
	if now-lastClassifyTime > 6*3600 {
		classifyAllProfiles(profile)
		lastClassifyTime = now
	}

	profile.LastUpdated = now

	// POST or PUT
	existing, _ := vendcommon.GetEntity(profileService, profileServiceArea,
		&vend.VendMachineProfile{ProfileId: m.MachineId}, nic)
	if existing != nil {
		vendcommon.PutEntity(profileService, profileServiceArea, profile, nic)
	} else {
		vendcommon.PostEntity(profileService, profileServiceArea, profile, nic)
	}
}

func getOrCreateProfile(m *vend.VendFleetMachine, nic ifs.IVNic) *vend.VendMachineProfile {
	existing, _ := vendcommon.GetEntity(profileService, profileServiceArea,
		&vend.VendMachineProfile{ProfileId: m.MachineId}, nic)
	if existing != nil {
		return existing.(*vend.VendMachineProfile)
	}
	return &vend.VendMachineProfile{
		ProfileId:           m.MachineId,
		MachineId:           m.MachineId,
		MachineName:         m.Name,
		CascadeThresholdPct: 30,
		TrendMultiplier:     1.0,
	}
}

func ensureArraySizes(p *vend.VendMachineProfile) {
	for len(p.DowDepletionRate) < 7 {
		p.DowDepletionRate = append(p.DowDepletionRate, 0)
	}
	for len(p.DowSampleCount) < 7 {
		p.DowSampleCount = append(p.DowSampleCount, 0)
	}
	for len(p.HodDepletionRate) < 24 {
		p.HodDepletionRate = append(p.HodDepletionRate, 0)
	}
	for len(p.HodSampleCount) < 24 {
		p.HodSampleCount = append(p.HodSampleCount, 0)
	}
}

func rollingAvg(current float64, count int32, newValue float64) float64 {
	if count <= 0 {
		return newValue
	}
	return (current*float64(count) + newValue) / float64(count+1)
}

func updateProductProfiles(profile *vend.VendMachineProfile, m *vend.VendFleetMachine, now int64) {
	productMap := make(map[string]*vend.VendProductProfile)
	for _, pp := range profile.TopProducts {
		productMap[pp.ProductName] = pp
	}

	for _, slot := range m.Inventory {
		if slot.ProductName == "" || slot.Capacity <= 0 {
			continue
		}
		pp := productMap[slot.ProductName]
		if pp == nil {
			pp = &vend.VendProductProfile{
				ProductName: slot.ProductName,
				Price:       getSlotPrice(slot),
				Capacity:    slot.Capacity,
			}
			productMap[slot.ProductName] = pp
		}

		// Track depletion from stock changes
		slotKey := fmt.Sprintf("%s:%d:prof", m.MachineId, slot.SlotNumber)
		if prev, ok := prevSlotStockForProfile[slotKey]; ok && slot.CurrentStock < prev {
			sold := prev - slot.CurrentStock
			// Assume 5-min interval
			rate := float64(sold) / (5.0 / 60.0) // units per hour
			pp.DepletionRatePerHour = rollingAvg(pp.DepletionRatePerHour, 100, rate)
		}
		prevSlotStockForProfile[slotKey] = slot.CurrentStock

		pp.AvgStock = slot.CurrentStock
		pp.Capacity = slot.Capacity
		if pp.DepletionRatePerHour > 0 {
			pp.TimeToEmptyHours = float64(slot.CurrentStock) / pp.DepletionRatePerHour
		}
	}

	// Rebuild top products list
	profile.TopProducts = make([]*vend.VendProductProfile, 0, len(productMap))
	for _, pp := range productMap {
		profile.TopProducts = append(profile.TopProducts, pp)
	}
}

func classifyAllProfiles(profile *vend.VendMachineProfile) {
	if len(profile.DowDepletionRate) < 7 {
		return
	}
	weekdaySum := profile.DowDepletionRate[1] + profile.DowDepletionRate[2] +
		profile.DowDepletionRate[3] + profile.DowDepletionRate[4] + profile.DowDepletionRate[5]
	weekendSum := profile.DowDepletionRate[0] + profile.DowDepletionRate[6]

	weekdayAvg := weekdaySum / 5
	weekendAvg := weekendSum / 2

	if weekdayAvg > 0 {
		profile.WeekendWeekdayRatio = weekendAvg / weekdayAvg
	}

	if profile.WeekendWeekdayRatio > 1.5 {
		profile.LocationClass = "retail"
	} else if profile.WeekendWeekdayRatio < 0.5 {
		profile.LocationClass = "office"
	} else {
		profile.LocationClass = "mixed"
	}
}
