# Threshold Events for Vending Machine Inventory

## Context

The Fleet section shows 7 vending machines with per-slot inventory data (stock levels, capacity, status). We need to generate threshold-based events when:
1. **Slot fill %** drops below warning/critical levels (per-slot)
2. **Total machine inventory %** drops below warning/critical levels (aggregate across all slots)

These events follow the L8 Ecosystem architecture: `l8events` for event records, `l8alarms` for alarm lifecycle (correlation, deduplication, auto-clear), and `l8notify` for notification delivery.

## Project Classification

This project is **hybrid**: observation/collection (pipeline from Nayax simulator) + business (Fleet CRUD, alerts management). The threshold evaluator sits at the boundary — it reads observed data from the collection pipeline and creates business events. This is the same pattern probler uses with its alarms service (`../probler/go/prob/alarms/`). The canonical reference for the collection side is probler; for the business/alerts side is l8alarms.

## Prime Object Classification

**VendAlert** is a **Prime Object**: it has independent identity (`AlertId`), its own lifecycle (active → acknowledged → resolved), and users query alerts across all machines ("show all critical alerts"). The existing Alerts service at `/10/Alert` already treats it as a Prime Object with its own CRUD service. Threshold-generated alerts are POSTed to this existing service.

## Duplication Audit

The threshold evaluator uses an **abstracted rule loop** — one `ThresholdRule` struct with config, one `evaluateRule()` function that handles both UPPER and LOWER thresholds. Rules are data (array of ThresholdRule), evaluation is shared behavior. No per-rule if-blocks.

## Threshold Rules

### Per-Slot Thresholds (fill % = currentStock / capacity × 100)

| Rule | Metric | Type | Warning | Critical | Auto-Clear |
|------|--------|------|---------|----------|------------|
| VEND_SLOT_LOW | slot fill % | LOWER | 30% | 10% | Yes (clears at 50%) |
| VEND_SLOT_EMPTY | slot fill % | LOWER | — | 0% | Yes (clears at 1 item) |

### Per-Machine Thresholds (total inventory % = sum(currentStock) / sum(capacity) × 100)

| Rule | Metric | Type | Warning | Critical | Auto-Clear |
|------|--------|------|---------|----------|------------|
| VEND_MACHINE_LOW_STOCK | total inventory % | LOWER | 40% | 20% | Yes (clears at 60%) |
| VEND_MACHINE_MANY_EMPTY | empty slot count | UPPER | 3 slots | 5 slots | Yes (clears at 1) |

## Architecture

Following the L8 Ecosystem flow:

```
VCache→Fleet bridge (inv_vend)
    │
    ├── Update VendFleetMachine with slot data (existing)
    │
    └── Evaluate thresholds on each update
        │
        ├── Per-slot: currentStock / capacity < threshold?
        ├── Per-machine: totalStock / totalCapacity < threshold?
        │
        └── If threshold crossed → Post event to l8events
            │
            └── l8alarms matches event → creates/updates Alarm
                │
                └── NotificationPolicy → l8notify delivery
```

## Phase 1: Threshold Evaluator

**New file:** `go/vend/inv_vend/threshold.go`

**Architecture note:** The L8 ecosystem has no built-in threshold evaluation — l8alarms processes alarms that are POSTed to it, but doesn't generate them. The threshold evaluator is custom project code that bridges observed data to the alarm system.

**Timing:** Runs as a **separate periodic goroutine** in `inv_vend`, NOT inside the VCache→Fleet bridge. The bridge creates VendFleetMachine with basic fields, but slot data arrives later from per-machine polling (parser patches VendFleetMachine asynchronously). The evaluator must read the **fully-populated** VendFleetMachine (with slot data already patched) from the Machine service.

```
Bridge creates VendFleetMachine (basic fields)           ← no slot data yet
    → Per-machine parser patches slot data (async)       ← slot data arrives later
        → Threshold evaluator reads complete record      ← separate timer, reads from Machine service
            → Evaluates slots with actual data
                → POSTs VendAlert to Alerts service
```

**New goroutine in `inv_vend/main.go`:**
```go
go evaluateThresholds(nic)  // separate from bridgeVCacheToFleet
```

Reads VendFleetMachine records from the Machine service (`/10/Machine`) every 5 minutes, evaluates each against threshold rules, and POSTs VendAlert to the existing Alerts service (`/10/Alert`).

