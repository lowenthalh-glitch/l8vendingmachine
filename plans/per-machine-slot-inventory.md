# Per-Machine Targets for Slot Inventory Collection

## Context

The Fleet section shows 7 vending machines collected from the Nayax management API (fleet-wide `/lynx/v1/machines` endpoint). But slot inventory data (products, stock levels, capacity per slot) requires per-machine endpoints (`/lynx/v1/machines/{machineId}/inventory`). The simulator has this data — 8 slots per machine with product names, SKUs, prices, stock levels, and status.

The approach: when the VCache→Fleet bridge creates VendFleetMachine entries, it also creates a per-machine `L8PTarget` for each. These targets use `$symbol` substitution so the collector polls `/lynx/v1/machines/$symbol/inventory` where `$symbol` is replaced with the machine ID (e.g., "M-100001"). The parser maps the slot data into each VendFleetMachine.

This fits the probler pattern — one target per device, each with its own polling cadence.

## Phase 1: Add Slot Fields to VendFleetMachine Proto

**File:** `proto/vend-fleet.proto`

Add to VendFleetMachine (reuse existing `VendMachineSlot` child type):
```protobuf
message VendFleetMachine {
  // ... existing fields 1-14, 20 ...
  repeated VendMachineSlot inventory = 30;
  int32 total_slots = 31;
  int32 empty_slots = 32;
  int32 low_stock_slots = 33;
  string inventory_last_updated = 34;
}
```

Run `cd proto && ./make-bindings.sh`.

## Phase 2: Create Per-Machine Pollaris Config and Links Routing

### Pollaris Config

**New file:** `go/vend/parser/boot/vend_per_machine_polls.go`

Create a second pollaris config. **Name must be `"VMach"`** (matching LinksID for boot stage filtering — REST-only targets use boot stage group matching, not sysOID detection):
```go
p.Name = "VMach"
p.Groups = []string{"vending", "vending-per-machine", "Boot_Stage_00"}
```

Per-machine polls using `$symbol`:
- `vendMachineInventory` — `GET::/lynx/v1/machines/$symbol/inventory::` — EVERY_5_MINUTES
- `vendMachineAttributes` — `GET::/lynx/v1/machines/$symbol/attributes::` — EVERY_15_MINUTES
- `vendMachineAlerts` — `GET::/lynx/v1/machines/$symbol/alerts::` — EVERY_5_MINUTES

Register via `boot.RegisterPollaris()` in `common/Links.go` init().

### Links Routing (switch on linkid — probler pattern)

**File:** `go/vend/common/Links.go`

Add constants:
```go
VendPerMachine_Links_ID         = "VMach"
VMach_Cache_Service_Name        = "Machine"     // Fleet CRUD service
VMach_Cache_Service_Area        = byte(10)
VMach_Parser_Service_Name       = "VMPars"      // Separate parser for per-machine
VMach_Parser_Service_Area       = byte(10)
VMach_Model_Name                = "vendfleetmachine"
```

