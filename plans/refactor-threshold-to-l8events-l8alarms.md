# Refactor Threshold Evaluator to Use l8events/l8alarms

## Context

The current `go/vend/inv_vend/threshold.go` reinvents functionality available in the L8 Ecosystem:
- Custom `ThresholdRule` struct → duplicates `AlarmDefinition`
- Custom dedup via deterministic AlertId → duplicates `AlarmDefinition.dedup_enabled`
- Custom auto-clear logic → duplicates `AlarmDefinition.auto_clear_enabled`
- Custom `PostEntity(VendAlert)` → should post `alm.Alarm` to the l8alarms Alarm service

The UI currently shows alerts under the Maintenance section using a custom VendAlert CRUD service. This should be replaced with a top-level **Alarms & Events** section using l8alarms UI files — the same pattern probler uses.

## Infrastructure Gap (to be addressed upstream later)

l8alarms has `AlarmDefinition` with event matching criteria (`EventPattern`, `ThresholdCount`, `DedupKeyExpression`, `ClearEventPattern`) but **no event processor** that reads incoming events and creates alarms from matching definitions. The event→alarm bridge is not yet built in l8alarms.

**Current workaround (this plan):** The threshold evaluator acts as the event processor — it evaluates thresholds, posts `alm.Event` for audit/history, and directly POSTs `alm.Alarm` to the l8alarms Alarm service. This gives us l8alarms correlation, notification, and escalation for free.

**Future (upstream):** Move the event→alarm matching logic into l8alarms as a generic event processor. At that point the threshold evaluator would only post events and l8alarms would handle alarm creation automatically.

## Architecture After Refactoring

```
Threshold evaluator (inv_vend)
    │
    ├── Evaluates slot fill % and machine fill %
    │
    ├── Posts alm.Event to l8alarms Event service (audit trail)
    │       eventType = "VEND_STOCK_LOW" / "VEND_SLOT_LOW"
    │       severity, sourceId=machineId, attributes
    │
    └── Posts alm.Alarm to l8alarms Alarm service (triggers lifecycle)
        │   definitionId = "VEND-DEF-002" (matches AlarmDefinition)
        │   nodeId = machineId
        │   dedupKey = "machineId:slotId:definitionId"
        │   severity, name, description
        │
        └── l8alarms After hooks fire automatically:
            ├── runCorrelation — suppresses symptoms
            ├── runNotification — dispatches via l8notify
            └── runEscalation — escalation timers

Auto-clear: evaluator POSTs Alarm with state=CLEARED when threshold recovers
```

## Dependencies

- `l8alarms/go/alm/services` — for `ActivateAlmServices()`
- `l8alarms/go/types/alm` — for `alm.Alarm`, `alm.Event`, `alm.AlarmDefinition`
- `l8events/go/types/l8events` — for `Severity`, `AlarmState` enums (already vendored)

Note: `l8alarms/go/alm/ui/RegisterAlmTypes()` imports l8topology which is NOT vendored. We register alarm types manually instead (only the types we need for web API).

## Risk: l8types Version Pin

Current go.mod pins l8types to `v0.0.0-20260502164503-192680a11be2` to avoid a breaking `Register()` change. Re-vendoring l8alarms service packages may pull in a newer l8types. If this happens:
- Keep the pin, verify that l8alarms compiles against the pinned version
- If incompatible, defer this refactoring until l8ql is updated upstream

## Phase 0: Re-vendor Dependencies

**Action:** Add imports for `l8alarms/go/alm/services` and re-run vendor refresh with the l8types pin.

```bash
cd go
rm -rf go.sum go.mod vendor
go mod init
GOPROXY=direct GOPRIVATE=github.com go get github.com/saichler/l8types@v0.0.0-20260502164503-192680a11be2 2>/dev/null
GOPROXY=direct GOPRIVATE=github.com go mod tidy
go mod vendor
```

**Verification:** `go build ./...` compiles with new imports.

**If l8types conflict:** Stop here and report incompatibility to user.

## Phase 1: Activate l8alarms Services in vend_demo

**File: `go/vend/main/main.go`**

Add after `services.ActivateAllServices()`:
```go
// Activate l8alarms services (alarm definitions, alarms, correlation, notification)
almServices.ActivateAlmServices(common.DB_CREDS, common.DB_NAME, nic)
```

