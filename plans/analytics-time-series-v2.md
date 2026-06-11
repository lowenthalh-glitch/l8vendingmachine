# Analytics v2: Real Time-Series Inventory History + Move Product Summary to Fleet

## Context

The current `VendFleetInventory` implementation is a current-state product summary (aggregated now), not time-series analytics. It shows "product X has Y units across Z machines right now" — useful, but belongs under Fleet, not Analytics.

True analytics shows **trends over time**: fill % declining between restocks, rising after restocks, depletion rate patterns, day-of-week seasonality.

## What Changes

1. **Move product summary** (`VendFleetInventory`) from Analytics to Fleet as a "Products" sub-tab
2. **Create time-series type** `VendInventorySnapshot` — one record per machine per 5-min cycle (never overwritten)
3. **Analytics shows historical chart** — line graph of fill % over time per machine, with visible depletion/restock patterns

## Traceability Matrix

| # | Section | Action Item | Phase |
|---|---------|-------------|-------|
| 1 | UI Fleet | Add Products service to Fleet config + section config | Phase 1 |
| 2 | UI Analytics | Remove fleet-inventory from Analytics config | Phase 1 |
| 3 | Proto | Create VendInventorySnapshot + List message | Phase 2 |
| 4 | Proto | Run make-bindings.sh | Phase 2 |
| 5 | Backend | Create snapshot service (InvSnap, area 10) | Phase 3 |
| 6 | Backend | Add snapshot service to activate_analytics.go | Phase 3 |
| 7 | Backend | Write inventory snapshots every 5 min (POST, never PUT) | Phase 4 |
| 8 | Backend | Keep computeFleetInventory for product summary | Phase 4 |
| 9 | UI Desktop | Add VendInventorySnapshot to analytics config | Phase 5 |
| 10 | UI Desktop | Add VendInventorySnapshot columns | Phase 5 |
| 11 | UI Desktop | Add VendInventorySnapshot form (all fields, readOnly) | Phase 5 |
| 12 | Backend | Register VendInventorySnapshot types in UI shared + inv_vend | Phase 6 |
| 13 | Backend | Add monthly data retention cleanup (configurable, default 30 days) | Phase 7 |
| 14 | UI Mobile | Add VendInventorySnapshot columns to mobile | Phase 5 |
| 15 | UI Mobile | Update mobile nav config for Fleet Products | Phase 1 |

## Phase 1: Move Product Summary to Fleet

**File: `go/vend/ui/web/vend-ui/fleet/fleet-config.js`**

Add "Products" service under the machines module:
```js
{ key: 'products', label: 'Products', icon: '📦', endpoint: '/10/FleetInv', model: 'VendFleetInventory', readOnly: true }
```

**File: `go/vend/ui/web/vend-ui/fleet/fleet-section-config.js`**

Add products service entry:
```js
{ key: 'products', label: 'Products', icon: '📦' }
```

**File: `go/vend/ui/web/vend-ui/analytics/analytics-config.js`**

Remove `fleet-inventory` service from Analytics (replace with time-series snapshot).

## Phase 2: Create VendInventorySnapshot Proto Type

**File: `proto/vend-analytics.proto`**

Add new message:
```protobuf
message VendInventorySnapshot {
  string snapshot_id = 1;
  string machine_id = 2;
  string machine_name = 3;
  int64 timestamp = 4;
  int32 total_stock = 5;
  int32 total_capacity = 6;
  int32 fill_pct = 7;
  int32 empty_slots = 8;
  int32 low_stock_slots = 9;
  int32 total_slots = 10;
}

message VendInventorySnapshotList {
  repeated VendInventorySnapshot list = 1;
  l8api.L8MetaData metadata = 2;
}
```

Run `make-bindings.sh` to generate Go types.

## Phase 3: Create Snapshot Service

**File: `go/vend/analytics/snapshots/SnapshotService.go`** (new)

```go
const (
    ServiceName = "InvSnap"
    ServiceArea = byte(10)
)
```

Standard CRUD service with `PrimaryKey = "SnapshotId"`, read-only (no user edits), non-unique key on `Timestamp` for time-range queries.

**File: `go/vend/services/activate_analytics.go`**

Add snapshot service activation alongside existing analytics services.

