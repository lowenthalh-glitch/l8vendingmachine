# Plan: Bring l8vendingmachine to Full Global Rules Compliance

**Date**: 2026-06-14
**Status**: Pending Approval

---

## 1. Audit Summary

An audit of the l8vendingmachine project against ALL global rules in `~/.claude/rules/` identified **8 violation categories**. The project is compliant in many areas (proto design, list conventions, enum zero values, l8ui submodule setup, login.json adaptation, VB-pattern service callbacks, run-local.sh). The violations below are the gaps that must be addressed.

---

## 2. Violations Found

### V1: Dockerfiles Use Wrong Base Images (CRITICAL)
**Rule**: `deployment-artifacts.md`, `l8pollaris-binary-deployment.md`
**Finding**: All 6 Dockerfiles use `saichler/erp-security:latest` or `saichler/erp-postgres:latest` instead of project-specific `saichler/vend-security:latest` / `saichler/vend-postgres:latest`.

| File | Current Base | Required Base |
|------|-------------|---------------|
| `go/vend/vnet/Dockerfile` | `erp-security` | `vend-security` |
| `go/vend/main/Dockerfile` | `erp-postgres` | `vend-postgres` |
| `go/vend/ui/main/Dockerfile` | `erp-security` | `vend-security` |
| `go/vend/collector/Dockerfile` | `erp-security` | `vend-security` |
| `go/vend/parser/Dockerfile` | `erp-security` | `vend-security` |
| `go/vend/inv_vend/Dockerfile` | `erp-security` | `vend-security` |

**Risk**: Protobuf namespace conflict panics at runtime — the `loader.so` security plugin in erp base images has a different dependency tree than what l8vendingmachine needs.

### V2: Missing 3 of 4 K8s Deployment Modes (CRITICAL)
**Rule**: `k8s-three-deployment-modes.md`
**Finding**: Only local-mode YAMLs exist (`go/k8s/vend*.yaml`). Missing:
- `k8s/vend-baremetal.yaml` (local-path PVCs)
- `k8s/vend-gke.yaml` (GCE PD storage)
- `k8s/vend-kind.yaml` (KIND built-in standard StorageClass)
- `k8s/kind-start.sh` and `k8s/kind-stop.sh`

### V3: Missing Log Services (CRITICAL)
**Rule**: `log-services-required.md`
**Finding**: No `go/vend/log-vnet/` or `go/vend/log-agent/` directories exist. These are required infrastructure services for every Layer 8 project.

### V4: Events Service Not Activated
**Rule**: `events-service-required.md`
**Finding**:
- `go/vend/main/main.go` does NOT call `evtservices.ActivateEvents()`
- `go/vend/ui/shared.go` does NOT register `EventRecord` type
- Note: l8alarms services ARE activated individually, which is a valid pattern for alarms. But the L8Events service (distinct from l8alarms `events.Activate`) is missing.

### V5: Missing login.html
**Rule**: `index-html-redirect.md`
**Finding**: `go/vend/ui/web/index.html` redirects to `l8ui/login/` (the l8ui login component directory), not to `login.html`. There is no `login.html` file at the web root. Per the rule, `index.html` must redirect to `login.html`, and `login.html` must exist.

### V6: Missing Security Config
**Rule**: `security-config-structure.md`, `security-provider-interface.md`
**Finding**: No `go/secure/plugin/vend/vend.json` security config file exists. The project has no defined roles, deny rules, or user provisioning configuration.

### V7: 4 Analytics Services Missing ServiceCallbacks
**Rule**: `maintainability.md` (ServiceCallback auto-generate ID)
**Finding**: These services have no `*ServiceCallback.go` — POST operations will fail with "XxxId is required" errors:
- `go/vend/analytics/restock/` (RestockService)
- `go/vend/analytics/snapshots/` (SnapshotService)
- `go/vend/analytics/profiles/` (ProfileService)
- `go/vend/analytics/topperformers/` (TopPerformerService)

