# Plan: Collect Inventory and Telemetry from Vending Cloud Platforms

## Objective
Add a new collection path to `l8vendingmachine` that ingests inventory and telemetry
from third-party vending **cloud platforms** (Nayax, Cantaloupe Seed, 365 Retail
Markets, Televend) instead of (or in addition to) polling each machine directly.

The current pipeline assumes one `L8PTarget` per machine, each exposing
`/api/v1/...` REST endpoints (matches `l8opensim`). Cloud platforms invert that:
**one tenant API, many machines, paginated/rate-limited, OAuth2-secured**, with
optional webhook push for real-time telemetry events.

This plan follows the **probler / l8collector / l8parser** family per
`canonical-project-selection.md` (observation/collection project, not ERP).

---

## 1. Architecture

### 1.1 Pipeline shape (new vs existing)

```
EXISTING (per-machine polling, l8opensim-style):
  L8PTarget(machine A) ──REST──▶ vend-collector ──▶ vend-parser ──▶ VendMachine cache (id=A)
  L8PTarget(machine B) ──REST──▶ vend-collector ──▶ vend-parser ──▶ VendMachine cache (id=B)

NEW (cloud-platform collection):
  L8PTarget(tenant X) ──REST(paginated)──▶ vend-cloud-collector
       │
       └─▶ batched JSON for N machines
              │
              ▼
       vend-cloud-parser  ── fan-out per machine_id ──▶ VendMachine cache (id=A,B,C,...)

  Webhook ingress (push) ──▶ vend-cloud-webhook ──▶ vend-cloud-parser ──▶ VendMachine cache
```

### 1.2 Reference: probler GPU stack

Concrete reference for every layer below — read these files alongside this
plan:

| Layer | Probler GPU file | What it does |
|-------|------------------|-------------|
| Target factory | `probler/go/prob/common/creates/CreateGPU.go` | Builds an `L8PTarget{LinksId="GPU", InventoryType=GPUS}` with multi-protocol `L8PHost` (SSH + SNMP + REST), each with its own `CredId` |
| Target registration | `probler/go/prob/common/commands/addGPU.go` | POSTs the target to the existing `targets.ServiceName` endpoint — no bespoke API |
| Routing | `probler/go/prob/common/Links.go` | Returns `(GCache, GPersist, GPars)` services for `LinksId="GPU"` at area 2; collector area is shared across all LinksIds (NetDev, K8s, GPU all use `Coll`/0) |
| Collector binary | `probler/go/prob/collector/main.go` | One binary; `service.Activate(NetworkDevice_Links_ID, nic)` covers all LinksIds because they share the collector area |
| Parser binary | `probler/go/prob/parser/main.go` | One binary; `service.Activate(GPU_Links_ID, &GpuDevice{}, false, nic, "Id")` per supported model — separate parser areas per LinksId |
| Cache (inventory) binary | `probler/go/prob/inv_gpu/main.go` | Separate binary that calls `inventory.Activate(GPU_Links_ID, &GpuDevice{}, &GpuDeviceList{}, nic, "Id")` from `l8inventory` and registers an `Online` metadata function |
| Polling profile | `l8parser/go/parser/boot/nvidia.go` + `nvidia_ssh_rest.go` | One `L8Pollaris` named `"nvidia-gpu"` with ~15 polls across SSH (`nvidia-smi -q -d UTILIZATION`), SNMP (OIDs), and REST (`GET::/api/v1/gpu/devices::`) |
| Per-card fan-out rule | `l8parser/go/parser/rules/RestGpuParse.go` | Iterates a JSON array; uses `pci_bus_id` (configurable via `key_field`) as the map key; writes per-card properties via path `gpudevice.gpus<{24}{key}>.{property}` |
| Per-card SSH fan-out | `l8parser/go/parser/rules/SshNvidiaSmiParse.go` | Same fan-out shape, applied to text output from `nvidia-smi` |
| Time-series fields | `gpu.proto` | Heavy use of `repeated l8api.L8TimeSeriesPoint` for metrics like `vram_used_mib`, `gpu_utilization_percent`, `temperature_celsius`, `power_draw_watts` |

**Key takeaways for our plan:**
1. The cache layer is a **separate binary** (`inv_gpu`-style) that calls
   `l8inventory.Activate(...)`. The current vending project has
   `vend.yaml` filling this role for direct collection. Per
   `plan-duplication-audit.md` we **try-then-fork**: first attempt to add
   a second `inventory.Activate(VendCloud_Links_ID, ...)` to the
   existing cache binary; only fork into a sibling binary if the
   framework rejects same-model registration at two areas.
2. The probler **collector binary handles multiple LinksIds in one
   process** because all use collector area `Coll`/0. We do the same —
   **no new collector binary**. The existing `vend/collector/main.go`
   gets a second `service.Activate(VendCloud_Links_ID, nic)` call.
3. The probler **parser uses a separate area per LinksId** (NetDev=0,
   K8s=1, GPU=2), but for *different models* per LinksId. Our case has
   the *same model* (`VendMachine`) for both LinksIds, so we again
   **try-then-fork**: first attempt a second
   `service.Activate(VendCloud_Links_ID, &vend.VendMachine{}, ...)` in
   the existing parser binary at a different area. Only fork if the
   framework rejects.
4. **`RestGpuParse` does fan-out into map<>/repeated *children of one
   parent*** (e.g., `gpudevice.gpus[pci_bus_id]`). Per-machine cloud
   fan-out is different — each array element should become a *separate
   top-level instance* keyed by `MachineId`. So `CloudListFanOut` is a
   genuinely new rule, not a rename of `RestGpuParse`. However,
   `RestGpuParse` is reusable as-is for **slot-level** updates inside a
   single `VendMachine` if we change `VendMachine.inventory` from
   `repeated VendMachineSlot` to `map<string, VendMachineSlot>` keyed by
   `slot_id` — see §5.4.
5. Time-series shape (`repeated l8api.L8TimeSeriesPoint`) is the
   canonical way to model trending metrics (temperature, cash level,
   vend rate). Worth adopting in `vend-temperature.proto` and
   `vend-payment.proto` for consistency with probler.

### 1.3b Framework constraints (verified against code)

Four constraints discovered by auditing the actual `l8collector`,
`l8parser`, and `l8pollaris` source. Each requires a design adjustment.

**C1. Parser rule registration is private.**
`_Parser.rules` is a private map in `l8parser/go/parser/service/Parser.go`
(line 38). All 16 rules are hardcoded in `newParser()` (lines 45-82).
There is no public `RegisterRule()` API. Downstream projects cannot add
rules at runtime.
→ **Decision**: add new vending cloud rules (`CloudListFanOut`,
`CloudTransactionAppend`, `CloudAlertNormalize`) directly to
`l8parser/go/parser/rules/` and register them in `newParser()`. This
follows the existing pattern: `RestGpuParse`, `SshNvidiaSmiParse`,
`SnmpGpuTable` all live inside `l8parser` even though only probler uses
them. Rules are only invoked when a poll's `L8PAttribute` references them
by name, so unused rules have zero runtime cost. A public
`RegisterRule()` method would be cleaner long-term but is out of scope
for this plan.

**C2. `CJob.Result` is singular.**
`CJob.Result` is a single `[]byte` field
(`l8pollaris/go/types/l8tpollaris/jobs.pb.go`, line ~50).
`RestCollector.Exec()` writes exactly one result. The plan originally
proposed "emitting one `CJob.Result` per page" for pagination — that is
impossible without a proto change.
→ **Decision**: the `Paginate` enhancement to `RestCollector.Exec()`
concatenates all pages into **one combined JSON array** stored in a
single `CJob.Result`. The parser receives one result containing the full
dataset. This means the cloud parser must handle arrays of up to
~10 000 elements in memory. For the cloud-platform use case (tens to low
thousands of machines per vendor tenant), this is acceptable. If a future
vendor returns millions of rows, the design can be revisited with a
`CJob.Results repeated bytes` proto change.

