# Plan: OptimoRoute Feature Parity

## Overview

Implement OptimoRoute-equivalent capabilities for the vending machine route optimizer. This plan compares OptimoRoute's feature set against our current implementation and defines phases to close the gaps.

---

## Feature Gap Analysis

### What We Already Have
| Feature | Our Implementation |
|---------|-------------------|
| Basic route optimization | Nearest-neighbor + 2-opt |
| Machine demand classification | List A (urgent) / List B (can wait) |
| Driver-truck matching | License class + schedule day |
| Dynamic facility reloads | Optimal facility placement by detour cost |
| Traffic integration | Google Maps Directions API (optional final refinement) |
| Fuel cost estimation | Distance / MPG × fuel price |
| 8-hour max route duration | Cluster cap during building |
| Interactive map | Leaflet with filters, route polylines |
| Generate Routes UI | Date + start time picker |

### Fundamental Engine Gaps (Must Fix First)

These are architectural flaws in the current optimizer that affect ALL routes. They must be fixed before adding new features.

**0.1 Driver-Aware Machine Assignment**
Currently: machines are clustered by proximity to EACH OTHER, then a driver is randomly assigned. A driver starting in North Austin may get a cluster of South Austin machines, passing through clusters assigned to other drivers.

Should be: assign machines to the driver whose start/end location is nearest. Each driver "owns" a geographic zone based on their position.

Algorithm change:
1. For each driver available on the planned day, compute their start location (from schedule) and end location (home address or home depot facility)
2. Assign each machine to the driver whose start location is nearest (weighted: 60% proximity to start, 40% proximity to end)
3. Within each driver's assigned machines, apply the existing nearest-neighbor + 2-opt ordering
4. This replaces the current "cluster first, assign driver later" approach with "assign driver first, build route per driver"

**0.2 End-of-Day Location in Route Optimization**
Currently: the route ends at the last machine stop. The driver still has to drive home or to a facility — that distance/fuel is unaccounted for. The 2-opt optimization doesn't know about the return trip, so it may leave the last stop far from home.

Should be: the driver's end-of-day location is a virtual final waypoint. The optimization minimizes total distance INCLUDING the drive home. See Phase 0 `route_builder.go` for implementation details.

**0.3 Facility Reloads Consider End Location**
Currently: facility selection minimizes detour to remaining machines' centroid only.

Should be: near end-of-route, prefer a facility on the way home. If the truck's home depot IS a facility, use it for the last reload (combining reload + end-of-day parking). See Phase 0 `route_builder.go` point 4 for implementation.

### What OptimoRoute Has That We Don't

#### Priority 1 — Core Optimization Gaps

**1.1 Time Windows**
OptimoRoute lets each order/stop have a time window (e.g., "deliver between 9am-12pm"). Machines in high-traffic locations (offices, hospitals) may only be accessible during certain hours.
- **Proto change:** Add `accessWindowStart` and `accessWindowEnd` (HH:MM strings) to `VendFleetMachine` or `VendLocation`
- **Optimizer change:** When ordering stops, respect time windows — don't schedule a machine before its window opens or after it closes. Reorder stops if needed to fit windows.

**1.2 Driver Working Hours & Overtime**
OptimoRoute tracks max working hours per driver and eliminates overtime. We cap routes at 8 hours but don't check against each driver's actual schedule end time.
- **Optimizer change:** Read driver's `schedule[day].startTime`, calculate end time (start + 8 hours or configurable shift length), reject routes that exceed the driver's actual working window.
- **Proto change:** Add `shiftDurationMinutes` to `VendDriverScheduleDay` (default 480 = 8 hours)

**1.3 Driver Breaks**
OptimoRoute plans lunch breaks and regulated breaks. A 6+ hour route needs a break.
- **Optimizer change:** After every 4 hours of driving/service, insert a 30-minute break stop. If a facility reload is near the break point, combine them.
- **Proto change:** Add `breakDurationMinutes` and `breakAfterMinutes` to `VendRouteOptRequest` (defaults: 30 and 240)
- Break appears as a stop with `stopType = "break"`

