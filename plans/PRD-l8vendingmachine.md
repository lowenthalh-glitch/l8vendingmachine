# PRD: L8VendingMachine -- Next-Gen AI Vending Machine Management System

## Overview

**L8VendingMachine** is a next-generation AI-powered vending machine management system built on the Layer 8 Ecosystem. It provides fleet management, real-time telemetry, inventory optimization, sales analytics, predictive maintenance, and AI-driven demand forecasting for vending machine operators managing fleets of smart vending machines.

The system integrates with the **l8opensim** vending machine simulator (TCN-ZK locker machines, Afen AF-60C beverage machines, Afen AF-D900 combo machines) for development and testing, and connects to real machines via REST API telemetry in production.

**Project location:** `../l8vendingmachine` (sibling to l8erp, l8ui, l8opensim)

---

## 1. Modules and Components

### 1.1 Fleet Management (fleet)
- Machine master data (model, serial, firmware, controller board, location)
- Machine status monitoring (online/offline, operational state, connectivity)
- Machine configuration (locale, currency, timezone, energy saving schedules)
- Machine grouping by location, region, operator, machine type
- Firmware version tracking and update management
- Machine lifecycle (installation, commissioning, decommissioning)
- GPS/location tracking with map visualization

### 1.2 Inventory Management (inventory)
- Real-time slot-level inventory tracking (stock quantity, capacity, par levels)
- Planogram management (product-to-slot assignment)
- Product catalog (SKU, name, category, price, expiration)
- Low-stock and sold-out alerting
- Restock order generation
- Product expiration tracking
- Slot configuration (spring motor, locker cell, conveyor)

### 1.3 Sales & Transactions (sales)
- Itemized transaction recording (product, slot, price, payment method, timestamp)
- Sales summary and aggregation (daily, weekly, monthly)
- Payment method breakdown (cash, card, NFC, QR, mobile wallet)
- Failed vend tracking and refund management
- Revenue reporting by machine, location, product, time period
- Receipt and settlement management
- **Immutable**: `VendTransaction` rejects PUT -- transactions are append-only (follows `l8events` `rejectPut()` pattern)

### 1.4 Payment Systems (payment)
- Payment peripheral status monitoring (coin acceptor, bill validator, card reader, QR scanner)
- Cash position tracking (coin tubes by denomination, bill stacker, cash box)
- Cash collection scheduling and recording
- Cashless transaction settlement tracking
- Change availability monitoring (exact change required detection)
- Payment device health and error tracking

### 1.5 Temperature & Refrigeration (temperature)
- Multi-zone temperature monitoring (ambient, refrigerated, frozen)
- Temperature setpoint management
- Compressor status and duty cycle tracking
- Defrost cycle monitoring
- Temperature alarm management (out-of-range alerts)
- Health department compliance logging
- Glass heater and door seal monitoring

### 1.6 Alerts & Maintenance (maintenance)
- Alert management (create, acknowledge, resolve)
- Alert categorization (INVENTORY, TEMPERATURE, PAYMENT, MECHANICAL, CONNECTIVITY)
- Alert severity levels (CRITICAL, WARNING, INFO)
- Predictive maintenance scoring per component
- Motor current draw trending (spring motor degradation)
- Compressor health trending (cycle frequency, power draw)
- Bill validator rejection rate trending
- Service visit scheduling and recording
- Work order creation and tracking

### 1.7 Route Optimization (route)
- Service route definition and management
- Route stop ordering by priority (revenue risk, stock urgency, cash box fullness)
- Driver assignment and vehicle capacity tracking
- Product load list generation per route
- Service visit duration estimation
- Route efficiency metrics (planned vs actual)

### 1.8 AI Analytics (analytics)
- Demand forecasting (by product, machine, hour, day, season)
- Sales velocity and product performance ranking
- Slot productivity analysis (revenue per slot per day)
- Customer traffic pattern analysis (approaches vs purchases, conversion rate)
- Dynamic pricing recommendations
- Product swap recommendations (replace slow movers)
- Predicted stockout time estimation
- Restock urgency scoring
- Energy consumption analysis and optimization
- Anomaly detection (unusual sales patterns, security events)
- Fleet inventory summary (aggregated product stock counts, pricing, and deployment across all machines)
  - Total units of each product across all machines fleet-wide
  - Total units by product category
  - Products with zero stock (sold out fleet-wide)
  - Products below fleet-wide par level threshold
  - Average price per product across all machines
  - Number of machines carrying each product
  - Warehouse stock vs machine stock comparison (total supply chain position)

### 1.9 Access & Security (access)
- Door open/close event logging
- Service visit correlation (door event + DEX read = service visit)
- Lock status monitoring (locker machines)
- Operator authentication and access control
- Tamper detection alerts
- Camera event integration (optional)
- **Immutable**: `VendAccessEvent` rejects PUT -- event logs are append-only

### 1.10 DEX Audit (dex)
- DEX/UCS audit data collection and storage
- EVA-DTS compliance reporting
- Audit data comparison (interval vs cumulative)
- Cash audit reconciliation
- Selection-level vend counts
- Event log aggregation
- **Immutable**: `VendDexAudit` rejects PUT -- audit records are append-only

### 1.11 Warehouse & Supply Chain (warehouse)
- Warehouse master data (name, address, capacity, operating hours, contact)
- Product stock levels per warehouse (quantity on hand, reorder point, reorder quantity)
- Stock movement tracking (receive from supplier, transfer to vehicle, return from vehicle, write-off/spoilage)
- Supplier management (name, contact, lead time, payment terms, product catalog)
- Purchase orders to suppliers (order lines with product/quantity/price, status lifecycle: DRAFT → SUBMITTED → CONFIRMED → RECEIVED → CLOSED)
- Vehicle inventory tracking (products loaded onto each route vehicle, pre-route load list vs post-route remaining)
- Low-stock warehouse alerts (product falls below reorder point)
- Stock reconciliation (expected vs actual after route completion)
- Demand-driven reorder suggestions (AI analytics predicts warehouse depletion based on fleet consumption rates)
- Multi-warehouse support (regional warehouses serving different route clusters)
- Expiration tracking at warehouse level (FIFO enforcement, near-expiry alerts)
- **Immutable**: `VendStockMovement` rejects PUT -- movement records are append-only ledger entries

### 1.12 Dashboard & KPIs (dashboard)
Integrates with the l8ui `Layer8DWidget` component and follows the l8erp BI module patterns (`../l8erp/go/erp/bi/dashboards/`, `../l8erp/go/erp/bi/kpis/`).

- **KPI Definitions** (following `BiKPI` pattern from l8erp, ServiceArea=35):
  - `VEND_KPI_REVENUE_TODAY` -- Total fleet revenue today (sum of VendTransaction amounts)
  - `VEND_KPI_VENDS_TODAY` -- Total vend count today
  - `VEND_KPI_MACHINES_ONLINE` -- Count of machines with status OPERATIONAL
  - `VEND_KPI_MACHINES_OFFLINE` -- Count of machines with status OFFLINE or OUT_OF_SERVICE
  - `VEND_KPI_CRITICAL_ALERTS` -- Count of active CRITICAL severity alerts
  - `VEND_KPI_RESTOCK_URGENT` -- Count of machines with restock urgency HIGH or CRITICAL
  - `VEND_KPI_CASH_TO_COLLECT` -- Total cash across all machines pending collection
  - `VEND_KPI_AVG_FILL_RATE` -- Average inventory fill rate (current qty / capacity) fleet-wide
  - `VEND_KPI_FAILED_VEND_RATE` -- Failed vend percentage over last 24 hours
  - `VEND_KPI_TEMP_VIOLATIONS` -- Count of machines with temperature out of compliance

- **KPI Thresholds**: Each KPI has ON_TARGET / AT_RISK / OFF_TARGET status computed via `BiKPIThreshold` rules (e.g., machines offline > 5% = AT_RISK, > 15% = OFF_TARGET)

- **Dashboard Layout** (following `BiDashboard` pattern):
  - Top row: 5 KPI cards with sparkline trends (revenue, vends, online, alerts, restock)
  - Middle row: Revenue chart (bar, last 7 days) + inventory heatmap by product category
  - Bottom row: Recent alerts table + machines needing service list

- **Widget Rendering**: Uses `Layer8DWidget.render()` for KPI cards with sparkline data, `Layer8DChart` for charts, standard `Layer8DTable` for alert/machine lists

### 1.13 Map Visualization (map)
Integrates with the topology map rendering from `../l8topology` (`topology-map.js`, SVG/WebGL rendering).

- **Fleet Map View**: Interactive map showing all machine locations with color-coded status markers:
  - Green = OPERATIONAL
  - Yellow = WARNING (has active WARNING alerts)
  - Red = OFFLINE or CRITICAL alert
  - Gray = DECOMMISSIONED or MAINTENANCE
  - Blue = IN_SERVICE (currently being restocked)

- **Map Interactions**:
  - Click machine marker -> popup with machine summary (model, status, stock level, revenue today, last heartbeat)
  - Click popup "Details" -> navigate to machine detail view
  - Cluster markers when zoomed out (show count + aggregate status)
  - Filter by status, machine type, route, operator

- **Route Overlay**: Show active service routes on the map with stop order, connecting lines between machines, driver position (if GPS-tracked)

- **Data Source**: Machine GPS coordinates from `VendLocation.coordinates` (VendGpsCoordinates: latitude, longitude). Uses SVG world map background from l8topology or integrates with tile-based map provider.

### 1.14 Compliance & Health Inspections (compliance)
Follows the compliance module patterns from `../l8erp/go/erp/comp/` (audit schedules, findings, controls, certifications).

- **Health Inspection Schedules** (following `CompAuditSchedule` pattern):
  - Scheduled food safety inspections per machine or location
  - Inspection types: HEALTH_DEPARTMENT, INTERNAL_AUDIT, FOOD_SAFETY, EQUIPMENT_SAFETY
  - Recurrence: monthly, quarterly, annually
  - Assigned inspector (internal or health department contact)
  - Planned vs actual dates, status tracking

- **Inspection Findings** (following `CompAuditFinding` pattern):
  - Finding severity: CRITICAL, HIGH, MEDIUM, LOW, INFORMATIONAL
  - Standard audit structure: condition (what was found), criteria (what should be), cause, effect, recommendation
  - Remediation actions with due dates and responsible party
  - Repeat finding tracking (links to prior finding)
  - Evidence document attachment (photo of violation, temperature log)

- **Compliance Controls** (following `CompControl` pattern):
  - Temperature logging compliance (HACCP requirements for refrigerated machines)
  - Product expiration tracking compliance
  - Cash handling procedures
  - Cleaning and sanitization schedules
  - Machine safety checks (electrical, mechanical)
  - Control effectiveness testing (EFFECTIVE, PARTIALLY_EFFECTIVE, INEFFECTIVE)

- **Certifications** (following `CompCertification` pattern):
  - Health department operating permits per location
  - Food handler certifications for service personnel
  - Equipment safety certifications
  - Expiry tracking with renewal alerts

- **Compliance Dashboard**: Inspection pass rate, open findings count, overdue remediation actions, upcoming inspections, certification expiry warnings

### 1.15 Scheduled Reports (reports)
Follows the BI reporting patterns from `../l8erp/go/erp/bi/reports/` (`BiReport`, `BiReportSchedule`, `report_scheduler.go`).

- **Report Definitions** (following `BiReport` pattern):
  - `VEND_RPT_DAILY_REVENUE` -- Daily revenue by machine, location, product
  - `VEND_RPT_WEEKLY_INVENTORY` -- Weekly inventory status across fleet
  - `VEND_RPT_MONTHLY_PNL` -- Monthly profit & loss per machine (revenue - product cost - service cost)
  - `VEND_RPT_CASH_COLLECTION` -- Cash collection summary (collected vs outstanding)
  - `VEND_RPT_MAINTENANCE_DUE` -- Machines with upcoming or overdue maintenance
  - `VEND_RPT_PRODUCT_PERFORMANCE` -- Product ranking by revenue, margin, velocity
  - `VEND_RPT_COMPLIANCE_STATUS` -- Open findings, upcoming inspections, certification expiries
  - `VEND_RPT_ROUTE_EFFICIENCY` -- Route performance (planned vs actual stops, duration, fill rates)
  - `VEND_RPT_WAREHOUSE_STOCK` -- Warehouse stock levels, reorder alerts, PO status
  - `VEND_RPT_FLEET_HEALTH` -- Machine health scores, predictive maintenance alerts, component trends

- **Report Schedules** (following `BiReportSchedule` pattern):
  - Frequency: DAILY, WEEKLY, MONTHLY, QUARTERLY
  - Delivery: email to configurable recipients (operators, managers, location contacts)
  - Output format: PDF, CSV, Excel
  - Time of day: configurable (e.g., daily revenue report at 6:00 AM)
  - Active/inactive toggle

