# High-Scale Analytics: Pre-Computed Machine Profiles

## Problem

The restock recommendation engine and analytics charts need to analyze historical patterns across hundreds of machines with millions of data points (300 machines × 288 snapshots/day × 30 days = 2.6M records/month). Fetching and processing raw snapshots at query time doesn't scale:

- Server returns 100 records per page — reading 2.6M records requires 26,000 page fetches
- Each recommendation refresh (every 30 min) would re-process the entire dataset
- Chart rendering with millions of data points is impractical client-side
- Adding more machines or longer history multiplies the problem linearly

## Solution: Pre-Computed Profiles

Instead of querying raw snapshots at analysis time, **compute profiles incrementally as data arrives**. The snapshot writer already iterates all machines every 5 minutes — it extends to update per-machine profile records in the same loop. Consumers (restock engine, charts, dashboards) read lightweight profile records instead of raw data.

### Architecture

```
Snapshot Writer (every 5 min, already exists)
    │
    ├── POST VendInventorySnapshot (raw data, as today)
    │
    └── UPDATE VendMachineProfile (incremental, in same loop)
            │
            ├── Restock Engine reads 300 profiles (not 2.6M snapshots)
            ├── Charts read 300 profiles (not 2.6M snapshots)
            └── Dashboard reads 300 profiles (not 2.6M snapshots)
```

### Key principle: Write-time aggregation

Every time a snapshot is written, the profile for that machine is updated incrementally:
- New data point is folded into rolling averages
- No need to re-read historical snapshots
- Profile stays current without expensive batch queries
- O(1) per machine per cycle, not O(N) where N is snapshot count

## VendMachineProfile Entity

One record per machine, updated every 5 minutes. Contains pre-computed analytics that would otherwise require scanning thousands of snapshots.

```protobuf
message VendMachineProfile {
    string profile_id = 1;              // = machineId
    string machine_id = 2;
    string machine_name = 3;
    int64 last_updated = 4;

    // Day-of-week depletion rates (units of fill% drop per hour)
    // Index 0=Sunday, 1=Monday, ..., 6=Saturday
    repeated double dow_depletion_rate = 10;
    repeated int32 dow_sample_count = 11;

    // Hour-of-day depletion rates (24 entries)
    repeated double hod_depletion_rate = 12;
    repeated int32 hod_sample_count = 13;

    // Overall metrics (rolling 30-day)
    double avg_hourly_depletion = 20;
    double avg_daily_revenue = 21;
    int64 total_revenue_30d = 22;
    int32 avg_fill_pct = 23;
    int32 restock_count_30d = 24;
    double avg_restock_interval_hours = 25;

    // Location classification (computed from depletion pattern)
    string location_class = 30;         // "office", "retail", "transit", "mixed"
    double weekend_weekday_ratio = 31;

    // Per-product depletion rates (top products)
    repeated VendProductProfile top_products = 40;

    // Trend (week-over-week)
    double trend_multiplier = 50;       // >1.0 = increasing demand, <1.0 = decreasing

    // Cascade threshold (fill% where revenue drops disproportionately)
    int32 cascade_threshold_pct = 60;   // default 30 if not enough data
}

message VendProductProfile {
    string product_name = 1;
    double depletion_rate_per_hour = 2;
    int64 price = 3;
    int32 avg_stock = 4;
    int32 capacity = 5;
    double time_to_empty_hours = 6;
}

message VendMachineProfileList {
    repeated VendMachineProfile list = 1;
    l8api.L8MetaData metadata = 2;
}
```

## Incremental Update Algorithm

The profile update runs inside the existing snapshot writer loop — zero additional data fetches.

### Per-snapshot update (every 5 min per machine):