```go
type ThresholdRule struct {
    Name           string
    MetricName     string
    ThresholdType  string  // "UPPER" or "LOWER"
    WarningValue   float64
    CriticalValue  float64
    ClearThreshold float64
    AutoClear      bool
    Category       string  // maps to VendAlertCategory
}
```

**Evaluation logic (single abstracted loop):**
```go
for _, machine := range machines {
    for _, rule := range rules {
        value := extractMetric(machine, rule.MetricName)
        severity := evaluateRule(rule, value)
        if severity != "" && !activeAlertExists(machine.MachineId, rule.Name) {
            postAlert(machine, rule, value, severity)
        } else if severity == "" && activeAlertExists(machine.MachineId, rule.Name) && rule.AutoClear {
            clearAlert(machine.MachineId, rule.Name)
        }
    }
}
```

**Metric extraction:**
- Per-slot fill %: `slot.CurrentStock / slot.Capacity × 100` (iterates `machine.Inventory`)
- Per-machine total fill %: `sum(CurrentStock) / sum(Capacity) × 100`
- Empty slot count: count of slots with `CurrentStock == 0`

## Phase 2: Alert Generation

The threshold evaluator POSTs `VendAlert` objects directly to the existing Alerts service (`/10/Alert`) — the same CRUD service that already exists in the Maintenance section. No l8events dependency needed for the initial implementation; l8events integration (for the full event→alarm→notify pipeline) is deferred to Phase 4.

```go
alert := &vend.VendAlert{
    MachineId:      machine.MachineId,
    Severity:       severity,  // WARNING or CRITICAL
    Category:       vend.VendAlertCategory_VEND_ALERT_CATEGORY_INVENTORY,
    Code:           rule.Name,
    Description:    fmt.Sprintf("Slot %d fill at %.0f%% (threshold: %.0f%%)", slotNum, fillPct, rule.WarningValue),
    CurrentValue:   fillPct,
    ThresholdValue: rule.WarningValue,
    Status:         vend.VendAlertStatus_VEND_ALERT_STATUS_ACTIVE,
    SlotId:         fmt.Sprintf("%d", slotNum),
    Timestamp:      time.Now().Unix(),
}
vendcommon.PostEntity("Alert", byte(10), alert, nic)
```

**Auto-clear:** When a threshold recovers, the evaluator queries active alerts for that machine+rule and PATCHes their status to `RESOLVED` with `ResolvedAt = now`.

**Deduplication:** Before POSTing, check if an active alert already exists for `machineId + rule.Name + slotId`. Use `vendcommon.GetEntities` with a filter. Only create if no active alert exists.

## Phase 3: l8events Integration (Optional Enhancement)

**Dependency:** Add `l8events` to go.mod

Optionally post l8events `EventRecord` alongside VendAlert for the full event→alarm→notify pipeline:
```go
l8events.PostEvent(vnic,
    evt.EVENT_CATEGORY_PERFORMANCE,
    "VEND_SLOT_LOW",
    evt.SEVERITY_WARNING,
    machineId, machineName, "VendFleetMachine",
    fmt.Sprintf("Slot %d fill at %.0f%% (threshold: %.0f%%)", slotNum, fillPct, threshold),
    map[string]string{
        "slotNumber": strconv.Itoa(slotNum),
        "currentStock": strconv.Itoa(stock),
        "capacity": strconv.Itoa(capacity),
        "fillPercent": fmt.Sprintf("%.1f", fillPct),
        "threshold": fmt.Sprintf("%.1f", threshold),
    })
```

## Phase 3: Alarm Definitions

**Data:** Create alarm definitions via the mock generator or programmatically:

| Definition | Event Pattern | Threshold Count | Window | Auto-Clear |
|-----------|---------------|-----------------|--------|------------|
| Low Slot Stock | VEND_SLOT_LOW | 1 | — | Yes (clear pattern: VEND_SLOT_CLEARED) |
| Empty Slot | VEND_SLOT_EMPTY | 1 | — | Yes |
| Low Machine Stock | VEND_MACHINE_LOW_STOCK | 1 | — | Yes |
| Many Empty Slots | VEND_MACHINE_MANY_EMPTY | 1 | — | Yes |

Deduplication key: `machineId:slotNumber` (per-slot) or `machineId` (per-machine).