**C3. Rule chaining does not support fan-out.**
The parser executes rules sequentially on a shared workspace with a
**single object instance** (`any`). A rule cannot invoke downstream
rules per element. `RestGpuParse` writes per-card data into the same
parent via indexed property paths (`gpudevice.gpus<{24}{key}>.prop`),
but there is no mechanism to create N separate top-level instances from
within a rule chain.
→ **Decision**: `CloudListFanOut` is the **final rule** in the
attribute chain (not a middle rule). It receives the full JSON array
from the workspace, iterates elements, and for each element: resolves
the `MachineId`, looks up or creates the top-level `VendMachine`
instance in the cache, and writes properties via property-path sets
(same mechanism `RestGpuParse` uses, but targeting top-level instances
keyed by `MachineId` instead of map children keyed by `pci_bus_id`).
No downstream rules are chained after it.

**C4. `L8PHost.Groups` is not accessible to `RestCollector`.**
`RestCollector.Init()` receives `L8PHostProtocol`, not `L8PHost`. The
`Groups` map is one level up on `L8PHost` and is never forwarded.
No code in `l8collector` or `l8parser` reads `L8PHost.Groups`.
→ **Decision**: pagination mode, page size, rate-limit RPS, and cursor
field names are stored as **`L8PRule` parameters** on each poll's
`L8PAttribute` (the `Paginate` and `RateLimit` rules in §4.2/§4.3
already describe this). The collector accesses poll rules via
`pollaris.Poll(jobName)` which it already calls. Vendor identity and
operator scope remain in `L8PHost.Groups` as metadata for UI display
only — they are never read by the collector or parser at runtime.

### 1.4 Pollaris alignment

The plan stays inside the Pollaris model. Specifically:

- **A "tenant" is an `L8PTarget`**, not a new Prime Object. Vendor and
  operator scope live in `L8PHost.Groups` (`map<string,string>`) as
  **UI-only metadata**. Collector/parser behavioral params (pagination
  mode, rate-limit) live in **`L8PRule` parameters** on each poll (per
  constraint C4). Auth lives in the existing `L8PHostProtocol.Ainfo`
  (`AuthInfo`) and `CredId`. There is no new Prime Object proto.
- **Routing uses `LinksId`**. We add `VendCloud_Links_ID = "VendCloud"` and
  extend `vend/common/Links.go` so `Collector("VendCloud")`,
  `Parser("VendCloud")`, `Cache("VendCloud")`, `Persist("VendCloud")`,
  `Model("VendCloud")` return the right (service, area). Following the
  probler GPU pattern, parser uses a **separate area** per LinksId
  (`VPars` for direct, `VCPars` for cloud). Cache + Persist also use
  separate areas (`VCCache`, `VCPersist`) — same convention as probler
  `(GCache, GPersist, GPars)` for GPU. **All ServiceNames are ≤10 chars
  per `maintainability.md`.** **Per
  `plan-duplication-audit.md` try-then-fork**: separate areas can be
  served by the same binary (with two `Activate` calls) until the
  framework forces a fork.
- **Polling profiles are `L8Pollaris` objects** registered via the existing
  Pollaris service, one per vendor (`NayaxCloud`, `CantaloupeCloud`,
  `M365Cloud`, `TelevendCloud`). Each is built in
  `vend/cloud/boot/{vendor}_pollaris.go` exactly like
  `vend/parser/boot/vend_pollaris.go` does today and exactly like
  `nvidia.go` + `nvidia_ssh_rest.go` do for the probler GPU stack.
- **Authentication uses `AuthInfo`**. OAuth2 client-credentials maps onto
  the existing `AuthPath` + `AuthBody` + `AuthResp` + `AuthToken` fields.
  API-key vendors set `IsApiKey = true` with `ApiUser` / `ApiKey`. Secrets
  are referenced via `CredId`, not embedded.

### 1.5 Components added

Per `plan-duplication-audit.md` Second Instance Rule, this plan **reuses
existing binaries** wherever the framework allows it. New binaries are
added only when behavior genuinely differs.

| Component | Path | Probler-GPU analog | Status | Purpose |
|-----------|------|--------------------|--------|---------|
| Existing collector binary | `vend/collector/main.go` (already exists) | `prob/collector/main.go` (one binary, multiple LinksIds) | **REUSE** | Add a second `service.Activate(VendCloud_Links_ID, nic)` call. Pollaris already routes both LinksIds to the same collector area (`Links.Collector` returns constant `(VColl, 0)`), so one binary covers both. No new collector binary. |
| Existing parser binary | `vend/parser/main.go` (already exists) | `prob/parser/main.go` (one binary, per-LinksId activation) | **REUSE if framework allows** | Try adding `service.Activate(VendCloud_Links_ID, &vend.VendMachine{}, false, nic, "MachineId")` next to the existing direct activation. If the framework rejects same-model registration at two areas, fork into `vend/cloud/parser/main.go` — verified in Phase 2.5, decided in Phase 4. |
| Existing inventory cache binary | (current `vend.yaml` deploys it) | `prob/inv_gpu/main.go` | **REUSE if framework allows** | Same try-then-fork pattern: add a second `inventory.Activate(VendCloud_Links_ID, ...)` to the existing cache binary. Verified in Phase 2.5, decided in Phase 4b. |
| Webhook receiver | `vend/cloud/webhook/` | (none — push is outside Pollaris) | **NEW** | Genuinely different behavior — HTTPS server, push not pull. Cannot be merged into the collector or parser. |
| Per-vendor Pollaris boot | `vend/cloud/boot/{vendor}_pollaris.go` | `nvidia.go` + `nvidia_ssh_rest.go` | **NEW (config-only)** | Built on top of `vend/cloud/boot/helpers.go` + `cadences.go` extracted in Phase 2.5 (see below). Each vendor file is pure data, no helpers redefined. |
| Shared boot helpers | `vend/cloud/boot/helpers.go` + `cadences.go` | `vend/parser/boot/helpers.go` + `cadences.go` | **NEW (Phase 2.5, before any vendor file)** | `createCloudPoll(...)`, `createCloudRestAttribute(...)`, `addParameter(...)`, cadence constants. Extracted **before** the first vendor profile per Second Instance Rule. |
| Shared target builder | `vend/common/targets/BuildVendTarget(opts)` | `prob/common/creates/CreateGPU.go` | **NEW (Phase 1, replaces both factories)** | Single constructor that produces an `L8PTarget` for direct OR cloud, driven by an `opts` struct. Replaces both `vend/common/commands/createMachine.go` AND a hypothetical `vend/cloud/common/createTarget.go` per Second Instance Rule. |
| List fan-out parsing rule | `vend/cloud/rules/CloudListFanOut.go` | (no probler analog — see §5.4) | **NEW** | Splits an array response into separate top-level `VendMachine` instances keyed by `machine_id`. Distinct from `RestGpuParse`. |
| RestCollector enhancements | upstream `l8collector/.../rest/` | (would benefit probler too) | **NEW (upstream)** | Pagination loop driven by new `L8PRule` params; `429`/`503` backoff; OAuth2 token refresh on `401`. |

### 1.6 Reuse, do NOT duplicate