```
function updateProfile(machine, prevSnapshot, currentSnapshot, profile):
    now = currentTime
    dayOfWeek = now.Weekday()    // 0-6
    hourOfDay = now.Hour()       // 0-23

    // 1. Depletion rate for this interval
    if prevSnapshot exists AND currentSnapshot.fillPct < prevSnapshot.fillPct:
        // Stock decreased = depletion (not a restock)
        deltaFill = prevSnapshot.fillPct - currentSnapshot.fillPct
        hoursElapsed = (now - prevSnapshot.timestamp) / 3600
        depletionRate = deltaFill / hoursElapsed

        // Update day-of-week rolling average
        profile.dowDepletionRate[dayOfWeek] = rollingAvg(
            profile.dowDepletionRate[dayOfWeek],
            profile.dowSampleCount[dayOfWeek],
            depletionRate
        )
        profile.dowSampleCount[dayOfWeek]++

        // Update hour-of-day rolling average
        profile.hodDepletionRate[hourOfDay] = rollingAvg(
            profile.hodDepletionRate[hourOfDay],
            profile.hodSampleCount[hourOfDay],
            depletionRate
        )
        profile.hodSampleCount[hourOfDay]++

        // Update overall average
        profile.avgHourlyDepletion = rollingAvg(...)

    else if currentSnapshot.fillPct > prevSnapshot.fillPct + 10:
        // Fill% jumped up = restock event
        profile.restockCount30d++

    // 2. Revenue tracking
    profile.totalRevenue30d += currentSnapshot.dailyRevenue delta
    profile.avgDailyRevenue = profile.totalRevenue30d / 30

    // 3. Per-product depletion (from slot data)
    for each slot in machine.Inventory:
        if prevSlotStock[slot] > slot.currentStock:
            sold = prevSlotStock[slot] - slot.currentStock
            rate = sold / hoursElapsed
            updateProductProfile(profile, slot.productName, rate, slot)

    // 4. Location classification (recompute periodically, not every cycle)
    if shouldReclassify(profile):  // e.g., every 6 hours
        weekdayAvg = avg(profile.dowDepletionRate[1..5])
        weekendAvg = avg(profile.dowDepletionRate[0, 6])
        profile.weekendWeekdayRatio = weekendAvg / weekdayAvg
        if ratio > 1.5: profile.locationClass = "retail"
        elif ratio < 0.5: profile.locationClass = "office"
        else: profile.locationClass = "mixed"

    // 5. Trend (compare recent week vs 4 weeks ago)
    // Tracked via separate weekly buckets, updated incrementally

    profile.lastUpdated = now
    PUT profile to service
```

### Rolling average formula (no history needed):

```
rollingAvg(currentAvg, sampleCount, newValue) =
    (currentAvg * sampleCount + newValue) / (sampleCount + 1)
```

This is an **exponentially weighted moving average** that doesn't require storing past values. Each new data point nudges the average. Older data naturally fades as more samples arrive.

### 30-day window decay:

To prevent stale data from dominating, apply a decay factor monthly:
```
// On the 1st of every month (same schedule as snapshot retention cleanup):
profile.dowSampleCount[i] = profile.dowSampleCount[i] / 2  // halve sample counts
profile.hodSampleCount[i] = profile.hodSampleCount[i] / 2
profile.restockCount30d = profile.restockCount30d / 2
profile.totalRevenue30d = profile.totalRevenue30d / 2
```

This gradually forgets old patterns while preserving recent trends. After 2 months, data from day 1 has been halved twice (25% weight), while recent data has full weight.

## Scale Analysis

### Current approach (raw snapshot queries):

| Metric | Value |
|--------|-------|
| Machines | 300 |
| Snapshots/day/machine | 288 (every 5 min) |
| Snapshots/month | 2,592,000 |
| Records to fetch for profiles | 2,592,000 (26,000 page fetches) |
| Time to fetch (100ms/page) | ~43 minutes |
| Memory for processing | ~500 MB |

### Profile approach:

| Metric | Value |
|--------|-------|
| Machines | 300 |
| Profile records | 300 (one per machine) |
| Records to fetch for restock engine | 300 (3 page fetches) |
| Time to fetch (100ms/page) | 0.3 seconds |
| Memory for processing | ~1 MB |
| Update cost per cycle | 300 PUT operations (already in the snapshot loop) |

### At 10,000 machines:

| Metric | Raw Snapshots | Profiles |
|--------|--------------|----------|
| Monthly records | 86,400,000 | 10,000 |
| Fetch for analysis | 864,000 pages | 100 pages |
| Fetch time | ~24 hours | 10 seconds |
| Profile update/cycle | N/A | 10,000 PUTs (parallelizable) |

## Proto Field-to-JS Mapping

| Proto field | JSON name | JS key |
|---|---|---|
| profile_id | profileId | profileId |
| machine_id | machineId | machineId |
| machine_name | machineName | machineName |
| last_updated | lastUpdated | lastUpdated |
| dow_depletion_rate | dowDepletionRate | dowDepletionRate |
| dow_sample_count | dowSampleCount | dowSampleCount |
| hod_depletion_rate | hodDepletionRate | hodDepletionRate |
| hod_sample_count | hodSampleCount | hodSampleCount |
| avg_hourly_depletion | avgHourlyDepletion | avgHourlyDepletion |
| avg_daily_revenue | avgDailyRevenue | avgDailyRevenue |
| total_revenue_30d | totalRevenue30d | totalRevenue30d |
| avg_fill_pct | avgFillPct | avgFillPct |
| restock_count_30d | restockCount30d | restockCount30d |
| avg_restock_interval_hours | avgRestockIntervalHours | avgRestockIntervalHours |
| location_class | locationClass | locationClass |
| weekend_weekday_ratio | weekendWeekdayRatio | weekendWeekdayRatio |
| top_products | topProducts | topProducts |
| trend_multiplier | trendMultiplier | trendMultiplier |
| cascade_threshold_pct | cascadeThresholdPct | cascadeThresholdPct |

VendProductProfile (child):
| product_name | productName | productName |
| depletion_rate_per_hour | depletionRatePerHour | depletionRatePerHour |
| price | price | price |
| avg_stock | avgStock | avgStock |
| capacity | capacity | capacity |
| time_to_empty_hours | timeToEmptyHours | timeToEmptyHours |

## Traceability Matrix

| # | Section | Action Item | Phase |
|---|---------|-------------|-------|
| 1 | VendMachineProfile Entity | Add VendMachineProfile + VendProductProfile to proto | Phase 1 |
| 2 | VendMachineProfile Entity | Run make-bindings.sh | Phase 1 |
| 3 | Architecture | Create ProfileService.go (ServiceName="MachProf") | Phase 1 |
| 4 | Architecture | ProfileServiceCallback: profileId = machineId (not auto-generated — deterministic ID) | Phase 1 |
| 5 | Architecture | Add to activate_analytics.go | Phase 1 |
| 6 | Architecture | Register types with AddPrimaryKeyDecorator in UI shared + inv_vend | Phase 1 |
| 7 | Incremental Update Algorithm | Create profiles.go with updateMachineProfile | Phase 2 |
| 8 | Incremental Update Algorithm | Add profile update to snapshot writer loop | Phase 2 |
| 9 | Incremental Update Algorithm | Implement rolling average (no history storage) | Phase 2 |
| 10 | Incremental Update Algorithm | Implement per-product depletion tracking | Phase 2 |
| 11 | Incremental Update Algorithm | Implement location classification from depletion pattern | Phase 2 |
| 12 | 30-day Window Decay | Add profile decay to monthly cleanup goroutine | Phase 3 |
| 13 | Restock Engine Uses Profiles | Update restock engine to read profiles instead of snapshots | Phase 4 |
| 14 | Top Performers Uses Profiles | Replace snapshot-based computation with profile query | Phase 5 |
| 15 | Refactor Forecasts | Use profile depletionRatePerHour instead of estimated velocity | Phase 5.5 |
| 16 | Refactor Slot Performance | Derive from profiles, remove in-memory slotStates | Phase 5.5 |
| 17 | Refactor Analytics | Remove hardcoded productPrices map, use profile prices | Phase 5.5 |
| 18 | UI Desktop | Add Machine Profiles to Fleet config + section config | Phase 6 |
| 19 | UI Desktop | Add desktop columns (all scalar fields + topProducts inline) | Phase 6 |
| 20 | UI Desktop | Add desktop form with all fields + topProducts inline table | Phase 6 |
| 21 | UI Desktop | Profile chart view (dow/hod pattern bar charts) | Phase 7 |
| 22 | UI Mobile | Add mobile columns with primary/secondary markers | Phase 6 |
| 23 | UI Mobile | Add mobile forms | Phase 6 |
| 24 | UI Mobile | Add mobile nav config entry | Phase 6 |
| 25 | Mocks | Create gen_profiles.go seeding sample profiles with non-zero values | Phase 6 |

