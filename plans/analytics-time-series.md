# Analytics: Time-Series Inventory & Performance Metrics

## Context

The Analytics section currently points at empty CRUD tables. Real data exists in `VendFleetMachine` (inventory snapshots every 5 min) and `VendTransaction` (sales events from simulator). The analytics proto types (`VendSlotPerformance`, `VendFleetInventory`) already define the right aggregated structures — they just need to be populated.

## Goal

Show meaningful time-series data in Analytics:
- **Fleet Inventory**: Fleet-wide fill % over time, per-product stock levels, stockout trends
- **Performance**: Per-slot vend velocity, revenue, stockout hours per period
- **Forecasts**: Predicted stockout times based on observed depletion rates

## Data Source

The `inv_vend` binary already reads all `VendFleetMachine` records every 5 minutes (bridge loop + threshold evaluator). It can also compute and POST aggregated analytics records on the same cycle — no new binary or polling mechanism needed.

## Phase 1: Fleet Inventory Snapshotter

**New file: `go/vend/inv_vend/analytics.go`**

A goroutine that runs every 5 minutes (same as bridge/threshold) and computes `VendFleetInventory` summaries.

**Logic:**
1. Read all `VendFleetMachine` records (same data threshold evaluator already fetches)
2. Group slots by `productName` (or `sku`) across all machines
3. For each product, compute:
   - `totalMachines` — how many machines stock it
   - `totalSlots` — total slot count for this product
   - `totalUnitsInMachines` — sum of `currentStock`
   - `totalCapacity` — sum of `capacity`
   - `fleetSoldOutCount` — machines where this product has 0 stock
   - `fleetLowStockCount` — machines where this product is below 30% fill
   - `lastUpdated` — current timestamp
4. POST each `VendFleetInventory` record to `/10/FleetInv`

**Dedup:** Use `SummaryId = productName` (one record per product, updated each cycle via PUT).

**Frequency:** Every 5 minutes, same timer as the bridge loop. Can share the fetched machine data.

```go
func computeFleetInventory(machines []*vend.VendFleetMachine, nic ifs.IVNic) {
    products := make(map[string]*vend.VendFleetInventory)
    for _, m := range machines {
        for _, slot := range m.Inventory {
            key := slot.ProductName
            if key == "" { continue }
            summary := products[key]
            if summary == nil {
                summary = &vend.VendFleetInventory{
                    SummaryId:   key,
                    ProductName: slot.ProductName,
                }
                products[key] = summary
            }
            summary.TotalSlots++
            summary.TotalCapacity += slot.Capacity
            summary.TotalUnitsInMachines += slot.CurrentStock
            // track per-machine counts...
        }
    }
    // POST/PUT each summary
}
```

## Phase 2: Slot Performance Aggregator

**New file: `go/vend/inv_vend/performance.go`**

Computes `VendSlotPerformance` records per slot per day. Since we don't have actual transaction data from the simulator (only `dailyTransactions` count on the machine), we derive what we can:

**Available metrics per slot (from inventory snapshots):**
- `vendCount` — derived from stock level changes between snapshots (previous - current = vended)
- `stockoutHours` — if `currentStock == 0`, accumulate time since last non-zero
- `velocity` — vendCount / time period

**Tracking state:** The aggregator needs to remember previous stock levels to compute deltas. Use a local `map[string]int32` of `machineId:slotNumber → previousStock`.

**Period:** Daily (midnight-to-midnight). At end of day, POST the completed `VendSlotPerformance` record.

