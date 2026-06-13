# Plan: On-Demand Route Optimizer

## Overview

A backend Go service that, when triggered, reads machine inventory levels, classifies machines into two priority lists, clusters them by proximity, matches them to available trucks/drivers/facilities, and generates optimized VendRoute entities.

---

## Core Algorithm

### Step 1: Build Restock Demand Lists

Read all VendFleetMachine records (from the Machine service, area 10). For each machine, examine its `inventory` slots:

**List A — Needs Restock Now:**
- Any slot has `currentStock == 0` (empty)
- Machine `emptySlots > 0`
- Machine profile `cascadeThresholdPct` is reached (overall fill % below threshold)

**List B — Can Wait 1 Day:**
- Slots with `currentStock > 0` but below 30% of `capacity` (low)
- Machine `lowStockSlots > 0` but no empty slots
- Not yet at cascade threshold but trending toward it

Each entry records: machineId, GPS (locationLat/locationLng), and a product demand map (productName → quantity needed to fill to capacity).

### Step 2: Cluster Machines by Proximity

Using haversine distance between GPS coordinates:

1. Take all List A machines
2. Pick the machine closest to any facility as a seed
3. Grow the cluster by adding the nearest unassigned List A machine, as long as:
   - Total product demand of the cluster doesn't exceed the truck's cargo capacity
   - Adding it doesn't increase total route distance beyond a max (configurable, e.g., 50 miles)
4. Once a cluster can't grow, start a new cluster with the next closest unassigned machine
5. Repeat until all List A machines are assigned to a cluster

### Step 3: Add List B Machines to Existing Routes

For each List B machine:
1. Calculate its distance to every leg of every cluster's route (insertion cost)
2. If inserting it between two existing stops adds less than a threshold (e.g., 3 miles detour), add it to that cluster
3. Mark the stop with `serviceUrgency = "low"` vs `"high"` for List A stops

### Step 4: Match Cluster → Truck → Driver

For each cluster:

1. **Truck**: Pick a truck that:
   - Has status Active
   - Has `cargoCapacityCuFt` sufficient for the estimated restock volume
   - `licenseClass` requirement can be met by an available driver

2. **Driver**: From drivers with `truckId` matching the selected truck:
   - Has `isActive == true`
   - Has a schedule entry for the planned day (`schedule[day].day` matches)
   - `licenseClass` matches truck requirements

The truck's current `stock` determines whether a facility stop is needed at all (see Step 5).

### Step 5: Build Route with Dynamic Facility Stops

The route is NOT "start at facility, visit machines, return." Instead, facility stops are inserted only when the truck runs out of product, at whatever point in the route that happens.

**Algorithm:**

1. Start from the driver's start location for the day (`schedule[day].startLocationId`, or home address if blank)
2. Order machines using nearest-neighbor from the start point
3. Walk through the ordered stops, tracking truck stock as it depletes:
   - Before each stop, check: does the truck have enough stock for this machine's demand?
   - **If yes**: visit the machine, deduct products from truck stock, continue
   - **If no**: insert a facility reload stop BEFORE this machine
     - Pick the nearest Active facility (by haversine from current position) that:
       - Has sufficient stock for at least the next N machines' demand
       - Is open on the planned date (`operatingDays`)
     - Add a facility stop (type: "reload") to the route
     - Replenish truck stock from facility stock
     - Add reload time to duration estimate
     - Continue to the machine
4. After all machines are visited, the route ends at the last machine (no mandatory return to facility)

**If the truck already has sufficient stock for ALL machines in the cluster, no facility stop is added at all.**

5. Apply 2-opt improvement on the machine stops (keeping facility reload stops anchored where stock runs out — they shift only if the reordering moves the depletion point)

**Calculate (initial — haversine-based):**
- `totalDistance`: sum of haversine distances between all stops (including facility detours)
- `totalDuration`: `totalDistance / avgSpeed + serviceTimePerStop * machineStops + reloadTime * facilityStops`
- `plannedArrival` per stop: cumulative time from driver's start time
- `estimatedFuelCost`: `totalDistance / truck.milesPerGallon * fuelPricePerGallon`