**1.4 Workload Balancing**
OptimoRoute distributes work evenly across drivers — by hours, order count, or mileage. Phase 0's `AssignMachinesToDrivers` assigns by geography with an 8-hour cap, but doesn't actively equalize. One driver may get 18 stops and another 8.
- **Optimizer change:** After Phase 0's assignment, run a balancing pass that moves machines between drivers to equalize workload. This is a POST-PROCESSING step on Phase 0's output, not a replacement.
  - Balance by stop count: move from most-loaded to least-loaded driver, preferring machines near the receiving driver
  - Balance by estimated duration: same, using duration instead of count
- **Config:** Add `balanceMode` to `VendRouteOptRequest`: `"stops"`, `"duration"`, or `"none"` (default)
- **Note:** Phase 0 handles the geographic assignment and cap. Phase 2 refines the distribution. No duplicate logic — Phase 2 calls a separate `BalanceWorkload()` function that takes Phase 0's output as input.

**1.5 Vehicle Capacity Constraints**
OptimoRoute respects vehicle weight/volume limits. We check `cargoCapacityCuFt` during clustering but don't enforce it per-product-volume.
- **Optimizer change:** Track cumulative product volume/weight as stops are added. If the truck can't physically carry enough product for the next machine, trigger a reload.
- This is partially handled by stock depletion already, but explicit volume tracking is more accurate.

**1.6 Priority Sequencing**
OptimoRoute supports priority levels (urgent, high, normal). We have urgent vs low but don't ORDER by priority — just classify.
- **Optimizer change:** Within a cluster, serve high-priority stops first (even if slightly farther), then lower priority. Add priority weight to the nearest-neighbor selection.

#### Priority 2 — Operational Features

**2.1 Multi-Day Planning**
OptimoRoute plans up to 5 weeks ahead. We only generate routes for a single day.
- **Optimizer change:** Accept a date range in `VendRouteOptRequest` (`plannedDateEnd`). Generate routes for each day in the range, factoring in which drivers work each day.
- Machines served on day 1 are excluded from day 2's demand list (their stock was restocked).

**2.2 Real-Time Reoptimization**
OptimoRoute allows mid-route changes: insert last-minute stops, drag-and-drop reordering, handle sick driver.
- **UI change:** Allow editing a route's stop order via drag-and-drop in the route detail popup
- **Backend:** Add a `PUT /vend/10/OptRoute` endpoint that takes an existing route ID and recalculates with a new/removed stop
- **Sick driver:** Reassign a route's driver+truck by finding another available pair

**2.3 Order/Job Duration**
OptimoRoute lets each stop have a different service duration. We use a fixed 20 min per stop.
- **Proto change:** Add `estimatedServiceMinutes` to `VendFleetMachine` (default 20)
- **Optimizer change:** Use per-machine service time instead of global constant
- Machines with more slots or harder access take longer

**2.4 Skills Matching**
OptimoRoute matches driver skills to job requirements. Some machines may need refrigeration-certified drivers or special training.
- **Proto change:** Add `requiredSkills` (repeated string) to `VendFleetMachine`, `skills` (repeated string) to `VendDriver`
- **Optimizer change:** When assigning driver to cluster, verify driver has all required skills for machines in the cluster

#### Priority 3 — Customer/Field Features

**3.1 Real-Time Driver Tracking**
OptimoRoute shows driver location in real-time with ETA per stop.
- **We already have:** Driver `currentLatitude/currentLongitude` updated on device connect
- **Gap:** No ETA calculation. The map shows driver position but not "arriving at machine X in 15 min"
- **Implementation:** When a route is in progress, calculate ETA per remaining stop from driver's current position using Google Maps

**3.2 Proof of Service**
OptimoRoute captures photos, signatures, and forms at each stop.
- **Proto change:** Add `completionPhoto`, `completionNotes`, `completedAt` to `VendRouteStop`
- **Mobile app:** At each stop, driver takes a photo of the restocked machine, adds notes
- **Backend:** PUT to update the route stop with completion data

**3.3 Notifications**
OptimoRoute sends SMS/email notifications with ETA to customers.
- For vending machines, the "customer" is the location manager
- **Implementation:** Use l8notify to send email to `VendLocation.contactEmail` when the truck is N minutes away