## Phase 1: Proto + Service

1. Add `VendMachineProfile`, `VendProductProfile` to proto
2. Run make-bindings.sh
3. Create `ProfileService.go` (ServiceName="MachProf", area 10) with `ProfileServiceCallback.go`. **Exemption from ServiceCallback Auto-Generate ID rule:** `profileId` is set deterministically to `machineId` by the writer (one profile per machine). Auto-generating would break the 1:1 mapping. The callback validates `ProfileId` is present on POST but does NOT auto-generate it.
4. Add to `activate_analytics.go`
5. Register types with `AddPrimaryKeyDecorator(&vend.VendMachineProfile{}, "ProfileId")` in UI shared + inv_vend

## Phase 2: Incremental Profile Writer

**File: `go/vend/inv_vend/profiles.go`** (new)

Add to the snapshot writer loop in `writeInventorySnapshots()`:
```go
// After writing snapshot, update profile incrementally
updateMachineProfile(m, prevSnapshots[m.MachineId], snapshot, nic)
prevSnapshots[m.MachineId] = snapshot
```

Implement:
- `updateMachineProfile(machine, prevSnapshot, currentSnapshot, nic)` — incremental update
- `rollingAvg(current, count, newValue)` — weighted average without history
- `updateProductProfiles(profile, machine, prevSlotStocks)` — per-product depletion
- `classifyLocation(profile)` — from weekend/weekday ratio
- Store `prevSnapshots` and `prevSlotStocks` in memory (maps keyed by machineId)

## Phase 3: Monthly Decay

Add to the existing `cleanOldSnapshots` goroutine (runs on 1st of month):
```go
// After cleaning old snapshots, decay profile sample counts
decayProfiles(nic)
```

## Phase 4: Restock Engine Uses Profiles

Update `smart-restock-recommendations.md` Phase 2:
- `computeRestockRecommendations` reads `VendMachineProfile` records (300 records, 3 pages)
- No snapshot queries needed — all depletion rates, location classes, trends are pre-computed
- Each scenario evaluator receives the profile, not raw snapshots

```go
profiles, _ := vendcommon.GetEntities("MachProf", 10, &vend.VendMachineProfile{}, nic)
for _, profile := range profiles {
    candidates = append(candidates, evaluateDayOfWeekDemand(profile))
    candidates = append(candidates, evaluateFastMovers(profile))
    candidates = append(candidates, evaluateCriticalPrediction(profile))
}
applyRevenuePriority(candidates, profiles)
```

## Phase 5: Top Performers Uses Profiles

Replace `computeTopPerformers` (which reads 10,000 snapshots) with a simple query:
```go
profiles, _ := vendcommon.GetEntities("MachProf", 10, &vend.VendMachineProfile{}, nic)
// Sort by totalRevenue30d, write top performers
```

## Phase 5.5: Refactor Existing Analytics to Use Profiles

The current analytics services compute from raw snapshots or in-memory state. With profiles available, refactor them to consume pre-computed data:

### Top Performers (already addressed in Phase 5)
- Before: `GetEntitiesByQuery("select * from VendInventorySnapshot limit 10000")` → aggregate in Go
- After: `GetEntities("MachProf", 10, ...)` → sort by `TotalRevenue30d` → write VendTopPerformer