Note: These analytics services are primarily populated by the `inv_vend` process via bridging, not by user CRUD. However, if any UI or API POST is attempted, it will fail without a callback to auto-generate the primary key.

### V8: Mobile Parity Gaps
**Rule**: `mobile-rules.md` (Desktop/Mobile Functional Parity)
**Finding**: 4 desktop sections have no mobile equivalent module directories:
- `alarms` — desktop has `sections/alarms.html`, no mobile JS module
- `map` — desktop has `sections/map.html`, no mobile JS module
- `nayax` — desktop has `sections/nayax.html`, no mobile JS module
- `dashboard` — desktop has `sections/dashboard.html`, no mobile JS module (mobile has `app-core.js` but no dedicated dashboard module)

---

## 3. Items NOT Violated (Compliant)

The following rules were checked and are already compliant:
- Proto enum zero values (`proto-enum-zero-value.md`) — all 16 proto files checked
- Proto list convention (`proto-list-convention.md`) — all use `repeated X list = 1` + metadata
- Cross-references use ID fields, not struct pointers (`prime-object-references.md`)
- l8ui is a git submodule (`l8ui-copy-to-new-project.md`) — `.gitmodules` confirmed
- `login.json` properly adapted (`login-json-adaptation.md`) — `apiPrefix: "/vend"`, title correct
- `run-local.sh` exists and is comprehensive (`run-local-script.md`)
- ServiceCallback validation uses VB pattern from l8common (functionally equivalent to `common.GenerateID`)
- No Go generics used (`no-go-generics.md`)
- No `l8secure` imports (`never-import-l8secure.md`)
- Test files are under `go/tests/` (`test-location-and-approach.md`)
- `app.html` body structure follows l8erp pattern (`app-html-body-from-l8erp.md`)
- Desktop modules use `Layer8DModuleFactory.create()` pattern
- Mobile modules use `Layer8MModuleRegistry.create()` pattern
- `css/base-core.css` and `css/responsive.css` exist

---

## 4. Traceability Matrix

| # | Violation | Gap | Phase |
|---|-----------|-----|-------|
| 1 | V1 | vnet Dockerfile uses erp-security | Phase 1 |
| 2 | V1 | main Dockerfile uses erp-postgres | Phase 1 |
| 3 | V1 | ui Dockerfile uses erp-security | Phase 1 |
| 4 | V1 | collector Dockerfile uses erp-security | Phase 1 |
| 5 | V1 | parser Dockerfile uses erp-security | Phase 1 |
| 6 | V1 | inv_vend Dockerfile uses erp-security | Phase 1 |
| 7 | V2 | Missing vend-baremetal.yaml | Phase 2 |
| 8 | V2 | Missing vend-gke.yaml | Phase 2 |
| 9 | V2 | Missing vend-kind.yaml | Phase 2 |
| 10 | V2 | Missing kind-start.sh | Phase 2 |
| 11 | V2 | Missing kind-stop.sh | Phase 2 |
| 12 | V3 | Missing log-vnet directory (main.go, Dockerfile, build.sh) | Phase 3 |
| 13 | V3 | Missing log-agent directory (main.go, Dockerfile, build.sh) | Phase 3 |
| 14 | V3 | log-vnet and log-agent not in build-all-images.sh | Phase 3 |
| 15 | V3 | log-vnet and log-agent not in K8s YAMLs | Phase 3 |
| 16 | V3 | log-vnet and log-agent not in deploy.sh/undeploy.sh | Phase 3 |
| 17 | V4 | evtservices.ActivateEvents() not called in main.go | Phase 4 |
| 18 | V4 | EventRecord type not registered in ui/shared.go | Phase 4 |
| 19 | V5 | No login.html at web root | Phase 5 |
| 20 | V5 | index.html redirects to l8ui/login/ instead of login.html | Phase 5 |
| 21 | V6 | No go/secure/plugin/vend/vend.json | Phase 6 |
| 22 | V7 | RestockService missing callback | Phase 7 |
| 23 | V7 | SnapshotService missing callback | Phase 7 |
| 24 | V7 | ProfileService missing callback | Phase 7 |
| 25 | V7 | TopPerformerService missing callback | Phase 7 |
| 26 | V8 | alarms module missing on mobile | Phase 8 |
| 27 | V8 | map module missing on mobile | Phase 8 |
| 28 | V8 | nayax module missing on mobile | Phase 8 |
| 29 | V8 | dashboard module missing on mobile | Phase 8 |