This activates all 10+ l8alarms services (AlmDef, Alarm, Event, CorrRule, NotifPolicy, EscPolicy, MainWin, AlarmFilter, ArchivedAlarm, ArchivedEvent, Enrichment) in the same process as vend services.

**File: `go/vend/ui/main/main.go`**

Register alarm types for web API introspection (manual registration, NOT `RegisterAlmTypes()` to avoid l8topology dependency):
```go
l8common.RegisterType(resources, &alm.AlarmDefinition{}, &alm.AlarmDefinitionList{}, "DefinitionId")
l8common.RegisterType(resources, &alm.Alarm{}, &alm.AlarmList{}, "AlarmId")
l8common.RegisterType(resources, &alm.Event{}, &alm.EventList{}, "EventId")
l8common.RegisterType(resources, &alm.CorrelationRule{}, &alm.CorrelationRuleList{}, "RuleId")
l8common.RegisterType(resources, &alm.NotificationPolicy{}, &alm.NotificationPolicyList{}, "PolicyId")
l8common.RegisterType(resources, &alm.EscalationPolicy{}, &alm.EscalationPolicyList{}, "PolicyId")
l8common.RegisterType(resources, &alm.MaintenanceWindow{}, &alm.MaintenanceWindowList{}, "WindowId")
l8common.RegisterType(resources, &alm.AlarmFilter{}, &alm.AlarmFilterList{}, "FilterId")
l8common.RegisterType(resources, &alm.ArchivedAlarm{}, &alm.ArchivedAlarmList{}, "AlarmId")
l8common.RegisterType(resources, &alm.ArchivedEvent{}, &alm.ArchivedEventList{}, "EventId")
```

## Phase 2: Refactor threshold.go to Post alm.Alarm

**File: `go/vend/inv_vend/threshold.go`**

Replace the entire custom alarm lifecycle:

1. **Remove**: `ThresholdRule` struct, `inventoryRules` var, `evaluateRule()`, `activeAlertExists()`, `postAlert()`, `clearAlert()`
2. **Remove**: imports for VendAlert types and alerts service
3. **Keep**: `evaluateThresholds()` goroutine structure and `evaluateMachine()` metric calculation
4. **Replace**: custom alert posting with `alm.Alarm` POSTs to l8alarms

**Threshold rules as data** (replaces custom ThresholdRule struct):

```go
type thresholdDef struct {
    definitionId string   // matches AlarmDefinition.DefinitionId
    eventType    string   // e.g. "VEND_STOCK_LOW"
    clearType    string   // e.g. "VEND_STOCK_NORMAL"
    thresholdType string  // "LOWER" or "UPPER"
    warningValue float64
    criticalValue float64
    clearValue   float64
}
```

These reference the seed data AlarmDefinitions by ID — the evaluator knows which definition it's creating alarms for.

**Alarm creation** — when threshold crossed:
```go
alarm := &alm.Alarm{
    AlarmId:      fmt.Sprintf("THR-%s-%s-%s", machineId, def.definitionId, slotId),
    DefinitionId: def.definitionId,
    NodeId:       machineId,
    Name:         def.eventType,
    State:        l8events.AlarmState_ALARM_STATE_ACTIVE,
    Severity:     severity,
    DedupKey:     fmt.Sprintf("%s:%s:%s", machineId, slotId, def.definitionId),
    Description:  message,
}
vendcommon.PostEntity("Alarm", 10, alarm, nic)
```

**Auto-clear** — when threshold recovers:
```go
alarm := &alm.Alarm{
    AlarmId: existingAlarmId,
    State:   l8events.AlarmState_ALARM_STATE_CLEARED,
}
vendcommon.PutEntity("Alarm", 10, alarm, nic)
```

**Dedup**: l8alarms Alarm service `protectSystemFields` ensures identity fields are immutable after creation. The evaluator checks if an alarm with that ID already exists before POSTing (same pattern as current `activeAlertExists` but querying `alm.Alarm` instead of `VendAlert`).

**Tracking "last state"**: Local `map[string]bool` of `{machineId+slotId+definitionId} → wasAlerting`. When a metric was alerting and now recovers, PUT the alarm to CLEARED state.

**File: `go/vend/inv_vend/main.go`**

- Remove VendAlert type registration (no longer needed)
- Add `alm.Alarm`, `alm.AlarmList`, `alm.Event`, `alm.EventList` type registration
- Keep VendFleetMachine and L8PTarget registrations

## Phase 3: Seed Alarm Definitions via Mocks