## Phase 4: Snapshot Writer in inv_vend

**File: `go/vend/inv_vend/analytics.go`**

Replace the current `computeFleetInventory` (which PUTs current state) with a **snapshot writer** that POSTs a new record per machine every cycle:

```go
func writeInventorySnapshots(fleetMachines []*vend.VendFleetMachine, nic ifs.IVNic) {
    now := time.Now().Unix()
    for _, m := range fleetMachines {
        totalStock, totalCapacity := 0, 0
        for _, slot := range m.Inventory {
            totalStock += int(slot.CurrentStock)
            totalCapacity += int(slot.Capacity)
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
            TotalStock:    int32(totalStock),
            TotalCapacity: int32(totalCapacity),
            FillPct:       fillPct,
            EmptySlots:    m.EmptySlots,
            LowStockSlots: m.LowStockSlots,
            TotalSlots:    m.TotalSlots,
        }
        vendcommon.PostEntity("InvSnap", 10, snapshot, nic)
    }
}
```

**Key difference from current:** Always POST (never PUT). Each 5-min cycle creates new records. Over 24h with 7 machines = 7 × 288 = 2,016 records/day. This is the time-series data.

**Keep `computeFleetInventory`** — it still runs to maintain the product summary in `VendFleetInventory` (now shown under Fleet > Products).

## Phase 5: Analytics UI — Time-Series Chart

**File: `go/vend/ui/web/vend-ui/analytics/analytics-config.js`**

Replace the fleet-inventory service with snapshot service (explicit chart config per probler pattern):
```js
svc('snapshots', 'Inventory History', '', '/10/InvSnap', 'VendInventorySnapshot',
    { supportedViews: ['table', 'chart'], readOnly: true,
      viewConfig: { chartType: 'line', categoryField: 'timestamp', valueField: 'fillPct',
                    aggregation: 'avg', pageSize: 2000 } })
```

The `viewConfig` is passed through the service registry to the chart factory:
- `chartType: 'line'` — renders time-series line chart (not bar/pie)
- `categoryField: 'timestamp'` — X-axis uses the snapshot timestamp
- `valueField: 'fillPct'` — Y-axis shows fill percentage
- `aggregation: 'avg'` — averages fill % when multiple machines land in same time bucket
- `pageSize: 2000` — fetches enough data points for ~1 day of all machines (vs default 100)

**File: `go/vend/ui/web/vend-ui/analytics/data/data-columns.js`**

Add columns for `VendInventorySnapshot`:
```js
VendInventorySnapshot: [
    ...col.id('snapshotId'),
    ...col.col('machineName', 'Machine'),
    ...col.date('timestamp', 'Time'),
    ...col.number('fillPct', 'Fill %'),
    ...col.number('totalStock', 'Stock'),
    ...col.number('totalCapacity', 'Capacity'),
    ...col.number('totalSlots', 'Slots'),
    ...col.number('emptySlots', 'Empty'),
    ...col.number('lowStockSlots', 'Low Stock')
]
```

**File: `go/vend/ui/web/vend-ui/analytics/data/data-forms.js`**

Add form for detail popup (all fields readOnly):
```js
VendInventorySnapshot: f.form('Inventory Snapshot', [
    f.section('Snapshot Details', [
        ...f.text('snapshotId', 'Snapshot ID', false, { readOnly: true }),
        ...f.text('machineName', 'Machine', false, { readOnly: true }),
        ...f.text('machineId', 'Machine ID', false, { readOnly: true }),
        ...f.date('timestamp', 'Timestamp', false, { readOnly: true }),
        ...f.number('fillPct', 'Fill %', false, { readOnly: true }),
        ...f.number('totalStock', 'Total Stock', false, { readOnly: true }),
        ...f.number('totalCapacity', 'Total Capacity', false, { readOnly: true }),
        ...f.number('totalSlots', 'Total Slots', false, { readOnly: true }),
        ...f.number('emptySlots', 'Empty Slots', false, { readOnly: true }),
        ...f.number('lowStockSlots', 'Low Stock Slots', false, { readOnly: true })
    ])
])
```

**File: `go/vend/ui/web/vend-ui/analytics/data/data-enums.js`**

Add `VendInventorySnapshot` to primaryKeys.

## Phase 6: Register Types