---

## 5. Phase Breakdown

### Phase 1: Fix All Dockerfiles — Replace Base Images
**Scope**: 6 Dockerfiles
**Effort**: Small (text replacement)

Replace base image references in all Dockerfiles:

| File | Change |
|------|--------|
| `go/vend/vnet/Dockerfile` | `erp-security` → `vend-security` |
| `go/vend/main/Dockerfile` | `erp-postgres` → `vend-postgres` |
| `go/vend/ui/main/Dockerfile` | `erp-security` → `vend-security` |
| `go/vend/collector/Dockerfile` | `erp-security` → `vend-security` |
| `go/vend/parser/Dockerfile` | `erp-security` → `vend-security` |
| `go/vend/inv_vend/Dockerfile` | `erp-security` → `vend-security` |

**Prerequisite**: The user must build `saichler/vend-security:latest` and `saichler/vend-postgres:latest` base images before these Dockerfiles can produce working images. This is outside the scope of this plan.

### Phase 2: Create Missing K8s Deployment Mode YAMLs + KIND Scripts
**Scope**: 5 new files in `go/k8s/`
**Effort**: Medium (copy from probler, adapt)
**Reference**: `../probler/k8s/`

1. Create `go/k8s/vend-baremetal.yaml` — convert all local-mode YAMLs:
   - DaemonSets → StatefulSets with podAntiAffinity
   - hostPath volumes → volumeClaimTemplates with `vend-local-storage` StorageClass
   - Add `rancher.io/local-path` StorageClass definition
2. Create `go/k8s/vend-gke.yaml` — convert all local-mode YAMLs:
   - Keep DaemonSets as DaemonSets
   - hostPath → shared PVC `vend-data` (50Gi) with `kubernetes.io/gce-pd` StorageClass
3. Create `go/k8s/vend-kind.yaml` — copy baremetal and:
   - Remove custom StorageClass definition
   - Replace `storageClassName: vend-local-storage` → `storageClassName: standard`
4. Create `go/k8s/kind-start.sh` — copy from probler, adapt cluster name and image list
5. Create `go/k8s/kind-stop.sh` — copy from probler, adapt cluster name

### Phase 3: Add Log Services (log-vnet + log-agent)
**Scope**: 6 new files + updates to build/deploy scripts
**Effort**: Medium (copy from probler, adapt)
**Reference**: `../probler/go/prob/log-vnet/`, `../probler/go/prob/log-agent/`

1. Create `go/vend/log-vnet/` directory:
   - `main.go` — copy from probler, change image name and binary name
   - `Dockerfile` — use `saichler/vend-security:latest` as base
   - `build.sh` — build `saichler/vend-logs-vnet:latest`
2. Create `go/vend/log-agent/` directory:
   - `main.go` — copy from probler, change image name, set `LOGPATH` default to `/data/logs/vend`
   - `Dockerfile` — use `saichler/vend-security:latest` as base
   - `build.sh` — build `saichler/vend-log-agent:latest`
3. Update `go/build-all-images.sh` — add log-vnet and log-agent (before other services)
4. Update `go/k8s/deploy.sh` — add log service YAMLs
5. Update `go/k8s/undeploy.sh` — add log service YAMLs
6. Add log-vnet and log-agent entries to ALL four K8s YAML files (local, baremetal, gke, kind)

### Phase 4: Activate Events Service
**Scope**: 2 files modified
**Effort**: Small

