# Restructure Inventory to Probler Pattern (Collection → Cache)

## Context

The l8vendingmachine inventory module currently uses CRUD services backed by PostgreSQL (Products, Planograms, RestockOrders via `l8common.ActivateService()`). Per the canonical project selection rule, this project's objective is **observation/collection** — its canonical reference is **probler**, not l8erp.

The inventory must be restructured to follow the probler pattern: observed live state cached in-memory via `l8inventory.Activate()`. Data flows from the Nayax management API's Lynx REST endpoints through collector → parser → inventory cache. Business CRUD services will be built later on top of this cache.

### Architecture: Single Management Target

Unlike probler (one target = one device), this project uses **one target = one Nayax management API** that knows about all vending machines. The collector polls fleet-wide Lynx endpoints that return data for all machines in a single response. The parser splits the response into individual VendMachine inventory entries.

**UI behavior**: The inventory tab shows individual vending machines as **root-level table rows** (not nested under a management parent). Clicking a machine row opens a popup with full detail: slot inventory, alerts, device info, statistics.

## Duplication Audit

The new `inv_vend/main.go` follows probler's `inv_box/main.go` exactly. It is **config-only** (~45 lines): resource creation, VNIC setup, one `inventory.Activate()` call, one metadata function, wait for signal. No behavioral code is duplicated — all behavior lives in the `l8inventory` library.

## Infrastructure Findings (from GPU REST Pattern)

| Concern | Status | Detail |
|---------|--------|--------|
| HTTPS + self-signed certs | **Supported** | `InsecureSkipVerify: true` in RestCollector |
| JSON response parsing | **Supported** | `RestJsonParse` rule with dot-notation field mapping |
| Array → multiple entries | **Supported** | `RestJsonParse` + `CTableToInstances` rule chain splits array into separate inventory entries |
| Vending Machine target type | **Already defined** | `L8PTargetType_Vending_Machine = 8` exists in targets.proto |
| Authentication | **Supported** | `AuthInfo` on `L8PHostProtocol` supports login-based auth |
| Bearer token header | **Gap** | RestCollector does not auto-inject `Authorization: Bearer` after login — deferred infrastructure item |

### 1:N Parsing: Single Response → Multiple Inventory Entries

The `/lynx/v1/machines` endpoint returns an `items` array of 25 machines. The parser uses the `RestJsonParse` + `CTableToInstances` rule chain:

1. `RestJsonParse` extracts the `items` array into a CTable structure
2. `CTableToInstances` iterates each row, creates a separate VendMachine proto instance per row, stores all in `workspace["instances"]`
3. `Parser.ParseMulti()` returns the 25 instances
4. `ParsingCenter.JobComplete()` loops through all 25 and calls `agg.AddElement()` for each — sending each as a separate PATCH to the cache service

This is the same mechanism used for K8s resource parsing in probler (verified in `ParsingCenter.go` lines 86-92).

### Single-Target Polling Strategy

Since this is a cloud management API (one IP serves all machines), the polls use **fleet-wide machine-state endpoints only** — no `$symbol` substitution needed:

- `/lynx/v1/machines` → returns all 25 machines with status, model, location, revenue
- `/lynx/v1/devices` → returns all payment terminals with connection, SIM, firmware
- `/lynx/v1/transactions` → returns recent transactions across all machines
- `/lynx/v1/reports/revenue` → fleet-wide revenue breakdown
- `/lynx/v1/reports/sales-by-period` → 7-day sales data
- `/lynx/v1/reports/machine-performance` → top machines by revenue/uptime

**Not collected** (business layer — not machine state):
- `/lynx/v1/routes` → operational logistics (drivers, schedules) — belongs in CRUD business layer
- `/lynx/v1/tasks` → work orders (refill/repair) — belongs in CRUD business layer
- `/lynx/v1/actors` → organizational master data (operators) — belongs in CRUD business layer

Per-machine detail endpoints (`/lynx/v1/machines/{id}/inventory`, `/attributes`, `/statistics`, `/alerts`) are **not polled directly** in this phase. The fleet-wide `/lynx/v1/machines` response contains enough data (status, model, location, revenue, device IDs) to populate the root table. Per-machine detail can be fetched on-demand when the user clicks a row, or added as a second polling phase later.