**File: `go/tests/mocks/main_phases.go`**

Add a new phase after business foundation:
```go
runPhase("Phase 2: Alarm Definitions", func() error {
    return seedAlarmDefinitions(client)
})
```

**File: `go/tests/mocks/phase_alarms.go`** (new)

Posts the existing seed data to l8alarms services:
- `POST /vend/10/AlmDef` with `AlarmDefinitionList` — from `seeddata.GetAlarmDefinitions()`
- `POST /vend/10/CorrRule` with `CorrelationRuleList` — from `seeddata.GetCorrelationRules()`

## Phase 4: Add Top-Level Alarms Section (like probler)

Replace the custom VendAlert UI under Maintenance with l8alarms' own UI files — a full top-level **Alarms & Events Management** section.

### Step 1: Copy l8alarms UI files

Copy the `alm/` directory from l8alarms into the vend web directory:

**Source:** `../l8alarms/go/alm/ui/web/alm/`
**Destination:** `go/vend/ui/web/alm/`

This brings in the complete alarms UI:
```
alm/
├── alm-section-config.js          # Section config for Layer8SectionGenerator
├── alm-config.js                  # Module config (endpoints, models, views)
├── alm-init.js                    # Layer8DModuleFactory.create() bootstrap
├── alm.css                        # Alarms section styles
├── alarms/
│   ├── alarms-enums.js            # Severity, AlarmState, DefinitionStatus
│   ├── alarms-columns.js          # Alarm, AlarmDefinition, AlarmFilter columns
│   ├── alarms-forms.js            # Alarm, AlarmDefinition, AlarmFilter forms
│   ├── alarms-correlation-tree.js # Correlation tree in alarm detail popup
│   └── alarms-correlation-tree.css
├── events/
│   ├── events-enums.js
│   ├── events-columns.js
│   └── events-forms.js
├── correlation/
│   ├── correlation-enums.js
│   ├── correlation-columns.js
│   └── correlation-forms.js
├── policies/
│   ├── policies-enums.js
│   ├── policies-columns.js
│   └── policies-forms.js
├── maintenance/                   # l8alarms maintenance windows (different from vend Maintenance)
│   ├── maintenance-enums.js
│   ├── maintenance-columns.js
│   └── maintenance-forms.js
└── archive/
    ├── archive-enums.js
    ├── archive-columns.js
    └── archive-forms.js
```

These are l8alarms' own files — configuration only, no project-specific code. They define columns, forms, and enums for the l8alarms protobuf types (`Alarm`, `AlarmDefinition`, `Event`, `CorrelationRule`, etc.).

### Step 2: Create section HTML

**File: `go/vend/ui/web/sections/alarms.html`** (new)

```html
<div id="alarms-section-placeholder"></div>
<script>
    (function() {
        var placeholder = document.getElementById('alarms-section-placeholder');
        if (placeholder && window.Layer8SectionGenerator) {
            placeholder.outerHTML = Layer8SectionGenerator.generate('alarms');
        }
    })();
</script>
```

### Step 3: Wire into sections.js

**File: `go/vend/ui/web/js/sections.js`**

Add section mapping:
```js
alarms: 'sections/alarms.html'
```

Add initializer:
```js
alarms: () => {
    if (typeof initializeAlm === 'function') {
        initializeAlm();
    }
}
```

### Step 4: Wire into app.html

**Sidebar nav** — add Alarms link (after Maintenance):
```html
<li><a href="#" data-section="alarms" class="nav-link sidebar-item">
    <span class="nav-icon">🔔</span><span>Alarms</span>
</a></li>
```

**CSS** — add before module-specific CSS:
```html
<link rel="stylesheet" href="alm/alm.css">
<link rel="stylesheet" href="alm/alarms/alarms-correlation-tree.css">
```

**JS** — add section config after other section configs (around line 170):
```html
<script src="alm/alm-section-config.js"></script>
```