- **Report Subscriptions** (following `BiReportSubscription` pattern):
  - Multiple subscribers per report schedule
  - Per-subscriber format preference
  - Include/exclude empty reports

- **Report Scheduler** (following `report_scheduler.go` pattern):
  - Background goroutine checking due schedules every 1 minute
  - Computes next run based on frequency (DAILY = +24h, WEEKLY = +7d, MONTHLY = +1mo)
  - Delegates report execution to generate data, format output, and deliver via email
  - Execution history tracking (start/end time, row count, status, errors)

### 1.16 Threshold Events (via l8events)
Integrates with `../l8events` (EventRecord service, ServiceArea=76) to generate structured events when vending machine metrics cross configurable thresholds.

- **Temperature threshold events**: Cabinet temperature exceeds safe range (e.g., > 8C for refrigerated, > 25C for ambient)
- **Inventory threshold events**: Stock level falls below par level, slot sells out, product approaching expiration
- **Payment threshold events**: Cash box exceeds capacity %, coin tube empty (exact change required), card reader offline
- **Mechanical threshold events**: Motor current draw exceeds baseline (degradation), compressor duty cycle > 85%, vend failure rate > threshold
- **Connectivity threshold events**: Machine offline > N minutes, signal strength below threshold, heartbeat missed
- **Sales threshold events**: Revenue anomaly (sudden drop/spike), unusual transaction pattern, refund rate exceeds threshold

Each threshold crossing generates an `l8events.EventRecord` with:
- `category`: `EVENT_CATEGORY_PERFORMANCE` (metrics) or `EVENT_CATEGORY_SYSTEM` (connectivity/state)
- `severity`: Computed from threshold proximity (WARNING at 80%, CRITICAL at 100%)
- `source_id`: Machine ID
- `attributes`: Map with `metricName`, `currentValue`, `thresholdValue`, `thresholdType` (UPPER/LOWER)

Threshold rules are configured per machine type via `AlarmDefinition` records (see 1.12).

### 1.17 Alarms & Correlation (via l8alarms)
Integrates with `../l8alarms` (Alarm service, ServiceArea=10) for intelligent alarm lifecycle management.

- **Alarm Definitions**: Pre-configured rules that match threshold events to alarms:
  - `VEND_TEMP_HIGH` -- Temperature above safe limit (severity: CRITICAL, auto-clear when temp drops)
  - `VEND_STOCK_LOW` -- Slot below par level (severity: WARNING, auto-clear on restock)
  - `VEND_SOLD_OUT` -- Slot completely empty (severity: MINOR)
  - `VEND_CASH_BOX_FULL` -- Cash box above 90% capacity (severity: WARNING)
  - `VEND_EXACT_CHANGE` -- Cannot make change (severity: WARNING)
  - `VEND_PAYMENT_OFFLINE` -- Card reader or NFC offline (severity: MAJOR)
  - `VEND_MOTOR_DEGRADED` -- Motor current > 130% of baseline (severity: WARNING)
  - `VEND_COMPRESSOR_FAIL` -- Compressor not running when temp > setpoint+5C (severity: CRITICAL)
  - `VEND_MACHINE_OFFLINE` -- No heartbeat for > 10 minutes (severity: CRITICAL)
  - `VEND_VEND_FAILURE_RATE` -- Failed vend rate > 5% over last hour (severity: MAJOR)

- **Alarm Lifecycle**: ACTIVE -> ACKNOWLEDGED -> CLEARED (with full state history audit trail)
- **Deduplication**: Same alarm on same machine within time window increments `occurrence_count` instead of creating duplicate
- **Auto-clear**: Temperature and stock alarms auto-clear when the condition resolves (via `clear_event_pattern`)

- **Correlation Rules** (using l8alarms correlation engine):
  - **Temporal**: "Compressor Failure" within 15 minutes before "Temperature High" -> compressor is root cause, temp alarm suppressed as symptom
  - **Pattern**: "Payment Offline" + "Vend Failure Rate" -> payment system is root cause
  - **Pattern**: "Machine Offline" suppresses all other alarms from that machine (connectivity is root cause)

- **Maintenance Windows**: Suppress alarms during scheduled service visits (restocking, cleaning, maintenance). Configured per machine or per location with optional recurrence (daily, weekly).

- **Alarm Filters**: Saved views for operators (e.g., "Critical Only", "My Route Machines", "Temperature Alarms", "Unacknowledged")

### 1.18 Notifications (via l8notify)
Integrates with `../l8notify` for multi-channel alert delivery when alarms fire or escalate.

- **Notification Policies**: Configurable rules determining who gets notified for which alarms:
  - Route driver receives WEBHOOK/EMAIL for WARNING+ alarms on their route machines
  - Operations manager receives EMAIL for all CRITICAL alarms
  - Maintenance team Slack channel receives all MECHANICAL category alarms
  - Location contact receives EMAIL when machine goes offline at their site