- `VendMachine`, `VendMachineSlot`, `VendTransaction`, `VendAlert` protos are
  unchanged. (Open question: switch `VendMachine.inventory` from `repeated`
  to `map<string, VendMachineSlot>` keyed by `slot_id` so we can reuse
  `RestGpuParse` directly for slot updates — see §5.4.)
- `VCache` / `VPersist` keep collecting direct-poll data; new
  `VCCache` / `VCPersist` collect cloud-poll data. The UI joins
  both sources at query time (or a downstream merger pushes both into a
  shared analytics store — out of scope here).
- `l8collector` is reused. The pagination + backoff + OAuth2-refresh
  enhancements are upstream changes to `RestCollector`, not a fork. They are
  driven by **new `L8PRule` parameter names** so existing direct-poll usage
  is unaffected (no new params = current behavior).
- `l8parser` is reused. Fan-out is a new `L8PRule` (`CloudListFanOut`)
  registered alongside the existing `RestJsonParse`. No parser binary fork
  beyond the standard `service.Activate(LinksId, ...)` entrypoint.

---

## 2. Tenant as a Pollaris target

There is **no new Prime Object proto**. A "tenant" is an `L8PTarget`
registered through the existing Pollaris targets service. The target is
shaped like this:

```go
target := &l8tpollaris.L8PTarget{
    TargetId:      "nayax-acme",                       // operator-chosen
    LinksId:       common.VendCloud_Links_ID,          // routes to cloud collector/parser
    InventoryType: l8tpollaris.L8PTargetType_Vending_Machine,
    Hosts: map[string]*l8tpollaris.L8PHost{
        "api": {
            HostId: "api",
            Configs: map[int32]*l8tpollaris.L8PHostProtocol{
                int32(l8tpollaris.L8PProtocol_L8PRESTAPI): {
                    Protocol:   l8tpollaris.L8PProtocol_L8PRESTAPI,
                    Addr:       "api.nayax.com",
                    Port:       443,
                    HttpPrefix: "",
                    CredId:     "nayax-acme-creds",     // resolved via security provider
                    Timeout:    30000,
                    Ainfo: &l8tpollaris.AuthInfo{
                        NeedAuth:      true,
                        AuthPath:      "/oauth/token",
                        AuthBody:      `{"grant_type":"client_credentials"}`,
                        AuthUserField: "client_id",
                        AuthPassField: "client_secret",
                        AuthResp:      "access_token",
                    },
                },
            },
            Polls: map[string]string{"NayaxCloud": ""},  // Pollaris profile name
            Groups: map[string]string{
                "vendor":         "nayax",
                "operator_scope": "OP-12345",
                "page_size":      "100",
                "rate_limit_rps": "5",
            },
        },
    },
}
```

Notes:
- `CredId` resolves through the **existing security provider** per
  `security-provisioning-channels.md`. No new secret-storage proto.
- API-key vendors (e.g., 365): set `Ainfo.IsApiKey = true`, `Ainfo.ApiUser`,
  `Ainfo.ApiKey`. Pollaris already supports this.
- `L8PHost.Polls` selects which `L8Pollaris` profile applies (one per vendor
  — `NayaxCloud`, `CantaloupeCloud`, `M365Cloud`, `TelevendCloud`).
- `L8PHost.Groups` carries **UI-only metadata** (`vendor`,
  `operator_scope`) for display in the Cloud Tenants table. Per
  constraint C4, the collector and parser do NOT read `Groups` at
  runtime — all behavioral params (`page_size`, `rate_limit_rps`,
  `cursor_field`, `pagination_mode`) live in **`L8PRule` parameters** on
  each poll's `L8PAttribute` (§4.2/§4.3).
- A small admin-only proto (`VendCloudTenantView`) MAY be added later if the
  UI needs a denormalized list view, but it is **not** a Prime Object and is
  not the source of truth — `L8PTarget` is. **If added**, it MUST follow:
  - `proto-list-convention.md`: `VendCloudTenantViewList { repeated
    VendCloudTenantView list = 1; l8api.L8MetaData metadata = 2; }`
  - `proto-enum-zero-value.md`: any enum has `*_UNSPECIFIED = 0` as its
    first value.
  - `protobuf-generation.md`: regenerate via `cd proto &&
    ./make-bindings.sh` after every edit.

---

## 3. What we collect (per vendor)

Mapping per vendor is documented as a **table only** — the actual mapping logic
lives in vendor adapter `boot/*_polls.go` files (config-only, no behavior).

| Domain | Internal field (`VendMachine.*`) | Nayax | Cantaloupe Seed | 365 Retail Markets | Televend |
|--------|----------------------------------|-------|-----------------|--------------------|----------|
| Machine identity | `machine_id`, `serial_number`, `model` | `/v1/devices` | `/v1/devices` | `/api/markets` | `/v3/machines` |
| Machine status / heartbeat | `status`, `last_heartbeat`, `connectivity` | `/v1/devices/status` | `/v1/devices/{id}/status` | `/api/devices/status` | `/v3/machines/state` |
| Inventory snapshot | `inventory[]` | `/v1/products/inventory` | `/v1/planograms` | `/api/products/inventory` | `/v3/planograms` |
| Vend transactions | (writes `VendTransaction`) | `/v1/transactions?since=` | `/v1/transactions?since=` | `/api/transactions` | `/v3/sales` |
| Cashbox / payment | `payment.*` | `/v1/cash/positions` | `/v1/cash` | `/api/cash` | `/v3/cash` |
| Temperature | `temperature.*` | `/v1/telemetry/temperature` | `/v1/telemetry/temperature` | n/a (micro-markets) | `/v3/telemetry/temp` |
| Alerts | (writes `VendAlert`) | `/v1/alerts` | `/v1/alerts` | `/api/alerts` | `/v3/events` |
| DEX audits | `vend-dex` | `/v1/dex/latest` | `/v1/dex` | n/a | `/v3/dex` |

Cells marked "n/a" mean the vendor does not expose that domain — the adapter
omits the corresponding poll rather than returning empty data.

---

## 4. Pagination, rate limiting, OAuth2 — upstream `RestCollector`

These three concerns are not vending-specific and are not cloud-specific.
Any REST-based collector (k8s API, GraphQL, vendor cloud APIs) hits the same
issues. They belong in `l8collector`'s `RestCollector`, not in a parallel
`vend/cloud/common/` library.

The plan upstream the changes as **opt-in** behavior driven by new `L8PRule`
parameter names so existing direct-poll usage is unaffected.

### 4.1 OAuth2 token caching (extend `RestCollector`)
- `RestCollector` already reads `Ainfo.AuthPath`, `AuthBody`, `AuthResp`,
  `AuthToken`. Today it is set up but the refresh lifecycle is incomplete.
  We extend `Connect()` / `Exec()` so:
  - A missing/expired `AuthToken` triggers a `POST AuthPath` with a body
    built from `AuthBody` + `AuthUserField` / `AuthPassField` substitution
    using the credential resolved via `CredId`.
  - The response is parsed and `AuthResp` is extracted; `AuthToken` is set
    on the in-memory `L8PHostProtocol` (not persisted).
  - On `401` from a regular call, invalidate `AuthToken` and retry **once**.
- API-key path (`IsApiKey == true`) is already supported and is a no-op
  here.

### 4.2 Pagination loop (extend `RestCollector`)
A poll opts in by adding an `L8PRule` named `Paginate` to its `L8PAttribute`,
with these `L8PParameter`s:

| Param | Values | Meaning |
|-------|--------|---------|
| `mode` | `cursor` \| `link_header` \| `page_number` | how the API exposes "next" |
| `cursor_field` | e.g. `next_cursor` | only for `mode=cursor`; JSON path on the response body |
| `cursor_param` | e.g. `cursor` | query-string param to send the cursor as |
| `page_param` | e.g. `page` | only for `mode=page_number` |
| `page_size_param` | e.g. `limit` | query-string param for page size |
| `page_size_value` | e.g. `100` | integer |
| `max_pages` | e.g. `1000` | safety cap |

`RestCollector.Exec()` loops while a "next" indicator is present,
**concatenating all pages into a single combined JSON array** stored in
one `CJob.Result` (per constraint C2 — `CJob.Result` is singular). The
parser receives one result containing the full dataset.

Memory bound: for the cloud-platform use case (tens to low thousands of
machines per vendor tenant), this is acceptable. If a future vendor
returns millions of rows, the design can be revisited with a
`CJob.Results repeated bytes` proto change. The `max_pages` safety cap
prevents unbounded fetching.

Polls without the `Paginate` rule keep the current single-shot behavior.

### 4.3 Rate-limit / `429`-`503` backoff (extend `RestCollector`)
A poll opts in via an `L8PRule` named `RateLimit`:

| Param | Meaning |
|-------|---------|
| `rps` | token-bucket rate per second (per host) |
| `burst` | token-bucket burst size |
| `backoff_max_ms` | cap on exponential backoff |

The collector keeps a per-`(TargetId, HostId)` token bucket and an
exponential-backoff state. On `429` / `503` it sleeps with jittered backoff
and retries; on persistent failure it records `job.Error` exactly as today
so the existing health machinery sees the failure.

### 4.4 Webhook ingress (separate binary, **not** Pollaris)
Pollaris is pull-only; webhooks are push. A small `vend/cloud/webhook/`
binary stands up an HTTPS server that:
- Resolves vendor + tenant from URL path: `POST /webhook/{vendor}/{tenant}`.
- Verifies the HMAC signature (header name per vendor) using a credential
  resolved via the security provider (the same credential store
  `RestCollector.CredId` uses).
- Dedupes on vendor `event_id` via a 24 h cache.
- Publishes the normalized event into the **parser service area** for
  `VendCloud_Links_ID` — same downstream as the polled path.

This is the only component that legitimately sits outside the Pollaris
collector path.

---

## 5. Parser fan-out

### 5.1 Service activation
In `vend/cloud/parser/main.go`:
```go
service.Activate(common.VendCloud_Links_ID,
    &vend.VendMachine{}, false, nic, "MachineId")
```
The cloud parser registers a **separate** parser service area for
`VendCloud_Links_ID` (so direct and cloud parsers can scale independently),
but its writes flow into the **shared `VCache` / `VPersist`** keyed by
`MachineId`. A machine collected via Nayax today and Cantaloupe tomorrow
keeps the same cache identity.

### 5.2 New parsing rules (added to `l8parser` core — see constraint C1)

Per constraint C1, the parser's rule map is private. New rules are added
directly to `l8parser/go/parser/rules/` and registered in `newParser()`
(upstream PR). This follows the existing pattern — `RestGpuParse`,
`SshNvidiaSmiParse`, and `SnmpGpuTable` all live inside `l8parser` even
though only probler uses them. Rules are invoked only when a poll's
`L8PAttribute` references them by name.

| Rule | Input | Output | Notes |
|------|-------|--------|-------|
| `CloudListFanOut` | combined JSON array (from paginated collector result) | one `VendMachine` cache update per element, keyed by configured id field | **Must be the FINAL rule** in the attribute chain (per constraint C3 — rules cannot fan out to downstream rules). Iterates the array, resolves each element's `MachineId`, looks up or creates the `VendMachine` instance in the cache, and writes properties via property-path sets (same mechanism as `RestGpuParse` but targeting top-level instances, not map children). |
| `CloudTransactionAppend` | JSON array of vends | append-only `VendTransaction` writes | **Final rule.** Per immutability rule for transactions. |
| `CloudAlertNormalize` | vendor alert JSON | mapped to internal `VendAlertCategory` / `VendAlertSeverity` enums | **Final rule.** |

**Removed from original plan:** `CloudInventorySlotMerge` is not needed
as a separate rule. `CloudListFanOut` handles slot merging via nested
property-path writes (e.g., `vendmachine.inventory<{24}{slotId}>.currentQuantity`)
when the §5.4 map<> migration is adopted. Without the migration,
`CloudListFanOut` writes slot arrays directly.

### 5.3 How a poll uses these rules

Per constraint C3, each poll's `L8PAttribute` has **one parser-side
rule** that is the terminal rule — no chaining to downstream rules.
Collector-side rules (`Paginate`, `RateLimit`) are consumed by
`RestCollector` before the result reaches the parser.

```go
// Machine identity/status poll — CloudListFanOut is the ONLY parser rule
attr.Rules = []*l8tpollaris.L8PRule{
    {Name: "Paginate",        Params: paginationParams},   // collector-side
    {Name: "RateLimit",       Params: rateLimitParams},    // collector-side
    {Name: "CloudListFanOut", Params: map[string]*l8tpollaris.L8PParameter{
        "id_field": {Value: "machine_id"},
        "mapping":  {Value: "serial_number:serialNumber,model:model,..."},
    }},                                                     // parser-side, FINAL
}

// Transaction poll — CloudTransactionAppend is the ONLY parser rule
attr.Rules = []*l8tpollaris.L8PRule{
    {Name: "Paginate",                 Params: paginationParams},
    {Name: "RateLimit",                Params: rateLimitParams},
    {Name: "CloudTransactionAppend",   Params: txnParams},   // FINAL
}
```

Collector-side rules are identified by name in `RestCollector.Exec()`.
The parser ignores them (they have no registered `ParsingRule`
implementation). Parser-side rules are identified by name in `_Parser`
and the collector ignores them (not in its known-rules set).

### 5.4 Reuse `RestGpuParse` for slot-level updates (decision needed)

`RestGpuParse` already does what we need for the **slot dimension** inside
one machine: iterate a JSON array, key each element by an ID field, and
write properties into a map<>/repeated child of the parent. Today
`VendMachine.inventory` is `repeated VendMachineSlot`, which works with
the property-path indexing but is less ergonomic than a map.

**Recommended:** change `vend-fleet.proto`:

```protobuf
message VendMachine {
    ...
    map<string, VendMachineSlot> inventory = 20;  // keyed by slot_id
    ...
}
```

Benefits:
- Reuse `RestGpuParse` as-is for slot updates with `key_field="slot_id"` —
  no new rule.
- Slot lookups by `slot_id` become O(1).
- Aligns exactly with probler's `GpuDevice.gpus map<string, Gpu>` shape.

Cost:
- Migration of any existing code that iterates `inventory` as a slice.
- Simulator and UI need to handle the map representation.

This change is independent of cloud collection but is a strong
architectural alignment with probler's reference pattern. If we do **not**
make this change, we add `RestSlotFanOut` (a near-copy of `RestGpuParse`
that handles `repeated`) to the parser rules.

**Per `protobuf-generation.md`**, the proto change MUST be followed by:
```bash
cd proto && ./make-bindings.sh
```
Before running, confirm `make-bindings.sh` uses `-i` (not `-it`) on its
`docker run` commands — the `-t` flag fails in non-interactive
environments. After regeneration, run `go build ./...` to verify the type
change compiled across the project.

---

## 6. Webhook ingress (push path)

### 6.1 Daemon (`vend/cloud/webhook/main.go`)
- Stands up an HTTPS server inside the cluster (DaemonSet, hostNetwork or
  Service+Ingress depending on cluster).
- Routes: `POST /webhook/{vendor}/{tenant_id}` → verify signature → normalize
  payload → publish onto the same internal channel the cloud parser consumes.
