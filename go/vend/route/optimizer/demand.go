/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 */
package optimizer

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/vend/fleet/machines"
)

const (
	LowStockThresholdPct = 30 // Below this % of capacity → List B
)

// MachineDemand represents a machine that needs restocking.
type MachineDemand struct {
	MachineId string
	Lat       float64
	Lng       float64
	Products  map[string]int32 // sku → quantity needed to fill
	Urgency   string           // "high" (List A) or "low" (List B)
}

// BuildDemandLists reads fleet machines and classifies them into
// List A (needs restock now) and List B (can wait 1 day).
func BuildDemandLists(nic ifs.IVNic) ([]MachineDemand, []MachineDemand, error) {
	results, err := vendcommon.GetEntities(machines.ServiceName, machines.ServiceArea, &vend.VendFleetMachine{}, nic)
	if err != nil {
		return nil, nil, err
	}

	var listA, listB []MachineDemand

	for _, elem := range results {
		fm, ok := elem.(*vend.VendFleetMachine)
		if !ok || fm.LocationLat == 0 && fm.LocationLng == 0 {
			continue
		}
		if fm.Status == "offline" || fm.Status == "decommissioned" {
			continue
		}

		demand := buildMachineDemand(fm)
		if demand == nil {
			continue
		}

		if fm.EmptySlots > 0 {
			demand.Urgency = "high"
			listA = append(listA, *demand)
		} else if fm.LowStockSlots > 0 {
			demand.Urgency = "low"
			listB = append(listB, *demand)
		}
	}

	return listA, listB, nil
}

func buildMachineDemand(fm *vend.VendFleetMachine) *MachineDemand {
	products := make(map[string]int32)
	for _, slot := range fm.Inventory {
		if slot.Capacity <= 0 {
			continue
		}
		needed := slot.Capacity - slot.CurrentStock
		if needed <= 0 {
			continue
		}
		fillPct := float64(slot.CurrentStock) / float64(slot.Capacity) * 100
		if fillPct < float64(LowStockThresholdPct) || slot.CurrentStock == 0 {
			sku := slot.Sku
			if sku == "" {
				sku = slot.ProductName
			}
			if sku != "" {
				products[sku] += needed
			}
		}
	}
	if len(products) == 0 {
		return nil
	}
	return &MachineDemand{
		MachineId: fm.MachineId,
		Lat:       fm.LocationLat,
		Lng:       fm.LocationLng,
		Products:  products,
	}
}
