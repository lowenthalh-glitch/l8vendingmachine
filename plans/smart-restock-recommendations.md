# Smart Restock Recommendations

## Context

The vending machine system collects inventory snapshots every 5 minutes with fill %, revenue, daily revenue, and per-slot stock levels. Historical data spans 30 days. Currently, restocking is reactive — operators see low-stock machines and respond. This plan adds proactive restock recommendations based on historical demand patterns, predicting WHEN to restock BEFORE machines hit critical levels.

## Recommendation Entity

**New proto type: `VendRestockRecommendation`**

```protobuf
message VendRestockRecommendation {
    string recommendation_id = 1;
    string machine_id = 2;
    string machine_name = 3;
    string location = 4;
    VendRestockPriority priority = 5;
    string reason = 6;
    VendRestockReasonCode reason_code = 7;
    int64 predicted_empty_time = 8;
    int32 current_fill_pct = 9;
    int32 projected_fill_pct = 10;
    int64 revenue_at_risk = 11;
    repeated VendRestockItem suggested_products = 12;
    string route_group_id = 13;
    double confidence = 14;
    int64 created_at = 15;
    int64 expires_at = 16;
    int32 revenue_rank = 17;
    int64 avg_daily_revenue = 18;
    string location_class = 19;
}

enum VendRestockPriority {
    VEND_RESTOCK_PRIORITY_UNSPECIFIED = 0;
    VEND_RESTOCK_PRIORITY_LOW = 1;
    VEND_RESTOCK_PRIORITY_MEDIUM = 2;
    VEND_RESTOCK_PRIORITY_HIGH = 3;
    VEND_RESTOCK_PRIORITY_CRITICAL = 4;
}

enum VendRestockReasonCode {
    VEND_RESTOCK_REASON_UNSPECIFIED = 0;
    VEND_RESTOCK_REASON_WEEKEND_DEMAND = 1;
    VEND_RESTOCK_REASON_WEEKDAY_DEMAND = 2;
    VEND_RESTOCK_REASON_RUSH_HOUR = 3;
    VEND_RESTOCK_REASON_FAST_MOVER_EMPTY = 4;
    VEND_RESTOCK_REASON_CASCADE_THRESHOLD = 5;
    VEND_RESTOCK_REASON_ROUTE_GROUPING = 6;
    VEND_RESTOCK_REASON_SEASONAL_TREND = 7;
    VEND_RESTOCK_REASON_EVENT_DRIVEN = 8;
    VEND_RESTOCK_REASON_REVENUE_AT_RISK = 9;
    VEND_RESTOCK_REASON_CRITICAL_PREDICTION = 10;
}

message VendRestockItem {
    string product_name = 1;
    int32 current_stock = 2;
    int32 capacity = 3;
    int32 units_to_add = 4;
    double depletion_rate = 5;
}

message VendRestockRecommendationList {
    repeated VendRestockRecommendation list = 1;
    l8api.L8MetaData metadata = 2;
}
```

## Scenarios

### Scenario 1: Day-of-Week Demand Patterns

**Problem:** A mall machine sells 2x on weekends. Thursday evening it's at 50%. By Saturday afternoon it'll be empty — but no one checks until Monday.

**Computation:**
1. Group 30 days of snapshots by day-of-week per machine
2. Calculate average depletion rate per day-of-week (fill% drop per hour, excluding restock jumps where fill% increases >10%)
3. Detect "weekend-heavy" machines (weekend depletion > 1.5x weekday) and "weekday-heavy" machines (inverse)
4. On Thursday evening, check weekend-heavy machines: project fill% through Sunday using weekend depletion rate
5. If projected fill% drops below 20% during the peak period → generate recommendation

**Output:**
- priority: HIGH
- reason: "Weekend demand 2.1x higher than weekday. Projected to hit 12% fill by Saturday 6 PM."
- reasonCode: WEEKEND_DEMAND
- confidence: 0.85 (4 weeks of data)

### Scenario 2: Rush Hour Depletion

**Problem:** An office machine drops 25 units between 11 AM-1 PM lunch rush. At 10 AM it has 30 units — it'll be nearly empty by 1 PM.

**Computation:**
1. Group snapshots by hour-of-day per machine across all 30 days
2. Identify "rush hours" where hourly depletion > 2x daily average
3. Common patterns: morning (7-9 AM), lunch (11 AM-1 PM), afternoon (2-4 PM), evening (5-7 PM)
4. At rush-minus-2-hours, check: can current stock survive the rush? If (currentStock - rushRate × rushDuration) < 15% capacity → recommend

