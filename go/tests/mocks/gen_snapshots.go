package mocks

import (
	"fmt"
	"math/rand"
	"time"

	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

// generateHistoricalSnapshots creates 30 days of inventory snapshots for a set of machines.
// Simulates realistic depletion (stock decreases over days) with periodic restocks.
func generateHistoricalSnapshots(machineCount int) []*vend.VendInventorySnapshot {
	snapshots := make([]*vend.VendInventorySnapshot, 0)
	now := time.Now()
	startTime := now.AddDate(0, 0, -30)

	// Snapshot every 2 hours (12 per day × 30 days = 360 per machine)
	interval := 2 * time.Hour

	for i := 0; i < machineCount; i++ {
		machineId := fmt.Sprintf("M-%06d", 400000+i)
		machineName := machineNames[i%len(machineNames)]
		totalSlots := int32(10 + rand.Intn(8)) // 10-17 slots
		totalCapacity := totalSlots * int32(8+rand.Intn(5)) // 8-12 per slot

		// Start near full
		currentStock := totalCapacity - int32(rand.Intn(int(totalCapacity/10)))
		restockDay := 3 + rand.Intn(4) // restock every 3-6 days
		daysSinceRestock := 0

		// Average price per item for this machine (cents)
		avgPrice := int64(175 + rand.Intn(150)) // $1.75 - $3.25

		var dailyRevenue int64
		var currentDay int64

		t := startTime
		for t.Before(now) {
			daysSinceRestock++
			ts := t.Unix()
			day := (ts / 86400) * 86400

			// Reset daily revenue at midnight
			if day != currentDay {
				dailyRevenue = 0
				currentDay = day
			}

			// Depletion: sell 3-8% of remaining stock per 2-hour interval
			depleteRate := 0.03 + rand.Float64()*0.05
			sold := int32(float64(currentStock) * depleteRate)
			if sold < 1 && currentStock > 0 {
				sold = 1
			}
			// Accumulate daily revenue from items sold
			dailyRevenue += int64(sold) * avgPrice

			currentStock -= sold
			if currentStock < 0 {
				currentStock = 0
			}

			// Restock check
			if daysSinceRestock >= restockDay*12 { // restockDay in 2-hour intervals
				currentStock = totalCapacity - int32(rand.Intn(int(totalCapacity/10)))
				daysSinceRestock = 0
			}

			fillPct := int32(0)
			if totalCapacity > 0 {
				fillPct = int32(float64(currentStock) / float64(totalCapacity) * 100)
			}

			emptySlots := int32(0)
			lowStockSlots := int32(0)
			if fillPct < 10 {
				emptySlots = totalSlots / 3
				lowStockSlots = totalSlots / 2
			} else if fillPct < 30 {
				lowStockSlots = totalSlots / 3
			}

			revenue := int64(totalCapacity-currentStock) * avgPrice

			snapshots = append(snapshots, &vend.VendInventorySnapshot{
				SnapshotId:    fmt.Sprintf("%s-%d", machineId, ts),
				MachineId:     machineId,
				MachineName:   machineName,
				Timestamp:     ts,
				TotalStock:    currentStock,
				TotalCapacity: totalCapacity,
				FillPct:       fillPct,
				EmptySlots:    emptySlots,
				LowStockSlots: lowStockSlots,
				TotalSlots:    totalSlots,
				Revenue:       revenue,
				DailyRevenue:  dailyRevenue,
			})

			t = t.Add(interval)
		}
	}

	return snapshots
}

var machineNames = []string{
	"Lobby Snacks", "Break Room Drinks", "Cafeteria Main", "Floor 2 East",
	"Main Entrance", "Gym Refreshments", "Library Corner", "Parking Garage",
	"Pool Area", "Office Wing A", "Office Wing B", "Visitor Center",
	"Hospital Lobby", "Airport Terminal", "Hotel Lobby", "School Cafe",
	"Train Station", "Bus Terminal", "Mall Food Court", "Stadium Level 1",
}