**JS** — add module scripts (after maintenance module, before routes module):
```html
<!-- JS: Alarms Module (l8alarms) -->
<script src="alm/alm-config.js"></script>
<script src="alm/alarms/alarms-enums.js"></script>
<script src="alm/alarms/alarms-columns.js"></script>
<script src="alm/alarms/alarms-forms.js"></script>
<script src="alm/events/events-enums.js"></script>
<script src="alm/events/events-columns.js"></script>
<script src="alm/events/events-forms.js"></script>
<script src="alm/correlation/correlation-enums.js"></script>
<script src="alm/correlation/correlation-columns.js"></script>
<script src="alm/correlation/correlation-forms.js"></script>
<script src="alm/policies/policies-enums.js"></script>
<script src="alm/policies/policies-columns.js"></script>
<script src="alm/policies/policies-forms.js"></script>
<script src="alm/maintenance/maintenance-enums.js"></script>
<script src="alm/maintenance/maintenance-columns.js"></script>
<script src="alm/maintenance/maintenance-forms.js"></script>
<script src="alm/archive/archive-enums.js"></script>
<script src="alm/archive/archive-columns.js"></script>
<script src="alm/archive/archive-forms.js"></script>
<script src="alm/alm-init.js"></script>
<script src="alm/alarms/alarms-correlation-tree.js"></script>
```

### Step 5: Update Maintenance section — remove Alerts

The Maintenance section currently has Alerts, Work Orders, and Service Visits under one "alerts" module. After this change, Alerts move to the new Alarms section. Maintenance keeps Work Orders and Service Visits.

**File: `go/vend/ui/web/vend-ui/maintenance/maintenance-config.js`**

Remove the `alerts` service entry, rename the module from `alerts` to `work-orders`:
```js
Layer8ModuleConfigFactory.create({
    namespace: 'Maintenance',
    modules: {
        'work-orders': mod('Work Orders', '', [
            svc('work-orders', 'Work Orders', '', '/10/WorkOrder', 'VendWorkOrder',
                { supportedViews: ['table', 'kanban', 'gantt'] }),
            svc('service-visits', 'Service Visits', '', '/10/SvcVisit', 'VendServiceVisit')
        ])
    },
    submodules: ['MaintenanceAlerts']
});
```

**File: `go/vend/ui/web/vend-ui/maintenance/maintenance-section-config.js`**

Update to remove Alerts service:
```js
Layer8SectionConfigs.register('maintenance', {
    title: 'Maintenance',
    subtitle: 'Work Orders, Service Visits',
    icon: '🔧',
    initFn: 'initializeMaintenance',
    modules: [{
        key: 'work-orders', label: 'Work Orders', icon: '🔧', isDefault: true,
        services: [
            { key: 'work-orders', label: 'Work Orders', icon: '📝', isDefault: true },
            { key: 'service-visits', label: 'Service Visits', icon: '🔧' }
        ]
    }]
});
```

**File: `go/vend/ui/web/vend-ui/maintenance/maintenance-init.js`**

Update `defaultModule` and `sectionSelector` to match new module key:
```js
Layer8DModuleFactory.create({
    namespace: 'Maintenance',
    defaultModule: 'work-orders',
    defaultService: 'work-orders',
    sectionSelector: 'work-orders',
    initializerName: 'initializeMaintenance',
    requiredNamespaces: ['MaintenanceAlerts']
});
```

Note: The `MaintenanceAlerts` enums/columns/forms files stay — they still define VendWorkOrder and VendServiceVisit. The VendAlert definitions in those files become unused but harmless. They can be cleaned up in Phase 5 (deferred).

### Step 6: Update Dashboard

**File: `go/vend/ui/web/vend-ui/dashboard/dashboard-init.js`**

Change the alerts query from:
```js
var alertQuery = encodeURIComponent(JSON.stringify({ text: 'select * from VendAlert where status=1' }));
fetch(prefix + '/10/Alert?body=' + alertQuery, {
```
To:
```js
var alertQuery = encodeURIComponent(JSON.stringify({ text: 'select * from Alarm where state=1' }));
fetch(prefix + '/10/Alarm?body=' + alertQuery, {
```

Update rendering to use l8alarms Alarm fields:
- `a.machineId` → `a.nodeId`
- `a.code` → `a.name`
- `a.severity === 3` → `a.severity === 5` (l8events SEVERITY_CRITICAL = 5)
- Severity labels: 1=INFO, 2=WARNING, 3=MINOR, 4=MAJOR, 5=CRITICAL

### Step 7: Copy alm/ to demo in run-local.sh

**File: `go/run-local.sh`**

The existing `cp -r vend/ui/web demo/.` already copies the entire web directory, which will include `alm/` after the files are added. No change needed to run-local.sh.

## Phase 5: Remove Custom VendAlert Service (Deferred)