**Output:**
- priority: HIGH (rush starts in < 2 hours) or MEDIUM (rush starts in 2-6 hours)
- reason: "Lunch rush (11 AM-1 PM) depletes 23 units/hour. Current stock supports only 1.2 hours. Restock before 10 AM."
- reasonCode: RUSH_HOUR
- suggestedProducts: prioritized by rush-specific depletion rate

### Scenario 3: Product-Specific Fast Movers

**Problem:** Coca-Cola sells 5x faster than Kombucha. The machine shows 60% fill overall, but the top 3 products will be empty in 3 hours.

**Computation:**
1. Track per-slot stock deltas between consecutive snapshots per machine
2. Compute average hourly depletion rate per product over past 7 days
3. Rank products: top 20% are "fast movers"
4. For each fast mover, calculate time-to-empty = currentStock / depletionRate
5. If any fast mover time-to-empty < 8 hours → recommend
6. Also flag "dead stock" (depletion < 0.5 units/day) — wasted capacity

**Output:**
- reason: "Coca-Cola (slot 3) depletes at 4.2 units/hour, empty in 3.1 hours. Doritos (slot 7) at 3.8/hr, empty in 4.0 hours."
- reasonCode: FAST_MOVER_EMPTY
- suggestedProducts: sorted by time-to-empty (shortest first), with unitsToAdd = capacity - currentStock
- revenueAtRisk: sum(depletionRate × price × hoursUntilEmpty) for products that will empty before next restock

### Scenario 4: Location-Based Patterns

**Problem:** Office machines and mall machines have opposite demand curves. A one-size-fits-all restock schedule wastes trips.

**Computation:**
1. Classify machines by their depletion signature (not metadata):
   - **Office**: high Mon-Fri, low Sat-Sun (weekday/weekend ratio > 2.0)
   - **Retail/Mall**: high Sat-Sun, moderate weekday (weekend/weekday ratio > 1.5)
   - **Transit**: consistent daily with morning/evening peaks
   - **Mixed**: no strong pattern
2. Apply class-specific restock schedules:
   - Office: heavy Monday AM + Wednesday PM
   - Mall: heavy Friday PM + Sunday AM
   - Transit: daily during off-peak
3. Generate recommendations based on class, not universal thresholds

**Output:**
- locationClass: "office" / "retail" / "transit" / "mixed"
- reason: "Office-pattern machine. Monday depletion 3.2x Sunday. Pre-stock by Sunday evening."

### Scenario 5: Low-Stock Cascade Effect

**Problem:** When 30%+ slots are empty, customers see a "picked-over" machine and walk away. Revenue drops disproportionately — a machine at 25% fill earns less than 25% of full revenue.

**Computation:**
1. From historical data, plot dailyRevenue vs fillPct for each machine
2. Find the "cascade threshold" — the fill% where revenue/fill% ratio drops sharply (typically 30-40%)
3. Detect numerically: find the fill% where (deltaRevenue / deltaFill) > 2x the average = inflection point
4. Generate CRITICAL recommendation when machine approaches its cascade threshold

**Output:**
- priority: CRITICAL (within 5% of cascade threshold)
- reason: "Machine at 33% fill (8 of 24 slots empty). Historical data shows revenue drops 47% below 30% fill. 2.1 hours to cascade threshold."
- reasonCode: CASCADE_THRESHOLD
- revenueAtRisk: (normalHourlyRevenue - postCascadeHourlyRevenue) × hoursUntilRestock

### Scenario 6: Route Optimization (Grouping)

**Problem:** Machine A needs critical restock. Machine B, 0.5 km away, is at MEDIUM. Sending a truck to A without servicing B means another trip tomorrow.

**Computation:**
1. Group machines by location zone (name prefix, city, or GPS cluster if available)
2. Apply "pull-forward" rule: if Machine A in a zone is HIGH/CRITICAL, and Machine B in same zone is MEDIUM and will need restock within 12 hours → upgrade B to HIGH
3. Assign shared routeGroupId to grouped machines
4. Sort machines within route by priority then predicted empty time