---

## Phase 0: Update VendMachine Proto

The VendMachine proto is missing fields returned by the Nayax Lynx `/machines` API. Add inline fields to hold the collected data (matching the GPU proto pattern of inline device info).

### Fields to add to `VendMachine` in `proto/vend-fleet.proto`

| Field | Type | Proto Name | Reason |
|-------|------|-----------|--------|
| `name` | string | `name` | Machine display name (e.g., "Office Lobby Snack") |
| `daily_transactions` | int32 | `daily_transactions` | Today's transaction count |
| `revenue_today` | int64 | `revenue_today` | Today's revenue in cents |
| `revenue_currency` | string | `revenue_currency` | Currency code (e.g., "USD") |
| `last_transaction_at` | string | `last_transaction_at` | ISO timestamp of last transaction |
| `device_id` | string | `device_id` | Associated payment terminal ID |
| `uptime_percent` | double | `uptime_percent` | Uptime percentage |
| `location_address` | string | `location_address` | Street address (inline, not reference) |
| `location_city` | string | `location_city` | City |
| `location_state` | string | `location_state` | State/province |
| `location_lat` | double | `location_lat` | Latitude |
| `location_lng` | double | `location_lng` | Longitude |

Keep existing `location_id` field for future business layer use. The inline location fields hold the observed/collected location from the management API (same pattern as `GpuDeviceInfo` having `latitude`/`longitude` directly).

After modifying the proto: run `cd proto && ./make-bindings.sh` to regenerate `.pb.go` files (per `protobuf-generation.md`). Verify `make-bindings.sh` uses `-i` not `-it` on docker run commands.

---

## Phase 1: Create `inv_vend` Binary

Create `go/vend/inv_vend/` with 3 files. **Copy from probler's `inv_box` and the existing collector, then adapt** (per `l8pollaris-binary-deployment.md`).

### `go/vend/inv_vend/main.go` (new, ~45 lines)

Copy from `../probler/go/prob/inv_box/main.go`, adapt:
- `vendcommon.CreateResources("inv-vend")`
- `inventory.Activate(vendcommon.VendMachine_Links_ID, &vend.VendMachine{}, &vend.VendMachineList{}, nic, "MachineId")`
- `invCenter.AddMetadata("Online", Online)` — checks `machine.Status == VendMachineStatus_VEND_MACHINE_STATUS_OPERATIONAL`
- `vendcommon.WaitForSignal(res)`

### `go/vend/inv_vend/build.sh` (new)

Copy from `go/vend/collector/build.sh`, change image name to `saichler/vendmachine-inv:latest`.

### `go/vend/inv_vend/Dockerfile` (new)

Copy from `go/vend/collector/Dockerfile`, change binary name to `vendmachine-inv`.

### Dependencies

Running `go mod tidy` will pull in `github.com/saichler/l8inventory`.

### Existing infrastructure (no changes needed)

- `go/vend/common/Links.go` — already implements `TargetLinks` with `VCache`/`byte(0)`
- `go/vend/common/defaults.go` — already has `VendMachine_Links_ID = "Vend"`, `Vend_Cache_Service_Name = "VCache"`

---

## Phase 2: Update Polling Configs to Nayax Lynx Paths

All files in `go/vend/parser/boot/`. Change endpoints to fleet-wide Lynx paths (no `$symbol` — single management target). Reduce from 14 polls to 8 fleet-wide polls. Remove per-machine polls that have no fleet-wide equivalent.

### Polls to update

