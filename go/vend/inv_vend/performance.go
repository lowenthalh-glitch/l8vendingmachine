/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 */

package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/saichler/l8types/go/ifs"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

const (
	slotPerfService     = "SlotPerf"
	slotPerfServiceArea = byte(10)
)

type slotState struct {
	previousStock int32
	lastSeen      int64
	stockoutStart int64
	totalVended   int32
	stockoutHours float64
}

var slotStates = make(map[string]*slotState)
var lastPerfPost int64

func computeSlotPerformance(fleetMachines []*vend.VendFleetMachine, nic ifs.IVNic) {
	now := time.Now().Unix()

	for _, m := range fleetMachines {
		for _, slot := range m.Inventory {
			if slot.ProductName == "" || slot.Capacity <= 0 {
				continue
			}
			key := fmt.Sprintf("%s:%d", m.MachineId, slot.SlotNumber)
			state := slotStates[key]

			if state == nil {
				// First observation — store baseline
				slotStates[key] = &slotState{
					previousStock: slot.CurrentStock,
					lastSeen:      now,
				}
				continue
			}

			// Calculate vended units (stock decrease = vend, increase = restock)
			delta := state.previousStock - slot.CurrentStock
			if delta > 0 {
				state.totalVended += delta
			}

			// Track stockout time
			if slot.CurrentStock == 0 {
				if state.stockoutStart == 0 {
					state.stockoutStart = now
				}
			} else if state.stockoutStart > 0 {
				hours := float64(now-state.stockoutStart) / 3600
				state.stockoutHours += hours
				state.stockoutStart = 0
			}

			state.previousStock = slot.CurrentStock
			state.lastSeen = now
		}
	}

	// Post performance records every hour (12 cycles × 5 min)
	if lastPerfPost == 0 {
		lastPerfPost = now
		return
	}
	if now-lastPerfPost < 3600 {
		return
	}

	postSlotPerformance(fleetMachines, nic, now)
	lastPerfPost = now
}

func postSlotPerformance(fleetMachines []*vend.VendFleetMachine, nic ifs.IVNic, now int64) {
	periodStart := lastPerfPost
	periodEnd := now
	periodHours := float64(periodEnd-periodStart) / 3600

	type slotPerf struct {
		machineId   string
		slotId      string
		productName string
		vendCount   int32
		velocity    float64
		stockout    float64
	}

	// Collect performance per machine for ranking
	machinePerfs := make(map[string][]*slotPerf)

	for _, m := range fleetMachines {
		for _, slot := range m.Inventory {
			if slot.ProductName == "" || slot.Capacity <= 0 {
				continue
			}
			key := fmt.Sprintf("%s:%d", m.MachineId, slot.SlotNumber)
			state := slotStates[key]
			if state == nil {
				continue
			}

			velocity := float64(0)
			if periodHours > 0 {
				velocity = float64(state.totalVended) / (periodHours / 24)
			}

			// Add pending stockout hours
			stockout := state.stockoutHours
			if state.stockoutStart > 0 {
				stockout += float64(now-state.stockoutStart) / 3600
			}

			sp := &slotPerf{
				machineId:   m.MachineId,
				slotId:      fmt.Sprintf("%d", slot.SlotNumber),
				productName: slot.ProductName,
				vendCount:   state.totalVended,
				velocity:    velocity,
				stockout:    stockout,
			}
			machinePerfs[m.MachineId] = append(machinePerfs[m.MachineId], sp)
		}
	}

	count := 0
	for _, perfs := range machinePerfs {
		// Rank by velocity (highest = rank 1)
		sort.Slice(perfs, func(i, j int) bool {
			return perfs[i].velocity > perfs[j].velocity
		})

		for rank, sp := range perfs {
			perfId := fmt.Sprintf("%s-s%s-%d", sp.machineId, sp.slotId, periodEnd)
			record := &vend.VendSlotPerformance{
				PerformanceId: perfId,
				MachineId:     sp.machineId,
				SlotId:        sp.slotId,
				ProductName:   sp.productName,
				PeriodStart:   periodStart,
				PeriodEnd:     periodEnd,
				VendCount:     sp.vendCount,
				Velocity:      sp.velocity,
				Rank:          int32(rank + 1),
				StockoutHours: sp.stockout,
			}
			vendcommon.PostEntity(slotPerfService, slotPerfServiceArea, record, nic)
			count++
		}
	}

	// Reset accumulators for next period
	for _, state := range slotStates {
		state.totalVended = 0
		state.stockoutHours = 0
		if state.stockoutStart > 0 {
			state.stockoutStart = now
		}
	}

	if count > 0 {
		fmt.Printf("[ANALYTICS] Posted %d slot performance records\n", count)
	}
}