**Output:**
- routeGroupId: shared zone identifier
- reason on pulled-forward machine: "Grouped with Machine A (CRITICAL). Currently MEDIUM but needs restock in 8 hours. Servicing together saves a trip."

### Scenario 7: Seasonal/Trend Detection

**Problem:** Summer is starting — beverage sales are climbing week over week. Current restock schedules based on last month's winter data are insufficient.

**Computation:**
1. Compute 7-day rolling average depletion rate per machine
2. Compare week 4 vs week 1 trend
3. If trend > 1.2x (20% increase over 4 weeks), project forward and inflate restock quantities
4. Apply calendar heuristics (configurable multipliers):
   - Summer: beverages +20%
   - Winter: hot drinks +20%
   - Back-to-school (Aug-Sep): snacks near campus +20%

**Output:**
- reason: "Depletion rate increased 35% over past 4 weeks. Increasing restock quantity by 25%."
- reasonCode: SEASONAL_TREND
- confidence: 0.5-0.6 (30 days is thin for seasonal inference)

### Scenario 8: Event-Driven Demand

**Problem:** July 4th is in 2 days. Mall machines will see 1.8x normal demand, but the restock schedule is business-as-usual.

**Computation:**
1. Detect past anomalies: days where depletion > 2x day-of-week average
2. Match against known holidays (configurable list)
3. Apply location-class adjustments:
   - Office machines: demand drops on holidays (×0.3)
   - Mall machines: demand spikes (×1.8)
   - Transit: moderate increase (×1.2)
4. Generate recommendations 2 days before known holidays for spike-pattern machines

**Output:**
- reason: "July 4th in 2 days. Mall-pattern machine historically sees 1.8x demand on holidays."
- reasonCode: EVENT_DRIVEN
- confidence: 0.7 for known holidays with data, 0.4 for first-time

### Scenario 9: Revenue-Based Priority Adjustment

**Problem:** Two machines both need restocking. Machine A earns $800/day, Machine B earns $100/day. They shouldn't have equal priority.

**Computation:**
1. Rank all machines by 30-day average daily revenue
2. Compute revenue-at-risk per machine: sum(depletionRate × price × hoursUntilEmpty) for all products
3. Revenue modifier: top 20% revenue machines get priority upgrade (+1 level), bottom 20% get downgrade (-1 level)
4. Applied as a modifier on top of all other scenario priorities

**Output:**
- revenueRank: percentile (e.g., top 10%)
- avgDailyRevenue: for context
- reason appended: "Revenue tier: Top 10% ($847/day). Revenue at risk: $234 if not restocked by 3 PM."

### Scenario 10: Critical Threshold Prediction

**Problem:** A machine is at 45% fill right now. Based on its specific depletion pattern (this day, this time), it will hit critical at 4:30 PM. Next scheduled restock is 8 PM — 3.5 hour gap.

**Computation:**
1. Take current fill% and current time
2. Look up machine's depletion rate for current day-of-week and upcoming hours (from Scenarios 1+2)
3. Project fill% forward hour by hour:
   ```
   for each future hour h:
       projected[h] = projected[h-1] - depletionRate[dayOfWeek][hourOfDay]
       if projected[h] < criticalThreshold: predictedEmptyTime = h; break
   ```
4. Critical threshold is machine-specific (from Scenario 5 cascade threshold, or default 15%)
5. Generate tiered recommendations:
   - < 2 hours: CRITICAL
   - < 6 hours: HIGH
   - < 12 hours: MEDIUM
   - < 24 hours: LOW

**Output:**
- predictedEmptyTime: unix timestamp
- projectedFillPct: fill% at next scheduled restock
- reason: "At current rate (adjusted for Thursday afternoon), hits 15% at 4:30 PM. Next restock: 8 PM. Gap: 3.5 hours critical."
- confidence: based on depletion rate stability (CV < 0.2 → 0.9, CV > 0.5 → 0.5)

## Shared Building Blocks

All scenarios share common computation that should be extracted into reusable functions:

| Building Block | Used By | Description |
|---|---|---|
| `hourlyDepletionRate(machineId)` | All | Avg fill% drop per hour from 7-day recent data |
| `dayOfWeekProfile(machineId)` | 1, 4, 8, 10 | Depletion rate per day-of-week (7 values) |
| `hourOfDayProfile(machineId)` | 2, 10 | Depletion rate per hour-of-day (24 values) |
| `productDepletionRate(machineId, slot)` | 3, 5, 9 | Per-product units/hour depletion |
| `detectRestockEvents(snapshots)` | 1, 2, 3, 10 | Filter out fill% jumps from restock events |
| `revenueFillCurve(machineId)` | 5, 9 | Revenue vs fill% correlation |
| `classifyLocation(machineId)` | 4, 6, 8 | Office/retail/transit/mixed classification |
| `trendMultiplier(machineId)` | 7 | 4-week rolling average trend factor |

## Architecture

### New files
```
go/vend/inv_vend/restock.go               — engine orchestration, merge, dedup, expiry (<200 lines)
go/vend/inv_vend/restock_profiles.go       — building blocks (depletion profiles, classification)
go/vend/inv_vend/restock_scenarios_core.go — scenarios 1, 3, 9, 10
go/vend/inv_vend/restock_scenarios_sec.go  — scenarios 2, 4, 5
go/vend/inv_vend/restock_scenarios_adv.go  — scenarios 6, 7, 8
go/vend/analytics/restock/RestockService.go         — CRUD service
go/vend/analytics/restock/RestockServiceCallback.go — auto-generate ID on POST
```

### Data flow
```
VendInventorySnapshot (30 days)
    │
    ├── Build depletion profiles (day-of-week, hour-of-day, per-product)
    ├── Classify machine location patterns
    │
    └── Run scenarios 1-10 per machine
        │
        └── Generate/update VendRestockRecommendation
            │
            ├── Merge: one machine may trigger multiple scenarios → take highest priority, combine reasons
            ├── Dedup: one recommendation per machine (update existing, don't create duplicates)
            └── Expire: delete recommendations older than 24 hours
```

### Computation schedule
- Full profile rebuild: every 6 hours (day-of-week, hour-of-day, location classification)
- Recommendation refresh: every 30 minutes (scenarios against current data + cached profiles)
- Expiry cleanup: every hour (delete stale recommendations)

## Proto-to-JS Field Name Mapping

| Proto field | JSON name | JS column key |
|---|---|---|
| recommendation_id | recommendationId | recommendationId |
| machine_id | machineId | machineId |
| machine_name | machineName | machineName |
| location | location | location |
| priority | priority | priority |
| reason | reason | reason |
| reason_code | reasonCode | reasonCode |
| predicted_empty_time | predictedEmptyTime | predictedEmptyTime |
| current_fill_pct | currentFillPct | currentFillPct |
| projected_fill_pct | projectedFillPct | projectedFillPct |
| revenue_at_risk | revenueAtRisk | revenueAtRisk |
| suggested_products | suggestedProducts | suggestedProducts |
| route_group_id | routeGroupId | routeGroupId |
| confidence | confidence | confidence |
| created_at | createdAt | createdAt |
| expires_at | expiresAt | expiresAt |
| revenue_rank | revenueRank | revenueRank |
| avg_daily_revenue | avgDailyRevenue | avgDailyRevenue |
| location_class | locationClass | locationClass |

## Traceability Matrix

