# Refactor: Add Pollaris Targets for Machine Connection Management

## Overview

Add Pollaris `L8PTarget` service alongside the existing `VendMachine` service, following the probler pattern where L8PTarget and NetworkDevice coexist. L8PTarget handles connection configuration (offline config, polling state), while VendMachine stores rich inventory data populated by the collector/parser.

**No PRD deviation**: VendMachine is retained. L8PTarget is added as infrastructure.

## Architecture

```
                    L8PTarget                          VendMachine
                (Connection Config)               (Inventory Data)
                ┌──────────────────┐           ┌──────────────────────┐
User creates →  │ targetId: IP     │           │ machineId: IP        │
in Targets UI   │ host: REST:8443  │           │ model: TCN-ZK...     │
                │ state: Down/Up   │           │ serial: SN2024...    │
                │ linksId: "Vend"  │           │ inventory: [slots]   │
                └────────┬─────────┘           │ firmware: 22SP-3.8.2 │
                         │                     │ status: OPERATIONAL  │
                    State = Up                 └──────────┬───────────┘
                         │                                ↑
                         ▼                                │
                ┌──────────────┐    parse    ┌────────────┴───┐
                │  Collector   │ ──────────→ │     Parser     │
                │  polls REST  │   CJob      │ populates      │
                │  /api/v1/*   │             │ VendMachine    │
                └──────────────┘             └────────────────┘
```

## What Changes

### Add
- Pollaris Targets service activation in vend backend `main.go`
- Targets UI in `l8ui/targets/` (shared component, Phase 0)
- Fleet section: Targets tab alongside existing Machines tab
- L8PTarget type registration in `ui/shared.go`
- `login.json`: add `api.targetsPath` and `api.credsPath` fields for targets config.js

### Keep (no changes)
- `VendMachine` service, proto, callback, activation -- unchanged
- `VendMachineSlot` inventory -- unchanged
- All other services referencing `machineId` -- unchanged
- Machine forms, columns, enums -- unchanged
- Collector/Parser -- already use Pollaris targets
- Mock data generators -- unchanged
- Mobile UI -- unchanged (targets are desktop-only initially)

### Modify
- Fleet section config: add a "Targets" tab that loads the targets iframe
- `vend/main/main.go`: activate Pollaris targets service
- `ui/shared.go`: register L8PTarget types
- `login.json`: add api section for targets config compatibility

---

## Implementation Phases

### Phase 0: Move Targets UI to l8ui (Shared, Zero Duplication)

The probler targets UI (`targets.js`, `targets-hosts.js`, `targets-detail.js`, `targets.css`, `config.js`, `index.html`) totals ~1,376 lines. It is fully config-driven -- reads `apiPrefix` from `login.json` at runtime.

Steps:
1. Copy `probler/go/prob/newui/web/targets/*` → `l8ui/targets/`
2. Update probler's `sections/inventory.html` iframe src to `l8ui/targets/index.html`
3. Verify probler still works
4. Commit to l8ui, update submodule in both projects

### Phase 1: Backend Changes
1. Activate Pollaris Targets service in `vend/main/main.go`:
   ```go
   import "github.com/saichler/l8pollaris/go/pollaris/targets"
   // After service activation:
   targets.Activate(common.DB_CREDS, common.DB_NAME, nic)
   ```
2. Set the targets Links to use our vending Links interface:
   ```go
   targets.Links = &common.Links{}
   ```
3. Register L8PTarget types in `ui/shared.go`
4. Update `login.json` to add api section for targets config.js compatibility:
   ```json
   {
       "app": { "apiPrefix": "/vend", "healthPath": "/0/Health" },
       "api": { "prefix": "/vend", "targetsPath": "/91/Targets", "credsPath": "/75/Creds" }
   }
   ```

### Phase 2: Desktop UI Changes
1. Update Fleet section to add a "Targets" module tab alongside "Machines":
   - Fleet > Machines (existing -- shows VendMachine inventory data)
   - Fleet > Targets (new -- iframe to `l8ui/targets/index.html` for connection config)
2. Update `fleet-section-config.js` to add targets module
3. Update `fleet-config.js` to add targets module entry
4. Update `sections/fleet.html` to include targets content area

### Phase 3: Verification

Desktop:
- [ ] Fleet > Machines tab shows existing machine inventory (unchanged)
- [ ] Fleet > Targets tab loads targets iframe
- [ ] Add Target: fill IP, add Host with REST protocol, port 8443, LinksId "Vend"
- [ ] Target saved with State=Down
- [ ] Toggle State to Up → verify collector log shows polling started
- [ ] After polling, VendMachine record created/updated by parser
- [ ] Click machine row → popup shows inventory (Qty, Cap, Fill %, Status)
- [ ] All other sections (Inventory, Sales, Maintenance, etc.) work unchanged

Live simulator:
- [ ] Add target with IP 192.168.200.1, port 8443, REST protocol
- [ ] Toggle to Up → collector connects to simulator
- [ ] Verify VendMachine record populated with simulator data

---

## Traceability Matrix

| # | Section | Gap / Action Item | Phase |
|---|---------|-------------------|-------|
| 1 | Shared UI | Move targets UI from probler to l8ui/targets/ | Phase 0 |
| 2 | Shared UI | Update probler iframe src | Phase 0 |
| 3 | Shared UI | Update l8ui submodule in both projects | Phase 0 |
| 4 | Backend | Activate Pollaris Targets service in main.go | Phase 1 |
| 5 | Backend | Set targets.Links to vending Links interface | Phase 1 |
| 6 | Backend | Register L8PTarget types in ui/shared.go | Phase 1 |
| 7 | Config | Update login.json with api section for targets config.js | Phase 1 |
| 8 | Desktop UI | Add Targets tab to Fleet section | Phase 2 |
| 9 | Desktop UI | Update fleet-section-config.js and fleet-config.js | Phase 2 |
| 10 | Verification | Desktop end-to-end (machines + targets) | Phase 3 |
| 11 | Verification | Live simulator integration | Phase 3 |
| 12 | Mobile | Add targets management to mobile UI | Deferred -- desktop-only initially |