**File: `go/vend/ui/shared.go`**

Add:
```go
common.RegisterType(resources, &vend.VendInventorySnapshot{}, &vend.VendInventorySnapshotList{}, "SnapshotId")
```

**File: `go/vend/inv_vend/main.go`**

Register snapshot types for POST:
```go
nic.Resources().Registry().Register(&vend.VendInventorySnapshot{})
nic.Resources().Registry().Register(&vend.VendInventorySnapshotList{})
nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendInventorySnapshot{}, "SnapshotId")
```

## Phase 7: Data Retention

Time-series data grows indefinitely. Add a cleanup goroutine that runs on the 1st of every month and deletes snapshots older than the configured retention period (default: 30 days).

Retention period is configurable via a constant in the analytics package:
```go
const DefaultRetentionDays = 30
```

This can be overridden by environment variable `VEND_SNAPSHOT_RETENTION_DAYS` at startup, allowing operators to adjust without rebuilding:
```go
func getRetentionDays() int {
    if val := os.Getenv("VEND_SNAPSHOT_RETENTION_DAYS"); val != "" {
        if days, err := strconv.Atoi(val); err == nil && days > 0 {
            return days
        }
    }
    return DefaultRetentionDays
}
```

Cleanup logic:
```go
func cleanOldSnapshots(nic ifs.IVNic) {
    // Wait until first of next month, then run monthly
    for {
        sleepUntilFirstOfMonth()
        cutoff := time.Now().AddDate(0, 0, -getRetentionDays()).Unix()
        // Delete snapshots where timestamp < cutoff
    }
}
```

With 7 machines at 5-min intervals, 30 days = ~60K records. Manageable for the ORM.

## Mobile Parity

Mobile analytics columns at `m/js/analytics/analytics-columns.js` need a `VendInventorySnapshot` entry. The existing `VendFleetInventory` columns stay (now served under Fleet > Products via the same nav config update).

**File: `go/vend/ui/web/m/js/analytics/analytics-columns.js`**

Add:
```js
VendInventorySnapshot: [
    ...col.id('snapshotId'),
    { key: 'machineName', label: 'Machine', primary: true, sortKey: 'machineName' },
    ...col.date('timestamp', 'Time'),
    { key: 'fillPct', label: 'Fill %', secondary: true },
    ...col.number('totalStock', 'Stock'),
    ...col.number('totalSlots', 'Slots'),
    ...col.number('emptySlots', 'Empty')
]
```

## Phase 8: Verification

- [ ] `go build ./...` passes
- [ ] Navigate to Fleet > Products — verify per-product summary table with fill bars
- [ ] Navigate to Analytics > Inventory History — verify time-series records accumulating
- [ ] Switch to chart view — verify line chart with time X-axis and fill % Y-axis
- [ ] Wait 15+ minutes — verify multiple data points per machine appear
- [ ] Click a snapshot row — verify detail popup with all fields
- [ ] Verify Analytics > Performance still works (slot performance records)
- [ ] Verify Analytics > Forecasts still works
- [ ] Mobile: verify Analytics shows snapshot data
- [ ] Mobile: verify Fleet > Products shows product summary

## Critical Files

| Action | File |
|--------|------|
| Modify | `proto/vend-analytics.proto` (add VendInventorySnapshot) |
| Create | `go/vend/analytics/snapshots/SnapshotService.go` |
| Modify | `go/vend/services/activate_analytics.go` (add snapshot service) |
| Modify | `go/vend/inv_vend/analytics.go` (add writeInventorySnapshots) |
| Modify | `go/vend/inv_vend/main.go` (register snapshot types) |
| Modify | `go/vend/ui/shared.go` (register snapshot types for web) |
| Modify | `go/vend/ui/web/vend-ui/fleet/fleet-config.js` (add Products) |
| Modify | `go/vend/ui/web/vend-ui/fleet/fleet-section-config.js` (add Products) |
| Modify | `go/vend/ui/web/vend-ui/analytics/analytics-config.js` (snapshots) |
| Modify | `go/vend/ui/web/vend-ui/analytics/data/data-columns.js` |
| Modify | `go/vend/ui/web/vend-ui/analytics/data/data-forms.js` |
| Modify | `go/vend/ui/web/m/js/analytics/analytics-columns.js` |