| # | Section | Action Item | Phase |
|---|---------|-------------|-------|
| 1 | Recommendation Entity | Add VendRestockRecommendation + enums + VendRestockItem to proto | Phase 1 |
| 2 | Recommendation Entity | Run make-bindings.sh | Phase 1 |
| 3 | Architecture | Create RestockService.go + RestockServiceCallback.go (auto-generate ID on POST) | Phase 1 |
| 4 | Architecture | Add to activate_analytics.go | Phase 1 |
| 5 | Architecture | Register types with AddPrimaryKeyDecorator in UI shared + inv_vend | Phase 1 |
| 6 | Shared Building Blocks | Create restock_profiles.go (depletion profiles, classification, restock detection) | Phase 1 |
| 7 | Scenario 1 | Day-of-week demand evaluation in restock_scenarios_core.go | Phase 2 |
| 8 | Scenario 3 | Product-specific fast movers in restock_scenarios_core.go | Phase 2 |
| 9 | Scenario 9 | Revenue-based priority adjustment in restock_scenarios_core.go | Phase 2 |
| 10 | Scenario 10 | Critical threshold prediction in restock_scenarios_core.go | Phase 2 |
| 11 | Architecture | Create restock.go (engine orchestration, merge, dedup) | Phase 2 |
| 12 | Computation schedule | Wire computeRestockRecommendations goroutine (every 30 min) | Phase 2 |
| 13 | Computation schedule | Profile rebuild timer (every 6 hours) | Phase 2 |
| 14 | Computation schedule | Expiry cleanup timer (every hour, delete >24h) | Phase 2 |
| 15 | Scenario 2 | Rush hour depletion in restock_scenarios_sec.go | Phase 3 |
| 16 | Scenario 4 | Location-based patterns in restock_scenarios_sec.go | Phase 3 |
| 17 | Scenario 5 | Low-stock cascade in restock_scenarios_sec.go | Phase 3 |
| 18 | Scenario 6 | Route optimization in restock_scenarios_adv.go | Phase 4 |
| 19 | Scenario 7 | Seasonal trend detection in restock_scenarios_adv.go | Phase 4 |
| 20 | Scenario 8 | Event-driven demand in restock_scenarios_adv.go | Phase 4 |
| 21 | UI | Add Restock Recommendations to Analytics config + section config | Phase 5 |
| 22 | UI | Add desktop columns (all proto fields listed or omission documented) | Phase 5 |
| 23 | UI | Add desktop forms with suggestedProducts inline table | Phase 5 |
| 24 | UI | Add desktop enums for VendRestockPriority + VendRestockReasonCode | Phase 5 |
| 25 | UI | Add dashboard "Urgent Restocks" widget | Phase 5 |
| 26 | Mobile Parity | Add mobile columns with primary/secondary markers | Phase 5 |
| 27 | Mobile Parity | Add mobile forms | Phase 5 |
| 28 | Mobile Parity | Add mobile enums | Phase 5 |
| 29 | Mobile Parity | Update mobile nav config for restock service | Phase 5 |
| 30 | Mock Data | Seed sample recommendations with non-zero values for all fields | Phase 6 |
| 31 | Mock Data | Populate VendRestockItem repeated field in seeded data | Phase 6 |

## Phase 1: Proto + Service + Building Blocks

1. Add `VendRestockRecommendation`, `VendRestockItem`, enums to proto
2. Run make-bindings.sh
3. Create `RestockService.go` (ServiceName="Restock", area 10) with `RestockServiceCallback.go` that auto-generates `RecommendationId` on POST via `common.GenerateID`
4. Add to `activate_analytics.go`
5. Register types with `Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendRestockRecommendation{}, "RecommendationId")` in UI shared.go and inv_vend main.go
6. Create `restock_profiles.go` with shared building blocks:
   - `buildDepletionProfiles(snapshots)` → returns per-machine profiles
   - `classifyLocations(profiles)` → returns location class per machine
   - `detectRestockEvents(snapshots)` → filters restock jumps

## Phase 2: Core Scenarios (1, 3, 9, 10)

