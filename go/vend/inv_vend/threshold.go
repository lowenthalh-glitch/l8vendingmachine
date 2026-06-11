/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 */

package main

import (
	"fmt"
	"time"

	"github.com/saichler/l8alarms/go/types/alm"
	evt "github.com/saichler/l8events/go/types/l8events"
	"github.com/saichler/l8types/go/ifs"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/vend/fleet/machines"
	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

// thresholdDef references an AlarmDefinition by ID and defines the threshold logic.
type thresholdDef struct {
	definitionId  string
	eventType     string
	clearType     string
	thresholdType string // "LOWER" or "UPPER"
	warningValue  float64
	criticalValue float64
	clearValue    float64
}

var slotLowDef = thresholdDef{
	definitionId:  "VEND-DEF-002",
	eventType:     "VEND_STOCK_LOW",
	clearType:     "VEND_STOCK_NORMAL",
	thresholdType: "LOWER",
	warningValue:  30, criticalValue: 10, clearValue: 50,
}

var machineLowDef = thresholdDef{
	definitionId:  "VEND-DEF-002",
	eventType:     "VEND_STOCK_LOW",
	clearType:     "VEND_STOCK_NORMAL",
	thresholdType: "LOWER",
	warningValue:  40, criticalValue: 20, clearValue: 60,
}

const (
	almAlarmService     = "Alarm"
	almAlarmServiceArea = byte(10)
)

// alertState tracks which machine+slot combinations are currently alerting.
var alertState = make(map[string]bool)

// evaluateThresholds runs as a separate goroutine, periodically reading
// VendFleetMachine records and evaluating thresholds.
func evaluateThresholds(nic ifs.IVNic) {
	time.Sleep(60 * time.Second)

	for {
		results, err := vendcommon.GetEntities(
			machines.ServiceName, machines.ServiceArea,
			&vend.VendFleetMachine{}, nic)
		if err != nil {
			time.Sleep(60 * time.Second)
			continue
		}

		alertCount := 0
		for _, elem := range results {
			machine, ok := elem.(*vend.VendFleetMachine)
			if !ok || len(machine.Inventory) == 0 {
				continue
			}
			alertCount += evaluateMachine(machine, nic)
		}

		if alertCount > 0 {
			fmt.Printf("[THRESHOLD] Processed %d alarm actions from %d machines\n", alertCount, len(results))
		}

		time.Sleep(5 * time.Minute)
	}
}

func evaluateMachine(machine *vend.VendFleetMachine, nic ifs.IVNic) int {
	count := 0

	// Per-slot evaluation
	for _, slot := range machine.Inventory {
		if slot.Capacity <= 0 {
			continue
		}
		fillPct := float64(slot.CurrentStock) / float64(slot.Capacity) * 100
		slotId := fmt.Sprintf("%d", slot.SlotNumber)
		count += checkThreshold(machine.MachineId, machine.Name, slotId, fillPct, slotLowDef, nic)
	}

	// Per-machine total inventory evaluation
	totalStock, totalCapacity := 0, 0
	for _, slot := range machine.Inventory {
		totalStock += int(slot.CurrentStock)
		totalCapacity += int(slot.Capacity)
	}
	if totalCapacity > 0 {
		machineFillPct := float64(totalStock) / float64(totalCapacity) * 100
		count += checkThreshold(machine.MachineId, machine.Name, "", machineFillPct, machineLowDef, nic)
	}

	return count
}

func checkThreshold(machineId, machineName, slotId string, value float64, def thresholdDef, nic ifs.IVNic) int {
	stateKey := fmt.Sprintf("%s:%s:%s", machineId, slotId, def.definitionId)
	severity := evalSeverity(def, value)

	if severity != evt.Severity_SEVERITY_UNSPECIFIED {
		// Threshold crossed — create alarm if not already alerting
		if !alertState[stateKey] {
			postAlarm(machineId, machineName, slotId, value, def, severity, nic)
			alertState[stateKey] = true
			return 1
		}
	} else if alertState[stateKey] {
		// Threshold recovered — clear the alarm
		clearAlarm(machineId, slotId, def, nic)
		alertState[stateKey] = false
		return 1
	}
	return 0
}

func evalSeverity(def thresholdDef, value float64) evt.Severity {
	if def.thresholdType == "LOWER" {
		if value <= def.criticalValue {
			return evt.Severity_SEVERITY_CRITICAL
		}
		if value <= def.warningValue {
			return evt.Severity_SEVERITY_WARNING
		}
	} else {
		if value >= def.criticalValue {
			return evt.Severity_SEVERITY_CRITICAL
		}
		if value >= def.warningValue {
			return evt.Severity_SEVERITY_WARNING
		}
	}
	return evt.Severity_SEVERITY_UNSPECIFIED
}

func postAlarm(machineId, machineName, slotId string, value float64, def thresholdDef, severity evt.Severity, nic ifs.IVNic) {
	alarmId := fmt.Sprintf("THR-%s-%s-%s", machineId, def.definitionId, slotId)
	dedupKey := fmt.Sprintf("%s:%s:%s", machineId, slotId, def.definitionId)

	desc := fmt.Sprintf("%s: %.0f%% (warning=%.0f%%, critical=%.0f%%)",
		def.eventType, value, def.warningValue, def.criticalValue)
	if slotId != "" {
		desc = fmt.Sprintf("Slot %s %s: %.0f%% fill", slotId, def.eventType, value)
	}

	alarm := &alm.Alarm{
		AlarmId:      alarmId,
		DefinitionId: def.definitionId,
		NodeId:       machineId,
		NodeName:     machineName,
		Name:         def.eventType,
		State:        evt.AlarmState_ALARM_STATE_ACTIVE,
		Severity:     severity,
		DedupKey:     dedupKey,
		Description:  desc,
	}

	_, err := vendcommon.PostEntity(almAlarmService, almAlarmServiceArea, alarm, nic)
	if err != nil {
		fmt.Printf("[THRESHOLD] Error posting alarm: %v\n", err)
	}
}

func clearAlarm(machineId, slotId string, def thresholdDef, nic ifs.IVNic) {
	alarmId := fmt.Sprintf("THR-%s-%s-%s", machineId, def.definitionId, slotId)

	// Fetch existing alarm to verify it's still active
	existing, err := vendcommon.GetEntity(almAlarmService, almAlarmServiceArea,
		&alm.Alarm{AlarmId: alarmId}, nic)
	if err != nil || existing == nil {
		return
	}
	alarm := existing.(*alm.Alarm)
	if alarm.State != evt.AlarmState_ALARM_STATE_ACTIVE {
		return
	}

	alarm.State = evt.AlarmState_ALARM_STATE_CLEARED
	alarm.ClearedAt = time.Now().Unix()
	vendcommon.PutEntity(almAlarmService, almAlarmServiceArea, alarm, nic)
}