Update Links methods with switch statements (following probler's `Links.go` pattern):
```go
func (this *Links) Cache(linkid string) (string, byte) {
    switch linkid {
    case VendPerMachine_Links_ID:
        return VMach_Cache_Service_Name, VMach_Cache_Service_Area
    }
    return Vend_Cache_Service_Name, Vend_Cache_Service_Area
}

func (this *Links) Parser(linkid string) (string, byte) {
    switch linkid {
    case VendPerMachine_Links_ID:
        return VMach_Parser_Service_Name, VMach_Parser_Service_Area
    }
    return Vend_Parser_Service_Name, Vend_Parser_Service_Area
}

func (this *Links) Model(linkid string) string {
    switch linkid {
    case VendPerMachine_Links_ID:
        return VMach_Model_Name
    }
    return Vend_Model_Name
}
```

This ensures:
- "Vend" targets → VCache (management system inventory cache)
- "VMach" targets → Machine service (per-machine Fleet CRUD, persisted in DB)

## Phase 3: Create Per-Machine Targets from VCache Bridge

**File:** `go/vend/inv_vend/main.go`

In `bridgeVCacheToFleet()`, after POSTing each VendFleetMachine, also POST an `L8PTarget`:
```go
target := &l8tpollaris.L8PTarget{
    TargetId:      info.MachineId,       // "M-100001"
    LinksId:       vendcommon.VendPerMachine_Links_ID,  // "VMach"
    InventoryType: l8tpollaris.L8PTargetType_Vending_Machine,
    State:         l8tpollaris.L8PTargetState_Up,
    Hosts: map[string]*l8tpollaris.L8PHost{
        info.MachineId: {
            HostId: info.MachineId,
            Configs: map[int32]*l8tpollaris.L8PHostProtocol{
                int32(l8tpollaris.L8PProtocol_L8PRESTAPI): {
                    Protocol: l8tpollaris.L8PProtocol_L8PRESTAPI,
                    Addr:     managementIp,  // same management API IP
                    Port:     8443,
                    Timeout:  30,
                    Ainfo:    &l8tpollaris.AuthInfo{},
                },
            },
        },
    },
}
```

Use `vendcommon.PostEntity(targets.ServiceName, targets.ServiceArea, target, nic)` to create the target. The collector picks it up and starts polling with `$symbol` = `"M-100001"`.

## Phase 4: Parser Mapping for Slot Data

The per-machine `/inventory` endpoint returns:
```json
{
  "machineId": "M-100001",
  "lastUpdated": "2026-04-17T06:00:00Z",
  "slots": [{"slotNumber": 1, "productName": "...", "sku": "...", "price": 175, "capacity": 10, "currentStock": 7, "status": "ok"}, ...],
  "totalSlots": 8, "emptySlots": 1, "lowStockSlots": 2
}
```

Use `RestJsonParse` rule with explicit mapping for the field name mismatch:
- JSON `"slots"` → proto `"inventory"` (repeated VendMachineSlot)
- JSON `"totalSlots"` → proto `"totalSlots"` (exact match)
- JSON `"emptySlots"` → proto `"emptySlots"` (exact match)
- JSON `"lowStockSlots"` → proto `"lowStockSlots"` (exact match)
- JSON `"lastUpdated"` → proto `"inventoryLastUpdated"` (name mismatch — use mapping)

Poll mapping string:
```
"slots:vendfleetmachine.inventory,totalSlots:vendfleetmachine.totalslots,emptySlots:vendfleetmachine.emptyslots,lowStockSlots:vendfleetmachine.lowstockslots,lastUpdated:vendfleetmachine.inventorylastupdated"
```

**Note:** `RestJsonParse.setRepeatedProperty()` handles JSON arrays by iterating items and setting indexed properties (`inventory<{2}0>.slotnumber`, `inventory<{2}0>.productname`, etc.). The sub-field names must match the `VendMachineSlot` proto JSON field names (verified after regeneration).

**JS field name verification:** After `make-bindings.sh`, verify generated JSON names in `vend-fleet.pb.go`:
- `total_slots` → JSON `totalSlots` ✓
- `empty_slots` → JSON `emptySlots` ✓
- `low_stock_slots` → JSON `lowStockSlots` ✓
- `inventory_last_updated` → JSON `inventoryLastUpdated` ✓

**Parser activation** in `parser/main.go` — add a second activation for the "VMach" LinksId:
```go
parserService.Activate(vendcommon.VendPerMachine_Links_ID,
    &vend.VendFleetMachine{}, false, nic, "MachineId")
```
This registers the parser to handle "VMach" jobs. Parsed VendFleetMachine data is forwarded to `Links.Cache("VMach")` = `("Machine", byte(10))` — the Fleet CRUD service.

**Note:** The parser forwards via PATCH, which updates existing VendFleetMachine records (created by the VCache bridge) with slot inventory data. Records are matched by `MachineId` primary key.

## Phase 5: UI — Add Inventory Sub-Tab to Fleet

Now that slot data is being collected, add the "Inventory" sub-tab.

**fleet-config.js** — add inventory service to the machines module:
```js
services: [
    { key: 'machines', label: 'Vending Machines', icon: '🏭', endpoint: '/10/Machine', model: 'VendFleetMachine', readOnly: true },
    { key: 'inventory', label: 'Inventory', icon: '📦', endpoint: '/10/Machine', model: 'VendFleetMachine', readOnly: true },
    { key: 'machine-groups', ... },
    { key: 'locations', ... }
]
```

**fleet-section-config.js** — add inventory service to section config:
```js
services: [
    { key: 'machines', label: 'Vending Machines', icon: '🏭', isDefault: true },
    { key: 'inventory', label: 'Inventory', icon: '📦' },
    { key: 'machine-groups', label: 'Groups', icon: '📁' },
    { key: 'locations', label: 'Locations', icon: '📍' }
]
```

**sections/fleet.html** — add container ID for inventory service view:
```html
<div class="l8-service-view" data-service="inventory">
    <div class="l8-table-container" id="machines-inventory-table-container"></div>
</div>
```

**machines-columns.js** — add VendFleetMachine inventory columns (slot-focused, uses same model key but different column set). Since both services use the same model `VendFleetMachine`, the Inventory sub-tab uses a different column key. If the framework doesn't support per-service columns for the same model, use the same VendFleetMachine columns (which include slot summary fields).

**machines-forms.js** — add read-only form with slot inline table:
```js
VendFleetMachine: f.form('Vending Machine', [
    // ... existing sections ...
    f.section('Slot Inventory', [
        ...f.inlineTable('inventory', 'Slots', [
            { key: 'slotNumber', label: 'Slot', type: 'number' },
            { key: 'productName', label: 'Product', type: 'text' },
            { key: 'sku', label: 'SKU', type: 'text' },
            { key: 'price', label: 'Price', type: 'number' },
            { key: 'currentStock', label: 'Stock', type: 'number' },
            { key: 'capacity', label: 'Capacity', type: 'number' },
            { key: 'status', label: 'Status', type: 'text' }
        ])
    ])
])
```
All forms are read-only (collected data from external API, not user-editable).

**No mobile changes** — mobile parity deferred (per `mobile-rules.md`, flag for future).

## Phase 6: Verify

1. `go build ./...` — compiles
2. Run `run-local.sh`, load mock data
3. Inventory section shows targets — management API + 7 per-machine targets
4. Fleet "Vending Machines" tab shows 7 machines with data
5. Fleet "Inventory" tab shows 7 machines with slot counts
6. Click a machine → detail shows slot inventory (8 slots per machine)
7. Management Systems popup still shows 7 machines

## Traceability Matrix

| # | Action Item | Phase |
|---|-------------|-------|
| 1 | Add slot fields to VendFleetMachine proto | Phase 1 |
| 2 | Run make-bindings.sh | Phase 1 |
| 3 | Create "VMach" pollaris config with $symbol polls (Name matches LinksID) | Phase 2 |
| 4 | Add VendPerMachine_Links_ID + VMach service constants to Links.go | Phase 2 |
| 5 | Update Links methods with switch on linkid (Cache/Parser/Model route "VMach" differently) | Phase 2 |
| 6 | Register per-machine pollaris via boot.RegisterPollaris | Phase 2 |
| 7 | Create per-machine targets in VCache bridge | Phase 3 |
| 8 | Add second parser activation for VMach LinksId (→ Machine service via PATCH) | Phase 4 |
| 9 | Map slot data via RestJsonParse with explicit field name mappings | Phase 4 |
| 10 | Verify JS field names against generated pb.go after make-bindings | Phase 4 |
| 11 | Add Inventory sub-tab to fleet-config.js and fleet-section-config.js | Phase 5 |
| 12 | Add inventory container ID to sections/fleet.html | Phase 5 |
| 13 | Add slot inline table to machines-forms.js (read-only) | Phase 5 |
| 14 | Add slot summary columns to machines-columns.js | Phase 5 |
| 15 | Verify per-machine targets appear in Inventory section | Phase 6 |
| 16 | Verify slot data in Fleet Inventory tab | Phase 6 |
| 17 | Mobile parity | Deferred |

## Critical Files

| Action | File |
|--------|------|
| Modify | `proto/vend-fleet.proto` |
| Regenerate | `go/types/vend/vend-fleet.pb.go` |
| Create | `go/vend/parser/boot/vend_per_machine_polls.go` |
| Modify | `go/vend/common/Links.go` |
| Modify | `go/vend/inv_vend/main.go` |
| Modify | `go/vend/parser/main.go` |
| Modify | `go/vend/ui/web/vend-ui/fleet/fleet-config.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/fleet-section-config.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-columns.js` |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-forms.js` |