### Step 6: Refine with Traffic Data (Google Maps)

After stops are ordered and facility reloads are inserted, call the Google Maps Directions API to get real driving distances and traffic-aware durations.

**How it works:**
1. Build a waypoints list from the ordered stops (driver start → stop 1 → stop 2 → ... → last stop)
2. Call the Directions API with:
   - `origin`: driver's start location
   - `destination`: last stop
   - `waypoints`: all intermediate stops in order
   - `departure_time`: the route's planned start time (from the request or driver's schedule start time)
   - `traffic_model`: `best_guess`
3. The API returns per-leg data:
   - `distance.value` — actual road distance in meters (not straight-line)
   - `duration_in_traffic.value` — traffic-aware travel time in seconds
4. Replace haversine estimates with real values:
   - `totalDistance`: sum of all leg road distances
   - Per-stop `plannedArrival`: cumulative traffic-aware duration + service time at each preceding stop
   - `totalDuration`: last arrival + last service time - start time
   - `estimatedFuelCost`: real road distance / MPG × fuel price (more accurate than haversine)

**If traffic makes the route too long:**
- If total duration exceeds the driver's working hours (end time - start time from schedule), split the cluster: move the last N stops to a new route assigned to a different driver/truck

**API key management:**
- Stored in L8 credentials via `ISecurityProvider` (same pattern as DB credentials)
- Credential key: `"google-maps"`, looked up via `nic.Resources().Security().Credential()`

**Fallback:**
- If no API key is configured or the API call fails, keep the haversine-based estimates
- Log a warning but don't fail route generation

**Cost optimization:**
- One Directions API call per generated route (not per machine pair)
- Clustering and ordering use free haversine — Google API only called once per final route
- For a typical day generating 5-10 routes, this is ~5-10 API calls

### Step 7: Generate VendRoute

Create a VendRoute for each cluster:
```
routeId:       auto-generated
name:          "Route <date>-<seq>"
status:        PLANNED
driverId:      selected driver
vehicleId:     selected truck
facilityId:    primary facility (first reload, or nearest if no reload needed)
plannedDate:   specified date
totalDistance:  calculated (includes facility detours)
totalDuration: calculated (includes reload time)
estimatedFuelCost: distance / MPG * fuel price
stops:         ordered VendRouteStop list (machines + facility reloads interleaved)
```

VendRouteStop entries:
- Machine stops: `serviceUrgency = "high"` (List A) or `"low"` (List B)
- Facility reload stops: `serviceUrgency = "reload"`, `machineId = facilityId`

POST each route to the Route service.

---

## Implementation

### Phase 1: Route Optimizer Package

Create `go/vend/route/optimizer/` with these files:

**demand.go** — Step 1
- `BuildDemandLists(nic) (listA, listB []MachineDemand, err)`
- Reads VendFleetMachine from Machine service
- Reads VendMachineProfile for cascade thresholds
- Classifies into List A (urgent) and List B (can wait)
- `MachineDemand` struct: machineId, lat, lng, products map[string]int32, urgency

**distance.go** — Distance + route metrics
- `Haversine(lat1, lng1, lat2, lng2 float64) float64` — returns miles
- `Centroid(points [][2]float64) (lat, lng float64)`
- `ComputeRouteMetrics(legs []RouteLeg, startTime int64, mpg float64, fuelPrice float64, serviceMinutes int32, reloadMinutes int32) RouteMetrics`
  - `RouteLeg` struct: distanceMiles float64, durationSeconds int64, isReload bool
  - `RouteMetrics` struct: totalDistance, totalDuration, estimatedFuelCost, plannedArrivals []int64
  - Single function used by BOTH haversine path and Google Maps path — the only difference is the source of leg distances/durations