| File | Poll Name | Current Path | New Path | Notes |
|------|-----------|-------------|----------|-------|
| `vend_machine_polls.go` | vendMachineIdentity | `/api/v1/machine` | `/lynx/v1/machines` | Fleet-wide, returns all machines with model/status/location |
| `vend_machine_polls.go` | vendMachineStatus | `/api/v1/machine/status` | **Remove** | Merged into vendMachineIdentity (machines endpoint includes status) |
| `vend_inventory_polls.go` | vendInventory | `/api/v1/inventory` | **Remove (deferred)** | Per-machine only (`/lynx/v1/machines/{id}/inventory`); add later |
| `vend_inventory_polls.go` | vendInventoryAlerts | `/api/v1/inventory/alerts` | **Remove (deferred)** | Per-machine only; add later |
| `vend_sales_polls.go` | vendTransactions | `/api/v1/transactions` | `/lynx/v1/transactions` | Fleet-wide |
| `vend_sales_polls.go` | vendTransactionSummary | `/api/v1/transactions/summary` | `/lynx/v1/reports/revenue` | Fleet-wide |
| `vend_monitor_polls.go` | vendTemperature | `/api/v1/temperature` | **Remove (deferred)** | Per-machine attribute; add later |
| `vend_monitor_polls.go` | vendAlerts | `/api/v1/alerts` | **Remove (deferred)** | Per-machine only; add later |
| `vend_monitor_polls.go` | vendEnergy | `/api/v1/energy` | **Remove (deferred)** | Per-machine attribute; add later |
| `vend_payment_polls.go` | vendCashbox | `/api/v1/payment/cashbox` | **Remove (deferred)** | Not in fleet-wide endpoints; add later |
| `vend_payment_polls.go` | vendPaymentStatus | `/api/v1/payment/status` | `/lynx/v1/devices` | Fleet-wide device/payment terminal status |
| `vend_analytics_polls.go` | vendTraffic | `/api/v1/analytics/traffic` | `/lynx/v1/reports/sales-by-period` | Fleet-wide |
| `vend_analytics_polls.go` | vendHealth | `/api/v1/analytics/health` | `/lynx/v1/reports/machine-performance` | Fleet-wide |
| `vend_dex_polls.go` | vendDexAudit | `/api/v1/dex/audit` | **Remove (deferred)** | Per-machine only; add later |

### Resulting 6 active polls

1. **vendMachines** — `GET::/lynx/v1/machines::` — EVERY_5_MINUTES, always=true
2. **vendDevices** — `GET::/lynx/v1/devices::` — EVERY_5_MINUTES
3. **vendTransactions** — `GET::/lynx/v1/transactions::` — EVERY_30_SECONDS, always=true
4. **vendRevenue** — `GET::/lynx/v1/reports/revenue::` — EVERY_5_MINUTES
5. **vendSalesByPeriod** — `GET::/lynx/v1/reports/sales-by-period::` — EVERY_1_HOUR
6. **vendMachinePerformance** — `GET::/lynx/v1/reports/machine-performance::` — EVERY_1_HOUR

### Field mapping updates

Update poll rule chains to match verified Nayax Lynx JSON keys:

**vendMachines poll** (`/lynx/v1/machines`) — returns `{"items": [...], "totalCount": 25}`:
- Rule 1: `RestJsonParse` — extracts `items` array into CTable
- Rule 2: `CTableToInstances` — splits into 25 VendMachine instances keyed by `machineId`
- Field mappings: `machineId`, `name`, `model`, `status`, `type`, `dailyTransactions`, `revenue.today` → `revenueToday`, `revenue.currency` → `revenueCurrency`, `location.address` → `locationAddress`, `location.city` → `locationCity`, `location.state` → `locationState`, `location.lat` → `locationLat`, `location.lng` → `locationLng`, `deviceId`, `lastTransactionAt`

**vendDevices poll** (`/lynx/v1/devices`) — returns `{"items": [...]}`:
- Maps payment terminal data. Each device has `machineId` which links it to the corresponding VendMachine entry.

**vendTransactions, vendRevenue, vendSalesByPeriod, vendMachinePerformance** — fleet-wide aggregate data. These update summary fields on VendMachine entries matched by `machineId`.

### Files to simplify

Remove the per-machine poll files that are now empty:
- `vend_inventory_polls.go` — remove both functions, keep file with comment noting deferred
- `vend_monitor_polls.go` — remove temperature/alerts/energy functions
- `vend_payment_polls.go` — remove cashbox function, keep payment status
- `vend_dex_polls.go` — remove entirely

Update `vend_pollaris.go` to remove deferred poll calls and update the `Order` slice.

---

## Phase 3: Remove Old CRUD Inventory Services

### Delete entire directory
- `go/vend/inventory/` (products/, planograms/, restockorders/ — 6 files)

### Delete activation file
- `go/vend/services/activate_inventory.go`

