# Create VendFleetMachine Prime Object and Populate Fleet

## Context

The collection pipeline is working — the Nayax simulator's vending machines are collected into VCache and displayed in the "Management Systems" section. Now we need to create individual vending machines as **Prime Objects** in the Fleet section, following the Layer 8 Ecosystem architecture.

Each vending machine from the management system's `machines` map should become a standalone `VendFleetMachine` entity with its own CRUD lifecycle, persisted in the database. This mirrors how probler has NetworkDevice in the inventory cache AND separate business entities derived from it.

## Duplication Audit

- **Go service**: reuses `common.ActivateService()`, `common.NewValidation()` (shared factory/builder). No behavioral duplication.
- **UI enums**: STATUS_MAP/TYPE_MAP are identical to Nayax Cloud's enums. Per probler pattern, each module has its own namespaced enums — no shared JS file (accepted duplication of configuration data).

## Notes on Key Decisions

- **GenerateID safety**: `common.GenerateID()` checks if field is empty before generating — pre-set IDs from VCache (e.g., "M-100001") won't be overwritten. POST with pre-set IDs is safe (same pattern as probler).
- **No Inventory sub-tab yet**: Per probler pattern, tabs are only added when backing code exists. Slot inventory requires per-machine polling (deferred).

---

## Phase 1: Proto — Add VendFleetMachine

**File:** `proto/vend-fleet.proto`

Add after VendMachineList:

```protobuf
// @PrimeObject
message VendFleetMachine {
  string machine_id = 1;
  string name = 2;
  string type = 3;
  string model = 4;
  string status = 5;
  string device_id = 6;
  int32 daily_transactions = 7;
  string last_transaction_at = 8;
  string management_ip = 9;
  string location_address = 10;
  string location_city = 11;
  string location_state = 12;
  double location_lat = 13;
  double location_lng = 14;
  l8common.AuditInfo audit_info = 20;
}

message VendFleetMachineList {
  repeated VendFleetMachine list = 1;
  l8api.L8MetaData metadata = 2;
}
```

Run `cd proto && ./make-bindings.sh` to regenerate.

## Phase 2: Go Service — Update FleetMachine Service

**Directory:** `go/vend/fleet/machines/`

**MachineService.go** — change types to `VendFleetMachine`:
- ServiceName stays `"Machine"`, ServiceArea stays `byte(10)`
- PrimaryKey: `"MachineId"`
- Types: `&vend.VendFleetMachine{}`, `&vend.VendFleetMachineList{}`
- Helper functions return `*vend.VendFleetMachine`

**MachineServiceCallback.go** — change to validate `VendFleetMachine`:
- Auto-generate ID on POST: `common.GenerateID(&entity.MachineId)` guarded by `if action == ifs.POST`
- Require `MachineId`

No new files. Activation in `activate_fleet.go` already calls `machines.Activate()`.

## Phase 3: UI Type Registration

**File:** `go/vend/ui/shared.go`

Add in `RegisterTypes()`:
```go
common.RegisterType(resources, &vend.VendFleetMachine{}, &vend.VendFleetMachineList{}, "MachineId")
```

Keep existing VendMachine registration (still needed for VCache/Management Systems).

## Phase 4: UI — Update Fleet Module Files

**fleet-config.js** — change model to `VendFleetMachine`:
```js
services: [
    { key: 'machines', label: 'Vending Machines', icon: '🏭',
      endpoint: '/10/Machine', model: 'VendFleetMachine', readOnly: true },
    { key: 'machine-groups', label: 'Groups', icon: '📁',
      endpoint: '/10/MachGrp', model: 'VendMachineGroup' },
    { key: 'locations', label: 'Locations', icon: '📍',
      endpoint: '/10/Location', model: 'VendLocation' }
]
```

**fleet-section-config.js** — update title/subtitle.

**fleet-init.js** — standard `Layer8DModuleFactory.create` (already correct).

**machines/machines-enums.js** — update primaryKeys to `VendFleetMachine: 'machineId'`.