**3.4 Route Analytics & Reporting**
OptimoRoute provides breadcrumbs (planned vs actual route), performance metrics, arrival accuracy.
- **Proto change:** Add `actualDistance`, `actualDuration`, `actualFuelCost` to `VendRoute` (populated after route completion)
- **UI:** Side-by-side comparison of planned vs actual on the map
- **Analytics:** Average arrival accuracy, time per stop, fuel efficiency trends

**3.5 Live Customer Tracking Link**
OptimoRoute provides a tracking link for customers to see the driver's location.
- Not directly applicable to vending (no end customer waiting), but location managers might want to see when their machine will be restocked
- **Deferred** — low priority for vending use case

---

## Implementation Phases

### Phase 0: Engine Fundamentals — Driver-Aware Assignment + End Location + Smart Reloads
**Estimated scope:** Proto changes + major optimizer rewrite

**Proto changes:**
1. Add `end_location_id` (string) to `VendDriverScheduleDay` — end-of-day location (blank = home address)
2. Add `stop_type = "end"` support in `VendRouteStop` — final leg to driver's end location
3. Run `cd proto && ./make-bindings.sh` after proto changes
4. Verify: `go build ./...`

**End-location resolution fallback chain** (same pattern as probler — no explicit chain in ecosystem, use parallel fields):
1. `schedule[day].endLocationId` → look up VendLocation coordinates
2. If blank → driver's `homeAddress` → geocode or use driver's `currentLatitude/currentLongitude`
3. If no coordinates → truck's `homeDepotId` → facility coordinates (facilities always have GPS)
4. If all fail → driver's start location (round trip)

**Mock data:** Update 5 driver schedules with specific Austin end locations (pre-defined, not random — follows l8erp config data pattern):
- drv-001: ends at Austin Central Depot (fac-001)
- drv-002: ends at home address
- drv-003: ends at Austin North Hub (fac-002)
- drv-004: ends at home address
- drv-005: ends at Austin South Hub (fac-003)

**Mobile parity:** Update mobile route forms to show `stopType = "end"` stops. Update mobile map if route visualization exists.

**Optimizer rewrite — assign.go:**
1. Replace `ClusterMachines` + `AssignResources` with `AssignMachinesToDrivers`:
   - For each driver, resolve start location and end location (using fallback chain above)
   - For each machine needing restock, score against each driver: `0.6 × distance(machine, driver.start) + 0.4 × distance(machine, driver.end)`
   - Assign each machine to the driver with the lowest score
   - Cap per-driver assignments at 8-hour estimated duration
   - Overflow machines go to the next-best driver
2. Result: one machine list per driver (not anonymous clusters)

**Optimizer rewrite — route_builder.go:**
1. `BuildRoute` receives driver's start location, end location, and assigned machines
2. Nearest-neighbor ordering: start from driver's start, end near driver's end location
   - Specifically: remove the machine nearest to end location from the pool, save it as the last stop. Then nearest-neighbor the rest from the start. Append the saved last stop.
3. 2-opt includes the end-location leg: when evaluating edge swaps, the "last leg" distance is from the last machine to the end location
4. Facility reloads: when fewer than 3 stops remain, weight facility selection toward end location (50% remaining-machines centroid, 50% end location). If truck's `homeDepotId` facility is within 2× the distance of the optimal facility, prefer it (driver parks truck at home depot).
5. Add final `stopType = "end"` with the driver's end location to the route
6. `totalDistance` and `estimatedFuelCost` include the end-location leg

**Optimizer rewrite — cluster.go:**
1. `ClusterMachines` is removed (replaced by driver-aware assignment)
2. List B insertion moves to `AssignMachinesToDrivers`: a List B machine is added to a driver's list only if it's within detour threshold of that driver's existing route

**Map:**
1. Route polylines include the leg to end location (shown as a different dash pattern)
2. End location marker shown as a house icon

### Phase 0.5: OSRM Offline Routing + Traffic Statistics
**Estimated scope:** New infrastructure + optimizer integration

Replace haversine with real road distances using a local OSRM (Open Source Routing Machine) instance. Add statistical traffic multipliers based on time-of-day. Keep Google Maps as an optional final refinement layer.