### Modify `go/vend/services/activate_all.go`
Remove line 17: `all = append(all, collectInventoryActivations(creds, dbname, nic)...)`

### Keep proto types
Do NOT delete `proto/vend-inventory.proto` — these types will be used later for the business CRUD layer.

---

## Phase 4: Update Build & Deploy Infrastructure

### `go/build-all-images.sh`
Add step 6/6:
```bash
echo "6/6 Building inventory cache..."
cd vend/inv_vend && ./build.sh && cd ../..
```

### New: `go/k8s/vend-inv.yaml`

**Copy from `go/k8s/vend-collector.yaml`** and adapt. Required K8s entries (per `k8s-yaml-required-entries.md`):

- Namespace metadata with **labels** (`name: vend-inv`)
- DaemonSet metadata with **labels** (`app: vend-inv`)
- Container **`env` section with NODE_IP** from `status.hostIP`
- Volume name **`hdata`** (not `data`), mountPath `/data`, hostPath type `DirectoryOrCreate`
- Image: `saichler/vendmachine-inv:latest`, `imagePullPolicy: Always`

### `go/k8s/deploy.sh`
Add after parser: `kubectl apply -f ./vend-inv.yaml`

### `go/k8s/undeploy.sh`
Add before parser: `kubectl delete -f ./vend-inv.yaml`

### `go/run-local.sh`

Build line (after parser build):
```bash
cd vend/inv_vend && go build -o ../../demo/inv_demo && cd ../../
```

Start line (after vend_demo, before collector/parser):
```bash
echo "Starting inventory cache..."
./inv_demo &
sleep 2
```

Per `demo-directory-sync.md` — only edit `run-local.sh`, never files in `demo/`.

---

## Traceability Matrix

| # | Gap / Action Item | Phase |
|---|-------------------|-------|
| 1 | Add missing fields to VendMachine proto (name, location, revenue, etc.) | Phase 0 |
| 2 | Run make-bindings.sh to regenerate .pb.go files | Phase 0 |
| 3 | Create inv_vend binary with `l8inventory.Activate()` | Phase 1 |
| 4 | Create inv_vend build.sh (copy from collector) | Phase 1 |
| 5 | Create inv_vend Dockerfile (copy from collector) | Phase 1 |
| 6 | Pull l8inventory dependency via go mod tidy | Phase 1 |
| 7 | Update vendMachineIdentity to `/lynx/v1/machines` with CTableToInstances rule | Phase 2 |
| 8 | Remove vendMachineStatus (merged into machines endpoint) | Phase 2 |
| 9 | Defer vendInventory (per-machine only) | Phase 2 |
| 10 | Defer vendInventoryAlerts (per-machine only) | Phase 2 |
| 11 | Update vendTransactions to `/lynx/v1/transactions` | Phase 2 |
| 12 | Update vendTransactionSummary to `/lynx/v1/reports/revenue` | Phase 2 |
| 13 | Defer vendTemperature (per-machine only) | Phase 2 |
| 14 | Defer vendAlerts (per-machine only) | Phase 2 |
| 15 | Defer vendEnergy (per-machine only) | Phase 2 |
| 16 | Defer vendCashbox (per-machine only) | Phase 2 |
| 17 | Update vendPaymentStatus to `/lynx/v1/devices` | Phase 2 |
| 18 | Update vendTraffic to `/lynx/v1/reports/sales-by-period` | Phase 2 |
| 19 | Update vendHealth to `/lynx/v1/reports/machine-performance` | Phase 2 |
| 20 | Defer vendDexAudit (per-machine only) | Phase 2 |
| 21 | Update field mappings to Nayax Lynx JSON keys | Phase 2 |
| 22 | Clean up empty/deferred poll files | Phase 2 |
| 23 | Update vend_pollaris.go Order slice | Phase 2 |
| 24 | Delete go/vend/inventory/ directory (6 files) | Phase 3 |
| 25 | Delete go/vend/services/activate_inventory.go | Phase 3 |
| 26 | Remove collectInventoryActivations from activate_all.go | Phase 3 |
| 27 | Add inv_vend to build-all-images.sh | Phase 4 |
| 28 | Create vend-inv.yaml with all required K8s entries | Phase 4 |
| 29 | Add vend-inv.yaml to deploy.sh | Phase 4 |
| 30 | Add vend-inv.yaml to undeploy.sh | Phase 4 |
| 31 | Add inv_vend build and start to run-local.sh | Phase 4 |
| 32 | Verify compilation with `go build ./...` | Phase 5 |
| 33 | Verify with `go vet ./...` | Phase 5 |
| 34 | End-to-end test with Nayax simulator | Phase 5 |
| 35 | Bearer token auth — infrastructure item | Deferred |
| 36 | Per-machine detail endpoints (inventory, attributes, alerts, statistics) | Deferred |
| 37 | Routes/tasks/actors — business layer, not inventory | Deferred — CRUD layer |