Implement the four most impactful scenarios first:
- **Scenario 1**: Day-of-week demand (the user's original request)
- **Scenario 3**: Product-specific fast movers (most actionable)
- **Scenario 9**: Revenue-based priority (highest business value)
- **Scenario 10**: Critical threshold prediction (prevents stockouts)

Create `restock.go` (~200 lines) with:
- `computeRestockRecommendations(nic)` — main engine, runs every 30 min
- `mergeRecommendations(candidates)` → one per machine, highest priority wins
- `expireOldRecommendations(nic)` → delete recommendations older than 24 hours
- Goroutine timers: 30-min refresh, 6-hour profile rebuild, 1-hour expiry

Create `restock_scenarios_core.go` (~300 lines) with:
- `evaluateDayOfWeekDemand(machine, profile)` → candidate recommendation
- `evaluateFastMovers(machine, slotProfiles)` → candidate recommendation
- `evaluateCriticalPrediction(machine, profile)` → candidate recommendation
- `applyRevenuePriority(recommendations, revenueRanks)` → adjusts priorities

## Phase 3: Secondary Scenarios (2, 4, 5)

- **Scenario 2**: Rush hour depletion
- **Scenario 4**: Location-based patterns
- **Scenario 5**: Low-stock cascade

## Phase 4: Advanced Scenarios (6, 7, 8)

- **Scenario 6**: Route optimization grouping
- **Scenario 7**: Seasonal trend detection
- **Scenario 8**: Event-driven demand

## Phase 5: UI — Restock Recommendations Tab

### Desktop Analytics Config
- New service: `{ key: 'restock', label: 'Restock Recommendations', endpoint: '/10/Restock', model: 'VendRestockRecommendation', readOnly: true, defaultSort: { column: 'priority', direction: 'desc' } }`

### Desktop Enums (`data-enums.js`)
```js
VendRestockPriority: { 0: 'Unspecified', 1: 'Low', 2: 'Medium', 3: 'High', 4: 'Critical' }
VendRestockReasonCode: { 0: 'Unspecified', 1: 'Weekend Demand', 2: 'Weekday Demand', 3: 'Rush Hour',
    4: 'Fast Mover Empty', 5: 'Cascade Threshold', 6: 'Route Grouping', 7: 'Seasonal Trend',
    8: 'Event Driven', 9: 'Revenue At Risk', 10: 'Critical Prediction' }
```

### Desktop Columns (`data-columns.js`)
All proto fields covered:
```js
VendRestockRecommendation: [
    col.status('priority', 'Priority', PRIORITY_VALUES, render.priority),
    col.col('machineName', 'Machine'),
    col.col('reason', 'Reason'),
    col.enum('reasonCode', 'Reason Code', null, render.reasonCode),
    col.date('predictedEmptyTime', 'Predicted Empty'),
    col.number('currentFillPct', 'Current Fill %'),
    col.number('projectedFillPct', 'Projected Fill %'),
    col.money('revenueAtRisk', 'Revenue At Risk'),
    col.number('confidence', 'Confidence'),
    col.col('locationClass', 'Location Type'),
    col.number('revenueRank', 'Revenue Rank'),
    col.money('avgDailyRevenue', 'Avg Daily Revenue'),
    col.col('routeGroupId', 'Route Group'),
    col.date('createdAt', 'Created'),
    col.date('expiresAt', 'Expires')
]
```
// Omitted from columns: recommendationId (primary key, shown in detail), machineId (machineName shown instead), location (locationClass shown instead), suggestedProducts (shown in form inline table)

### Desktop Forms (`data-forms.js`)
```js
VendRestockRecommendation: f.form('Restock Recommendation', [
    f.section('Recommendation', [
        f.text('recommendationId', 'ID', false, { readOnly: true }),
        f.text('machineName', 'Machine', false, { readOnly: true }),
        f.text('location', 'Location', false, { readOnly: true }),
        f.select('priority', 'Priority', PRIORITY_ENUM, false, { readOnly: true }),
        f.text('reason', 'Reason', false, { readOnly: true }),
        f.select('reasonCode', 'Reason Code', REASON_CODE_ENUM, false, { readOnly: true }),
        f.text('locationClass', 'Location Type', false, { readOnly: true })
    ]),
    f.section('Prediction', [
        f.date('predictedEmptyTime', 'Predicted Empty', false, { readOnly: true }),
        f.number('currentFillPct', 'Current Fill %', false, { readOnly: true }),
        f.number('projectedFillPct', 'Projected Fill %', false, { readOnly: true }),
        f.money('revenueAtRisk', 'Revenue At Risk', false, { readOnly: true }),
        f.number('confidence', 'Confidence', false, { readOnly: true }),
        f.number('revenueRank', 'Revenue Rank', false, { readOnly: true }),
        f.money('avgDailyRevenue', 'Avg Daily Revenue', false, { readOnly: true })
    ]),
    f.section('Suggested Products', [
        f.inlineTable('suggestedProducts', 'Products to Restock', [
            { key: 'productName', label: 'Product', type: 'text' },
            { key: 'currentStock', label: 'Current', type: 'number' },
            { key: 'capacity', label: 'Capacity', type: 'number' },
            { key: 'unitsToAdd', label: 'Units to Add', type: 'number' },
            { key: 'depletionRate', label: 'Depletion/hr', type: 'number' }
        ])
    ]),
    f.section('Metadata', [
        f.text('routeGroupId', 'Route Group', false, { readOnly: true }),
        f.date('createdAt', 'Created', false, { readOnly: true }),
        f.date('expiresAt', 'Expires', false, { readOnly: true })
    ])
])
```

### Dashboard
- "Urgent Restocks" widget showing CRITICAL + HIGH recommendations count with link to Analytics

### Mobile Parity

**Mobile columns** (`m/js/analytics/analytics-columns.js`):
```js
VendRestockRecommendation: [
    { key: 'machineName', label: 'Machine', primary: true },
    { key: 'priority', label: 'Priority', secondary: true, render: render.priority },
    col.col('reason', 'Reason'),
    col.date('predictedEmptyTime', 'Predicted Empty'),
    col.number('currentFillPct', 'Fill %'),
    col.money('revenueAtRisk', 'Revenue At Risk')
]
```

**Mobile forms** (`m/js/analytics/analytics-forms.js`): Same field coverage as desktop form.

**Mobile enums** (`m/js/analytics/analytics-enums.js`): Same VendRestockPriority and VendRestockReasonCode enums.

**Mobile nav config** (`m/js/nav-configs/layer8m-nav-config-vend.js`): Add restock service under analytics:
```js
{ key: 'restock', label: 'Restock', icon: 'analytics', endpoint: '/10/Restock',
  model: 'VendRestockRecommendation', idField: 'recommendationId', readOnly: true }
```

## Phase 6: Mock Data

Create `gen_restock.go` to seed sample recommendations:
- Generate 15-20 recommendations across different machines
- Cover all priority levels (CRITICAL, HIGH, MEDIUM, LOW)
- Cover multiple reason codes (at least 4 different codes)
- Populate ALL fields with non-zero values including:
  - `suggestedProducts` repeated field with 2-4 VendRestockItem entries per recommendation
  - `revenueAtRisk`, `avgDailyRevenue` with realistic cent values
  - `confidence` between 0.5-0.95
  - `predictedEmptyTime` and `expiresAt` as future timestamps
  - `locationClass` with "office", "retail", "transit" values

## Phase 7: Verification

### Build
- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] All scenario files under 500 lines