1. **`go/vend/main/main.go`**: Add import `evtservices "github.com/saichler/l8events/go/services"` and call `evtservices.ActivateEvents(common.DB_CREDS, common.DB_NAME, nic)` after `services.ActivateAllServices()`.
2. **`go/vend/ui/shared.go`**: Add import `l8events "github.com/saichler/l8types/go/types/l8events"` and register `common.RegisterType(resources, &l8events.EventRecord{}, &l8events.EventRecordList{}, "EventId")`.

**Note**: After adding the import, the user will need to re-vendor to pull in `l8events`.

### Phase 5: Add login.html and Fix index.html Redirect
**Scope**: 2 files (1 new, 1 modified)
**Effort**: Small

1. Create `go/vend/ui/web/login.html` — copy from `../l8erp/go/erp/ui/web/login.html` or use the l8ui login component pattern. The l8ui shared login page is at `l8ui/login/index.html` — if the project uses that directly, create a `login.html` that either redirects to it or embeds its content.
2. Update `go/vend/ui/web/index.html` — change `<meta http-equiv="refresh" content="0; url=l8ui/login/">` to `url=login.html` and the JS `window.location.href` to `'login.html'`.

### Phase 6: Create Security Config
**Scope**: 1 new directory + 1 new JSON file
**Effort**: Medium

1. Create directory `go/secure/plugin/vend/`
2. Create `go/secure/plugin/vend/vend.json` with:
   - `credentials` section (postgres connection)
   - `key` and `secret` (AES encryption key, shared secret)
   - `roles` section with at minimum:
     - `admin` role — full access to all services
     - `operator` role — CRUD access to fleet, inventory, sales, maintenance, routes
     - `viewer` role — read-only access
   - `users` section with at minimum:
     - `admin` user with admin role
   - `sysconfig` section (dataStoreType, dataStoreName, webPort, etc.)

**Reference**: `../l8erp/go/secure/plugin/erp/erp.json` or `l8secure/go/secure/plugin/phy/phy.json`

### Phase 7: Add ServiceCallbacks for 4 Analytics Services
**Scope**: 4 new files
**Effort**: Small

Create `*ServiceCallback.go` for each service using the VB pattern:

1. `go/vend/analytics/restock/RestockServiceCallback.go` — auto-generate `RecommendationId`
2. `go/vend/analytics/snapshots/SnapshotServiceCallback.go` — auto-generate `SnapshotId`
3. `go/vend/analytics/profiles/ProfileServiceCallback.go` — auto-generate `ProfileId`
4. `go/vend/analytics/topperformers/TopPerformerServiceCallback.go` — auto-generate `PerformerId`

Each callback follows the existing project pattern:
```go
func (cb *XxxServiceCallback) Before(elements ifs.IElements, action ifs.Action, vnic ifs.IVNic) error {
    entity := elements.Element().(*vend.XxxType)
    if action == ifs.POST {
        return common.NewValidation(entity).Require("XxxId").Build()
    }
    return nil
}
```

### Phase 8: Mobile Parity — Add Missing Mobile Modules
**Scope**: 4 desktop sections need mobile equivalents
**Effort**: Large
**Platform**: Mobile (`go/vend/ui/web/m/`)

For each missing mobile module, follow the `adding-module-mobile.md` pattern:

1. **Dashboard** (`m/js/dashboard/`):
   - Mobile dashboard with KPI cards (similar to desktop `sections/dashboard.html`)
   - Use `Layer8MWidget.render()` for stats
   - Add nav config entry

2. **Alarms** (`m/js/alarms/`):
   - Create enums, columns, forms for l8alarms types (Alarm, Event, AlarmDefinition, etc.)
   - Reuse existing l8alarms JS from l8ui SYS module or create vend-specific wrappers
   - Add nav config entry

3. **Map** (`m/js/map/`):
   - Mobile map view may be simplified (map components require visible container dimensions)
   - Could be deferred if map functionality requires desktop-only libraries
   - Mark as "Deferred — requires platform-specific map SDK evaluation" if not feasible on mobile