---

## Phase 5: End-to-End Verification

### Compilation
1. `go build ./...` — zero errors
2. `go vet ./...` — no issues

### Local run
1. Run `run-local.sh`
2. Verify all services start (vnet, vend_demo, inv_demo, collector, parser, ui)
3. Verify inv_demo logs show `inventory.Activate` registration for VendMachine

### Data flow (with Nayax simulator at 192.168.200.1:8443)
1. Create a target: IP `192.168.200.1`, port `8443`, REST protocol, `L8PTargetType_Vending_Machine`
2. Configure AuthInfo for Lynx SignIn
3. Verify collector polls `/lynx/v1/machines` and other fleet-wide endpoints
4. Verify parser splits machines array into individual VendMachine cache entries
5. Verify inv_vend caches the data in dcache
6. Verify querying VCache returns individual machine entries (not one management blob)

### UI verification
- [ ] Inventory tab root table shows individual vending machines as rows
- [ ] Each row shows: machine name, model, status, location, revenue
- [ ] Clicking a row opens popup with full machine detail
- [ ] Old CRUD endpoints (Product, Planogram, RstockOrd) no longer exist

### Sections to verify
- [ ] inv_vend binary starts and registers VCache service
- [ ] Collector connects and polls fleet-wide Lynx endpoints
- [ ] Parser splits array responses into per-machine entries
- [ ] Inventory cache serves queries with individual machine data
- [ ] K8s deployment succeeds with new vend-inv.yaml

### Deferred items (not blocking)
- Bearer token injection (may need RestCollector enhancement)
- Per-machine detail endpoints (inventory slots, attributes, statistics, alerts)
- UI popup detail data (depends on per-machine polling or on-demand fetch)

---

## Critical Files Summary

| Action | File |
|--------|------|
| **Modify** | `proto/vend-fleet.proto` (add missing Lynx fields to VendMachine) |
| **Regenerate** | `go/types/vend/vend-fleet.pb.go` (via `cd proto && ./make-bindings.sh`) |
| **Create** | `go/vend/inv_vend/main.go` (copy from `../probler/go/prob/inv_box/main.go`) |
| **Create** | `go/vend/inv_vend/build.sh` (copy from `go/vend/collector/build.sh`) |
| **Create** | `go/vend/inv_vend/Dockerfile` (copy from `go/vend/collector/Dockerfile`) |
| **Create** | `go/k8s/vend-inv.yaml` (copy from `go/k8s/vend-collector.yaml`) |
| **Modify** | `go/vend/parser/boot/vend_pollaris.go` |
| **Modify** | `go/vend/parser/boot/vend_machine_polls.go` |
| **Modify** | `go/vend/parser/boot/vend_sales_polls.go` |
| **Modify** | `go/vend/parser/boot/vend_payment_polls.go` |
| **Modify** | `go/vend/parser/boot/vend_analytics_polls.go` |
| **Simplify/Remove** | `go/vend/parser/boot/vend_inventory_polls.go` |
| **Simplify/Remove** | `go/vend/parser/boot/vend_monitor_polls.go` |
| **Remove** | `go/vend/parser/boot/vend_dex_polls.go` |
| **Modify** | `go/vend/services/activate_all.go` |
| **Modify** | `go/build-all-images.sh` |
| **Modify** | `go/k8s/deploy.sh` |
| **Modify** | `go/k8s/undeploy.sh` |
| **Modify** | `go/run-local.sh` |
| **Delete** | `go/vend/inventory/` (entire directory) |
| **Delete** | `go/vend/services/activate_inventory.go` |