### Forecasts
- Before: estimates stockout from `(capacity - currentStock) / 7 days` (rough guess)
- After: uses profile's `avgHourlyDepletion` and `hodDepletionRate` for per-hour projection
```go
func computeForecasts(fleetMachines []*vend.VendFleetMachine, profiles map[string]*vend.VendMachineProfile, nic ifs.IVNic) {
    for _, m := range fleetMachines {
        profile := profiles[m.MachineId]
        if profile == nil { continue }
        for _, slot := range m.Inventory {
            // Use profile's per-product depletion rate (not estimated)
            productProfile := findProductProfile(profile, slot.ProductName)
            if productProfile == nil { continue }
            velocity := productProfile.DepletionRatePerHour * 24  // daily
            hoursToEmpty := float64(slot.CurrentStock) / productProfile.DepletionRatePerHour
            // ... rest of forecast using accurate velocity
        }
    }
}
```

### Slot Performance
- Before: tracks deltas in memory (`slotStates` map), lost on restart, posts hourly
- After: profile's `topProducts` already has per-product `depletionRatePerHour`, `avgStock`, `timeToEmptyHours`
- Refactor `computeSlotPerformance` to read from profiles instead of maintaining its own state
- The `slotStates` map and `prevStockLevels` in `performance.go` become redundant — profiles track this persistently
- Keep `performance.go` for posting `VendSlotPerformance` records (periodic summary), but derive values from profiles

### Fleet Inventory
- Keep as-is — it aggregates by product across machines (different dimension than profiles)
- Can optionally add `unitPrice` from profile's `topProducts` price data instead of the hardcoded price map

### What gets removed
- `prevStockLevels` map in `performance.go` → profiles track this
- `slotStates` map in `performance.go` → profiles track this
- Hardcoded `productPrices` map in `analytics.go` → profiles have per-product prices
- `GetEntitiesByQuery` with `limit 10000` in `computeTopPerformers` → replaced by profile query

### What stays
- `writeInventorySnapshots` → still writes raw snapshots (needed for Inventory History chart)
- `computeFleetInventory` → still aggregates by product (different view)
- `cleanOldSnapshots` → still manages retention

## Phase 6: UI — Machine Profile Tab

### Desktop Config
Add to Fleet section:
```js
{ key: 'profiles', label: 'Machine Profiles', endpoint: '/10/MachProf',
  model: 'VendMachineProfile', readOnly: true,
  supportedViews: ['table', 'chart'],
  viewConfig: { chartType: 'bar', categoryField: 'machineName', valueField: 'avgDailyRevenue', aggregation: 'sum' } }
```

### Desktop Columns
```js
VendMachineProfile: [
    col.id('profileId'),
    col.col('machineName', 'Machine'),
    col.col('locationClass', 'Location Type'),
    col.number('avgHourlyDepletion', 'Depletion/hr'),
    col.money('avgDailyRevenue', 'Avg Daily Revenue'),
    col.money('totalRevenue30d', '30-Day Revenue'),
    col.number('avgFillPct', 'Avg Fill %'),
    col.number('trendMultiplier', 'Trend'),
    col.number('restockCount30d', 'Restocks (30d)'),
    col.number('avgRestockIntervalHours', 'Restock Interval (hrs)'),
    col.number('weekendWeekdayRatio', 'Weekend/Weekday'),
    col.number('cascadeThresholdPct', 'Cascade Threshold'),
    col.date('lastUpdated', 'Updated')
]
```
// Omitted from table columns: profileId (PK, shown in detail only), machineId (= profileId, redundant — the profile IS the machine), dowDepletionRate/hodDepletionRate/dowSampleCount/hodSampleCount (repeated arrays — visualized in Phase 7 chart, not tabular), topProducts (shown in form inline table)