4. **Nayax** (`m/js/nayax/`):
   - Nayax payment integration section
   - Create mobile enums, columns, forms matching desktop definitions
   - Add nav config entry

For each implemented module:
- Create `*-enums.js`, `*-columns.js`, `*-forms.js` (with `primary`/`secondary` annotations)
- Create registry index file using `Layer8MModuleRegistry.create()`
- Add nav config entries to `m/js/nav-configs/layer8m-nav-config-vend.js`
- Register module in `layer8m-nav-data.js` lookup arrays
- Add `<script>` tags to `m/app.html`

---

## 6. Phase Dependencies

```
Phase 1 (Dockerfiles)     — no dependencies, can start immediately
Phase 2 (K8s modes)       — depends on Phase 3 (log services must be in the YAMLs)
Phase 3 (Log services)    — no dependencies
Phase 4 (Events)          — no dependencies (requires re-vendor after)
Phase 5 (login.html)      — no dependencies
Phase 6 (Security config) — no dependencies
Phase 7 (Callbacks)       — no dependencies
Phase 8 (Mobile parity)   — no dependencies

Recommended order: 1 → 3 → 2 → 4 → 5 → 6 → 7 → 8
(Phase 3 before 2 so log services are included in the new K8s YAMLs)
```

---

## 7. End-to-End Verification (Phase 9)

After all phases are implemented:

1. **Build verification**:
   - [ ] `cd go && go build ./...` — all Go code compiles
   - [ ] `cd go && go vet ./...` — no vet errors

2. **Dockerfile verification**:
   - [ ] `grep "erp-security\|erp-postgres" go/vend/*/Dockerfile go/vend/*/*/Dockerfile` returns 0 results
   - [ ] All Dockerfiles reference `vend-security` or `vend-postgres`

3. **K8s verification**:
   - [ ] All 4 YAML files exist: `go/k8s/vend-local.yaml`, `vend-baremetal.yaml`, `vend-gke.yaml`, `vend-kind.yaml`
   - [ ] KIND scripts exist and are executable: `go/k8s/kind-start.sh`, `go/k8s/kind-stop.sh`
   - [ ] Same images across all 4 modes
   - [ ] NODE_IP env var present in all services across all 4 modes

4. **Log services verification**:
   - [ ] `go/vend/log-vnet/` has main.go, Dockerfile, build.sh
   - [ ] `go/vend/log-agent/` has main.go, Dockerfile, build.sh
   - [ ] Both use `vend-security` base image (not `erp-security`)
   - [ ] Both in build-all-images.sh, deploy.sh, undeploy.sh

5. **Events verification**:
   - [ ] `grep "ActivateEvents" go/vend/main/main.go` returns a match
   - [ ] `grep "EventRecord" go/vend/ui/shared.go` returns a match

6. **Login verification**:
   - [ ] `go/vend/ui/web/login.html` exists
   - [ ] `go/vend/ui/web/index.html` redirects to `login.html`

7. **Security config verification**:
   - [ ] `go/secure/plugin/vend/vend.json` exists with valid JSON structure

8. **Callback verification**:
   - [ ] All 4 analytics services have `*ServiceCallback.go` files
   - [ ] Each callback auto-generates the primary key on POST

9. **Mobile parity verification**:
   - [ ] Desktop: all sections load data and render tables
   - [ ] Mobile: dashboard, alarms, nayax have mobile module directories
   - [ ] Mobile: nav config includes entries for new modules
   - [ ] Mobile: all new module scripts included in `m/app.html`

---

## 8. Out of Scope

The following are NOT addressed in this plan:
- Building `saichler/vend-security` and `saichler/vend-postgres` base Docker images (user-managed)
- Re-vendoring dependencies after adding `l8events` import (user-managed)
- Mock data generators for any new entities
- Consolidating the separate local-mode YAML files (`vend.yaml`, `vend-collector.yaml`, etc.) into a single `vend-local.yaml` (the current split-file approach works but differs from probler's single-file convention)