## Phase 4: Alarms Service Activation

**File:** `go/vend/main/main.go`

Add l8alarms service activation (same pattern as probler's `go/prob/alarms/main.go`):
- Activate alarm definitions service
- Activate alarms service
- Activate correlation engine
- Activate notification policies

## Phase 5: UI — Alerts Section

Reuse the existing VendAlert infrastructure (already a Prime Object with CRUD service):
- **Maintenance section** already has Alerts service (`/10/Alert`) with:
  - Config: `maintenance-config.js` (service entry for alerts)
  - Enums: `VendAlertSeverity`, `VendAlertCategory`, `VendAlertStatus` in `vend-common.proto`
  - Columns: `alertId`, `machineId`, `severity`, `category`, `code`, `description`, `status`
  - Forms: alert detail with `currentValue`, `thresholdValue`, `slotId`
  - Type registration: `VendAlert`/`VendAlertList` in `shared.go`
- Threshold-generated alerts POST to this existing service — no new UI files needed
- Dashboard KPI cards can show active alarm count

**Data completeness:** VendAlert proto has 16 fields. Existing columns/forms cover the display fields. Threshold evaluator populates: `alertId` (generated), `machineId`, `severity`, `category` (INVENTORY), `code` (rule name), `description` (human-readable), `currentValue`, `thresholdValue`, `status` (ACTIVE), `slotId` (for per-slot alerts). Remaining fields (`acknowledgedBy`, `acknowledgedAt`, `resolvedAt`) are populated on user action.

## Phase 6: Notification Policies (Optional)

Create notification policies that fire when:
- Critical slot empty → immediate webhook/email
- Warning low stock → batched daily summary
- Machine-wide low stock → escalation chain (driver → supervisor → operations)

## Deferred

- **Mobile parity** — mobile alerts section must mirror desktop (per `mobile-rules.md`). Flag for next iteration.
- **Security/auth** — alarm acknowledge/clear must use `ISecurityProvider` (per `security-provider-interface.md`). Phase 6 must address this.
- l8notify channel configuration (email/webhook/Slack setup)
- Escalation chains with time-based re-notification
- Correlation rules (e.g., multiple empty slots → single "machine needs restock" root cause)
- Maintenance window suppression during scheduled restocking

## Traceability Matrix

| # | Action Item | Phase |
|---|-------------|-------|
| 1 | Create threshold evaluator as separate goroutine in inv_vend | Phase 1 |
| 2 | Define threshold rules for slot fill % and machine total % | Phase 1 |
| 3 | Read VendFleetMachine from Machine service (with slot data) | Phase 1 |
| 4 | Register VendAlert type in inv_vend for serialization | Phase 1 |
| 5 | POST VendAlert to Alerts service when thresholds crossed | Phase 2 |
| 6 | Deduplication: check active alerts before creating new ones | Phase 2 |
| 7 | Auto-clear: PATCH alert status to RESOLVED when threshold recovers | Phase 2 |
| 8 | l8events EventRecord integration (optional) | Phase 3 (optional) |
| 9 | Alarm definitions for l8alarms pipeline | Phase 3 (optional) |
| 10 | Activate l8alarms services in vend_demo | Phase 4 (deferred) |
| 11 | Verify alerts in Maintenance → Alerts table | Phase 5 |
| 12 | Dashboard active alarm count | Phase 5 |
| 13 | Notification policies | Phase 6 (deferred) |
| 14 | Mobile alerts parity | Deferred |
| 15 | Security/auth for alarm acknowledge/clear | Deferred |
| 16 | Verify threshold events fire on low stock | Verify |
| 17 | Verify auto-clear when stock recovers | Verify |
| 18 | Verify deduplication (no duplicate alerts) | Verify |

## Critical Files

| Action | File |
|--------|------|
| Create | `go/vend/inv_vend/threshold.go` |
| Modify | `go/vend/inv_vend/main.go` (add separate evaluator goroutine + register VendAlert type) |
| Modify | `go/vend/main/main.go` (activate l8alarms) |
| Modify | `go.mod` (add l8events, l8alarms dependencies) |
| Modify | `go/tests/mocks/phase_setup.go` (create alarm definitions) |
| Modify | `go/vend/ui/web/vend-ui/` (dashboard alarm count) |