- Ack with `2xx` only after the message is durably enqueued.

### 6.2 Idempotency
- Vendors retry on non-2xx. Every event carries a vendor `event_id`; the
  webhook daemon caches `(tenant_id, event_id)` for 24h and dedupes before
  publishing.

### 6.3 Failover to polling
- Webhooks are best-effort. The polling path remains the system of record for
  inventory/state. Webhook events are merged in opportunistically for low
  latency.

---

## 7. Configuration & deployment

### 7.1 Adding a tenant
- Operator creates an `L8PTarget` with `LinksId = "VendCloud"` via the
  existing Pollaris targets service (POST to `/0/Device` per current
  `addMachine.go` pattern; no new endpoint).
- Credentials referenced by `CredId` are provisioned through the security
  provider — **never** a bespoke users/secrets endpoint per
  `security-provisioning-channels.md`.
- The Pollaris targets service notifies the cloud collector via the standard
  target-callback path (same as direct-poll registration today).
- UI for tenant creation is a thin wrapper that writes the `L8PTarget`
  shape — see Phase 7.

### 7.2 New k8s manifests
Following `k8s-yaml-required-entries.md` and `deployment-artifacts.md`:

| File | Kind | Probler-GPU analog | Notes |
|------|------|--------------------|-------|
| `k8s/vend-cloud-collector.yaml` | StatefulSet | `probler/k8s/collector.yaml` | One pod per N tenants, sharded by hash(target_id) |
| `k8s/vend-cloud-parser.yaml`    | DaemonSet   | `probler/k8s/parser.yaml`    | Stateless |
| `k8s/vend-cloud-inventory.yaml` | StatefulSet | `probler/k8s/gpu.yaml`       | In-memory cache via `l8inventory.Activate(...)` |
| `k8s/vend-cloud-webhook.yaml`   | DaemonSet, hostNetwork | (none — push not in probler) | Public ingress |

**Required entries in EVERY new YAML (per `k8s-yaml-required-entries.md`):**
- `Namespace.metadata.labels.name` — namespace label (do not omit)
- Resource `metadata.labels.app` — app label (do not omit)
- Container `env.NODE_IP` from `fieldRef: status.hostIP` (do not omit)
- Volume named `hdata` (NOT `data`) mounted at `/data` from
  `hostPath: /data, type: DirectoryOrCreate`
- `imagePullPolicy: Always`

Verification after writing each YAML:
```bash
grep -A2 "kind: Namespace" <file> | grep "labels:"
grep -A4 "kind: StatefulSet\|kind: DaemonSet" <file> | grep "labels:"
grep "NODE_IP" <file>
grep "name: hdata" <file>
```

### 7.3 `build.sh` and `Dockerfile` per binary

Each new binary follows the canonical multi-stage pattern from
`probler/go/prob/inv_gpu/Dockerfile` and `probler/go/prob/inv_gpu/build.sh`:

```dockerfile
FROM saichler/builder:latest AS build
COPY main.go /home/src/github.com/saichler/build/
RUN go mod init && GOPROXY=direct GOPRIVATE=github.com go mod tidy
RUN go build -o vend-cloud-<name>
FROM saichler/erp-security:latest AS final
COPY --from=build /home/src/github.com/saichler/build/vend-cloud-<name> /home/run/
ENTRYPOINT ["/home/run/vend-cloud-<name>"]
```

```bash
#!/usr/bin/env bash
set -e
docker build --no-cache --platform=linux/amd64 -t saichler/vend-cloud-<name>:latest .
docker push saichler/vend-cloud-<name>:latest
```

### 7.4 `build-all-images.sh`
Add four lines:
```
cd ../vend/cloud/collector && ./build.sh
cd ../vend/cloud/parser    && ./build.sh
cd ../vend/cloud/inventory && ./build.sh
cd ../vend/cloud/webhook   && ./build.sh
```

### 7.5 `run-local.sh`
Per `run-local-script.md`, this is an observation/collection project so the
canonical reference is `../probler/go/run-local.sh` (NOT l8erp). **Copy and
adapt** rather than write from scratch:
1. Copy `../probler/go/run-local.sh` template structure.
2. Replace probler binary names (`box_demo`, `gpu_demo`) with vending names
   (`vend_cloud_collector_demo`, `vend_cloud_parser_demo`,
   `vend_cloud_inventory_demo`, `vend_cloud_webhook_demo`).
3. Point the cloud collector at the simulator's base URL via env var. The
   simulator itself is started by its own project's `run-local.sh`; this
   project consumes the URL.
4. Preserve the `kill_demo.sh` cleanup generation flow.
5. Per `cleanup-test-binaries.md`, build verification uses
   `go build ./...` (not `go build ./path/to/main/`) to avoid leftover
   binaries.

---

## 8. Multi-vendor duplication audit

Per `plan-duplication-audit.md`, before writing four vendor configs that look
95% identical:

| Concern | Where it lives | Per-vendor work |
|---------|---------------|-----------------|
| OAuth2 / API-key wiring | `RestCollector` + existing `AuthInfo` | set `Ainfo.AuthPath` / `IsApiKey` fields per vendor |
| Pagination | `RestCollector` + new `Paginate` `L8PRule` | choose `mode` + field names |
| Rate-limit / backoff | `RestCollector` + new `RateLimit` `L8PRule` | choose `rps` |
| List fan-out | `CloudListFanOut` parser rule | set ID-field name |
| JSON → VendMachine mapping | existing `RestJsonParse` | per-poll `mapping` string |
| Endpoint URLs + cadences | `vend/cloud/boot/{vendor}_pollaris.go` | the only per-vendor file |

Per-vendor file = ~80 lines of pure `L8Pollaris` config (poll names,
endpoints, cadence selectors, rule param maps). All behavior lives in:
- `l8collector` (extended `RestCollector`) — upstream
- `vend/cloud/rules/` — parser rules
- `vend/cloud/boot/helpers.go` + `cadences.go` — shared boot helpers
  (extracted in Phase 2.5 **before** the first vendor file)

If a vendor file exceeds ~120 lines or contains an `if vendor == X` branch
in Go code anywhere, it is a code-smell — extend the shared helper or add a
new `L8PRule` parameter instead.

### Binary-level duplication

`Dockerfile` and `build.sh` per binary are near-identical (differing
only in binary name). This is a known trade-off — k8s deployment
conventions require per-image files. The probler project has the same
duplication (`prob/collector/Dockerfile`, `prob/parser/Dockerfile`,
`prob/inv_box/Dockerfile`, `prob/inv_gpu/Dockerfile` are structurally
identical). If this plan's try-then-fork outcome results in more than
one forked binary, and a future plan adds a 5th, a templated
`build_image.sh <name>` helper should be introduced. For now (at most
one fork + webhook) it remains acceptable.

---

## 9. Traceability matrix