### Desktop Form
```js
VendMachineProfile: f.form('Machine Profile', [
    f.section('Machine', [
        f.text('profileId', 'Profile ID', false, { readOnly: true }),
        f.text('machineName', 'Machine', false, { readOnly: true }),
        f.text('machineId', 'Machine ID', false, { readOnly: true }),
        f.text('locationClass', 'Location Type', false, { readOnly: true }),
        f.number('weekendWeekdayRatio', 'Weekend/Weekday Ratio', false, { readOnly: true }),
        f.date('lastUpdated', 'Last Updated', false, { readOnly: true })
    ]),
    f.section('Depletion & Revenue', [
        f.number('avgHourlyDepletion', 'Avg Depletion/hr', false, { readOnly: true }),
        f.money('avgDailyRevenue', 'Avg Daily Revenue', false, { readOnly: true }),
        f.money('totalRevenue30d', '30-Day Revenue', false, { readOnly: true }),
        f.number('avgFillPct', 'Avg Fill %', false, { readOnly: true }),
        f.number('trendMultiplier', 'Trend Multiplier', false, { readOnly: true }),
        f.number('cascadeThresholdPct', 'Cascade Threshold %', false, { readOnly: true })
    ]),
    f.section('Restock History', [
        f.number('restockCount30d', 'Restocks (30 days)', false, { readOnly: true }),
        f.number('avgRestockIntervalHours', 'Avg Interval (hours)', false, { readOnly: true })
    ]),
    f.section('Top Products', [
        f.inlineTable('topProducts', 'Product Depletion', [
            { key: 'productName', label: 'Product', type: 'text' },
            { key: 'depletionRatePerHour', label: 'Units/hr', type: 'number' },
            { key: 'price', label: 'Price', type: 'money' },
            { key: 'avgStock', label: 'Avg Stock', type: 'number' },
            { key: 'capacity', label: 'Capacity', type: 'number' },
            { key: 'timeToEmptyHours', label: 'Time to Empty (hrs)', type: 'number' }
        ])
    ])
])
```
// Note: dowDepletionRate (7 values) and hodDepletionRate (24 values) are repeated double arrays. These are best visualized in the Phase 7 chart view, not as form fields. dowSampleCount and hodSampleCount are internal accuracy metrics, not shown in UI.

### Mobile

**Mobile columns** (`m/js/fleet/fleet-columns.js`):
```js
VendMachineProfile: [
    col.id('profileId'),
    { key: 'machineName', label: 'Machine', primary: true, sortKey: 'machineName' },
    { key: 'locationClass', label: 'Type', secondary: true },
    col.money('avgDailyRevenue', 'Revenue/Day'),
    col.number('avgFillPct', 'Avg Fill %'),
    col.number('trendMultiplier', 'Trend'),
    col.number('restockCount30d', 'Restocks')
]
```

**Mobile forms** (`m/js/fleet/fleet-forms.js`):
```js
VendMachineProfile: f.form('Machine Profile', [
    f.section('Machine', [
        f.text('machineName', 'Machine', false, { readOnly: true }),
        f.text('locationClass', 'Location Type', false, { readOnly: true }),
        f.number('weekendWeekdayRatio', 'Weekend/Weekday', false, { readOnly: true })
    ]),
    f.section('Depletion & Revenue', [
        f.number('avgHourlyDepletion', 'Depletion/hr', false, { readOnly: true }),
        f.money('avgDailyRevenue', 'Avg Daily Revenue', false, { readOnly: true }),
        f.money('totalRevenue30d', '30-Day Revenue', false, { readOnly: true }),
        f.number('avgFillPct', 'Avg Fill %', false, { readOnly: true }),
        f.number('trendMultiplier', 'Trend', false, { readOnly: true }),
        f.number('cascadeThresholdPct', 'Cascade %', false, { readOnly: true })
    ]),
    f.section('Restock', [
        f.number('restockCount30d', 'Restocks (30d)', false, { readOnly: true }),
        f.number('avgRestockIntervalHours', 'Interval (hrs)', false, { readOnly: true })
    ]),
    f.section('Top Products', [
        f.inlineTable('topProducts', 'Products', [
            { key: 'productName', label: 'Product', type: 'text' },
            { key: 'depletionRatePerHour', label: 'Units/hr', type: 'number' },
            { key: 'price', label: 'Price', type: 'money' },
            { key: 'capacity', label: 'Capacity', type: 'number' },
            { key: 'timeToEmptyHours', label: 'Empty (hrs)', type: 'number' }
        ])
    ])
])
```