**Simplified initial version:** Since delta tracking requires persistence across restarts (which we don't have), start with a simpler approach:
- Compute `velocity` from `(capacity - currentStock) / hoursSinceLastRestock` (approximation)
- Set `rank` by sorting slots within a machine by `currentStock/capacity` ratio
- `stockoutHours` = 0 if currentStock > 0, else estimate from fill rate

```go
type slotState struct {
    previousStock int32
    lastSeen      int64
    stockoutStart int64
}
var slotStates = make(map[string]*slotState) // "machineId:slotNum" -> state
```

## Phase 3: Simple Forecast Generator

**Add to: `go/vend/inv_vend/analytics.go`**

Generate `VendForecast` records based on observed depletion:

**Logic:**
1. For each machine+product combo, compare current stock to capacity
2. If we have delta tracking (Phase 2), use observed vend rate to predict stockout
3. Otherwise, use `velocity = (capacity - currentStock) / daysSinceRestock` as estimate
4. `predictedStockoutTime = now + (currentStock / velocity) * 86400`
5. `restockUrgency`:
   - HIGH: predicted stockout < 24h
   - MEDIUM: predicted stockout < 72h
   - LOW: predicted stockout > 72h

**Frequency:** Once per hour (not every 5 min — forecasts don't change that fast).

## Phase 4: Analytics UI — Chart Default View

**File: `go/vend/ui/web/vend-ui/analytics/analytics-config.js`**

Update Fleet Inventory to default to chart view and add proper columns/forms for the populated data.

**File: `go/vend/ui/web/vend-ui/analytics/data/data-columns.js`**

Update columns to match the actual populated fields of `VendFleetInventory`:
- `productName` — Product
- `totalMachines` — Machines
- `totalSlots` — Slots
- `totalUnitsInMachines` — Units in Field
- `totalCapacity` — Total Capacity
- Fill % (custom column with fill bar)
- `fleetSoldOutCount` — Sold Out
- `fleetLowStockCount` — Low Stock
- `lastUpdated` — Last Updated

For `VendSlotPerformance`:
- `machineId` — Machine
- `productName` — Product
- `vendCount` — Vends
- `velocity` — Vends/Day
- `rank` — Rank
- `stockoutHours` — Stockout Hours

For `VendForecast`:
- `machineId` — Machine
- `productId` — Product
- `predictedDailyVends` — Predicted Vends/Day
- `predictedStockoutTime` — Predicted Stockout (date)
- `restockUrgency` — Urgency
- `confidenceScore` — Confidence

## Phase 5: Wire Analytics Goroutine

**File: `go/vend/inv_vend/main.go`**

Add new goroutine:
```go
go computeAnalytics(nic)
```

The `computeAnalytics` function:
1. Waits 90s for fleet data to be populated
2. Every 5 minutes: reads fleet machines, computes fleet inventory summaries, computes slot performance deltas
3. Every hour: generates forecasts

## Implementation Details

### No new binary needed
All analytics computation runs in the existing `inv_vend` binary alongside the bridge and threshold evaluator. It reads the same fleet machine data.

### No new proto types needed
`VendFleetInventory`, `VendSlotPerformance`, and `VendForecast` are already defined with appropriate fields. Services already exist (`/10/FleetInv`, `/10/SlotPerf`, `/10/Forecast`).

### Data flow
```
VendFleetMachine (collected every 5 min)
    │
    ├── bridgeVCacheToFleet (existing)
    ├── evaluateThresholds (existing)
    │
    └── computeAnalytics (new)
        ├── VendFleetInventory — per-product fleet aggregates (PUT every 5 min)
        ├── VendSlotPerformance — per-slot daily metrics (POST daily)
        └── VendForecast — predictive stockout times (PUT hourly)
```

### Analytics services endpoints
| Service | Endpoint | Model | View |
|---------|----------|-------|------|
| Fleet Inventory | `/10/FleetInv` | VendFleetInventory | table + chart |
| Performance | `/10/SlotPerf` | VendSlotPerformance | table + chart |
| Forecasts | `/10/Forecast` | VendForecast | table |

## Phase 4 Form Definitions

Desktop forms (`data-forms.js`) must be updated to show detail popups with the populated fields:

**VendFleetInventory form:**
- Section "Product Summary": productName (text, readOnly), totalMachines (number, readOnly), totalSlots (number, readOnly), totalUnitsInMachines (number, readOnly), totalCapacity (number, readOnly), fleetSoldOutCount (number, readOnly), fleetLowStockCount (number, readOnly), lastUpdated (date, readOnly)

**VendSlotPerformance form:**
- Section "Performance": machineId (text, readOnly), slotId (text, readOnly), productName (text, readOnly), vendCount (number, readOnly), velocity (number, readOnly), rank (number, readOnly), stockoutHours (number, readOnly), periodStart (date, readOnly), periodEnd (date, readOnly)

**VendForecast form:**
- Section "Forecast": machineId (text, readOnly), productId (text, readOnly), forecastDate (date, readOnly), horizonDays (number, readOnly), predictedDailyVends (number, readOnly), predictedStockoutTime (date, readOnly), restockUrgency (text, readOnly), confidenceScore (number, readOnly)

## Mobile Parity

Mobile analytics columns already exist at `m/js/analytics/analytics-columns.js` and cover the key fields for all three types (`VendFleetInventory`, `VendSlotPerformance`, `VendForecast`). No changes needed — the mobile columns already match the protobuf fields that the backend will populate.

Mobile forms at `m/js/analytics/analytics-forms.js` should be verified to include the same fields as desktop. If gaps exist, update in Phase 4.

## Traceability Matrix

| # | Section | Action Item | Phase |
|---|---------|-------------|-------|
| 1 | Backend | Create analytics.go with fleet inventory snapshotter | Phase 1 |
| 2 | Backend | Compute per-product aggregates from VendFleetMachine slots | Phase 1 |
| 3 | Backend | POST/PUT VendFleetInventory records every 5 min | Phase 1 |
| 4 | Backend | Create performance.go with slot performance tracker | Phase 2 |
| 5 | Backend | Track stock level deltas between snapshots | Phase 2 |
| 6 | Backend | Compute vendCount, velocity, stockoutHours, rank | Phase 2 |
| 7 | Backend | POST VendSlotPerformance records daily | Phase 2 |
| 8 | Backend | Add forecast computation (depletion-rate based) | Phase 3 |
| 9 | Backend | POST/PUT VendForecast records hourly | Phase 3 |
| 10 | UI Desktop | Update analytics columns for VendFleetInventory | Phase 4 |
| 11 | UI Desktop | Update analytics columns for VendSlotPerformance | Phase 4 |
| 12 | UI Desktop | Update analytics columns for VendForecast | Phase 4 |
| 13 | UI Desktop | Update analytics forms for all three types | Phase 4 |
| 14 | UI Desktop | Revert analytics-config endpoint back to FleetInv | Phase 4 |
| 15 | UI Mobile | Verify mobile columns/forms cover populated fields | Phase 4 |
| 16 | Backend | Wire computeAnalytics goroutine in inv_vend main.go | Phase 5 |

## Phase 6: Verification

- [ ] `go build ./...` passes
- [ ] Navigate to Analytics > Fleet Inventory — verify per-product rows with fill data
- [ ] Click a row — verify detail popup shows all fields correctly
- [ ] Switch to chart view — verify bar/pie chart shows product distribution
- [ ] Navigate to Analytics > Performance — verify per-slot records appear after 5+ min
- [ ] Click a performance row — verify detail popup fields
- [ ] Navigate to Analytics > Forecasts — verify forecast records appear after 1 hour
- [ ] Verify no performance impact (analytics shares data fetch with bridge/threshold)
- [ ] Mobile: Navigate to Analytics > Fleet Inventory — verify card data displays
- [ ] Mobile: Click a card — verify detail popup

## Critical Files

| Action | File |
|--------|------|
| Create | `go/vend/inv_vend/analytics.go` (fleet inventory + forecast computation) |
| Create | `go/vend/inv_vend/performance.go` (slot performance tracking) |
| Modify | `go/vend/inv_vend/main.go` (add computeAnalytics goroutine) |
| Modify | `go/vend/ui/web/vend-ui/analytics/analytics-config.js` (revert to FleetInv) |
| Modify | `go/vend/ui/web/vend-ui/analytics/data/data-columns.js` (real columns) |
| Modify | `go/vend/ui/web/vend-ui/analytics/data/data-forms.js` (real forms) |
| Verify | `go/vend/ui/web/m/js/analytics/analytics-columns.js` (mobile parity) |
| Verify | `go/vend/ui/web/m/js/analytics/analytics-forms.js` (mobile parity) |