**Infrastructure — OSRM Docker container:**
1. Download regional OSM extract from Geofabrik (default: Texas ~500MB, configurable for other regions)
2. Pre-process into OSRM format (one-time, ~10 min)
3. Run OSRM as a Docker container alongside the other services:
   ```bash
   docker run -d -p 5000:5000 -v /data/osrm:/data osrm/osrm-backend \
     osrm-routed --algorithm mld /data/texas-latest.osrm
   ```
4. Add OSRM container start to `run-local.sh` — idempotent (skip download if `.osm.pbf` exists, skip pre-process if `.osrm` exists, skip container start if already running). No flags needed — follows L8 ecosystem pattern of unconditional but idempotent setup.

**OSRM URL configuration — in login.json** (same pattern as `apiPrefix`):
```json
{
    "app": {
        "apiPrefix": "/vend",
        "osrmUrl": "http://localhost:5000",
        "osrmRegion": "texas"
    }
}
```
The optimizer reads OSRM URL from config. If not set, defaults to `http://localhost:5000`. If OSRM is unreachable, falls back to haversine silently.

**OSRM APIs used:**
- **Route API**: `GET {osrmUrl}/route/v1/driving/{lng1},{lat1};{lng2},{lat2}` → road distance + duration between two points
- **Table API**: `GET {osrmUrl}/table/v1/driving/{lng1},{lat1};{lng2},{lat2};...` → distance/duration matrix for N points in one call
- **Important:** The Table API is called **per-driver** (15-20 machines per route after assignment), NOT for all 300 machines. A 20×20 matrix = 400 elements, <10ms. Never build a 300×300 matrix.

**Distance abstraction — consolidate `distance.go` + new `routing.go` into a unified interface:**

Rename `distance.go` → keep only `ComputeRouteMetrics` and `Centroid` (pure math, no I/O).