**Mobile nav config** (`m/js/nav-configs/layer8m-nav-config-vend.js`):
```js
{ key: 'profiles', label: 'Profiles', icon: 'fleet',
  endpoint: '/10/MachProf', model: 'VendMachineProfile', idField: 'profileId', readOnly: true }
```

### Mock Data
Create `gen_profiles.go` seeding 20 sample profiles with:
- All scalar fields populated with non-zero values
- `dowDepletionRate` with 7 entries (higher on weekdays for office machines, weekends for retail)
- `hodDepletionRate` with 24 entries (peaks at lunch hour)
- `topProducts` with 3-5 VendProductProfile entries per machine
- `locationClass` with mix of "office", "retail", "transit", "mixed"
- `trendMultiplier` between 0.8-1.3
- `cascadeThresholdPct` between 25-40

## Phase 7: Chart Uses Profiles

The Inventory History chart can offer a "Profile View" that shows the day-of-week and hour-of-day patterns from the profile instead of raw snapshots. No thousands of records needed — just the 7 dow_depletion_rate values and 24 hod_depletion_rate values from the profile.

## Phase 8: Verification

### Build
- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] profiles.go under 500 lines

### Functional
- [ ] After 2+ snapshot cycles, VendMachineProfile records appear (one per machine)
- [ ] Profile dowDepletionRate has 7 entries (one per day-of-week)
- [ ] Profile hodDepletionRate has 24 entries (one per hour)
- [ ] Profile locationClass populated ("office"/"retail"/"transit"/"mixed")
- [ ] Profile totalRevenue30d accumulates over time
- [ ] Profile topProducts shows per-product depletion rates
- [ ] Monthly decay halves sample counts without losing patterns

### Desktop UI
- [ ] Navigate to Fleet > Machine Profiles — verify table loads with data
- [ ] Click a profile row — verify detail popup opens with all sections
- [ ] Verify topProducts inline table shows product rows in detail popup
- [ ] Verify locationClass shows "office"/"retail"/"transit"/"mixed"
- [ ] Switch to chart view — verify bar chart renders

### Integration
- [ ] Restock engine reads profiles (300 records) — not snapshots (millions)
- [ ] Top Performers reads profiles — fast query
- [ ] Charts use profile data for pattern visualization

### Scale
- [ ] Profile update adds < 50ms per cycle (300 machines)
- [ ] Restock engine completes in < 5 seconds (reads 300 profiles)
- [ ] No pagination issues (profiles fit in 3 pages)

### Mobile
- [ ] Navigate to Fleet > Machine Profiles — verify card list loads
- [ ] Verify primary (machineName) and secondary (locationClass) show on cards
- [ ] Click a card — verify detail popup with all fields and topProducts

## Critical Files

| Action | File |
|--------|------|
| Modify | `proto/vend-analytics.proto` (add profile types) |
| Create | `go/vend/analytics/profiles/ProfileService.go` |
| Create | `go/vend/analytics/profiles/ProfileServiceCallback.go` |
| Modify | `go/vend/services/activate_analytics.go` |
| Create | `go/vend/inv_vend/profiles.go` (incremental updater) |
| Modify | `go/vend/inv_vend/analytics.go` (wire profile update, remove productPrices map) |
| Modify | `go/vend/inv_vend/performance.go` (refactor to use profiles, remove slotStates) |
| Modify | `go/vend/inv_vend/main.go` (register profile types) |
| Modify | `go/vend/inv_vend/retention.go` (add profile decay) |
| Modify | `go/vend/ui/shared.go` (register profile types) |
| Modify | `go/vend/ui/web/vend-ui/fleet/fleet-config.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/fleet-section-config.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-columns.js` (add VendMachineProfile) |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-forms.js` (add VendMachineProfile) |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-enums.js` (add primaryKey) |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-columns.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-forms.js` |
| Modify | `go/vend/ui/web/m/js/fleet/fleet-columns.js` |
| Modify | `go/vend/ui/web/m/js/fleet/fleet-forms.js` |
| Modify | `go/vend/ui/web/m/js/nav-configs/layer8m-nav-config-vend.js` |
| Create | `go/tests/mocks/gen_profiles.go` |