### Desktop
- [ ] Navigate to Analytics > Restock Recommendations — verify table loads with data
- [ ] Verify CRITICAL/HIGH recommendations sorted to top
- [ ] Click a recommendation — verify detail popup with all sections (Recommendation, Prediction, Suggested Products, Metadata)
- [ ] Verify suggestedProducts inline table shows product rows
- [ ] Verify priority and reasonCode render as status badges (not raw numbers)
- [ ] Verify dashboard shows "Urgent Restocks" widget with count

### Scenario Logic
- [ ] Verify day-of-week pattern detection (weekend-heavy machines flagged by Thursday)
- [ ] Verify fast mover detection (products with high depletion rate flagged)
- [ ] Verify critical prediction (machines with <6h to critical threshold are HIGH/CRITICAL)
- [ ] Verify revenue priority adjustment (high-revenue machines get priority boost)
- [ ] Verify recommendations expire after 24 hours (old entries removed)

### Mobile
- [ ] Navigate to Analytics > Restock — verify card list loads
- [ ] Verify primary (machineName) and secondary (priority) fields show on cards
- [ ] Click a card — verify detail popup with all fields
- [ ] Verify suggestedProducts shows in detail

## Critical Files

| Action | File |
|--------|------|
| Modify | `proto/vend-analytics.proto` (add recommendation types) |
| Create | `go/vend/analytics/restock/RestockService.go` |
| Modify | `go/vend/services/activate_analytics.go` |
| Create | `go/vend/inv_vend/restock_profiles.go` (building blocks) |
| Create | `go/vend/inv_vend/restock.go` (recommendation engine) |
| Modify | `go/vend/inv_vend/main.go` (wire goroutine + register types) |
| Modify | `go/vend/ui/shared.go` (register types with AddPrimaryKeyDecorator) |
| Modify | `go/vend/ui/web/vend-ui/analytics/analytics-config.js` |
| Modify | `go/vend/ui/web/vend-ui/analytics/analytics-section-config.js` |
| Modify | `go/vend/ui/web/vend-ui/analytics/data/data-columns.js` |
| Modify | `go/vend/ui/web/vend-ui/analytics/data/data-forms.js` |
| Modify | `go/vend/ui/web/vend-ui/analytics/data/data-enums.js` |
| Modify | `go/vend/ui/web/vend-ui/dashboard/dashboard-init.js` |
| Modify | `go/vend/ui/web/m/js/analytics/analytics-columns.js` |
| Modify | `go/vend/ui/web/m/js/analytics/analytics-forms.js` |
| Modify | `go/vend/ui/web/m/js/analytics/analytics-enums.js` |
| Modify | `go/vend/ui/web/m/js/nav-configs/layer8m-nav-config-vend.js` |
| Create | `go/tests/mocks/gen_restock.go` |