New `routing.go` provides a single `Router` that all optimizer code calls:
```go
type Router struct {
    osrmURL    string
    httpClient *http.Client  // persistent, created once (l8collector RestCollector pattern)
    available  bool          // set to false on first OSRM failure, retried periodically
}

func NewRouter(osrmURL string) *Router  // creates persistent http.Client with connection pooling
func (r *Router) Distance(lat1, lng1, lat2, lng2 float64) (distMiles float64, durationSecs int64)
func (r *Router) Matrix(points [][2]float64) (distMatrix [][]float64, durMatrix [][]int64)
```
- `httpClient` created once in `NewRouter()`, reused across all calls (same lifecycle as l8collector's `RestCollector.httpClient`)
- Tries OSRM first via persistent client
- Falls back to haversine if OSRM unavailable (sets `available = false`, retries after 5 min)
- Single source for all distance calls — `assign.go`, `route_builder.go`, `optimize.go` all use `Router`, never call haversine directly
- Eliminates duplicate distance logic between `distance.go` and `routing.go`
- `Router` is created once per optimization run in `GenerateRoutes()` and passed to all functions
- **Note:** If this file exceeds 500 lines, split OSRM HTTP client into `osrm_client.go`

**Traffic statistics — `go/vend/route/optimizer/traffic_stats.go`:**
1. Time-of-day multipliers applied to OSRM base durations:
   ```
   06:00-07:00  → 1.0× (early, light traffic)
   07:00-09:00  → 1.4× (morning rush)
   09:00-11:30  → 1.1× (mid-morning)
   11:30-13:30  → 1.2× (lunch)
   13:30-16:00  → 1.1× (afternoon)
   16:00-18:30  → 1.4× (evening rush)
   18:30-22:00  → 0.9× (evening, light)
   ```
2. Day-of-week factors: weekday = 1.0×, Saturday = 0.7×, Sunday = 0.6×
3. `ApplyTrafficStats(durationSecs int64, arrivalTime int64) int64` — adjusts duration based on what time the driver arrives at each leg
4. Over time, build a learning model: after each completed route, compare planned vs actual durations per time slot and adjust multipliers.

**VendTrafficProfile — full Prime Object** (follows l8learn `GrowthRecord`/`CohortSnapshot` pattern for persisted statistical data):
```protobuf
// @PrimeObject
message VendTrafficProfile {
    string profile_id = 1;               // e.g., "austin-weekday"
    string region = 2;
    repeated VendTrafficSlot time_slots = 3;  // 24 slots (one per hour) or custom
    repeated double day_of_week_factors = 4;  // 7 values (Sun=0..Sat=6)
    int32 sample_count = 5;              // routes used to compute
    int64 last_updated = 6;
    l8common.AuditInfo audit_info = 7;
}

message VendTrafficSlot {
    int32 hour = 1;                      // 0-23
    double multiplier = 2;              // e.g., 1.4 for rush hour
    int32 sample_count = 3;
}

message VendTrafficProfileList {
    repeated VendTrafficProfile list = 1;
    l8api.L8MetaData metadata = 2;
}
```
- Needs: proto message + List type + service (`TrafProf`, area 10) + activation + type registration + mock data (seeded with default multipliers above)
- The learning function reads completed routes (`actualDuration` vs `totalDuration` per leg), buckets by hour-of-day, and updates the profile's multipliers

**Optimizer integration — all distance calls go through `Router`:**
1. **generator.go** — creates `Router` once via `NewRouter(osrmURL)`, passes to all functions
2. **assign.go** — `AssignMachinesToDrivers`: calls `router.Distance()` for driver-to-machine scoring. For 300 machines × 5 drivers, this is 1,500 calls — use `router.Matrix()` with all driver+machine points in one call instead (305 points = small matrix)
3. **route_builder.go** — `nearestNeighborOrder`: calls `router.Matrix()` for the driver's 15-20 assigned machines (one call per route). Uses cached matrix for all ordering decisions
4. **route_builder.go** — `findOptimalFacility`: calls `router.Distance()` for each facility candidate
5. **optimize.go** — `twoOpt`: uses the cached matrix from step 3 (no additional OSRM calls)
6. **generator.go** — `ComputeRouteMetrics`: legs use OSRM road distances + traffic-stat-adjusted durations
7. **traffic.go** — Google Maps becomes an optional FINAL layer: if API key exists AND user requests it, call Directions API for the final ordered route. Otherwise, OSRM + traffic stats is the production path

**Distance calculation hierarchy (highest to lowest priority):**
1. Google Maps Directions API (if API key exists AND requested) — real-time traffic
2. OSRM + traffic statistics (default) — road distance + statistical traffic
3. Haversine + fixed avg speed (fallback if OSRM unavailable) — straight-line estimate

**Cost comparison:**
| Method | Cost | Accuracy | Latency |
|--------|------|----------|---------|
| Google Maps | ~$5-10 per optimization run | High (real-time traffic) | 200-500ms per call |
| OSRM + traffic stats | $0 (self-hosted) | Good (real roads + statistical traffic) | <10ms per call |
| Haversine | $0 | Poor (straight-line, no roads) | <1ms |

**run-local.sh changes:**
1. Download Texas OSM data if not present: `wget https://download.geofabrik.de/north-america/us/texas-latest.osm.pbf`
2. Pre-process if `.osrm` file doesn't exist: `docker run osrm/osrm-backend osrm-extract/contract/partition`
3. Start OSRM container before other services
4. Add to `kill_demo.sh` cleanup

**Deployment artifacts (for K8s):**
1. K8s YAML: `k8s/osrm.yaml` — DaemonSet or StatefulSet running `osrm/osrm-backend`, volume mount for OSM data at `/data/osrm`
2. Add to `k8s/deploy.sh` — start OSRM before vend services
3. Add to `k8s/undeploy.sh`
4. No `build.sh` or `Dockerfile` needed — uses upstream `osrm/osrm-backend` image directly

**All optimizer files that use Router must stay under 500 lines.** If `routing.go` grows large, split OSRM HTTP parsing into `osrm_client.go`. If `traffic_stats.go` grows with the learning model, split into `traffic_learn.go`.

### Phase 1: Time Windows + Driver Hours + Breaks
**Estimated scope:** Proto changes + optimizer logic

1. Add `accessWindowStart`, `accessWindowEnd` to `VendLocation`
2. Add `shiftDurationMinutes` to `VendDriverScheduleDay`
3. Add `breakDurationMinutes`, `breakAfterMinutes` to `VendRouteOptRequest`
4. Update optimizer: respect time windows when ordering stops
5. Update optimizer: calculate driver end time from schedule, reject over-shift routes
6. Update optimizer: insert break stops after 4 hours
7. Update route forms to show break stops and time windows
8. Regenerate bindings, build, test

### Phase 2: Workload Balancing + Priority Sequencing
**Estimated scope:** Optimizer logic only (no proto changes)

1. Add `balanceMode` to `VendRouteOptRequest`
2. After clustering, redistribute machines to balance workload across drivers
3. Add priority weight to nearest-neighbor selection (serve urgent stops first within distance tolerance)
4. Test with 5 drivers, verify balanced stop counts and durations

### Phase 3: Per-Machine Service Duration + Skills Matching
**Estimated scope:** Proto changes + optimizer logic

1. Add `estimatedServiceMinutes` to `VendFleetMachine`
2. Add `requiredSkills` to `VendFleetMachine`, `skills` to `VendDriver`
3. Update optimizer: use per-machine service time
4. Update optimizer: filter drivers by required skills per cluster
5. Update mock data: vary service times (15-30 min), add skills to some machines/drivers
6. Update UI forms

### Phase 4: Multi-Day Planning
**Estimated scope:** Optimizer logic + UI

1. Add `plannedDateEnd` to `VendRouteOptRequest`
2. Update optimizer: loop over days, exclude machines served on previous days
3. UI: add end date picker to Generate Routes panel
4. UI: routes table shows day column, map shows routes per day

### Phase 5: Real-Time Reoptimization
**Estimated scope:** New endpoint + UI

1. Add `PUT /vend/10/OptRoute` for route modification (add/remove stop, change driver)
2. UI: drag-and-drop stop reordering in route detail popup
3. UI: "Reassign Driver" button on route
4. Recalculate metrics after modification

### Phase 6: Proof of Service + Analytics
**Estimated scope:** Proto changes + mobile + UI

1. Add completion fields to `VendRouteStop` (photo, notes, completedAt)
2. Add actual metrics to `VendRoute` (actualDistance, actualDuration)
3. Mobile: completion form at each stop
4. UI: planned vs actual comparison on map
5. Analytics dashboard: arrival accuracy, fuel efficiency, time per stop

---

## Traceability Matrix

| # | Feature | Phase | Priority |
|---|---------|-------|----------|
| 0a | Driver-aware machine assignment (replace blind clustering) | Phase 0 | P0 |
| 0b | End-of-day location in route optimization | Phase 0 | P0 |
| 0c | Facility reloads consider end location + home depot preference | Phase 0 | P0 |
| 0d | Route totalDistance/fuelCost includes drive to end location | Phase 0 | P0 |
| 0e | End-location fallback chain (schedule → home → facility → start) | Phase 0 | P0 |
| 0f | Mock data: driver end locations (3 at facilities, 2 at home) | Phase 0 | P0 |
| 0g | OSRM Docker container for offline road distance/duration | Phase 0.5 | P0 |
| 0h | OSRM URL configurable in login.json (not hardcoded) | Phase 0.5 | P0 |
| 0i | Router struct with persistent http.Client (l8collector pattern) | Phase 0.5 | P0 |
| 0j | OSRM Table API per-driver (15-20 points, not 300×300) | Phase 0.5 | P0 |
| 0k | Traffic statistics (time-of-day + day-of-week multipliers) | Phase 0.5 | P0 |
| 0l | VendTrafficProfile Prime Object (proto + service + mock data) | Phase 0.5 | P0 |
| 0m | Learning traffic model from completed route actuals | Phase 0.5 | P0 |
| 0n | Graceful fallback: OSRM → haversine, Google Maps optional top layer | Phase 0.5 | P0 |
| 0o | Idempotent OSRM setup in run-local.sh (no flags) | Phase 0.5 | P0 |
| 1 | Time windows | Phase 1 | P1 |
| 2 | Driver working hours/overtime | Phase 1 | P1 |
| 3 | Driver breaks (lunch, regulated) | Phase 1 | P1 |
| 4 | Workload balancing across drivers | Phase 2 | P1 |
| 5 | Priority sequencing | Phase 2 | P1 |
| 6 | Vehicle capacity enforcement | Phase 2 | P1 |
| 7 | Variable service duration per stop | Phase 3 | P2 |
| 8 | Skills matching (driver ↔ machine) | Phase 3 | P2 |
| 9 | Multi-day/week planning | Phase 4 | P2 |
| 10 | Real-time reoptimization (drag-drop) | Phase 5 | P2 |
| 11 | Sick driver reassignment | Phase 5 | P2 |
| 12 | Proof of service (photos, notes) | Phase 6 | P3 |
| 13 | Planned vs actual analytics | Phase 6 | P3 |
| 14 | Notifications to location managers | Phase 6 | P3 |
| 15 | Real-time ETA per stop | Phase 6 | P3 |
| 16 | Live tracking link | Deferred | P3 |
| 17 | Multi-technician jobs | N/A | — |
| 18 | 5-week advance planning | Phase 4 | P2 |
| 19 | Commercial truck routing (road restrictions) | Deferred | P3 |

---

## Rule Compliance Notes

**`events-service-required.md`** — The l8vendingmachine project does NOT currently call `evtservices.ActivateEvents()` in main.go. This is a pre-existing gap, not introduced by this plan. The route optimizer does not add or remove events activation. If events service is added to the project later, no optimizer changes are needed.

**`framework-interface-boundaries.md`** — The optimizer does NOT add methods or interfaces to `l8types/go/ifs/`. All optimizer code is in `go/vend/route/optimizer/` (project-specific implementation layer). The `OptimizerService` implements the existing `IServiceHandler` interface — it does not extend it.

**`never-import-l8secure.md`** — The optimizer does NOT import `l8secure`. Google Maps API key is read via the Credentials service (`GetEntities("Creds", 75, ...)`) not via l8secure.

**`single-owner-database-table.md`** — The optimizer does NOT activate its own ORM for VendRoute. It writes via `PostEntity(routes.ServiceName, routes.ServiceArea, ...)` which sends to the Route service's owning process. The optimizer's `OptimizerService` uses the ExecuteService pattern (no ORM, no cache).

**`k8s-three-deployment-modes.md`** — The plan adds OSRM as a new Docker container. The plan currently only mentions a single `k8s/osrm.yaml`. **Must provide four YAML variants** (local, baremetal, gke, kind) per the rule, plus KIND scripts if the project uses them. The project currently has NO k8s directory — this is a pre-existing gap.

**`plan-platform-completeness.md`** — Phase 0 mentions mobile parity for route forms. Phase 0.5 (OSRM) is backend-only. Phase 1+ UI changes must include both desktop and mobile. Each phase's verification includes "verify on both desktop and mobile."

**`log-services-required.md`** — The project does NOT have `log-vnet` or `log-agent` directories. This is a pre-existing gap, not introduced by this plan.

---

## What We Don't Need (Vending-Specific)

- **Pickup & delivery pairs**: Vending is restock-only, no customer returns
- **Multi-technician coordination**: One driver per machine
- **Customer-facing tracking links**: No end customer waiting for delivery
- **Reverse logistics**: Not applicable to vending restock
- **Commercial truck hazmat routing**: Our trucks are standard delivery vehicles

---

## Verification

### Phase 0 verification:
1. Each driver's assigned machines are geographically near their start/end location (no cross-city assignments)
2. The last machine stop is the one nearest to the driver's end location
3. Route totalDistance and estimatedFuelCost include the end-location leg
4. Map shows the end-location leg as a distinct line segment
5. Facility reload near end-of-route prefers driver's home depot when within 2× optimal distance
6. `stopType = "end"` appears in route detail popup

### Phase 0.5 verification:
1. OSRM container starts with `run-local.sh`
2. Route distances are road distances (larger than haversine equivalents)
3. Durations reflect traffic stats (rush hour routes take longer)
4. When OSRM is stopped, optimizer falls back to haversine without errors
5. Google Maps still works as optional refinement layer if API key exists
6. `go build ./...` passes with no new vendored dependencies

### General verification (all phases):
1. Generate routes and verify the new constraints are respected
2. Compare route quality metrics (total distance, duration, fuel cost) before and after
3. Verify on both desktop and mobile
4. Test edge cases: no machines in time window, all drivers over hours, no skills match
5. Tests in `go/tests/` exercise the optimizer via system API (not direct function calls)