**cluster.go** — Steps 2-3
- `ClusterMachines(listA, listB []MachineDemand, maxDistanceMiles, maxDetourMiles float64) []Cluster`
- `Cluster` struct: machines []MachineDemand, centroidLat, centroidLng, totalDemand map[string]int32

**assign.go** — Step 4
- `AssignResources(cluster *Cluster, trucks, drivers) (*Assignment, error)`
- `Assignment` struct: truckId, driverId, startLat, startLng (from driver's schedule)

**route_builder.go** — Step 5 (core of the new logic)
- `BuildRoute(cluster *Cluster, assignment *Assignment, truckStock map[string]int32, facilities []*vend.VendStockingFacility, config *RouteConfig) *BuiltRoute`
- Walks through ordered stops tracking truck stock depletion
- Inserts facility reload stops dynamically when stock runs out
- Picks nearest active/open facility at the point of depletion
- Applies 2-opt on machine stops, re-evaluates facility insertion points after reordering
- Builds haversine-based `RouteLeg` list, calls `ComputeRouteMetrics()` for totals
- `BuiltRoute` struct: stops []RouteStop, legs []RouteLeg, metrics RouteMetrics, facilityStops []string
- **Note:** If this file exceeds 500 lines, split 2-opt into `optimize.go`

**traffic.go** — Step 6 (Google Maps integration, raw net/http — no external SDK, same as l8collector REST pattern)
- `RefineWithTraffic(route *BuiltRoute, startTime int64, mpg float64, config *RouteConfig, nic ifs.IVNic) error`
- Looks up Google Maps API key from L8 credentials (`"google-maps"`)
- Builds Directions API request with ordered waypoints + `departure_time` via raw `net/http` GET
- Converts API response legs into `[]RouteLeg` (road distance + traffic-aware duration)
- Calls the SAME `ComputeRouteMetrics()` to recalculate totals — no duplicate calculation logic
- Replaces `route.legs` and `route.metrics` with traffic-refined values
- If total duration exceeds driver's working hours, returns a split indicator
- Falls back silently (keeps existing haversine metrics) if no API key or API failure

**generator.go** — Step 7
- `GenerateRoutes(nic ifs.IVNic, req *vend.VendRouteOptRequest) ([]*vend.VendRoute, error)`
- Orchestrates all steps, calls `RefineWithTraffic` after route building
- If traffic refinement triggers a split, re-clusters the overflow stops
- POSTs routes to Route service
- Returns the generated routes

### Phase 2: Proto + Service Endpoint

**Proto changes:**
- Add `VendRouteOptRequest` to `vend-route.proto` — command message, NO List type needed (same as `CJob` in l8collector)
- Add result fields to `VendRouteOptRequest` itself — L8 pattern is to return the mutated request with results populated (no separate Response type)
- Add `facility_id` (field 15) and `estimated_fuel_cost` (field 16) to `VendRoute`
- Add `stop_type` (field 8) and `facility_id` (field 9) to `VendRouteStop`
- Run `cd proto && ./make-bindings.sh` to regenerate bindings
- Verify: `go build ./...`

**Service — follows l8collector ExecuteService pattern (NOT standard ActivateService):**

Create `go/vend/route/optimizer/OptimizerService.go`:
- Does NOT use `common.ActivateService` — this is a command endpoint, not CRUD
- Implements custom `Post()` method that calls `GenerateRoutes()` and returns the mutated request with result fields
- Auto-generates `RouteId` via `common.GenerateID` on each generated route before POST to Route service
- Registers via `WebService()` with `AddEndpoint(&VendRouteOptRequest{}, ifs.POST, &VendRouteOptRequest{})`
- ServiceName: `"OptRoute"` (8 chars), ServiceArea: byte(10)

```go
// Pattern from l8collector/ExecuteService:
func (this *OptimizerService) WebService() ifs.IWebService {
    ws := web.New("OptRoute", this.serviceArea, 0)
    ws.AddEndpoint(&vend.VendRouteOptRequest{}, ifs.POST, &vend.VendRouteOptRequest{})
    return ws
}

func (this *OptimizerService) Post(pb ifs.IElements, vnic ifs.IVNic) ifs.IElements {
    req := pb.Element().(*vend.VendRouteOptRequest)
    routes, err := GenerateRoutes(vnic, req)
    // Populate result fields on req
    req.GeneratedRouteIds = routeIds
    req.GeneratedCount = int32(len(routes))
    return object.New(nil, req)
}
```

**Synchronous execution** — the optimizer runs within the POST request (same as l8agent chat tool loops). For 5-10 routes with one Google Maps API call each, this takes a few seconds.

**Concurrency note:** If two users trigger route generation simultaneously, both read the same truck/facility stock snapshots. L8 ecosystem does not use capacity reservation (confirmed: l8collector, l8alarms use mutex for thread safety but no slot allocation). The optimizer should use a sync.Mutex to serialize concurrent requests. Document that concurrent generation could produce overlapping facility reload plans.

**Type registration:**
- Register `VendRouteOptRequest` in `go/vend/ui/shared.go` (`RegisterTypes` function)

**Activation:**
- Register the optimizer service with the vnic in `go/vend/services/activate_route.go` (same pattern as l8collector's ExecuteService registration)

### Phase 3: UI Button (Desktop + Mobile)

Add a "Generate Routes" button to the Routes section on both platforms:

**Desktop:**
- Toolbar button above the routes table
- Calls `POST /vend/10/OptRoute` with the planned date and start time
- On success, refreshes the routes table to show new routes
- Shows notification with count: "Generated N routes"

**Mobile:**
- Same button in the routes card view header
- Calls same endpoint via `Layer8MAuth.post()`
- Refreshes table and shows notification

### Phase 4: Map Integration

When a route is selected/viewed:
- Draw the route path on the map as a polyline connecting all stops in order
- Color-code stops: red = urgent machine, yellow = low-priority machine, blue = facility reload
- Facility reload stops show as a different marker shape (square vs circle)
- Popup on each stop shows: machine name, products to restock, or "Reload at <facility name>"

---

## Proto Changes

Add a command message for the optimizer (NO List type — this is a command, not a CRUD entity, same pattern as CJob in l8collector):

```protobuf
message VendRouteOptRequest {
    // Input
    int64 planned_date = 1;              // Unix timestamp, default = tomorrow
    int64 start_time = 2;               // Unix timestamp for departure (for traffic query)
    double max_route_distance = 3;       // Max miles per route, default = 50
    double max_detour_distance = 4;      // Max detour for List B machines, default = 3
    int32 reload_time_minutes = 5;       // Time to reload at facility, default = 30

    // Output (populated by optimizer, returned in POST response)
    repeated string generated_route_ids = 10;
    int32 generated_count = 11;
    int32 list_a_count = 12;             // Machines that needed immediate restock
    int32 list_b_added = 13;             // "Can wait" machines added opportunistically
    string error = 14;                   // Error message if generation failed
}
```

Add fields to VendRoute:

```protobuf
// Add to existing VendRoute message:
string facility_id = 15;                 // cross-ref: VendStockingFacility (primary facility)
double estimated_fuel_cost = 16;         // distance / MPG * fuel price
```

Add `stop_type` to VendRouteStop to distinguish machine visits from facility reloads:

```protobuf
// Add to existing VendRouteStop message:
string stop_type = 8;                    // "machine" or "reload"
string facility_id = 9;                  // Only set for reload stops
```

---

## Configuration

Tuning constants (could be in a config or hardcoded initially):

| Parameter | Default | Description |
|-----------|---------|-------------|
| Low stock threshold | 30% | Slot fill % below which machine goes to List B |
| Max route distance | 50 miles | Won't add machines beyond this total distance |
| Max detour for List B | 3 miles | Max extra distance to pick up a "can wait" machine |
| Average speed | 25 mph | Haversine fallback when no Google Maps API key |
| Service time per stop | 20 min | Time to restock a machine |
| Facility reload time | 30 min | Time to reload truck at facility |
| Fuel price per gallon | $3.50 | For cost estimation |

---

## Traceability Matrix

| # | Item | Phase |
|---|------|-------|
| 1 | MachineDemand struct + BuildDemandLists (List A / List B) | Phase 1 (demand.go) |
| 2 | Haversine distance function | Phase 1 (distance.go) |
| 3 | Cluster machines by proximity + capacity | Phase 1 (cluster.go) |
| 4 | Insert List B machines into existing routes (detour check) | Phase 1 (cluster.go) |
| 5 | Match cluster to truck + driver (by schedule + license) | Phase 1 (assign.go) |
| 6 | Nearest-neighbor ordering from driver start location | Phase 1 (route_builder.go) |
| 7 | Dynamic facility reload insertion when truck stock depletes | Phase 1 (route_builder.go) |
| 8 | 2-opt improvement with facility re-insertion | Phase 1 (route_builder.go) |
| 9 | Fuel cost calculation using MPG | Phase 1 (route_builder.go) |
| 10 | Google Maps Directions API integration with departure_time | Phase 1 (traffic.go) |
| 11 | Traffic-aware duration/distance replaces haversine estimates | Phase 1 (traffic.go) |
| 12 | Route split when traffic makes duration exceed driver hours | Phase 1 (traffic.go + generator.go) |
| 13 | Fallback to haversine when no API key or API failure | Phase 1 (traffic.go) |
| 14 | Generate and POST VendRoute entities | Phase 1 (generator.go) |
| 15 | VendRouteOptRequest proto (with start_time) | Phase 2 |
| 16 | Regenerate bindings (make-bindings.sh) | Phase 2 |
| 17 | OptimizerService — custom Post() (ExecuteService pattern, NOT ActivateService) | Phase 2 |
| 18 | WebService().AddEndpoint() registration (no List type) | Phase 2 |
| 19 | Auto-generate RouteId on each generated route (GenerateID) | Phase 2 |
| 20 | VendRouteOptRequest result fields (returned in POST response) | Phase 2 |
| 21 | Proto: facility_id + estimatedFuelCost on VendRoute | Phase 2 |
| 22 | Proto: stopType + facilityId on VendRouteStop | Phase 2 |
| 23 | Register VendRouteOptRequest in ui/shared.go | Phase 2 |
| 24 | Register optimizer with vnic in activate_route.go | Phase 2 |
| 25 | sync.Mutex to serialize concurrent route generation | Phase 2 |
| 26 | Google Maps API key in L8 credentials | Phase 2 |
| 27 | Raw net/http for Google Maps API (no external SDK) | Phase 1 (traffic.go) |
| 28 | "Generate Routes" button — desktop UI | Phase 3 |
| 29 | "Generate Routes" button — mobile UI (parity) | Phase 3 |
| 30 | Route path visualization on map (polyline + color-coded stops) | Phase 4 |

## Verification

1. Trigger route generation via POST
2. Verify routes are created with correct driver/truck assignments
3. Verify List A machines all appear in routes
4. Verify List B machines are only added when the detour is small
5. Verify stop order is reasonable (no backtracking)
6. Verify facility reload stops appear only when truck stock is insufficient
7. Verify no facility stop when truck already has enough stock for all machines
8. Verify facility reload picks the nearest active/open facility at the depletion point
9. Verify fuel cost includes distance to/from facility reload stops
10. With Google Maps API key configured: verify distances are road distances (not haversine), durations reflect traffic, plannedArrival times account for traffic
11. With start_time set to rush hour: verify longer durations than off-peak
12. Without Google Maps API key: verify fallback to haversine estimates, no errors
13. Verify route split when traffic-aware duration exceeds driver working hours
14. Verify routes appear in the Routes table and on the map with correct stop types