**machines/machines-columns.js** — VendFleetMachine columns with ALL proto fields using correct JSON names:
`machineId`, `name`, `type`, `model`, `status`, `deviceId`, `dailyTransactions`, `lastTransactionAt`, `managementIp`, `locationAddress`, `locationCity`, `locationState`, `locationLat`, `locationLng`

**machines/machines-forms.js** — VendFleetMachine form with all fields. Remove machines map `transformData`.

Remove `machines-detail.js` from Fleet (Management Systems only). Update `app.html`.

**Section HTML** — verify `sections/fleet.html` container IDs match `{moduleKey}-{serviceKey}-table-container`.

**sections.js** — already has fleet mapping. No changes needed.

## Phase 5: Populate Fleet from VCache

**File:** `go/vend/inv_vend/main.go`

Add goroutine after `inventory.Activate()`:
- Periodically check VCache for new/updated machine data
- For each machine in the `VendMachine.Machines` map, POST a `VendFleetMachine` to the Machine service via vnic
- Set `ManagementIp` from `VendMachine.MachineId` (the management API IP)
- All 14 fields populated from VendMachineInfo + managementIp

## Phase 6: Verify

1. `go build ./...` — compiles
2. `go vet ./...` — passes
3. Run `run-local.sh`, load mock data
4. Management Systems section shows 1 management entry with 7 machines in popup
5. Fleet section shows 7 individual VendFleetMachine rows with full data
6. Clicking a Fleet row opens standard detail popup with all fields

## Deferred: Per-Machine Slot Inventory

An "Inventory" sub-tab under Fleet will be added when per-machine slot polling is wired. Requires:
1. Add slot fields to `VendFleetMachine` proto (repeated `VendMachineSlot`)
2. Wire per-machine polling: `/lynx/v1/machines/{machineId}/inventory` via `$symbol`
3. Map slot data into `VendFleetMachine` entries
4. Add Inventory sub-tab with slot-focused columns

## Traceability Matrix

| # | Action Item | Phase |
|---|-------------|-------|
| 1 | Add VendFleetMachine + VendFleetMachineList to proto | Phase 1 |
| 2 | Run make-bindings.sh | Phase 1 |
| 3 | Update MachineService.go to use VendFleetMachine | Phase 2 |
| 4 | Update MachineServiceCallback.go with GenerateID + VendFleetMachine | Phase 2 |
| 5 | Add VendFleetMachine type registration in shared.go | Phase 3 |
| 6 | Update fleet-config.js model to VendFleetMachine | Phase 4 |
| 7 | Update fleet-section-config.js title | Phase 4 |
| 8 | Update machines-enums.js primaryKeys | Phase 4 |
| 9 | Update machines-columns.js for all VendFleetMachine fields | Phase 4 |
| 10 | Update machines-forms.js for VendFleetMachine fields | Phase 4 |
| 11 | Remove machines-detail.js from Fleet + app.html | Phase 4 |
| 12 | Verify section HTML container IDs | Phase 4 |
| 13 | Add VCache→Fleet bridge in inv_vend/main.go | Phase 5 |
| 14 | Verify compilation | Phase 6 |
| 15 | Verify Fleet shows 7 machines | Phase 6 |
| 16 | Per-machine slot Inventory sub-tab | Deferred |

## Critical Files

| Action | File |
|--------|------|
| Modify | `proto/vend-fleet.proto` |
| Regenerate | `go/types/vend/vend-fleet.pb.go` |
| Modify | `go/vend/fleet/machines/MachineService.go` |
| Modify | `go/vend/fleet/machines/MachineServiceCallback.go` |
| Modify | `go/vend/ui/shared.go` |
| Modify | `go/vend/ui/web/vend-ui/fleet/fleet-config.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/fleet-section-config.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-enums.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-columns.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-forms.js` |
| Delete | `go/vend/ui/web/vend-ui/fleet/machines/machines-detail.js` |
| Modify | `go/vend/ui/web/app.html` |
| Modify | `go/vend/inv_vend/main.go` |
