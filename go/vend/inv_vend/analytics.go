/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 */

package main

import (
	"fmt"
	"time"

	l8common "github.com/saichler/l8common/go/types/l8common"
	"github.com/saichler/l8types/go/ifs"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/vend/fleet/machines"
	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

// productPrices maps product names to their price in cents (from simulator data).
// Used as fallback when slot.Price is not populated by the parser.
var productPrices = map[string]int64{
	"Coca-Cola 12oz":       175,
	"Pepsi 12oz":           175,
	"Snickers Bar":         200,
	"Lay's Classic Chips":  225,
	"Monster Energy 16oz":  350,
	"Bottled Water 16oz":   150,
	"KitKat Bar":           200,
	"Doritos Nacho":        225,
	"Red Bull 8.4oz":       325,
	"Gatorade Fruit Punch": 250,
	"M&M's Peanut":         200,
	"Pringles Original":    275,
	"Tropicana OJ 10oz":    225,
	"Reese's Cups":         200,
	"Smartwater 20oz":      250,
	"Nature Valley Granola": 175,
}

func getSlotPrice(slot *vend.VendMachineSlot) int64 {
	if slot.Price > 0 {
		return slot.Price
	}
	if p, ok := productPrices[slot.ProductName]; ok {
		return p
	}
	return 200 // default $2.00
}

const (
	fleetInvService     = "FleetInv"
	fleetInvServiceArea = byte(10)
	forecastService     = "Forecast"
	forecastServiceArea = byte(10)
)

const (
	snapshotService     = "InvSnap"
	snapshotServiceArea = byte(10)
	topPerfService      = "TopPerf"
	topPerfServiceArea  = byte(10)
)

func computeAnalytics(nic ifs.IVNic) {
	time.Sleep(90 * time.Second)

	cycles := 0
	for {
		fleetMachines := fetchFleetMachines(nic)
		if len(fleetMachines) > 0 {
			computeFleetInventory(fleetMachines, nic)
			writeInventorySnapshots(fleetMachines, nic)
			computeSlotPerformance(fleetMachines, nic)

			computeTopPerformers(nic)

			// Generate forecasts every 12 cycles (hourly at 5-min intervals)
			if cycles%12 == 0 {
				computeForecasts(fleetMachines, nic)
			}
		}
		cycles++
		time.Sleep(5 * time.Minute)
	}
}

// prevStockLevels tracks previous stock per machine+slot for delta revenue calculation.
var prevStockLevels = make(map[string]int32)

// dailyRevenueAccum accumulates revenue per machine for the current day.
var dailyRevenueAccum = make(map[string]int64)
var lastRevenueDay int64

func writeInventorySnapshots(fleetMachines []*vend.VendFleetMachine, nic ifs.IVNic) {
	now := time.Now().Unix()
	today := (now / 86400) * 86400

	// Reset daily accumulators at midnight
	if lastRevenueDay != today {
		dailyRevenueAccum = make(map[string]int64)
		lastRevenueDay = today
	}

	count := 0
	for _, m := range fleetMachines {
		totalStock, totalCapacity := int32(0), int32(0)
		var revenue int64
		for _, slot := range m.Inventory {
			totalStock += slot.CurrentStock
			totalCapacity += slot.Capacity
			sold := slot.Capacity - slot.CurrentStock
			if sold > 0 {
				revenue += int64(sold) * getSlotPrice(slot)
			}

			// Delta revenue: stock decrease = items sold (ignore restocks which increase stock)
			slotKey := fmt.Sprintf("%s:%d", m.MachineId, slot.SlotNumber)
			if prev, ok := prevStockLevels[slotKey]; ok && slot.CurrentStock < prev {
				itemsSold := prev - slot.CurrentStock
				dailyRevenueAccum[m.MachineId] += int64(itemsSold) * getSlotPrice(slot)
			}
			prevStockLevels[slotKey] = slot.CurrentStock
		}
		if totalCapacity == 0 {
			continue
		}
		fillPct := int32(float64(totalStock) / float64(totalCapacity) * 100)
		snapshot := &vend.VendInventorySnapshot{
			SnapshotId:    fmt.Sprintf("%s-%d", m.MachineId, now),
			MachineId:     m.MachineId,
			MachineName:   m.Name,
			Timestamp:     now,
			TotalStock:    totalStock,
			TotalCapacity: totalCapacity,
			FillPct:       fillPct,
			EmptySlots:    m.EmptySlots,
			LowStockSlots: m.LowStockSlots,
			TotalSlots:    m.TotalSlots,
			Revenue:       revenue,
			DailyRevenue:  dailyRevenueAccum[m.MachineId],
		}
		vendcommon.PostEntity(snapshotService, snapshotServiceArea, snapshot, nic)

		// Update machine profile incrementally
		updateMachineProfile(m, fillPct, dailyRevenueAccum[m.MachineId], nic)
		count++
	}
	if count > 0 {
		fmt.Printf("[ANALYTICS] Wrote %d inventory snapshots + updated profiles\n", count)
	}
}

func fetchFleetMachines(nic ifs.IVNic) []*vend.VendFleetMachine {
	results, err := vendcommon.GetEntities(
		machines.ServiceName, machines.ServiceArea,
		&vend.VendFleetMachine{}, nic)
	if err != nil {
		return nil
	}
	out := make([]*vend.VendFleetMachine, 0, len(results))
	for _, elem := range results {
		m, ok := elem.(*vend.VendFleetMachine)
		if ok && m != nil {
			out = append(out, m)
		}
	}
	return out
}

func computeFleetInventory(fleetMachines []*vend.VendFleetMachine, nic ifs.IVNic) {
	type productAgg struct {
		summary        *vend.VendFleetInventory
		machinesSeen   map[string]bool
		soldOutMachines map[string]bool
		lowStockMachines map[string]bool
	}

	products := make(map[string]*productAgg)

	for _, m := range fleetMachines {
		for _, slot := range m.Inventory {
			if slot.ProductName == "" {
				continue
			}
			key := slot.ProductName
			agg := products[key]
			if agg == nil {
				agg = &productAgg{
					summary: &vend.VendFleetInventory{
						SummaryId:   key,
						ProductName: slot.ProductName,
					},
					machinesSeen:    make(map[string]bool),
					soldOutMachines: make(map[string]bool),
					lowStockMachines: make(map[string]bool),
				}
				products[key] = agg
			}
			agg.summary.TotalSlots++
			agg.summary.TotalCapacity += slot.Capacity
			agg.summary.TotalUnitsInMachines += slot.CurrentStock
			if agg.summary.UnitPrice == nil {
				agg.summary.UnitPrice = &l8common.Money{
					Amount:     getSlotPrice(slot),
					CurrencyId: "USD",
				}
			}
			agg.machinesSeen[m.MachineId] = true

			if slot.CurrentStock == 0 {
				agg.soldOutMachines[m.MachineId] = true
			} else if slot.Capacity > 0 && float64(slot.CurrentStock)/float64(slot.Capacity) < 0.3 {
				agg.lowStockMachines[m.MachineId] = true
			}
		}
	}

	count := 0
	now := time.Now().Unix()
	for _, agg := range products {
		agg.summary.TotalMachines = int32(len(agg.machinesSeen))
		agg.summary.FleetSoldOutCount = int32(len(agg.soldOutMachines))
		agg.summary.FleetLowStockCount = int32(len(agg.lowStockMachines))
		agg.summary.LastUpdated = now

		existing, _ := vendcommon.GetEntity(fleetInvService, fleetInvServiceArea,
			&vend.VendFleetInventory{SummaryId: agg.summary.SummaryId}, nic)
		if existing != nil {
			vendcommon.PutEntity(fleetInvService, fleetInvServiceArea, agg.summary, nic)
		} else {
			vendcommon.PostEntity(fleetInvService, fleetInvServiceArea, agg.summary, nic)
		}
		count++
	}

	if count > 0 {
		fmt.Printf("[ANALYTICS] Updated %d fleet inventory summaries\n", count)
	}
}

func computeForecasts(fleetMachines []*vend.VendFleetMachine, nic ifs.IVNic) {
	// Load profiles for accurate depletion rates
	profileResults, _ := vendcommon.GetEntities(profileService, profileServiceArea,
		&vend.VendMachineProfile{}, nic)
	profileMap := make(map[string]*vend.VendMachineProfile)
	for _, elem := range profileResults {
		p, ok := elem.(*vend.VendMachineProfile)
		if ok && p != nil {
			profileMap[p.MachineId] = p
		}
	}

	now := time.Now().Unix()
	count := 0

	for _, m := range fleetMachines {
		profile := profileMap[m.MachineId]
		for _, slot := range m.Inventory {
			if slot.ProductName == "" || slot.Capacity <= 0 {
				continue
			}

			// Use profile's per-product depletion rate if available
			velocity := float64(0)
			confidence := 0.5
			if profile != nil {
				pp := findProductProfile(profile, slot.ProductName)
				if pp != nil && pp.DepletionRatePerHour > 0 {
					velocity = pp.DepletionRatePerHour * 24 // daily
					confidence = 0.85
				}
			}
			// Fallback to estimate if no profile data
			if velocity < 0.1 {
				consumed := slot.Capacity - slot.CurrentStock
				if consumed <= 0 {
					continue
				}
				velocity = float64(consumed) / 7.0
				if velocity < 0.1 {
					velocity = 0.1
				}
			}

			var stockoutTime int64
			if slot.CurrentStock > 0 {
				daysUntilEmpty := float64(slot.CurrentStock) / velocity
				stockoutTime = now + int64(daysUntilEmpty*86400)
			} else {
				stockoutTime = now
			}

			urgency := "LOW"
			hoursUntil := float64(stockoutTime-now) / 3600
			if hoursUntil <= 24 {
				urgency = "HIGH"
			} else if hoursUntil <= 72 {
				urgency = "MEDIUM"
			}

			forecastId := fmt.Sprintf("%s-slot%d", m.MachineId, slot.SlotNumber)
			forecast := &vend.VendForecast{
				ForecastId:            forecastId,
				MachineId:             m.MachineId,
				ProductId:             slot.ProductName,
				ForecastDate:          now,
				HorizonDays:           7,
				PredictedDailyVends:   velocity,
				PredictedStockoutTime: stockoutTime,
				RestockUrgency:        urgency,
				ConfidenceScore:       confidence,
			}

			existing, _ := vendcommon.GetEntity(forecastService, forecastServiceArea,
				&vend.VendForecast{ForecastId: forecastId}, nic)
			if existing != nil {
				vendcommon.PutEntity(forecastService, forecastServiceArea, forecast, nic)
			} else {
				vendcommon.PostEntity(forecastService, forecastServiceArea, forecast, nic)
			}
			count++
		}
	}

	if count > 0 {
		fmt.Printf("[ANALYTICS] Updated %d forecasts\n", count)
	}
}

func findProductProfile(profile *vend.VendMachineProfile, productName string) *vend.VendProductProfile {
	for _, pp := range profile.TopProducts {
		if pp.ProductName == productName {
			return pp
		}
	}
	return nil
}

func computeTopPerformers(nic ifs.IVNic) {
	// Read profiles (lightweight — one per machine)
	results, err := vendcommon.GetEntities(profileService, profileServiceArea,
		&vend.VendMachineProfile{}, nic)
	if err != nil || len(results) == 0 {
		return
	}

	// Sort by total revenue descending
	type ranked struct {
		profile *vend.VendMachineProfile
	}
	sorted := make([]ranked, 0, len(results))
	for _, elem := range results {
		p, ok := elem.(*vend.VendMachineProfile)
		if ok && p != nil {
			sorted = append(sorted, ranked{p})
		}
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].profile.TotalRevenue_30D > sorted[i].profile.TotalRevenue_30D {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	now := time.Now().Unix()
	count := 0
	for i, r := range sorted {
		perf := &vend.VendTopPerformer{
			PerformerId:    r.profile.MachineId,
			MachineId:      r.profile.MachineId,
			MachineName:    r.profile.MachineName,
			Revenue_30D:    r.profile.TotalRevenue_30D,
			AvgFillPct:     r.profile.AvgFillPct,
			TotalSnapshots: r.profile.DowSampleCount[0] + r.profile.DowSampleCount[1] + r.profile.DowSampleCount[2] + r.profile.DowSampleCount[3] + r.profile.DowSampleCount[4] + r.profile.DowSampleCount[5] + r.profile.DowSampleCount[6],
			Rank:           int32(i + 1),
			LastUpdated:    now,
		}
		existing, _ := vendcommon.GetEntity(topPerfService, topPerfServiceArea,
			&vend.VendTopPerformer{PerformerId: r.profile.MachineId}, nic)
		if existing != nil {
			vendcommon.PutEntity(topPerfService, topPerfServiceArea, perf, nic)
		} else {
			vendcommon.PostEntity(topPerfService, topPerfServiceArea, perf, nic)
		}
		count++
	}

	if count > 0 {
		fmt.Printf("[ANALYTICS] Updated %d top performers from profiles\n", count)
	}
}