| # | Concern | Where addressed |
|---|---------|----------------|
| 1 | `VendCloud_Links_ID` + `Links.go` routing | §1.2, §2 / Phase 1 |
| 2 | `L8PTarget` shape (no new Prime Object) | §2 / Phase 1 |
| 3 | Shared `BuildVendTarget(opts)` constructor (Second Instance Rule) | §1.5 / Phase 1 |
| 4 | OAuth2 lifecycle on existing `AuthInfo` fields | §4.1 / Phase 2 PR 1 (upstream `RestCollector`) |
| 5 | `Paginate` (3 modes, concatenated result per C2) | §4.2 / Phase 2 PR 1 (upstream `RestCollector`) |
| 6 | `RateLimit` + 429/503 backoff + 401 retry | §4.3 / Phase 2 PR 1 (upstream `RestCollector`) |
| 6b | New parser rules added to `l8parser` core (per C1) | §5.2 / Phase 2 PR 2 (upstream `l8parser`) |
| 7 | Framework double-activation verification | Phase 2.5(A) |
| 8 | Shared boot `helpers.go` + `cadences.go` (Second Instance Rule) | §1.5 / Phase 2.5(B) |
| 9 | Cloud collector activation (reuse existing binary) | §1.5 / Phase 3 |
| 10 | Per-vendor `L8Pollaris` boot file (config-only, first: Nayax) | §3, §8 / Phase 3 |
| 11 | `CloudListFanOut` as FINAL rule (per C3, no chaining) | §5 / Phase 2 PR 2 + Phase 4 |
| 12 | `CloudTransactionAppend` as FINAL rule (per C3) | §5.2 / Phase 2 PR 2 + Phase 4 |
| 13 | `CloudAlertNormalize` as FINAL rule (per C3) | §5.2 / Phase 2 PR 2 + Phase 4 |
| 14 | Parser activation (reuse or fork per Phase 2.5(A)) | §1.5 / Phase 4 |
| 15 | Inventory cache activation (reuse or fork per Phase 2.5(A)) | §1.5 / Phase 4b |
| 16 | Webhook receiver (genuinely new binary) + k8s manifest | §4.4, §6, §7.2 / Phase 5 |
| 17 | Three remaining vendor profiles (Cantaloupe, 365, Televend) | §8 / Phase 6 |
| 18 | `build-all-images.sh` + `run-local.sh` updates | §7.4, §7.5 / Phase 6 |
| 19 | (Optional) inventory map<> migration for slot reuse | §5.4 / decision before Phase 4 |
| 20 | UI surface for cloud target CRUD (desktop + mobile parity) | Phase 7 |
| 21 | End-to-end verification on all four vendors | Phase 8 |

Every gap above lands in exactly one phase. No orphans.

Synthetic data is **out of scope of this plan**. A separate simulator project
(see `plans/vending-machine-simulation.md` and follow-on work) is responsible
for standing up endpoints that mimic the four vendor cloud APIs. This plan
consumes whatever the simulator exposes; it does not generate fixtures, mocks,
or seed data of its own.

---

## 9b. PRD compliance checklist

Per `prd-compliance.md` Rule 1, this plan must comply with all rules at
`../l8book/rules` and the l8ui rules at `../l8ui/rules`. Each item below
ties a rule to the section/phase that satisfies it.

### Project structure & architecture
- [x] Observation/collection family — uses probler / l8collector /
      l8parser per `canonical-project-selection.md` (§1.1, §1.2)
- [x] No new Prime Object — `L8PTarget` reused per
      `prime-object-references.md` (§2)

### Service design
- [x] All ServiceNames ≤10 chars per `maintainability.md`: `VCPars`(6),
      `VCCache`(7), `VCPersist`(9). Collector area shared with `VColl`.
- [x] Same ServiceArea per LinksId (cloud area distinct from direct
      area, both internally consistent) per `maintainability.md`
- [x] No `ServiceCallback` auto-ID needed — `L8PTarget` is owned by
      Pollaris, not by this plan
- [x] No new types to register in UI `main.go` — `L8PTarget` already
      registered by Pollaris

### Protobuf design (only if §5.4 inventory migration is adopted, or if
`VendCloudTenantView` admin proto is added)
- [ ] Enum zero-values are `*_UNSPECIFIED = 0` per
      `proto-enum-zero-value.md`
- [ ] List types use `repeated X list = 1; l8api.L8MetaData metadata =
      2;` per `proto-list-convention.md`
- [ ] No direct struct references between Prime Objects per
      `prime-object-references.md`
- [ ] Bindings regenerated via `cd proto && ./make-bindings.sh` per
      `protobuf-generation.md` (Phase 2 prerequisite if §5.4 adopted)

### UI design (Phase 7)
- [x] Desktop module follows `adding-module-desktop.md`
- [x] Mobile module follows `adding-module-mobile.md`
- [x] Desktop/mobile parity per `mobile-rules.md`
- [x] Reference registry entries on both platforms per
      `reference-registry-completeness.md`
- [x] All factory methods used exist per `enum.md` /
      `enum-renderer-column-cascade.md`
- [x] Field names verified against `.pb.go` per
      `js-protobuf-field-names.md` (especially nav config `idField`)
- [x] `sectionSelector == defaultModule` per
      `module-init-section-selector.md`
- [x] Theme tokens (`--layer8d-*`) only per `l8ui-theme-compliance.md`
- [x] l8ui added as submodule, not copied, per
      `l8ui-copy-to-new-project.md`

### Mock data
- [x] Out of scope — simulator owns synthetic data (§9 note)

### Deployment
- [x] Try-then-fork: reuse existing collector/parser/inventory binaries
      unless framework forces a fork (Phase 2.5(A) per
      `plan-duplication-audit.md` Second Instance Rule)
- [x] Any forked binary + webhook get `build.sh`, `Dockerfile`, k8s
      YAML per `deployment-artifacts.md` (§7.2, §7.3)
- [x] All k8s YAMLs include namespace label, app label, `NODE_IP`,
      `hdata` volume per `k8s-yaml-required-entries.md` (§7.2)
- [x] `build-all-images.sh` updated (§7.4)
- [x] `run-local.sh` adapted from `../probler/go/run-local.sh` per
      `run-local-script.md` (§7.5)
- [x] Dockerfile/build.sh duplication acknowledged as k8s convention
      trade-off (§8)

### Configuration & security
- [x] `login.json` adaptation NOT NEEDED (no UI relocation)
- [x] No `Layer8DModuleFilter.load()` failure trap — vending project
      already handles its own ModConfig (`modconfig-failure-no-logout.md`)
- [x] Secrets only via `ISecurityProvider` per
      `security-provider-interface.md` (§2, §4.4, §7.1)
- [x] No bespoke users/credentials endpoints per
      `security-provisioning-channels.md` (§7.1)

### Code quality
- [x] No Go generics per `no-go-generics.md`
- [x] No file added by this plan exceeds 500 lines per
      `maintainability.md` (per-vendor boot files capped at ~120 lines
      per §8; rule files split if needed; `RestCollector` enhancements
      are localized inside its existing file or split)
- [x] Build verification uses `go build ./...` not `go build
      ./path/main/` per `cleanup-test-binaries.md`
- [x] Tests in `go/tests/cloud/` exercising the system API per
      `test-location-and-approach.md` (Phase 8)
- [x] Demo directory `go/demo/` is not touched — it is auto-generated
      by `run-local.sh` per `demo-directory-sync.md`
- [x] Third-party deps consumed from `go/vendor/` per
      `vendor-third-party-code.md`

### Plan discipline
- [x] Written to `./plans/` and not directly approved per
      `plan-approval-workflow.md`
- [x] Traceability matrix present per
      `plan-traceability-and-verification.md` (§9)
- [x] Final verification phase present (Phase 8) per
      `plan-traceability-and-verification.md`
- [x] Duplication audit present per `plan-duplication-audit.md` (§8)

---

## 10. Phases

### Phase 1 — Routing and shared target builder
- Add `VendCloud_Links_ID = "VendCloud"` and dedicated cloud
  service-name/area constants to `vend/common/defaults.go`. **All
  ServiceName values ≤10 chars** (`VCPars`, `VCCache`, `VCPersist`).
- Extend `vend/common/Links.go` so `Collector / Parser / Cache / Persist /
  Model` each handle the new `LinksId` via a `switch` (mirrors probler's
  `Links.go` GPU branch). `Collector(...)` keeps returning the existing
  `(VColl, 0)` for both LinksIds.