Once l8alarms is working:
- Remove `go/vend/maintenance/alerts/` service directory
- Remove VendAlert proto definitions from `proto/vend-fleet.proto`
- Remove alerts-enums.js, alerts-columns.js, alerts-forms.js (VendAlert parts only)
- Clean up MaintenanceAlerts namespace to only contain WorkOrder and ServiceVisit

Defer this to avoid scope creep — the unused VendAlert definitions are harmless.

## Traceability Matrix

| # | Gap / Action Item | Phase |
|---|-------------------|-------|
| 1 | Re-vendor l8alarms service packages | Phase 0 |
| 2 | Activate l8alarms services in vend main | Phase 1 |
| 3 | Register alarm types for web API (manual, no l8topology) | Phase 1 |
| 4 | Rewrite threshold.go to POST alm.Alarm | Phase 2 |
| 5 | Replace custom dedup with l8alarms dedup (AlarmId-based) | Phase 2 |
| 6 | Replace custom auto-clear with PUT state=CLEARED | Phase 2 |
| 7 | Track last-state for clear detection | Phase 2 |
| 8 | Register alm types in inv_vend binary | Phase 2 |
| 9 | Seed AlarmDefinitions via mock generator | Phase 3 |
| 10 | Seed CorrelationRules via mock generator | Phase 3 |
| 11 | Copy l8alarms UI files (alm/ directory) | Phase 4 Step 1 |
| 12 | Create sections/alarms.html | Phase 4 Step 2 |
| 13 | Wire alarms section into sections.js | Phase 4 Step 3 |
| 14 | Wire alm CSS and JS into app.html | Phase 4 Step 4 |
| 15 | Add Alarms sidebar nav link | Phase 4 Step 4 |
| 16 | Remove Alerts from Maintenance config | Phase 4 Step 5 |
| 17 | Update Maintenance section/init (work-orders default) | Phase 4 Step 5 |
| 18 | Update dashboard to query l8alarms Alarm | Phase 4 Step 6 |
| 19 | Remove VendAlert service and proto (optional) | Phase 5 (deferred) |
| 20 | Flag missing event processor as upstream gap | Documented above |

## Phase 6: End-to-End Verification

1. `go build ./...` — all packages compile
2. `run-local.sh` — all services start
3. Upload mocks — alarm definitions and correlation rules seeded
4. Navigate to **Alarms** section — verify section loads with 6 module tabs
5. Click **Definitions** tab — verify 10 alarm definitions from seed data
6. Click **Correlation > Rules** tab — verify 3 correlation rules
7. Wait for threshold evaluator cycle (60s initial + 5min)
8. Check `demo.log` for alarm POST log messages
9. Click **Alarms > Active Alarms** — verify alarms created by threshold evaluator
10. Navigate to **Maintenance** — verify only Work Orders and Service Visits (no Alerts)
11. Navigate to **Dashboard** — verify active alarm count from l8alarms
12. Modify simulator stock levels → verify auto-clear fires (alarm state changes to CLEARED)

## Critical Files

| Action | File |
|--------|------|
| Modify | `go/vend/inv_vend/threshold.go` (rewrite to POST alm.Alarm) |
| Modify | `go/vend/inv_vend/main.go` (swap VendAlert → alm type registrations) |
| Modify | `go/vend/main/main.go` (add l8alarms activation) |
| Modify | `go/vend/ui/main/main.go` (register alarm types for web) |
| Create | `go/tests/mocks/phase_alarms.go` (seed alarm definitions) |
| Modify | `go/tests/mocks/main_phases.go` (add alarm phase) |
| Copy   | `go/vend/ui/web/alm/` (entire l8alarms UI directory from ../l8alarms) |
| Create | `go/vend/ui/web/sections/alarms.html` (section placeholder) |
| Modify | `go/vend/ui/web/js/sections.js` (add alarms section) |
| Modify | `go/vend/ui/web/app.html` (sidebar link, CSS, JS script tags) |
| Modify | `go/vend/ui/web/vend-ui/maintenance/maintenance-config.js` (remove alerts) |
| Modify | `go/vend/ui/web/vend-ui/maintenance/maintenance-section-config.js` (remove alerts) |
| Modify | `go/vend/ui/web/vend-ui/maintenance/maintenance-init.js` (work-orders default) |
| Modify | `go/vend/ui/web/vend-ui/dashboard/dashboard-init.js` (query Alarm) |
| Modify | `go/go.mod` + `go/vendor/` (re-vendor) |