- **Notification Channels**:
  - **EMAIL**: Via SMTP to operators, managers, location contacts
  - **WEBHOOK**: HTTP POST to external systems (ticketing, ERP, fleet management)
  - **SLACK**: Incoming webhook to team channels (e.g., #vending-alerts, #maintenance)
  - **CUSTOM**: Pluggable for SMS gateways, push notification services, PagerDuty

- **Throttling**: Per-policy cooldown (e.g., max 1 notification per machine per 15 minutes) and hourly rate limits to prevent notification storms

- **Escalation Policies**: Time-based escalation if alarms are not acknowledged:
  - Step 1 (0 min): Notify route driver via webhook
  - Step 2 (15 min): Notify area supervisor via email
  - Step 3 (60 min): Notify operations manager via email + Slack
  - Escalation cancelled when alarm is acknowledged or cleared

- **Template Variables**: Notifications use templates with `{{alarm.name}}`, `{{alarm.severity}}`, `{{alarm.nodeName}}`, `{{alarm.location}}`, `{{alarm.description}}` placeholders

### 1.19 Data Retention & Archival (retention)
High-volume entities (`VendTransaction`, `VendTempReading`, `VendAccessEvent`, `VendDexAudit`, `VendStockMovement`) accumulate rapidly from collector polling. Uses the archival infrastructure from `../l8events/go/archive/` (Archive → save to archived table → delete from active table).

- **Archived Entity Services** (following `../l8alarms` `ArchivedAlarm`/`ArchivedEvent` pattern):
  - `VendArchivedTransaction` -- immutable copy of old transactions with `archivedAt`, `archivedBy`, `archiveReason`
  - `VendArchivedTempReading` -- immutable copy of old temperature readings
  - `VendArchivedAccessEvent` -- immutable copy of old access events

- **Retention Policies**: Configurable per entity type:
  - `VendTransaction`: Archive after 90 days, delete archived after 2 years
  - `VendTempReading`: Archive after 30 days, delete archived after 1 year
  - `VendAccessEvent`: Archive after 60 days, delete archived after 1 year
  - `VendDexAudit`: Archive after 180 days, keep archived indefinitely (compliance)
  - `VendStockMovement`: Archive after 90 days, delete archived after 2 years

- **Retention Scheduler** (`retention_scheduler.go`): Background goroutine (follows `report_scheduler.go` pattern) that runs daily at configurable time:
  1. For each entity type: query records older than retention threshold
  2. Copy to archived service (POST to `VendArchived*`)
  3. Delete from active service
  4. Log: archived count, deleted count, errors

### 1.20 AI Agent Chat (via l8agent)
Integrates with `../l8agent` for AI-powered chat interface. Requires full backend activation (`Initialize()` + `InitializeChat()` from `../l8agent/go/init.go`), not just frontend script includes.

- **Backend**: Activate l8agent services (Conversations, Messages, Prompts, Chat orchestration) in the main vend service startup. `InitializeChat()` must be called AFTER the introspector is populated with all VendMachine types.
- **Desktop**: Include `l8agent-chat.js`, `l8agent-bubble.js`, `l8agent-enums.js`, `l8agent-columns.js`, `l8agent-forms.js` in `app.html`. Include `l8agent-chat.css`, `l8agent-bubble.css` in CSS section.
- **Mobile**: Include `l8agent/m/l8agent-chat-m.js` in `m/app.html`.
- **Capabilities**: Operators can ask the AI agent questions like "Which machines need restocking today?", "What's the revenue for Building A this week?", "Show me all critical alerts", and the agent queries the vend services via L8Query.

### 1.21 Fleet Topology Map (via l8topology)
Registers vending machines as topology nodes for map visualization. Follows the `Layer1` discovery pattern from `../l8topology/go/topo/discover/Layer1.go`.

- **Topology Registration**:
  1. `ActivateVendTopology()` -- registers `VendMachine` type with topology service
  2. Adds primary key decorator: `AddPrimaryKeyDecorator(&VendMachine{}, "MachineId")`
  3. Discovery query: `select * from VendMachine` to fetch all machines

- **Node Conversion** (`ConvertToTopologyNode`):
  - `machine.MachineId` → `node.NodeId`
  - `machine.Model` → `node.Name`
  - `location.Name` → `node.Location`
  - `machine.MachineType` → `node.Type` (custom vending machine node type)
  - `location.Coordinates` (lat/lng) → SVG coordinates via Robinson projection

- **Custom SVG Icon**: Vending machine icon registered via `Layer8SvgFactory.registerTemplate('vendingMachine', ...)` -- rectangle with product grid pattern, status-colored based on `VendMachineStatus`

- **Node Type Mapping**:
  - LOCKER → "Locker Vending" (purple icon)
  - REFRIGERATED_BEVERAGE → "Beverage Vending" (blue icon)
  - COMBO → "Combo Vending" (green icon)

### 1.22 VNet Architecture
All 5 binaries (main, ui, vnet, collector, parser) share the same VNet ID (`VEND_VNET = 49010`), following the probler pattern where collector, parser, and main service are all part of the same interconnected mesh.

- **Service Routing**: The collector sends completed CJob results to the parser via VNet service routing (Links ID → parser service name/area). The parser writes parsed entities to the main vend services via VNet.
- **UI Service**: The UI binary creates a vnic on the same VNet to proxy REST API calls to backend services.
- **No Separate VNet**: Unlike probler which uses a separate VNet for log aggregation (`LOGS_VNET`), the vending system has no log collection component -- all 5 binaries share `VEND_VNET`.

---

## 2. Global Rules Compliance

### Project Structure & Architecture
- [x] Project structure follows l8erp layout
- [x] Directory names and file naming conventions match l8erp patterns

### Protobuf Design
- [x] Enum zero values are UNSPECIFIED (`proto-enum-zero-value`)
- [x] List types use `repeated X list = 1` convention (`proto-list-convention`)
- [x] No direct struct references between Prime Objects -- ID fields only (`prime-object-references`)
- [x] Child entities are embedded `repeated` fields, not separate services (`prime-object-references`)

### Service Design
- [x] ServiceName is 10 characters or less (`maintainability`)
- [x] ServiceArea is consistent within a module (`maintainability`)
- [x] ServiceCallback auto-generates primary key on POST (`maintainability`)
- [x] Types are registered in UI main.go (`maintainability`)

### UI Design
- [x] All UI module integration steps are planned
- [x] Desktop and mobile parity is addressed (`mobile-rules`)
- [x] Immutable entities/fields have read-only UI (`immutability-ui-alignment`)
- [x] Child types use inline tables, not standalone UI (`prime-object-references`)
- [x] UI components follow l8ui guide

### Mock Data
- [x] All services have mock data generators planned (`data-completeness-pipeline`)
- [x] Phase ordering accounts for cross-module dependencies (`mock-phase-ordering`)

### Deployment
- [x] Deployment artifacts included: build.sh, Dockerfile, K8s YAML (`deployment-artifacts`)
- [x] run-local.sh section included (`run-local-script`)
- [x] K8s YAMLs include all required entries (`k8s-yaml-required-entries`)

### Configuration
- [x] login.json adaptation planned (`login-json-adaptation`)
- [x] ModConfig handling addressed (`modconfig-failure-no-logout`)

---

## 3. Project Structure

```
l8vendingmachine/
├── go/
│   ├── go.mod
│   ├── go.sum
│   ├── vendor/
│   ├── run-local.sh
│   ├── build-all-images.sh
│   │
│   ├── vend/                              # Main service module
│   │   ├── common/
│   │   │   └── defaults.go               # PREFIX="/vend/", VNET, DB constants
│   │   │
│   │   ├── fleet/                         # Fleet Management services
│   │   │   ├── machines/
│   │   │   │   ├── MachineService.go
│   │   │   │   └── MachineServiceCallback.go
│   │   │   ├── machinegroups/
│   │   │   │   ├── MachineGroupService.go
│   │   │   │   └── MachineGroupServiceCallback.go
│   │   │   └── locations/
│   │   │       ├── LocationService.go
│   │   │       └── LocationServiceCallback.go
│   │   │
│   │   ├── inventory/                     # Inventory Management services
│   │   │   ├── products/
│   │   │   │   ├── ProductService.go
│   │   │   │   └── ProductServiceCallback.go
│   │   │   ├── planograms/
│   │   │   │   ├── PlanogramService.go
│   │   │   │   └── PlanogramServiceCallback.go
│   │   │   └── restockorders/
│   │   │       ├── RestockOrderService.go
│   │   │       └── RestockOrderServiceCallback.go
│   │   │
│   │   ├── sales/                         # Sales & Transactions services
│   │   │   ├── transactions/
│   │   │   │   ├── TransactionService.go
│   │   │   │   └── TransactionServiceCallback.go
│   │   │   └── settlements/
│   │   │       ├── SettlementService.go
│   │   │       └── SettlementServiceCallback.go
│   │   │
│   │   ├── payment/                       # Payment Systems services
│   │   │   ├── cashpositions/
│   │   │   │   ├── CashPositionService.go
│   │   │   │   └── CashPositionServiceCallback.go
│   │   │   └── collections/
│   │   │       ├── CollectionService.go
│   │   │       └── CollectionServiceCallback.go
│   │   │
│   │   ├── temperature/                   # Temperature & Refrigeration
│   │   │   └── tempreadings/
│   │   │       ├── TempReadingService.go
│   │   │       └── TempReadingServiceCallback.go
│   │   │
│   │   ├── maintenance/                   # Alerts & Maintenance services
│   │   │   ├── alerts/
│   │   │   │   ├── AlertService.go
│   │   │   │   └── AlertServiceCallback.go
│   │   │   ├── workorders/
│   │   │   │   ├── WorkOrderService.go
│   │   │   │   └── WorkOrderServiceCallback.go
│   │   │   └── servicevisits/
│   │   │       ├── ServiceVisitService.go
│   │   │       └── ServiceVisitServiceCallback.go
│   │   │
│   │   ├── route/                         # Route Optimization
│   │   │   ├── routes/
│   │   │   │   ├── RouteService.go
│   │   │   │   └── RouteServiceCallback.go
│   │   │   └── drivers/
│   │   │       ├── DriverService.go
│   │   │       └── DriverServiceCallback.go
│   │   │
│   │   ├── analytics/                     # AI Analytics
│   │   │   ├── forecasts/
│   │   │   │   ├── ForecastService.go
│   │   │   │   └── ForecastServiceCallback.go
│   │   │   └── performance/
│   │   │       ├── PerformanceService.go
│   │   │       └── PerformanceServiceCallback.go
│   │   │
│   │   ├── access/                        # Access & Security
│   │   │   └── accessevents/
│   │   │       ├── AccessEventService.go
│   │   │       └── AccessEventServiceCallback.go
│   │   │
│   │   ├── dex/                           # DEX Audit
│   │   │   └── audits/
│   │   │       ├── DexAuditService.go
│   │   │       └── DexAuditServiceCallback.go
│   │   │
│   │   ├── warehouse/                     # Warehouse & Supply Chain
│   │   │   ├── warehouses/
│   │   │   │   ├── WarehouseService.go
│   │   │   │   └── WarehouseServiceCallback.go
│   │   │   ├── warehousestock/
│   │   │   │   ├── WarehouseStockService.go
│   │   │   │   └── WarehouseStockServiceCallback.go
│   │   │   ├── suppliers/
│   │   │   │   ├── SupplierService.go
│   │   │   │   └── SupplierServiceCallback.go
│   │   │   ├── purchaseorders/
│   │   │   │   ├── PurchaseOrderService.go
│   │   │   │   └── PurchaseOrderServiceCallback.go
│   │   │   ├── stockmovements/
│   │   │   │   ├── StockMovementService.go
│   │   │   │   └── StockMovementServiceCallback.go
│   │   │   └── vehicleloads/
│   │   │       ├── VehicleLoadService.go
│   │   │       └── VehicleLoadServiceCallback.go
│   │   │
│   │   ├── dashboard/                     # Dashboard & KPIs
│   │   │   ├── kpis/
│   │   │   │   ├── VendKPIService.go
│   │   │   │   ├── VendKPIServiceCallback.go
│   │   │   │   └── kpi_threshold.go       # KPI status computation
│   │   │   └── dashboards/
│   │   │       ├── VendDashboardService.go
│   │   │       └── VendDashboardServiceCallback.go
│   │   │
│   │   ├── compliance/                    # Compliance & Health Inspections
│   │   │   ├── inspections/
│   │   │   │   ├── VendInspectionService.go
│   │   │   │   └── VendInspectionServiceCallback.go
│   │   │   ├── findings/
│   │   │   │   ├── VendInspectionFindingService.go
│   │   │   │   └── VendInspectionFindingServiceCallback.go
│   │   │   └── certifications/
│   │   │       ├── VendCertificationService.go
│   │   │       └── VendCertificationServiceCallback.go
│   │   │
│   │   ├── reports/                       # Scheduled Reports
│   │   │   └── reports/
│   │   │       ├── VendReportService.go
│   │   │       ├── VendReportServiceCallback.go
│   │   │       └── report_scheduler.go    # Background scheduler goroutine
│   │   │
│   │   ├── services/                      # Service activation
│   │   │   ├── activate_all.go
│   │   │   ├── activate_fleet.go
│   │   │   ├── activate_inventory.go
│   │   │   ├── activate_sales.go
│   │   │   ├── activate_payment.go
│   │   │   ├── activate_temperature.go
│   │   │   ├── activate_maintenance.go
│   │   │   ├── activate_route.go
│   │   │   ├── activate_analytics.go
│   │   │   ├── activate_access.go
│   │   │   ├── activate_dex.go
│   │   │   ├── activate_warehouse.go
│   │   │   ├── activate_dashboard.go
│   │   │   ├── activate_compliance.go
│   │   │   ├── activate_reports.go
│   │   │   └── activate_retention.go
│   │   │
│   │   ├── retention/                     # Data Retention & Archival
│   │   │   ├── archivedtxns/
│   │   │   │   ├── VendArchivedTransactionService.go
│   │   │   │   └── VendArchivedTransactionServiceCallback.go
│   │   │   ├── archivedtemps/
│   │   │   │   ├── VendArchivedTempReadingService.go
│   │   │   │   └── VendArchivedTempReadingServiceCallback.go
│   │   │   ├── archivedaccess/
│   │   │   │   ├── VendArchivedAccessEventService.go
│   │   │   │   └── VendArchivedAccessEventServiceCallback.go
│   │   │   └── scheduler/
│   │   │       └── retention_scheduler.go  # Background archival goroutine
│   │   │
│   │   ├── topology/                      # Fleet Topology Map
│   │   │   └── discovery.go               # VendMachine → L8TopologyNode conversion
│   │   │
│   │   ├── main/                          # Backend server entry point
│   │   │   ├── main.go
│   │   │   ├── build.sh
│   │   │   └── Dockerfile
│   │   │
│   │   ├── vnet/                          # Virtual network entry point
│   │   │   ├── main.go
│   │   │   ├── build.sh
│   │   │   └── Dockerfile
│   │   │
│   │   └── ui/                            # Web UI
│   │       ├── main.go                    # UI server + type registration
│   │       ├── build.sh
│   │       ├── Dockerfile
│   │       └── web/                       # Static web assets
│   │           ├── app.html
│   │           ├── login.html
│   │           ├── login.json
│   │           ├── l8ui/                  # Shared UI library (git submodule)
│   │           ├── js/
│   │           │   ├── sections.js
│   │           │   ├── app.js
│   │           │   └── reference-registry-vend.js
│   │           ├── sections/
│   │           │   ├── fleet.html
│   │           │   ├── inventory.html
│   │           │   ├── sales.html
│   │           │   ├── maintenance.html
│   │           │   ├── routes.html
│   │           │   ├── analytics.html
│   │           │   ├── warehouse.html
│   │           │   ├── dashboard.html
│   │           │   ├── compliance.html
│   │           │   └── reports.html
│   │           ├── vend-ui/               # Project-specific UI
│   │           │   ├── fleet/
│   │           │   │   ├── fleet-config.js
│   │           │   │   ├── fleet-init.js
│   │           │   │   └── machines/
│   │           │   │       ├── machines-enums.js
│   │           │   │       ├── machines-columns.js
│   │           │   │       └── machines-forms.js
│   │           │   ├── inventory/
│   │           │   ├── sales/
│   │           │   ├── maintenance/
│   │           │   ├── routes/
│   │           │   └── analytics/
│   │           └── m/                     # Mobile web assets
│   │               ├── app.html
│   │               └── js/
│   │
│   ├── types/                             # Generated protobuf types
│   │   └── vend/
│   │
│   ├── tests/
│   │   └── mocks/
│   │       ├── cmd/
│   │       ├── data.go
│   │       ├── store.go
│   │       ├── main_phases.go
│   │       ├── gen_fleet.go
│   │       ├── gen_inventory.go
│   │       ├── gen_sales.go
│   │       ├── gen_payment.go
│   │       ├── gen_maintenance.go
│   │       ├── gen_route.go
│   │       ├── gen_analytics.go
│   │       ├── gen_warehouse.go
│   │       ├── gen_temperature.go
│   │       ├── gen_access.go
│   │       ├── gen_dex.go
│   │       ├── gen_dashboard.go
│   │       ├── gen_compliance.go
│   │       └── gen_reports.go
│   │
│   └── k8s/
│       ├── deploy.sh
│       ├── undeploy.sh
│       ├── vend.yaml
│       ├── vend-web.yaml
│       ├── vend-vnet.yaml
│       ├── vend-collector.yaml
│       └── vend-parser.yaml
│
├── proto/
│   ├── make-bindings.sh
│   ├── vend-common.proto
│   ├── vend-fleet.proto
│   ├── vend-inventory.proto
│   ├── vend-sales.proto
│   ├── vend-payment.proto
│   ├── vend-temperature.proto
│   ├── vend-maintenance.proto
│   ├── vend-route.proto
│   ├── vend-analytics.proto
│   ├── vend-access.proto
│   ├── vend-dex.proto
│   ├── vend-warehouse.proto
│   ├── vend-dashboard.proto
│   ├── vend-compliance.proto
│   ├── vend-reports.proto
│   └── vend-retention.proto
│
└── plans/
    └── PRD-l8vendingmachine.md
```

---

## 4. Constants & Configuration

### defaults.go

```go
package common

const (
    VEND_VNET = 49010               // Virtual network ID
    PREFIX    = "/vend/"            // All API endpoints: /vend/{area}/{service}
)

var DB_CREDS = "postgres"
var DB_NAME  = "vendmachine"
```

### login.json

```json
{
    "login": {
        "appTitle": "L8 Vending",
        "appDescription": "Next-Gen AI Vending Machine Management",
        "authEndpoint": "/auth",
        "redirectUrl": "/app.html",
        "sessionTimeout": 30,
        "tfaEnabled": false
    },
    "app": {
        "dateFormat": "mm/dd/yyyy",
        "apiPrefix": "/vend",
        "healthPath": "/0/Health"
    }
}
```

---

## 5. Service Registry

All services use `ServiceArea = byte(10)` (single-module project).

| # | Submodule | Service Dir | ServiceName | Model (Proto Type) | Primary Key |
|---|-----------|-------------|-------------|-------------------|-------------|
| 1 | fleet | machines | `Machine` | `VendMachine` | `machineId` |
| 2 | fleet | machinegroups | `MachGrp` | `VendMachineGroup` | `groupId` |
| 3 | fleet | locations | `Location` | `VendLocation` | `locationId` |
| 4 | inventory | products | `Product` | `VendProduct` | `productId` |
| 5 | inventory | planograms | `Planogram` | `VendPlanogram` | `planogramId` |
| 6 | inventory | restockorders | `RstockOrd` | `VendRestockOrder` | `orderId` |
| 7 | sales | transactions | `Txn` | `VendTransaction` | `transactionId` |
| 8 | sales | settlements | `Settlemnt` | `VendSettlement` | `settlementId` |
| 9 | payment | cashpositions | `CashPos` | `VendCashPosition` | `positionId` |
| 10 | payment | collections | `CashColl` | `VendCashCollection` | `collectionId` |
| 11 | temperature | tempreadings | `TempRead` | `VendTempReading` | `readingId` |
| 12 | maintenance | alerts | `Alert` | `VendAlert` | `alertId` |
| 13 | maintenance | workorders | `WorkOrder` | `VendWorkOrder` | `workOrderId` |
| 14 | maintenance | servicevisits | `SvcVisit` | `VendServiceVisit` | `visitId` |
| 15 | route | routes | `Route` | `VendRoute` | `routeId` |
| 16 | route | drivers | `Driver` | `VendDriver` | `driverId` |
| 17 | analytics | forecasts | `Forecast` | `VendForecast` | `forecastId` |
| 18 | analytics | performance | `SlotPerf` | `VendSlotPerformance` | `performanceId` |
| 19 | access | accessevents | `AccEvent` | `VendAccessEvent` | `eventId` |
| 20 | dex | audits | `DexAudit` | `VendDexAudit` | `auditId` |
| 21 | warehouse | warehouses | `Warehouse` | `VendWarehouse` | `warehouseId` |
| 22 | warehouse | warehousestock | `WhseStock` | `VendWarehouseStock` | `stockId` |
| 23 | warehouse | suppliers | `Supplier` | `VendSupplier` | `supplierId` |
| 24 | warehouse | purchaseorders | `PurchOrd` | `VendPurchaseOrder` | `orderId` |
| 25 | warehouse | stockmovements | `StockMove` | `VendStockMovement` | `movementId` |
| 26 | warehouse | vehicleloads | `VehLoad` | `VendVehicleLoad` | `loadId` |
| 27 | analytics | fleetinventory | `FleetInv` | `VendFleetInventory` | `summaryId` |
| 28 | dashboard | kpis | `VendKpi` | `VendKPI` | `kpiId` |
| 29 | dashboard | dashboards | `VendDash` | `VendDashboard` | `dashboardId` |
| 30 | compliance | inspections | `Inspction` | `VendInspection` | `inspectionId` |
| 31 | compliance | findings | `InspFind` | `VendInspectionFinding` | `findingId` |
| 32 | compliance | certifications | `VendCert` | `VendCertification` | `certificationId` |
| 33 | reports | reports | `VendRpt` | `VendReport` | `reportId` |
| 34 | retention | archivedtxns | `ArchTxn` | `VendArchivedTransaction` | `transactionId` |
| 35 | retention | archivedtemps | `ArchTemp` | `VendArchivedTempReading` | `readingId` |
| 36 | retention | archivedaccess | `ArchAcces` | `VendArchivedAccessEvent` | `eventId` |

**Total: 36 Prime Object services**

**Immutable entities** (reject PUT via `rejectPut()` callback): `VendTransaction`, `VendAccessEvent`, `VendDexAudit`, `VendStockMovement`, `VendArchivedTransaction`, `VendArchivedTempReading`, `VendArchivedAccessEvent`

---

## 6. Protobuf Specifications

### vend-common.proto -- Shared Enums & Types

```protobuf
syntax = "proto3";
package vend;
option go_package = "./types/vend";

import "l8common.proto";

// === ENUMS ===

enum VendMachineType {
    VEND_MACHINE_TYPE_UNSPECIFIED = 0;
    VEND_MACHINE_TYPE_LOCKER = 1;
    VEND_MACHINE_TYPE_REFRIGERATED_BEVERAGE = 2;
    VEND_MACHINE_TYPE_COMBO = 3;
    VEND_MACHINE_TYPE_SNACK = 4;
    VEND_MACHINE_TYPE_FROZEN = 5;
    VEND_MACHINE_TYPE_FRESH_FOOD = 6;
}

enum VendMachineStatus {
    VEND_MACHINE_STATUS_UNSPECIFIED = 0;
    VEND_MACHINE_STATUS_OPERATIONAL = 1;
    VEND_MACHINE_STATUS_OUT_OF_SERVICE = 2;
    VEND_MACHINE_STATUS_MAINTENANCE = 3;
    VEND_MACHINE_STATUS_OFFLINE = 4;
    VEND_MACHINE_STATUS_DECOMMISSIONED = 5;
}

enum VendSlotStatus {
    VEND_SLOT_STATUS_UNSPECIFIED = 0;
    VEND_SLOT_STATUS_STOCKED = 1;
    VEND_SLOT_STATUS_LOW_STOCK = 2;
    VEND_SLOT_STATUS_SOLD_OUT = 3;
    VEND_SLOT_STATUS_EMPTY = 4;
    VEND_SLOT_STATUS_DISABLED = 5;
}

enum VendSlotMechanism {
    VEND_SLOT_MECHANISM_UNSPECIFIED = 0;
    VEND_SLOT_MECHANISM_SPRING_MOTOR = 1;
    VEND_SLOT_MECHANISM_ELECTRONIC_LOCK = 2;
    VEND_SLOT_MECHANISM_CONVEYOR = 3;
    VEND_SLOT_MECHANISM_GRAVITY = 4;
}

enum VendPaymentMethod {
    VEND_PAYMENT_METHOD_UNSPECIFIED = 0;
    VEND_PAYMENT_METHOD_CASH = 1;
    VEND_PAYMENT_METHOD_CREDIT_CARD = 2;
    VEND_PAYMENT_METHOD_NFC_CONTACTLESS = 3;
    VEND_PAYMENT_METHOD_QR_CODE = 4;
    VEND_PAYMENT_METHOD_MOBILE_WALLET = 5;
    VEND_PAYMENT_METHOD_PREPAID = 6;
    VEND_PAYMENT_METHOD_FREE = 7;
}

enum VendTransactionStatus {
    VEND_TRANSACTION_STATUS_UNSPECIFIED = 0;
    VEND_TRANSACTION_STATUS_COMPLETED = 1;
    VEND_TRANSACTION_STATUS_FAILED = 2;
    VEND_TRANSACTION_STATUS_REFUNDED = 3;
    VEND_TRANSACTION_STATUS_PENDING = 4;
}

enum VendAlertSeverity {
    VEND_ALERT_SEVERITY_UNSPECIFIED = 0;
    VEND_ALERT_SEVERITY_INFO = 1;
    VEND_ALERT_SEVERITY_WARNING = 2;
    VEND_ALERT_SEVERITY_CRITICAL = 3;
}

enum VendAlertCategory {
    VEND_ALERT_CATEGORY_UNSPECIFIED = 0;
    VEND_ALERT_CATEGORY_INVENTORY = 1;
    VEND_ALERT_CATEGORY_TEMPERATURE = 2;
    VEND_ALERT_CATEGORY_PAYMENT = 3;
    VEND_ALERT_CATEGORY_MECHANICAL = 4;
    VEND_ALERT_CATEGORY_CONNECTIVITY = 5;
    VEND_ALERT_CATEGORY_SECURITY = 6;
}

enum VendAlertStatus {
    VEND_ALERT_STATUS_UNSPECIFIED = 0;
    VEND_ALERT_STATUS_ACTIVE = 1;
    VEND_ALERT_STATUS_ACKNOWLEDGED = 2;
    VEND_ALERT_STATUS_RESOLVED = 3;
}

enum VendWorkOrderStatus {
    VEND_WORK_ORDER_STATUS_UNSPECIFIED = 0;
    VEND_WORK_ORDER_STATUS_OPEN = 1;
    VEND_WORK_ORDER_STATUS_ASSIGNED = 2;
    VEND_WORK_ORDER_STATUS_IN_PROGRESS = 3;
    VEND_WORK_ORDER_STATUS_COMPLETED = 4;
    VEND_WORK_ORDER_STATUS_CANCELLED = 5;
}

enum VendProductCategory {
    VEND_PRODUCT_CATEGORY_UNSPECIFIED = 0;
    VEND_PRODUCT_CATEGORY_COLD_BEVERAGE = 1;
    VEND_PRODUCT_CATEGORY_ENERGY_DRINK = 2;
    VEND_PRODUCT_CATEGORY_SNACK = 3;
    VEND_PRODUCT_CATEGORY_FRESH_FOOD = 4;
    VEND_PRODUCT_CATEGORY_CANDY = 5;
    VEND_PRODUCT_CATEGORY_WATER = 6;
    VEND_PRODUCT_CATEGORY_JUICE = 7;
    VEND_PRODUCT_CATEGORY_COFFEE = 8;
    VEND_PRODUCT_CATEGORY_HEALTH = 9;
}

enum VendTemperatureZoneType {
    VEND_TEMPERATURE_ZONE_TYPE_UNSPECIFIED = 0;
    VEND_TEMPERATURE_ZONE_TYPE_AMBIENT = 1;
    VEND_TEMPERATURE_ZONE_TYPE_REFRIGERATED = 2;
    VEND_TEMPERATURE_ZONE_TYPE_FROZEN = 3;
}

enum VendRouteStatus {
    VEND_ROUTE_STATUS_UNSPECIFIED = 0;
    VEND_ROUTE_STATUS_PLANNED = 1;
    VEND_ROUTE_STATUS_IN_PROGRESS = 2;
    VEND_ROUTE_STATUS_COMPLETED = 3;
    VEND_ROUTE_STATUS_CANCELLED = 4;
}

enum VendAccessEventType {
    VEND_ACCESS_EVENT_TYPE_UNSPECIFIED = 0;
    VEND_ACCESS_EVENT_TYPE_DOOR_OPEN = 1;
    VEND_ACCESS_EVENT_TYPE_DOOR_CLOSE = 2;
    VEND_ACCESS_EVENT_TYPE_SERVICE_VISIT = 3;
    VEND_ACCESS_EVENT_TYPE_VEND_DISPENSE = 4;
    VEND_ACCESS_EVENT_TYPE_TAMPER = 5;
}

// === SHARED TYPES (not Prime Objects) ===

message VendConnectivity {
    string connection_type = 1;        // "4G_LTE", "WIFI", "ETHERNET"
    string carrier = 2;
    int32 signal_strength = 3;         // dBm
    string ip_address = 4;
    bool wifi_backup = 5;
}

message VendGpsCoordinates {
    double latitude = 1;
    double longitude = 2;
}
```

### vend-fleet.proto -- Fleet Management

```protobuf
syntax = "proto3";
package vend;
option go_package = "./types/vend";

import "vend-common.proto";
import "l8common.proto";
import "api.proto";

// @PrimeObject
message VendMachine {
    string machine_id = 1;
    string serial_number = 2;
    string model = 3;
    string manufacturer = 4;
    VendMachineType machine_type = 5;
    VendMachineStatus status = 6;
    string firmware_version = 7;
    string controller_board = 8;
    int32 total_slots = 9;
    string location_id = 10;           // cross-ref: VendLocation
    string group_id = 11;             // cross-ref: VendMachineGroup
    VendConnectivity connectivity = 12;
    repeated string capabilities = 13;
    int64 installed_date = 14;
    int64 last_heartbeat = 15;
    int64 uptime = 16;
    string route_id = 17;             // cross-ref: VendRoute
    map<string, string> custom_fields = 18;
    l8common.AuditInfo audit_info = 19;
}

message VendMachineList {
    repeated VendMachine list = 1;
    l8api.L8MetaData metadata = 2;
}

// @PrimeObject
message VendMachineGroup {
    string group_id = 1;
    string name = 2;
    string description = 3;
    string region = 4;
    string operator_id = 5;
    int32 machine_count = 6;
    map<string, string> custom_fields = 7;
    l8common.AuditInfo audit_info = 8;
}

message VendMachineGroupList {
    repeated VendMachineGroup list = 1;
    l8api.L8MetaData metadata = 2;
}

// @PrimeObject
message VendLocation {
    string location_id = 1;
    string name = 2;
    l8common.Address address = 3;
    VendGpsCoordinates coordinates = 4;
    string location_type = 5;          // OFFICE, GYM, HOSPITAL, SCHOOL, etc.
    string timezone = 6;
    string contact_name = 7;
    string contact_phone = 8;
    string contact_email = 9;
    map<string, string> custom_fields = 10;
    l8common.AuditInfo audit_info = 11;
}

message VendLocationList {
    repeated VendLocation list = 1;
    l8api.L8MetaData metadata = 2;
}
```

### vend-inventory.proto -- Inventory Management

```protobuf
syntax = "proto3";
package vend;
option go_package = "./types/vend";

import "vend-common.proto";
import "l8common.proto";
import "api.proto";

// @PrimeObject
message VendProduct {
    string product_id = 1;            // SKU
    string name = 2;
    VendProductCategory category = 3;
    l8common.Money price = 4;
    string upc = 5;                   // UPC/EAN barcode
    string supplier_id = 6;
    int32 shelf_life_days = 7;
    bool is_active = 8;
    string image_url = 9;
    string description = 10;
    map<string, string> custom_fields = 11;
    l8common.AuditInfo audit_info = 12;
}

message VendProductList {
    repeated VendProduct list = 1;
    l8api.L8MetaData metadata = 2;
}

// @PrimeObject -- Product-to-slot assignment for a machine
message VendPlanogram {
    string planogram_id = 1;
    string machine_id = 2;            // cross-ref: VendMachine
    string name = 3;
    bool is_active = 4;
    int64 effective_date = 5;
    repeated VendSlotAssignment slots = 6;  // Child: embedded slot assignments
    map<string, string> custom_fields = 7;
    l8common.AuditInfo audit_info = 8;
}

// Child type (NOT a Prime Object) -- embedded in VendPlanogram
message VendSlotAssignment {
    string slot_id = 1;               // e.g., "A01", "B05"
    string row = 2;
    int32 column = 3;
    string product_id = 4;            // cross-ref: VendProduct
    int32 capacity = 5;
    int32 par_level = 6;
    int32 current_quantity = 7;
    VendSlotStatus status = 8;
    VendSlotMechanism mechanism = 9;
    l8common.Money price = 10;        // Price override (null = use product price)
    string motor_status = 11;         // OK, ATTENTION, FAILED
    string lock_status = 12;          // LOCKED, UNLOCKED, FAULT (locker only)
    int64 last_vend_time = 13;
    int32 vend_count = 14;
    int64 sold_out_since = 15;
    int64 expiration_date = 16;
}

// @PrimeObject
message VendRestockOrder {
    string order_id = 1;
    string machine_id = 2;            // cross-ref: VendMachine
    string route_id = 3;              // cross-ref: VendRoute
    VendWorkOrderStatus status = 4;
    int64 created_date = 5;
    int64 due_date = 6;
    string urgency = 7;               // LOW, MODERATE, HIGH, CRITICAL
    repeated VendRestockLine lines = 8;  // Child: items to restock
    string notes = 9;
    map<string, string> custom_fields = 10;
    l8common.AuditInfo audit_info = 11;
}

// Child type -- embedded in VendRestockOrder
message VendRestockLine {
    string slot_id = 1;
    string product_id = 2;            // cross-ref: VendProduct
    int32 quantity_needed = 3;
    int32 quantity_loaded = 4;
}

message VendPlanogramList {
    repeated VendPlanogram list = 1;
    l8api.L8MetaData metadata = 2;
}

message VendRestockOrderList {
    repeated VendRestockOrder list = 1;
    l8api.L8MetaData metadata = 2;
}
```

### vend-sales.proto -- Sales & Transactions

```protobuf
syntax = "proto3";
package vend;
option go_package = "./types/vend";

import "vend-common.proto";
import "l8common.proto";
import "api.proto";

// @PrimeObject
message VendTransaction {
    string transaction_id = 1;
    string machine_id = 2;            // cross-ref: VendMachine
    int64 timestamp = 3;
    string slot_id = 4;
    string product_id = 5;            // cross-ref: VendProduct
    string product_name = 6;
    VendProductCategory category = 7;
    l8common.Money price = 8;
    VendPaymentMethod payment_method = 9;
    VendTransactionStatus status = 10;
    string card_type = 11;            // VISA, MASTERCARD, etc.
    string card_last_four = 12;
    l8common.Money cash_inserted = 13;
    l8common.Money change_given = 14;
    double dispense_duration = 15;    // seconds
    string error_code = 16;
    string location_id = 17;          // cross-ref: VendLocation
    map<string, string> custom_fields = 18;
    l8common.AuditInfo audit_info = 19;
}

message VendTransactionList {
    repeated VendTransaction list = 1;
    l8api.L8MetaData metadata = 2;
}

// @PrimeObject
message VendSettlement {
    string settlement_id = 1;
    string machine_id = 2;            // cross-ref: VendMachine
    int64 settlement_date = 3;
    int64 period_start = 4;
    int64 period_end = 5;
    int32 transaction_count = 6;
    l8common.Money total_amount = 7;
    l8common.Money card_amount = 8;
    l8common.Money nfc_amount = 9;
    l8common.Money qr_amount = 10;
    string processor_reference = 11;
    string status = 12;               // PENDING, SETTLED, FAILED
    map<string, string> custom_fields = 13;
    l8common.AuditInfo audit_info = 14;
}

message VendSettlementList {
    repeated VendSettlement list = 1;
    l8api.L8MetaData metadata = 2;
}
```

### Remaining Proto Files (summarized)

**vend-payment.proto**: `VendCashPosition` (coin tubes, bill stacker, cash box totals per machine), `VendCashCollection` (cash collection event with amounts by denomination).

**vend-temperature.proto**: `VendTempReading` (zone, current temp, setpoint, compressor status, duty cycle, min/max, compliance status per machine per zone).

**vend-maintenance.proto**: `VendAlert` (severity, category, code, description, threshold, current value, status), `VendWorkOrder` (type, priority, assigned driver, parts needed, estimated duration), `VendServiceVisit` (machine, driver, activities performed, duration, parts used).

**vend-route.proto**: `VendRoute` (name, status, driver, vehicle, stops as embedded `VendRouteStop` children, total distance/duration), `VendDriver` (name, phone, license, vehicle, home base location).

**vend-analytics.proto**: `VendForecast` (machine, product, forecast horizon, predicted daily vends, predicted stockout time, restock urgency, confidence score), `VendSlotPerformance` (machine, slot, product, vend count, revenue, velocity, rank, margin, stockout hours), `VendFleetInventory` (product ref, product name, category, unit price, total machines carrying, total slots assigned, total units in machines, total capacity across machines, total units in warehouses, total supply chain position, fleet-wide sold out count, fleet-wide low stock count, last updated timestamp). Computed/refreshed periodically by aggregating across all VendPlanogram slot assignments and VendWarehouseStock records.

**vend-access.proto**: `VendAccessEvent` (machine, event type, door/cell ID, timestamp, duration, operator ID, activities).

**vend-dex.proto**: `VendDexAudit` (machine, audit timestamp, DEX version, interval/cumulative vend data, cash audit totals, embedded `VendDexSelectionAudit` children, event log).

**vend-dashboard.proto**: `VendKPI` (code, name, category, unit, calculation formula, current/target/previous value, status ON_TARGET/AT_RISK/OFF_TARGET, trend UP/DOWN/FLAT, refresh interval, embedded `VendKPIThreshold` children with operator/value/severity). `VendDashboard` (name, description, layout config JSON, refresh interval, owner, embedded `VendDashboardWidget` children with widget type, chart type, position, size, query, data source reference).

**vend-compliance.proto**: `VendInspection` (machine or location ref, inspection type enum HEALTH_DEPARTMENT/INTERNAL_AUDIT/FOOD_SAFETY/EQUIPMENT_SAFETY, status, planned/actual dates, inspector name, scope, embedded `VendInspectionReport` children). `VendInspectionFinding` (inspection ref, finding number, title, severity, status, condition/criteria/cause/effect/recommendation/management response fields, responsible party, due date, evidence document, repeat finding flag, embedded `VendRemediationAction` children). `VendCertification` (machine or location ref, certification type, standard, certifying body, certificate number, issue/expiry/renewal dates, status ACTIVE/PENDING/EXPIRED/REVOKED, scope).

**vend-reports.proto**: `VendReport` (code, name, description, report type, category, owner, query, default format PDF/CSV/EXCEL, is public, execution count, embedded `VendReportSchedule` children with frequency DAILY/WEEKLY/MONTHLY/QUARTERLY, run time, day of week/month, delivery email, output format, is active, next/last run timestamps. Embedded `VendReportExecution` children with status, start/end time, row count, file size, output path, error message. Embedded `VendReportSubscription` children with subscriber, format, delivery email).

**vend-retention.proto**: `VendArchivedTransaction` (copy of VendTransaction fields + `archivedAt`, `archivedBy`, `archiveReason`; immutable -- rejects PUT). `VendArchivedTempReading` (copy of VendTempReading fields + archive metadata; immutable). `VendArchivedAccessEvent` (copy of VendAccessEvent fields + archive metadata; immutable). `VendRetentionPolicy` (entity type, archive after days, delete archived after days, is active, last run timestamp -- embedded child of system config, NOT a separate service).

**vend-warehouse.proto**: `VendWarehouse` (name, address, GPS, capacity sqft, operating hours, contact, region). `VendWarehouseStock` (warehouse ref, product ref, quantity on hand, reorder point, reorder quantity, last counted date, expiration tracking). `VendSupplier` (name, contact, address, lead time days, payment terms, status). `VendPurchaseOrder` (supplier ref, warehouse ref, status lifecycle DRAFT→SUBMITTED→CONFIRMED→SHIPPED→RECEIVED→CLOSED, order date, expected delivery, embedded `VendPurchaseOrderLine` children with product/quantity/unit price/received quantity). `VendStockMovement` (warehouse ref, product ref, movement type enum RECEIVE_FROM_SUPPLIER/TRANSFER_TO_VEHICLE/RETURN_FROM_VEHICLE/WRITE_OFF/ADJUSTMENT, quantity, reference ID linking to PO or vehicle load, timestamp, performed by). `VendVehicleLoad` (route ref, driver ref, vehicle ID, load date, status LOADING→IN_TRANSIT→COMPLETED, embedded `VendVehicleLoadLine` children with product/quantity loaded/quantity returned/quantity dispensed).

---

## 7. UI Module Organization

### Desktop Modules

| Module Tab | Sub-modules | Services |
|-----------|-------------|----------|
| Fleet | Machines | machines, machine-groups, locations |
| Inventory | Products, Planograms | products, planograms, restock-orders |
| Sales | Transactions, Settlements | transactions, settlements |
| Maintenance | Alerts, Work Orders, Service Visits | alerts, work-orders, service-visits |
| Routes | Routes, Drivers | routes, drivers |
| Analytics | Forecasts, Performance | forecasts, slot-performance |
| Warehouse | Stock, Suppliers, Orders | warehouses, stock, suppliers, purchase-orders, movements, vehicle-loads |
| Dashboard | KPIs, Overview | kpis, dashboards (home page with KPI cards + charts) |
| Compliance | Inspections, Findings, Certs | inspections, findings, certifications |
| Reports | Report Definitions | reports (with schedule management) |

### Desktop Section HTML IDs

Following the `{moduleKey}-{serviceKey}-table-container` pattern:

| Container ID | Module | Service |
|-------------|--------|---------|
| `machines-machines-table-container` | fleet/machines | Machine |
| `machines-machine-groups-table-container` | fleet/machines | MachineGroup |
| `machines-locations-table-container` | fleet/machines | Location |
| `products-products-table-container` | inventory/products | Product |
| `products-planograms-table-container` | inventory/products | Planogram |
| `products-restock-orders-table-container` | inventory/products | RestockOrder |
| `transactions-transactions-table-container` | sales/transactions | Transaction |
| `alerts-alerts-table-container` | maintenance/alerts | Alert |
| `alerts-work-orders-table-container` | maintenance/work-orders | WorkOrder |
| `routes-routes-table-container` | route/routes | Route |
| `forecasts-forecasts-table-container` | analytics/forecasts | Forecast |

### View Types by Service

| Service | Views | Notes |
|---------|-------|-------|
| Machine | table, chart | Chart: machines by type/status |
| Transaction | table, chart, timeline | Chart: revenue by day; Timeline: transaction history |
| Alert | table, kanban | Kanban: by severity (INFO/WARNING/CRITICAL) |
| WorkOrder | table, kanban, gantt | Kanban: by status; Gantt: scheduled work |
| Route | table, gantt | Gantt: route schedule |
| Forecast | table, chart | Chart: predicted vs actual demand |
| SlotPerformance | table, chart | Chart: revenue by product |
| FleetInventory | table, chart | Chart: stock by category, machines per product |
| VendKPI | (dashboard cards) | Layer8DWidget KPI cards with sparklines |
| VendInspection | table, calendar | Calendar: scheduled inspections by date |
| VendInspectionFinding | table, kanban | Kanban: by status (OPEN/IN_PROGRESS/CLOSED) |
| VendReport | table | With schedule management inline |

---

## 8. Mock Data Generation

### Phase Ordering

```
Phase 1: Foundation (no dependencies)
    - Locations (10 locations)
    - Machine Groups (5 groups)
    - Products (30 products)
    - Drivers (8 drivers)
    - Warehouses (3 warehouses)
    - Suppliers (5 suppliers)

Phase 2: Core Entities (depend on Phase 1)
    - Machines (10 machines, reference locations + groups)
    - Routes (12 routes, reference drivers)
    - Warehouse Stock (90 records: 30 products x 3 warehouses)

Phase 3: Configuration (depend on Phase 2)
    - Planograms (10 planograms, one per machine, reference products)
    - Purchase Orders (8 POs, reference suppliers + warehouses)

Phase 4: Operational Data (depend on Phase 2-3)
    - Transactions (500 transactions across machines)
    - Cash Positions (10 cash positions, one per machine)
    - Temp Readings (100 readings across machines)
    - Access Events (50 events)
    - DEX Audits (10 audits, one per machine)
    - Stock Movements (40 movements: receives, transfers, returns)
    - Vehicle Loads (12 loads, one per route)

Phase 5: Analytics & Maintenance (depend on Phase 4)
    - Alerts (30 alerts across machines)
    - Work Orders (10 work orders)
    - Service Visits (15 visits)
    - Cash Collections (10 collections)
    - Settlements (20 settlements)
    - Forecasts (50 forecasts)
    - Slot Performance (100 performance records)
    - Restock Orders (8 orders)
    - Fleet Inventory (30 summaries, one per product, computed from planograms + warehouse stock)
    - KPIs (10 KPI records with computed values)
    - Dashboard (1 default dashboard with 8 widgets)
    - Inspections (5 inspection records, mix of completed and scheduled)
    - Inspection Findings (8 findings across inspections)
    - Certifications (10 certifications, some near expiry)
    - Reports (10 report definitions with schedules)
```

### Data Arrays (data.go)

```go
var MachineModels = []string{
    "TCN-ZK(22SP)+BLH-64S", "TCN-ZK(22SP)+BLH-40S",
    "AF-60C(22SP)", "AF-D900-54C(22SP)",
}

var LocationTypes = []string{
    "OFFICE_LOBBY", "GYM", "HOSPITAL", "SCHOOL", "AIRPORT",
    "HOTEL", "MALL", "FACTORY", "UNIVERSITY", "TRAIN_STATION",
}

var ProductNames = []string{
    "Coca-Cola 355ml", "Pepsi 355ml", "Sprite 355ml", "Dr Pepper 355ml",
    "Aquafina 500ml", "Dasani 500ml", "Gatorade 591ml",
    "Red Bull 250ml", "Monster Energy 473ml", "Celsius 355ml",
    "Snickers Bar 52g", "Doritos 28g", "Lay's Classic 28g",
    "Kit Kat 42g", "M&M's 49g", "Trail Mix 40g",
    "Turkey Sandwich", "Caesar Salad", "Chicken Wrap",
    // ...
}
```

---

## 9. Deployment

### Docker Images (5 images)

| Image | Directory | Base Image | K8s Kind |
|-------|-----------|------------|----------|
| `saichler/vendmachine` | `go/vend/main/` | `saichler/erp-postgres` | StatefulSet |
| `saichler/vendmachine-web` | `go/vend/ui/` | `saichler/erp-security` | DaemonSet |
| `saichler/vendmachine-vnet` | `go/vend/vnet/` | `saichler/erp-security` | DaemonSet |
| `saichler/vendmachine-collector` | `go/vend/collector/` | `saichler/erp-security` | DaemonSet |
| `saichler/vendmachine-parser` | `go/vend/parser/` | `saichler/erp-security` | DaemonSet |

### K8s Manifests

Each manifest follows the l8erp pattern with:
- Namespace with labels
- Resource labels (`app: vendmachine`)
- NODE_IP env var from `status.hostIP`
- Volume name `hdata` mounting `/data`

### build.sh (per image)

```bash
#!/usr/bin/env bash
set -e
docker build --no-cache --platform=linux/amd64 -t saichler/vendmachine:latest .
docker push saichler/vendmachine:latest
```

---

## 10. Local Development Setup

### run-local.sh

```bash
#!/bin/bash
set -e

rm -rf go.sum go.mod vendor
go mod init
GOPROXY=direct GOPRIVATE=github.com go mod tidy
go mod vendor

# Start database
docker rm -f unsecure-postgres 2>/dev/null || true
docker run -d --name unsecure-postgres -p 5432:5432 \
    -v /data/:/data/ saichler/unsecure-postgres:latest admin admin admin 5432

# Build binaries
rm -rf demo && mkdir -p demo
cd tests/mocks/cmd && go build -o ../../../demo/mocks_demo && cd ../../../
cd vend/vnet && go build -o ../../demo/vnet_demo && cd ../../
cd vend/main && go build -o ../../demo/vend_demo && cd ../../
cd vend/ui/main && go build -o ../../../demo/ui_demo && cd ../../../
cd vend/collector && go build -o ../../demo/collector_demo && cd ../../
cd vend/parser && go build -o ../../demo/parser_demo && cd ../../
cp -r vend/ui/web demo/.

# Generate kill script
cd demo
cat > kill_demo.sh <<'EOF'
cd ..
rm -rf demo
rm -rf /data/postgres/admin
pkill -9 demo
EOF
chmod +x kill_demo.sh

# Start services
./vnet_demo &
sleep 1
./vend_demo local &
./ui_demo &
./collector_demo &
./parser_demo &
sleep 8

# Upload mock data
EXTERNAL_IP=$(ip route get 1.1.1.1 | grep -oP 'src \K[0-9.]+')
read -p "Press Enter to upload mocks"
./mocks_demo --address https://${EXTERNAL_IP}:2773 --user admin --password admin --insecure

read -p "Press Enter to kill the demo"
./kill_demo.sh
```

---

## 11. Collector & Parser Architecture

The l8vendingmachine project uses the existing Layer 8 collector/parser framework (`l8collector`, `l8parser`, `l8pollaris`) to poll vending machines via REST API and parse responses into protobuf entities. This follows the same architecture as **probler** (network device collection) and uses the existing `RestCollector` protocol and `RestJsonParse` / `RestGpuParse` parsing rules.

### Reference Projects

| Project | Role | What We Reuse |
|---------|------|---------------|
| `../l8collector` | Job-based polling engine | `RestCollector` protocol, `HostCollector` lifecycle, `JobsQueue` cadence |
| `../l8parser` | Rule-based response parsing | `RestJsonParse` rule, `ParsingService` activation |
| `../l8pollaris` | Poll configuration registry | `L8Pollaris` configs, `L8Poll` job definitions, cadence plans |
| `../probler` | Reference implementation | `Links.go` pattern, `collector/main.go`, `parser/main.go`, `addPollConfig` command |

### Data Flow

```
l8opensim (100 vending machines on 192.168.100.x:8443)
    │
    │  HTTPS REST API
    ▼
┌─────────────────────────────────────────────┐
│  Vending Collector (go/vend/collector/)      │
│                                              │
│  l8collector.RestCollector                   │
│    ├── Poll: GET::/api/v1/machine::          │
│    ├── Poll: GET::/api/v1/inventory::        │
│    ├── Poll: GET::/api/v1/transactions::     │
│    ├── Poll: GET::/api/v1/temperature::      │
│    ├── Poll: GET::/api/v1/alerts::           │
│    ├── Poll: GET::/api/v1/payment/cashbox::  │
│    ├── Poll: GET::/api/v1/payment/status::   │
│    ├── Poll: GET::/api/v1/analytics/traffic::│
│    ├── Poll: GET::/api/v1/analytics/health:: │
│    └── Poll: GET::/api/v1/dex/audit::        │
│                                              │
│  Per-machine HostCollector with JobsQueue    │
│  Cadence: 30s (transactions), 5min (temp),   │
│           15min (inventory), 1hr (analytics), │
│           24hr (DEX audit)                    │
└──────────────────┬──────────────────────────┘
                   │ CJob results (raw JSON)
                   ▼
┌─────────────────────────────────────────────┐
│  Vending Parser (go/vend/parser/)            │
│                                              │
│  l8parser.RestJsonParse rules                │
│    ├── machine → VendMachine properties      │
│    ├── inventory.slots → VendPlanogram slots  │
│    ├── transactions → VendTransaction records │
│    ├── temperature.zones → VendTempReading    │
│    ├── alerts → VendAlert records             │
│    ├── cashbox → VendCashPosition             │
│    ├── payment.status → VendMachine.payment*  │
│    ├── analytics.traffic → VendSlotPerformance│
│    ├── analytics.health → VendMachine.health* │
│    └── dex.audit → VendDexAudit records       │
│                                              │
│  Parsed entities stored via vend services     │
└──────────────────┬──────────────────────────┘
                   │ Parsed VendMachine / VendTempReading / etc.
                   ▼
┌─────────────────────────────────────────────┐
│  Threshold Evaluator (in parser After hook)  │
│                                              │
│  Compares parsed values against thresholds:  │
│    temp > 8C → PostEvent(PERFORMANCE,        │
│                VEND_TEMP_HIGH, CRITICAL)      │
│    stock < par → PostEvent(PERFORMANCE,      │
│                  VEND_STOCK_LOW, WARNING)     │
│    cashbox > 90% → PostEvent(PERFORMANCE,    │
│                    VEND_CASH_BOX_FULL, WARN)  │
│    heartbeat gap > 10m → PostEvent(SYSTEM,   │
│                    VEND_MACHINE_OFFLINE, CRIT)│
│                                              │
│  Uses l8events.PostEvent() API               │
└──────────────────┬──────────────────────────┘
                   │ EventRecord
                   ▼
┌─────────────────────────────────────────────┐
│  l8alarms (../l8alarms)                      │
│                                              │
│  AlarmDefinition matching → Alarm creation   │
│  Correlation engine → root cause analysis    │
│  Notification engine → l8notify dispatch     │
│  Escalation scheduler → time-based escalate  │
└──────────────────┬──────────────────────────┘
                   │ NotifyTarget
                   ▼
┌─────────────────────────────────────────────┐
│  l8notify (../l8notify)                      │
│                                              │
│  EMAIL → operator/manager inboxes            │
│  WEBHOOK → ticketing systems, external APIs  │
│  SLACK → #vending-alerts channel             │
│  CUSTOM → SMS gateway, push notifications    │
└─────────────────────────────────────────────┘
```

### Project Structure (Collector/Parser additions)

```
go/vend/
├── collector/                         # Vending machine collector
│   └── main.go                        # Entry point (follows probler/go/prob/collector/main.go)
│
├── parser/                            # Vending machine parser
│   ├── main.go                        # Entry point (follows probler/go/prob/parser/main.go)
│   ├── boot/                          # Poll configurations
│   │   ├── cadences.go                # Cadence plan constants
│   │   ├── vend_machine_polls.go      # Machine identity polls
│   │   ├── vend_inventory_polls.go    # Inventory polls
│   │   ├── vend_sales_polls.go        # Transaction polls
│   │   ├── vend_monitor_polls.go      # Temperature, energy, alerts polls
│   │   ├── vend_payment_polls.go      # Cash position, payment status polls
│   │   ├── vend_analytics_polls.go    # Traffic, health, performance polls
│   │   └── vend_dex_polls.go          # DEX audit polls
│   └── threshold/                     # Threshold evaluation
│       └── evaluator.go               # Compares parsed values, posts l8events
│
├── common/
│   ├── defaults.go                    # PREFIX, VNET constants (existing)
│   ├── Links.go                       # VendMachine_Links_ID, service name/area mappings
│   ├── resources.go                   # CreateResources helper
│   └── commands/
│       ├── addMachine.go              # Register machine as L8PTarget for polling
│       └── addPollConfig.go           # Register vending pollaris configurations
```

### Links.go (following probler pattern)

```go
package common

const (
    Collector_Service_Name = "VColl"
    Collector_Service_Area = byte(0)

    VendMachine_Links_ID          = "Vend"
    Vend_Cache_Service_Name       = "VCache"
    Vend_Cache_Service_Area       = byte(0)
    Vend_Persist_Service_Name     = "VPersist"
    Vend_Persist_Service_Area     = byte(0)
    Vend_Parser_Service_Name      = "VPars"
    Vend_Parser_Service_Area      = byte(0)
    Vend_Model_Name               = "vendmachine"
)

type Links struct{}

func (this *Links) Collector(linkid string) (string, byte) {
    return Collector_Service_Name, Collector_Service_Area
}

func (this *Links) Parser(linkid string) (string, byte) {
    return Vend_Parser_Service_Name, Vend_Parser_Service_Area
}

func (this *Links) Cache(linkid string) (string, byte) {
    return Vend_Cache_Service_Name, Vend_Cache_Service_Area
}

func (this *Links) Persist(linkid string) (string, byte) {
    return Vend_Persist_Service_Name, Vend_Persist_Service_Area
}

func (this *Links) Model(linkid string) string {
    return Vend_Model_Name
}
```

### Collector main.go (following probler/go/prob/collector/main.go)

```go
package main

import (
    "github.com/saichler/l8bus/go/overlay/vnic"
    "github.com/saichler/l8collector/go/collector/common"
    "github.com/saichler/l8collector/go/collector/service"
    "github.com/saichler/l8pollaris/go/pollaris"
    "github.com/saichler/l8types/go/ifs"
    vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
)

func main() {
    common.SmoothFirstCollection = true
    res := vendcommon.CreateResources("vend-collector")
    ifs.SetNetworkMode(ifs.NETWORK_K8s)
    nic := vnic.NewVirtualNetworkInterface(res, nil)
    nic.Start()
    nic.WaitForConnection()

    pollaris.Activate(nic)
    service.Activate(vendcommon.VendMachine_Links_ID, nic)
    res.Logger().SetLogLevel(ifs.Error_Level)
    vendcommon.WaitForSignal(res)
}
```

### Parser main.go (following probler/go/prob/parser/main.go)

```go
package main

import (
    "github.com/saichler/l8bus/go/overlay/vnic"
    "github.com/saichler/l8parser/go/parser/service"
    "github.com/saichler/l8pollaris/go/pollaris"
    "github.com/saichler/l8types/go/ifs"
    vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
    vendtypes "github.com/saichler/l8vendingmachine/go/types/vend"
)

func main() {
    resources := vendcommon.CreateResources("vend-parser")
    ifs.SetNetworkMode(ifs.NETWORK_K8s)
    nic := vnic.NewVirtualNetworkInterface(resources, nil)
    nic.Start()
    nic.WaitForConnection()

    pollaris.Activate(nic)

    // Activate vending machine parser
    service.Activate(vendcommon.VendMachine_Links_ID,
        &vendtypes.VendMachine{}, false, nic, "MachineId")

    vendcommon.WaitForSignal(resources)
}
```

### Pollaris Configuration (Poll Definitions)

Each poll uses the `RestCollector` protocol with `"GET::/endpoint::"` format and `RestJsonParse` rule for mapping JSON fields to protobuf properties.

#### cadences.go

```go
package boot

import "github.com/saichler/l8pollaris/go/types/l8tpollaris"

var EVERY_30_SECONDS = &l8tpollaris.L8PCadencePlan{
    Enabled: true, Cadences: []int64{30},
}
var EVERY_5_MINUTES = &l8tpollaris.L8PCadencePlan{
    Enabled: true, Cadences: []int64{300},
}
var EVERY_15_MINUTES = &l8tpollaris.L8PCadencePlan{
    Enabled: true, Cadences: []int64{900},
}
var EVERY_1_HOUR = &l8tpollaris.L8PCadencePlan{
    Enabled: true, Cadences: []int64{3600},
}
var EVERY_24_HOURS = &l8tpollaris.L8PCadencePlan{
    Enabled: true, Cadences: []int64{86400},
}
var DEFAULT_TIMEOUT int64 = 30
```

#### Poll Definitions

| Poll Name | Endpoint | Cadence | Target Entity | JSON→Property Mapping |
|-----------|----------|---------|---------------|----------------------|
| `vendMachineIdentity` | `GET::/api/v1/machine::` | 15 min | VendMachine | `serialNumber:vendmachine.serialnumber`, `model:vendmachine.model`, `firmwareVersion:vendmachine.firmwareversion`, `totalSlots:vendmachine.totalslots`, `connectivity.signalStrength:vendmachine.connectivity.signalstrength` |
| `vendMachineStatus` | `GET::/api/v1/machine/status::` | 30 sec | VendMachine | `status:vendmachine.status`, `doorStatus:vendmachine.doorstatus`, `lastHeartbeat:vendmachine.lastheartbeat` |
| `vendInventory` | `GET::/api/v1/inventory::` | 15 min | VendPlanogram | `slots:vendplanogram.slots` (array parse) |
| `vendInventoryAlerts` | `GET::/api/v1/inventory/alerts::` | 5 min | VendAlert | Creates VendAlert records from inventory alerts |
| `vendTransactions` | `GET::/api/v1/transactions::` | 30 sec | VendTransaction | `transactions:vendtransaction` (array parse, each element → new record) |
| `vendTransactionSummary` | `GET::/api/v1/transactions/summary::` | 5 min | VendMachine | `totalVends:vendmachine.totalvendstoday`, `totalRevenue:vendmachine.totalrevenuetoday` |
| `vendTemperature` | `GET::/api/v1/temperature::` | 5 min | VendTempReading | `zones:vendtempreading` (one reading per zone) |
| `vendAlerts` | `GET::/api/v1/alerts::` | 30 sec | VendAlert | `activeAlerts:vendalert` (array parse) |
| `vendCashbox` | `GET::/api/v1/payment/cashbox::` | 15 min | VendCashPosition | `totalCashToCollect:vendcashposition.totalcash`, `coinTubes:vendcashposition.cointubes`, `billStacker:vendcashposition.billstacker` |
| `vendPaymentStatus` | `GET::/api/v1/payment/status::` | 5 min | VendMachine | `overallStatus:vendmachine.paymentstatus`, peripheral status fields |
| `vendTraffic` | `GET::/api/v1/analytics/traffic::` | 1 hr | VendSlotPerformance | `totalApproaches:vendslotperformance.approaches`, `conversionRate:vendslotperformance.conversionrate` |
| `vendHealth` | `GET::/api/v1/analytics/health::` | 1 hr | VendMachine | `overallHealthScore:vendmachine.healthscore`, component health arrays |
| `vendDexAudit` | `GET::/api/v1/dex/audit::` | 24 hr | VendDexAudit | Full DEX data mapping |

#### Example Poll Creation (vend_sales_polls.go)

```go
package boot

import "github.com/saichler/l8pollaris/go/types/l8tpollaris"

func CreateVendTransactionsPoll(p *l8tpollaris.L8Pollaris) {
    poll := &l8tpollaris.L8Poll{}
    poll.Name = "vendTransactions"
    poll.What = "GET::/api/v1/transactions::"
    poll.Protocol = l8tpollaris.L8PProtocol_L8PRESTAPI
    poll.Cadence = EVERY_30_SECONDS
    poll.Timeout = DEFAULT_TIMEOUT
    poll.Always = true  // Always send, even if hash matches (transactions are append-only)
    poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
    poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
    poll.Attributes = append(poll.Attributes, createVendRestAttribute(
        "vendtransaction",
        "transactions:vendtransaction"))
    p.Polling[poll.Name] = poll
}

func CreateVendTransactionSummaryPoll(p *l8tpollaris.L8Pollaris) {
    poll := &l8tpollaris.L8Poll{}
    poll.Name = "vendTransactionSummary"
    poll.What = "GET::/api/v1/transactions/summary::"
    poll.Protocol = l8tpollaris.L8PProtocol_L8PRESTAPI
    poll.Cadence = EVERY_5_MINUTES
    poll.Timeout = DEFAULT_TIMEOUT
    poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
    poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
    poll.Attributes = append(poll.Attributes, createVendRestAttribute(
        "vendmachine.salessummary",
        "totalVends:vendmachine.totalvendstoday,"+
            "totalRevenue:vendmachine.totalrevenuetoday,"+
            "averageTransactionValue:vendmachine.avgtransactionvalue"))
    p.Polling[poll.Name] = poll
}

// createVendRestAttribute creates a RestJsonParse attribute for vending data.
func createVendRestAttribute(propertyId, mapping string) *l8tpollaris.L8PAttribute {
    attr := &l8tpollaris.L8PAttribute{}
    attr.PropertyId = map[string]string{"vendmachine": propertyId}
    attr.Rules = make([]*l8tpollaris.L8PRule, 0)
    rule := &l8tpollaris.L8PRule{}
    rule.Name = "RestJsonParse"
    rule.Params = make(map[string]*l8tpollaris.L8PParameter)
    rule.Params["mapping"] = &l8tpollaris.L8PParameter{Value: mapping}
    attr.Rules = append(attr.Rules, rule)
    return attr
}
```

#### Master Pollaris Builder (boot/vend_pollaris.go)

```go
package boot

import "github.com/saichler/l8pollaris/go/types/l8tpollaris"

func GetVendingMachinePollaris() *l8tpollaris.L8Pollaris {
    p := &l8tpollaris.L8Pollaris{}
    p.Name = "VendingMachine"
    p.Polling = make(map[string]*l8tpollaris.L8Poll)

    // Machine identity & status (15min / 30sec)
    CreateVendMachineIdentityPoll(p)
    CreateVendMachineStatusPoll(p)

    // Inventory (15min / 5min for alerts)
    CreateVendInventoryPoll(p)
    CreateVendInventoryAlertsPoll(p)

    // Sales (30sec / 5min for summary)
    CreateVendTransactionsPoll(p)
    CreateVendTransactionSummaryPoll(p)

    // Temperature (5min)
    CreateVendTemperaturePoll(p)

    // Alerts (30sec)
    CreateVendAlertsPoll(p)

    // Payment (15min / 5min)
    CreateVendCashboxPoll(p)
    CreateVendPaymentStatusPoll(p)

    // Analytics (1hr)
    CreateVendTrafficPoll(p)
    CreateVendHealthPoll(p)

    // DEX Audit (24hr)
    CreateVendDexAuditPoll(p)

    return p
}
```

### addMachine Command (following probler addDevice pattern)

```go
package commands

import (
    "github.com/saichler/l8pollaris/go/types/l8tpollaris"
    "github.com/saichler/l8web/go/web/client"
    vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
)

// AddVendingMachine registers a machine IP as an L8PTarget for REST API polling.
func AddVendingMachine(rc *client.RestClient, ip string, port int32) {
    target := &l8tpollaris.L8PTarget{}
    target.Host = ip
    target.LinksId = vendcommon.VendMachine_Links_ID
    target.State = l8tpollaris.L8PTargetState_L8P_TARGET_POLL_ACTIVE

    hostProto := &l8tpollaris.L8PHostProtocol{}
    hostProto.Addr = ip
    hostProto.Port = port
    hostProto.Protocol = l8tpollaris.L8PProtocol_L8PRESTAPI
    hostProto.Ainfo = &l8tpollaris.L8PAuthInfo{}  // No auth for simulator

    target.Protocols = []*l8tpollaris.L8PHostProtocol{hostProto}
    target.PollarisName = "VendingMachine"

    rc.POST("0/Target", "L8PTarget", "", "", target)
}
```

### Deployment (Collector + Parser)

Two additional Docker images:

| Image | Directory | K8s Kind | Notes |
|-------|-----------|----------|-------|
| `saichler/vendmachine-collector` | `go/vend/collector/` | DaemonSet | Polls vending machine REST APIs |
| `saichler/vendmachine-parser` | `go/vend/parser/` | DaemonSet | Parses JSON responses into VendMachine entities |

These are deployed alongside the existing 3 images (vend, vend-web, vend-vnet), bringing the total to **5 Docker images**.

### Integration with l8opensim

1. Start l8opensim with 100 vending machines:
   ```bash
   cd ../l8opensim/go/simulator
   sudo ./simulator -auto-start-ip 192.168.100.1 -auto-count 10
   ```

2. Start l8vendingmachine services (vnet, main, ui, collector, parser)

3. Register poll configs:
   ```go
   // In mock data setup or CLI command
   commands.AddPollConfigs(rc, resources)  // Registers VendingMachine pollaris
   ```

4. Register each machine as an L8PTarget:
   ```go
   for i := 1; i <= 10; i++ {
       ip := fmt.Sprintf("192.168.100.%d", i)
       commands.AddVendingMachine(rc, ip, 8443)
   }
   ```

5. The collector automatically starts polling each machine at the configured cadences. The parser transforms JSON responses into protobuf entities and stores them via the vend services.

---

## 12. Implementation Phases

### Phase 1: Project Bootstrap
- Create project directory structure
- Set up `go.mod` with Layer 8 dependencies
- Add l8ui as git submodule (`setup-l8ui-submodule.sh`)
- Create `defaults.go` with PREFIX and VNET constants
- Create `login.json` with correct `apiPrefix`
- Create `make-bindings.sh`

### Phase 2: Protobuf Definitions
- Create all 16 proto files (`vend-common.proto` through `vend-retention.proto`)
- Run `make-bindings.sh` to generate Go types
- Verify generated `.pb.go` files in `go/types/vend/`

### Phase 3: Foundation Services
- Implement `VendLocation`, `VendMachineGroup`, `VendProduct`, `VendDriver`
- Implement `VendWarehouse`, `VendSupplier`
- Service files + Callback files for each
- Type registration in `ui/main.go`
- Activation functions

### Phase 4: Core Services
- Implement `VendMachine`, `VendRoute`, `VendPlanogram`
- Implement `VendWarehouseStock`, `VendPurchaseOrder`
- These reference Phase 3 entities

### Phase 5: Operational Services
- Implement `VendTransaction`, `VendCashPosition`, `VendTempReading`
- Implement `VendAccessEvent`, `VendDexAudit`
- Implement `VendStockMovement`, `VendVehicleLoad`

### Phase 6: Management Services
- Implement `VendAlert`, `VendWorkOrder`, `VendServiceVisit`
- Implement `VendCashCollection`, `VendSettlement`, `VendRestockOrder`

### Phase 7: Analytics, Dashboard, Compliance & Reports
- Implement `VendForecast`, `VendSlotPerformance`, `VendFleetInventory`
- `VendFleetInventory` uses `.Compute(computeFleetAggregates)` callback (follows `../l8erp/go/erp/bi/kpis/kpi_threshold.go` pattern) to aggregate across all VendPlanogram slot assignments + VendWarehouseStock records on POST. Refresh triggered by retention scheduler or manual API call.
- Implement `VendKPI` with `kpi_threshold.go` status computation (follows `../l8erp/go/erp/bi/kpis/kpi_threshold.go`)
- Implement `VendDashboard` with widget layout persistence
- Implement `VendInspection`, `VendInspectionFinding`, `VendCertification` (follows `../l8erp/go/erp/comp/` patterns)
- Implement `VendReport` with `report_scheduler.go` background goroutine (follows `../l8erp/go/erp/bi/reports/report_scheduler.go`)
- Seed 10 KPI definitions, default dashboard layout, 10 report definitions with schedules
- Implement `VendArchivedTransaction`, `VendArchivedTempReading`, `VendArchivedAccessEvent` (all immutable -- `rejectPut()` callback, follows `../l8alarms` archived entity pattern)
- Implement `retention_scheduler.go` background goroutine with configurable retention policies
- Implement `topology/discovery.go` -- `ActivateVendTopology()` + `ConvertToTopologyNode()` (follows `../l8topology/go/topo/discover/Layer1.go`)
- Register vending machine SVG icon via `Layer8SvgFactory.registerTemplate('vendingMachine', ...)`
- Activate l8agent backend: call `l8agent.Initialize()` + `l8agent.InitializeChat()` after introspector populated
- Verify immutable entities reject PUT: `VendTransaction`, `VendAccessEvent`, `VendDexAudit`, `VendStockMovement`

### Phase 8: Collector & Parser
- Create `go/vend/common/Links.go` with VendMachine_Links_ID and service mappings
- Create `go/vend/common/resources.go` with CreateResources helper
- Create `go/vend/parser/boot/cadences.go` with cadence plan constants
- Create 7 poll definition files (`vend_machine_polls.go` through `vend_dex_polls.go`)
- Create `go/vend/parser/boot/vend_pollaris.go` master builder
- Create `go/vend/collector/main.go` (follows `probler/go/prob/collector/main.go`)
- Create `go/vend/parser/main.go` (follows `probler/go/prob/parser/main.go`)
- Create `go/vend/common/commands/addMachine.go` and `addPollConfig.go`
- Create `go/vend/parser/threshold/evaluator.go` -- threshold evaluation logic
- Test: start simulator + collector + parser, register machines, verify data flows

### Phase 8b: Threshold Events, Alarms & Notifications
- Add `l8events`, `l8alarms`, `l8notify` as dependencies in `go.mod`
- Wire threshold evaluator into parser After hook to call `l8events.PostEvent()` on threshold crossings
- Create seed data for 10 `AlarmDefinition` records (VEND_TEMP_HIGH, VEND_STOCK_LOW, etc.)
- Create seed data for `CorrelationRule` records (compressor→temp, payment→vend failure, offline→suppress all)
- Create seed data for `NotificationPolicy` records (route driver webhook, ops manager email, Slack channel)
- Create seed data for `EscalationPolicy` records (0min→driver, 15min→supervisor, 60min→manager)
- Create seed data for `MaintenanceWindow` templates (restock window, quarterly maintenance)
- Create seed data for `AlarmFilter` records (Critical Only, My Route, Temperature, Unacknowledged)
- Activation: Add `l8alarms.ActivateAlmServices()` and `l8events` activation to `run-local.sh` startup
- Test: trigger a threshold (e.g., high temp from simulator), verify event→alarm→notification flow

### Phase 9: Desktop UI
- Module configs, section HTML, enums, columns, forms for all services
- View configurations (kanban for alerts, gantt for routes, charts for analytics)
- Reference registry
- Map visualization: fleet map view with l8topology map rendering, color-coded machine markers
- Dashboard: KPI cards with Layer8DWidget, sparkline trends, charts
- Remove `Layer8DModuleFilter.load()` from `app.js` (l8vendingmachine does not have ModConfig service -- see `modconfig-failure-no-logout` rule)
- Include l8agent chat JS/CSS in `app.html` (l8agent-chat.js, l8agent-bubble.js, l8agent-enums.js, l8agent-columns.js, l8agent-forms.js + CSS)
- Fleet topology map: integrate l8topology map rendering with custom vending machine SVG icons

### Phase 10: Mobile UI
- Mobile module data files (enums, columns, forms)
- Mobile nav config
- Mobile reference registry

### Phase 11: Mock Data Generators
- `data.go` with product/location/model arrays
- `store.go` with ID slices
- Generator files per phase
- Phase orchestration in `main_phases.go`

### Phase 12: Deployment Artifacts
- `build.sh` + `Dockerfile` for each image (5 total: vend, vend-web, vend-vnet, vend-collector, vend-parser)
- K8s manifests (5 YAMLs)
- `build-all-images.sh`
- `run-local.sh` (starts all 5 services)

### Phase 13: End-to-End Verification
- Build all binaries
- Run `run-local.sh`
- Upload mock data
- Verify all sections load in desktop and mobile
- Verify CRUD operations on each service
- Verify view types (kanban, gantt, chart) render correctly
- Verify search, filter, sort, pagination work

---

## Traceability Matrix

| # | Gap / Action Item | Phase |
|---|-------------------|-------|
| 1 | Project directory structure | Phase 1 |
| 2 | go.mod + l8ui submodule | Phase 1 |
| 3 | defaults.go (PREFIX, VNET) | Phase 1 |
| 4 | login.json adaptation | Phase 1 |
| 5 | 15 proto files (incl. warehouse, dashboard, compliance, reports) | Phase 2 |
| 6 | make-bindings.sh | Phase 2 |
| 7 | Foundation services (6: Location, MachineGroup, Product, Driver, Warehouse, Supplier) | Phase 3 |
| 8 | Core services (5: Machine, Route, Planogram, WarehouseStock, PurchaseOrder) | Phase 4 |
| 9 | Operational services (7: Transaction, CashPosition, TempReading, AccessEvent, DexAudit, StockMovement, VehicleLoad) | Phase 5 |
| 10 | Management services (6) | Phase 6 |
| 11 | Analytics services (3: Forecast, SlotPerformance, FleetInventory) | Phase 7 |
| 11b | Dashboard services (2: VendKPI with thresholds, VendDashboard) | Phase 7 |
| 11c | Compliance services (3: Inspection, Finding, Certification) | Phase 7 |
| 11d | Reports service (VendReport with scheduler) | Phase 7 |
| 11e | Map visualization (fleet map view with status markers) | Phase 9 |
| 11f | Archived entity services (3: Transaction, TempReading, AccessEvent) | Phase 7 |
| 11g | Retention scheduler (retention_scheduler.go) | Phase 7 |
| 11h | Topology discovery (VendMachine → L8TopologyNode) | Phase 7 |
| 11i | L8AgentChat backend activation (Initialize + InitializeChat) | Phase 7 |
| 11j | L8AgentChat frontend includes in app.html + m/app.html | Phase 9/10 |
| 11k | Immutable entity callbacks (rejectPut on 4 active + 3 archived entities) | Phase 3-7 |
| 11l | VNet architecture: all 5 binaries share VEND_VNET=49010 | Phase 1 |
| 11m | FleetInventory Compute callback for aggregate refresh | Phase 7 |
| 11n | Missing section HTML files (warehouse, dashboard, compliance, reports) | Phase 9 |
| 12 | Links.go (VendMachine_Links_ID, service mappings) | Phase 8 |
| 13 | Pollaris poll definitions (13 polls, 7 boot files) | Phase 8 |
| 14 | Collector main.go (RestCollector activation) | Phase 8 |
| 15 | Parser main.go (RestJsonParse activation) | Phase 8 |
| 16 | addMachine + addPollConfig commands | Phase 8 |
| 17 | Collector→Parser→Service data flow verification | Phase 8 |
| 17b | Threshold evaluator (parser After hook) | Phase 8b |
| 17c | AlarmDefinition seed data (10 rules) | Phase 8b |
| 17d | CorrelationRule seed data (3 rules) | Phase 8b |
| 17e | NotificationPolicy seed data | Phase 8b |
| 17f | EscalationPolicy seed data | Phase 8b |
| 17g | MaintenanceWindow templates | Phase 8b |
| 17h | AlarmFilter seed data | Phase 8b |
| 17i | l8events + l8alarms + l8notify activation in startup | Phase 8b |
| 17j | End-to-end threshold→event→alarm→notification test | Phase 8b |
| 18 | Desktop UI (all modules) | Phase 9 |
| 19 | Mobile UI (all modules) | Phase 10 |
| 20 | Mock data generators | Phase 11 |
| 21 | Docker images (5: vend, web, vnet, collector, parser) | Phase 12 |
| 22 | K8s manifests (5) | Phase 12 |
| 23 | run-local.sh (starts all 5 services) | Phase 12 |
| 24 | build-all-images.sh | Phase 12 |
| 25 | End-to-end verification | Phase 13 |
| 26 | Desktop/mobile parity check | Phase 13 |
| 27 | l8opensim integration test (10 machines) | Phase 13 |

---

## End-to-End Verification (Phase 13)

For every section:
1. Navigate to the section in desktop
2. Verify table data loads (not blank)
3. Verify row click opens detail/modal
4. Verify detail content is populated
5. Verify Add/Edit/Delete operations work
6. Verify on both desktop and mobile

Sections to verify:
- [ ] Fleet > Machines
- [ ] Fleet > Machine Groups
- [ ] Fleet > Locations
- [ ] Inventory > Products
- [ ] Inventory > Planograms
- [ ] Inventory > Restock Orders
- [ ] Sales > Transactions
- [ ] Sales > Settlements
- [ ] Maintenance > Alerts (+ kanban view)
- [ ] Maintenance > Work Orders (+ kanban + gantt views)
- [ ] Maintenance > Service Visits
- [ ] Routes > Routes (+ gantt view)
- [ ] Routes > Drivers
- [ ] Analytics > Forecasts (+ chart view)
- [ ] Analytics > Slot Performance (+ chart view)
- [ ] Analytics > Fleet Inventory (+ chart view: stock by category)
- [ ] Warehouse > Warehouses
- [ ] Warehouse > Warehouse Stock
- [ ] Warehouse > Suppliers
- [ ] Warehouse > Purchase Orders
- [ ] Warehouse > Stock Movements
- [ ] Warehouse > Vehicle Loads
- [ ] Dashboard > KPI cards (verify sparklines, trends, ON_TARGET/AT_RISK/OFF_TARGET status)
- [ ] Dashboard > Overview layout (KPI row + charts + tables)
- [ ] Map > Fleet map view (markers with color-coded status)
- [ ] Map > Click marker shows machine summary popup
- [ ] Compliance > Inspections (+ calendar view)
- [ ] Compliance > Inspection Findings (+ kanban view by status)
- [ ] Compliance > Certifications (verify expiry alerts)
- [ ] Reports > Report definitions
- [ ] Reports > Schedule management (create/edit/toggle schedules)
- [ ] Reports > Verify scheduler fires at configured time and delivers email

### Immutability Verification
- [ ] PUT to VendTransaction returns error "immutable"
- [ ] PUT to VendAccessEvent returns error "immutable"
- [ ] PUT to VendDexAudit returns error "immutable"
- [ ] PUT to VendStockMovement returns error "immutable"
- [ ] PUT to VendArchivedTransaction returns error "immutable"
- [ ] PUT to VendArchivedTempReading returns error "immutable"
- [ ] PUT to VendArchivedAccessEvent returns error "immutable"
- [ ] UI for immutable entities shows no Edit/Save buttons (read-only table mode)

### Retention & Archival Verification
- [ ] Retention scheduler runs at configured time
- [ ] Transactions older than 90 days are archived to VendArchivedTransaction
- [ ] Archived transactions are queryable via VendArchivedTransaction service
- [ ] Original records are deleted from VendTransaction after archival

### Topology & Agent Verification
- [ ] Fleet map shows vending machine markers at correct GPS coordinates
- [ ] Map markers are color-coded by machine status
- [ ] Click marker shows machine summary popup
- [ ] L8AgentChat bubble appears in bottom-right corner (desktop + mobile)
- [ ] Agent can answer "How many machines are online?" via L8Query

### Collector & Parser Verification
- [ ] Start l8opensim with 100 vending machines
- [ ] Start collector + parser services
- [ ] Register VendingMachine pollaris config
- [ ] Register 5 machines as L8PTargets
- [ ] Verify collector connects and starts polling (check logs)
- [ ] Verify parser receives jobs and creates VendMachine records
- [ ] Verify VendTransaction records appear after 30s poll cycle
- [ ] Verify VendTempReading records appear after 5min poll cycle
- [ ] Verify VendAlert records appear from alert polling
- [ ] Verify VendCashPosition records appear from cashbox polling
- [ ] Verify change detection works (unchanged data not re-parsed)
- [ ] Scale to 10 machines, verify no thundering herd (SmoothFirstCollection)

### Threshold, Alarm & Notification Verification
- [ ] Verify threshold evaluator detects simulated temperature > 8C and posts EventRecord
- [ ] Verify AlarmDefinition VEND_TEMP_HIGH matches the event and creates an Alarm
- [ ] Verify alarm state is ACTIVE with correct severity (CRITICAL)
- [ ] Verify correlation engine links "Compressor Failure" as root cause to "Temperature High"
- [ ] Verify symptom alarm (Temperature High) is suppressed when correlated
- [ ] Verify NotificationPolicy dispatches email/webhook for CRITICAL alarm
- [ ] Verify EscalationPolicy fires Step 2 after 15 minutes if alarm not acknowledged
- [ ] Acknowledge alarm via API, verify escalation cancels
- [ ] Verify auto-clear: when temperature drops back to normal, alarm state changes to CLEARED
- [ ] Verify MaintenanceWindow suppresses alarms during scheduled restock window
- [ ] Verify AlarmFilter "Critical Only" shows only CRITICAL severity alarms
- [ ] Verify throttling: same alarm on same machine within cooldown does not trigger duplicate notification