- **Per `plan-duplication-audit.md` Second Instance Rule**: replace
  `vend/common/commands/createMachine.go` with a shared
  `vend/common/targets/BuildVendTarget(opts BuildOpts) *L8PTarget`
  constructor that handles both direct and cloud target shapes (driven
  by `opts.Vendor`, `opts.Auth`, `opts.LinksId`). Existing call sites
  for `createMachine.go` are migrated to the new constructor; old file
  deleted in the same commit.
- Document the `L8PTarget` shape for cloud tenants (§2) — example shows
  `BuildVendTarget(BuildOpts{LinksId: "VendCloud", Vendor: "nayax",
  ...})`.

### Phase 2 — Upstream enhancements (two companion PRs)

**PR 1: `l8collector` — `RestCollector` enhancements**
Done in `l8collector/go/collector/protocols/rest/`. Adds:
- OAuth2 token-refresh lifecycle on the existing `AuthInfo` fields
  (`AuthPath`, `AuthBody`, `AuthUserField`, `AuthPassField`, `AuthResp`,
  `AuthToken`); 401 → invalidate + retry once. Currently these fields
  exist but are never read in `Exec()` (constraint C1 verification).
- `Paginate` behavior: `RestCollector.Exec()` reads `Paginate` rule
  params from the poll's `L8PAttribute` (accessed via `pollaris.Poll()`
  which the collector already calls). Loops pages and **concatenates
  all pages into one combined JSON array** in the single `CJob.Result`
  (per constraint C2 — `CJob.Result` is singular, not repeated).
- `RateLimit` behavior: reads `RateLimit` rule params the same way.
  Per-`(TargetId, HostId)` token bucket; jittered exponential backoff
  on `429` / `503`.
- Unit tests against the simulator endpoints.
- Existing direct-poll usage is unaffected (no Paginate/RateLimit params
  = current single-shot behavior).

**PR 2: `l8parser` — new vending cloud rules**
Done in `l8parser/go/parser/rules/`. Adds (per constraint C1 — parser
rule map is private, so rules must live in l8parser core):
- `CloudListFanOut` in `l8parser/go/parser/rules/CloudListFanOut.go`
- `CloudTransactionAppend` in `CloudTransactionAppend.go`
- `CloudAlertNormalize` in `CloudAlertNormalize.go`
- Registration calls added to `newParser()` in `Parser.go`
- Rules follow the `ParsingRule` interface (same as `RestGpuParse`):
  `Name() string`, `ParamNames() []string`,
  `Parse(resources, workSpace, params, any, pollWhat) error`

### Phase 2.5 — Framework verification + shared boot helpers

Two preconditions must be met before any vendor profile is written.

**A. Verify same-model double-activation in the framework**
Write a tiny test in `go/tests/cloud/double_activate_test.go` that:
1. Activates the existing parser binary against `VendMachine_Links_ID`
   (area `VPars`/0).
2. Activates the *same* binary against `VendCloud_Links_ID` (area
   `VCPars`/N).
3. Posts a `VendMachine` to the cloud area and verifies it lands in the
   cache without colliding with a direct-area write of the same machine.

If the framework supports it: **no parser fork** in Phase 4; **no
inventory cache fork** in Phase 4b. The existing `vend/parser/main.go`
and the existing inventory cache binary each get a second `Activate`.

If it rejects: fork into `vend/cloud/parser/main.go` and
`vend/cloud/inventory/main.go` with their own `build.sh` / `Dockerfile`
/ k8s YAML. Document the rejection reason in this section so the next
reader knows why two binaries exist.

**B. Extract shared boot helpers** (per Second Instance Rule, *before*
the first vendor profile is written)
- Add `vend/cloud/boot/cadences.go` with `EVERY_30_SECONDS_CLOUD`,
  `EVERY_5_MINUTES_CLOUD`, `EVERY_15_MINUTES_CLOUD`, etc., copying the
  pattern from `vend/parser/boot/cadences.go`.
- Add `vend/cloud/boot/helpers.go` with `createCloudPoll(name,
  endpoint, cadence, propertyId, mapping, paginateMode, rateLimitRps)`
  and `createCloudRestAttribute(propertyId, mapping)`, copying the
  pattern from `vend/parser/boot/helpers.go`. These helpers know how
  to attach the `Paginate` + `RateLimit` + `CloudListFanOut` +
  `RestJsonParse` rule chain so each vendor file does not redeclare
  it.
- **No vendor file may redefine its own `createPoll` or cadence
  constants.** If a vendor needs a one-off (e.g., a non-paginated
  endpoint), pass an option to the shared helper rather than fork it.

### Phase 3 — Activate cloud collector + first vendor (Nayax)
- **Reuse `vend/collector/main.go`** (no new binary). Add a second
  activation:
  ```go
  service.Activate(vendcommon.VendMachine_Links_ID, nic)  // existing
  service.Activate(vendcommon.VendCloud_Links_ID, nic)    // new
  ```
- Add `vend/cloud/boot/nayax_pollaris.go` — pure data, built entirely
  on the helpers from Phase 2.5. Per-domain polls (devices, status,
  inventory, transactions, cash, telemetry, alerts, dex). Target file
  size: ~80 lines. If it grows beyond ~120, add another shared helper
  rather than vendor-specific code.
- The existing `k8s/vend-collector.yaml` and `Dockerfile` are reused —
  no new k8s manifest for the collector.
- Verify the collector polls Nayax-shaped simulator endpoints and
  produces multi-page results with no parser-side fan-out yet.

### Phase 4 — Activate cloud parser + fan-out rules
- **If Phase 2.5(A) passed**: reuse `vend/parser/main.go`. Add a second
  activation:
  ```go
  service.Activate(common.VendMachine_Links_ID,
      &vend.VendMachine{}, false, nic, "MachineId")  // existing
  service.Activate(common.VendCloud_Links_ID,
      &vend.VendMachine{}, false, nic, "MachineId")  // new
  ```
- **If Phase 2.5(A) failed**: fork into `vend/cloud/parser/main.go` with
  its own `build.sh` / `Dockerfile` / `k8s/vend-cloud-parser.yaml`.
  Document the rejection reason in §1.5.
- Rules (`CloudListFanOut`, `CloudTransactionAppend`,
  `CloudAlertNormalize`) were added to `l8parser` core in Phase 2 PR 2.
  Verify they are available by name in the parser. The key design
  constraints:
  - `CloudListFanOut` is the **FINAL rule** in each poll's attribute
    chain (per constraint C3 — no downstream chaining). It receives the
    full combined JSON array, iterates elements, and writes per-instance
    property paths into the cache.
  - `CloudTransactionAppend` is a **FINAL rule** — append-only writes.
  - `CloudAlertNormalize` is a **FINAL rule** — maps vendor codes to
    internal enums.
  - For slot updates within a machine: reuse `RestGpuParse` directly
    (after the §5.4 inventory `map<>` migration) or have
    `CloudListFanOut` handle nested slot arrays inline.
- Add `build.sh`, `Dockerfile`, `k8s/vend-cloud-parser.yaml` (DaemonSet).
- Verify a 50-machine Nayax response produces 50 parser invocations.

### Phase 4b — Activate cloud inventory cache
- **If Phase 2.5(A) passed**: reuse the existing inventory cache binary
  (the one currently launched by `k8s/vend.yaml`). Add a second
  activation:
  ```go
  inventory.Activate(common.VendMachine_Links_ID,
      &vend.VendMachine{}, &vend.VendMachineList{}, nic, "MachineId")  // existing
  inventory.Activate(common.VendCloud_Links_ID,
      &vend.VendMachine{}, &vend.VendMachineList{}, nic, "MachineId")  // new
  s, a := targets.Links.Cache(common.VendCloud_Links_ID)
  invCenter := inventory.Inventory(res, s, a)
  invCenter.AddMetadata("Online", Online)
  ```
- **If Phase 2.5(A) failed**: fork into `vend/cloud/inventory/main.go`
  modeled on `probler/go/prob/inv_gpu/main.go`, with its own
  `build.sh` / `Dockerfile` / `k8s/vend-cloud-inventory.yaml`
  (StatefulSet). Document the rejection reason in §1.5.
- Define an `Online(any) (bool, string)` function that returns true when
  `last_heartbeat` is within the cadence-derived staleness threshold
  (mirrors probler's GPU `DEVICE_STATUS_ONLINE` check).
- Verify the cache holds a `VendMachine` per machine collected and that
  the `Online` metadata flips to true when fresh polls arrive.

### Phase 5 — Webhook receiver
- `vend/cloud/webhook/main.go` — small HTTPS server (Pollaris-independent,
  push not pull).
- Route: `POST /webhook/{vendor}/{tenant}` → HMAC verify → 24 h replay
  dedupe → publish into the cloud parser service area.
- HMAC secrets resolved through the security provider (same store as
  `CredId`).
- `build.sh`, `Dockerfile`, `k8s/vend-cloud-webhook.yaml` (DaemonSet,
  hostNetwork or Service+Ingress).

### Phase 6 — Three remaining vendor profiles + build/run-local wiring
- Add `vend/cloud/boot/cantaloupe_pollaris.go`,
  `vend/cloud/boot/m365_pollaris.go`,
  `vend/cloud/boot/televend_pollaris.go` — each a pure config file
  built entirely on `helpers.go` + `cadences.go` from Phase 2.5.
- Update `go/build-all-images.sh` with any new binaries (webhook is
  always new; collector/parser/inventory may or may not be forked per
  Phase 2.5(A) outcome).
- Update `go/run-local.sh`, pointing binaries at the
  simulator's URL via env var.
- If any vendor file exceeds ~120 lines or needs Go-level branching, stop
  and extend the shared helper instead.

### Phase 7 — UI surface (l8ui, config-only, desktop + mobile parity)
Per `prd-compliance.md` Rule 3, follow these l8ui rules exactly:
- `adding-module-desktop.md` — desktop module file layout
- `adding-module-mobile.md` — mobile module file layout
- `mobile-rules.md` — desktop/mobile **functional parity** (not just UI;
  every cloud-tenant operation must work identically on both)
- `factory-components.md` — use `Layer8EnumFactory`, `Layer8RefFactory`,
  `Layer8ColumnFactory`, `Layer8FormFactory`
- `enum.md` + `enum-renderer-column-cascade.md` — only valid factory
  methods (`create`, `simple`, `withValues`); no `createEnum` /
  `createStatus`
- `layer8d-table.md` + `layer8m-table.md` — table construction
- `layer8d-forms.md` + `layer8m-forms.md` — form construction
- `js-protobuf-field-names.md` — every JS field name verified against the
  `L8PTarget` `.pb.go` (e.g., `targetId`, NOT `targetID` or `TargetId`)
- `reference-registry-completeness.md` — every `lookupModel` referenced by
  the form MUST be registered in BOTH desktop
  (`reference-registry-vend.js`) and mobile
  (`layer8m-reference-registry-vend.js`)
- `module-init-section-selector.md` — `sectionSelector` must equal
  `defaultModule`
- `desktop-script-loading-order.md` + `mobile-script-loading-order.md` —
  load order discipline

Surface:
- "Cloud Tenants" service under the Fleet module.
- Columns surfaced from `L8PTarget`: target_id, vendor (from
  `Hosts.api.Groups[vendor]`), state, last poll timestamp, machine count
  (derived from cache).
- Form: target_id, vendor select (enum, with `UNSPECIFIED=0`), base URL,
  credential select (reference picker populated from the security
  provider — never free-text), Pollaris profile select, page size,
  rate-limit rps, webhook on/off.
- Read-only Health tab showing recent poll history per target.

Mobile parity:
- Same five columns rendered as a card list with `primary`/`secondary`.
- Same form rendered via `Layer8MForms`.
- Reference registry entries duplicated in
  `layer8m-reference-registry-vend.js`.
- Nav config entry in `layer8m-nav-config-fleet.js` with correct
  `idField` (lowercase `targetId` per `js-protobuf-field-names.md` —
  this is the 5x-regression trap).

### Phase 8 — End-to-end verification

**Test location and approach** (per `test-location-and-approach.md`):
- All tests live under `go/tests/cloud/`. Do NOT place `_test.go` files
  alongside source in `go/vend/cloud/...`.
- Tests exercise the **system API** (HTTP calls into the cloud collector,
  parser service area, inventory cache). Do NOT call unexported internal
  functions.
- Per `test-data-field-verification.md`, every JSON key in test request
  bodies MUST be verified against the corresponding `.pb.go`
  `protobuf:"...,json=..."` tag.

**Verification scenarios** — for each of the four vendors, against the
simulator's vendor-shaped endpoints:
1. Create a tenant `L8PTarget` pointing at the simulator URL.
2. Verify the simulator's machines appear in the `VendMachine` cache
   within 60 s.
3. Verify inventory slots populate and update on the next poll cycle.
4. Verify a vend event emitted by the simulator generates a
   `VendTransaction` row.
5. Verify a webhook POST from the simulator (Nayax + Cantaloupe only)
   shortcuts the poll.
6. Kill the collector pod; verify state survives (cache + persist).
7. Rotate the OAuth2 secret in the security provider; verify next poll
   picks it up without restart.
8. **UI smoke test** (per `mobile-rules.md` parity rule): navigate to the
   Cloud Tenants section on **both** desktop and mobile, verify the
   tenant list loads, click a row, verify the detail popup opens with
   reference fields resolved (display names, not IDs).

Then on **two real vendor accounts** (whichever the user has — Nayax +
one other) repeat steps 1–4. Mark the others "deferred — pending account
access" if unavailable.

---

## 11. Out of scope

- Pushing config **down** to vendor clouds (planogram pushes, price changes).
  This plan is read-only ingress. Bidirectional sync is a separate plan.
- Migrating the existing direct-poll path. Both paths coexist; an operator
  chooses per machine which one to use by registering it as a `VendMachine`
  target (`LinksId="Vend"`) or by including it inside the operator scope of a
  cloud `L8PTarget` (`LinksId="VendCloud"`).
- Vendor SDKs. We deliberately use only documented HTTPS REST APIs to avoid
  taking on vendor SDK dependencies.

---

## 12. Risks and mitigations

| Risk | Mitigation |
|------|-----------|
| Vendor API schema drift breaks parsing | Per-vendor `mapping` strings live in config (`boot/`); a schema change is a 1-line edit, not a code change |
| Rate limits cause data gaps | Backoff + per-tenant token bucket; alerting if `last_poll_ok > 2× cadence` |
| Webhook flood (vendor outage replay storm) | Replay cache + bounded queue with drop-oldest |
| Secret leakage via logs | All secret reads go through security provider; no `Sprintf("%v", tenant)` in logs — covered by code review checklist |
| Two collectors poll the same tenant | Sharding by `hash(tenant_id) % replicas`; each pod owns a slice |
| Cache contention with direct-poll path | Same `VendMachine` writes — last-writer-wins on the parser timestamp; document that an operator should not run both paths against the same machine |
