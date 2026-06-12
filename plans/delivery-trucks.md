# Plan: Add VendDeliveryTruck Entity

## Overview

Add `VendDeliveryTruck` as a Prime Object to the Routes module. Trucks are the vehicles that restock machines, collect cash, and follow planned routes. The `milesPerGallon` field will feed route optimization cost calculations.

---

## Phase 1: Protobuf

### 1.1 Add enums to `proto/vend-common.proto`

```protobuf
enum VendTruckStatus {
    VEND_TRUCK_STATUS_UNSPECIFIED = 0;
    VEND_TRUCK_STATUS_ACTIVE = 1;
    VEND_TRUCK_STATUS_MAINTENANCE = 2;
    VEND_TRUCK_STATUS_EN_ROUTE = 3;
    VEND_TRUCK_STATUS_DECOMMISSIONED = 4;
}

enum VendFuelType {
    VEND_FUEL_TYPE_UNSPECIFIED = 0;
    VEND_FUEL_TYPE_GASOLINE = 1;
    VEND_FUEL_TYPE_DIESEL = 2;
    VEND_FUEL_TYPE_ELECTRIC = 3;
    VEND_FUEL_TYPE_HYBRID = 4;
}

enum VendTruckType {
    VEND_TRUCK_TYPE_UNSPECIFIED = 0;
    VEND_TRUCK_TYPE_BOX_TRUCK = 1;
    VEND_TRUCK_TYPE_CARGO_VAN = 2;
    VEND_TRUCK_TYPE_REFRIGERATED = 3;
    VEND_TRUCK_TYPE_SPRINTER = 4;
    VEND_TRUCK_TYPE_PICKUP = 5;
}
```

### 1.2 Add message to `proto/vend-route.proto`

Append after `VendDriverList`:

```protobuf
// @PrimeObject
message VendDeliveryTruck {
    string truck_id = 1;
    string plate_number = 2;
    string vin = 3;
    string name = 4;

    // Vehicle Specs
    string make = 5;
    string model = 6;
    int32 year = 7;
    VendTruckType type = 8;
    int32 cargo_capacity_cu_ft = 9;
    int32 max_payload_lbs = 10;
    VendFuelType fuel_type = 11;
    int32 mileage = 12;
    double miles_per_gallon = 13;

    // Status
    VendTruckStatus status = 14;
    string current_driver_id = 15;       // cross-ref: VendDriver
    string current_route_id = 16;        // cross-ref: VendRoute
    string home_depot_id = 17;           // cross-ref: VendWarehouse

    // Location
    double last_latitude = 18;
    double last_longitude = 19;
    int64 last_location_update = 20;

    // Maintenance
    int64 last_maintenance_date = 21;
    int64 next_maintenance_date = 22;
    int32 next_maintenance_mileage = 23;
    int64 insurance_expiry = 24;
    int64 registration_expiry = 25;

    // Capacity Tracking
    bool refrigeration_equipped = 26;
    bool cash_collection_equipped = 27;
    bool coin_changer_equipped = 28;

    // Audit
    map<string, string> custom_fields = 29;
    l8common.AuditInfo audit_info = 30;
}

message VendDeliveryTruckList {
    repeated VendDeliveryTruck list = 1;
    l8api.L8MetaData metadata = 2;
}
```

### 1.3 Regenerate bindings

```bash
cd proto && ./make-bindings.sh
```

Verify: `grep "type VendDeliveryTruck struct" go/types/vend/*.pb.go`

---

## Phase 2: Backend Service

### 2.1 Create service directory

`go/vend/route/trucks/`

### 2.2 TruckService.go

Follow the `LocationService.go` pattern exactly:
- `ServiceName = "Truck"` (5 chars, under 10)
- `ServiceArea = byte(10)` (same as routes/drivers)
- `PrimaryKey = "TruckId"`
- Functions: `Activate()`, `Trucks()`, `Truck(truckId)`

### 2.3 TruckServiceCallback.go

Follow the `LocationServiceCallback.go` pattern:
- Auto-generate ID on POST: `common.GenerateID(&entity.TruckId)`
- Require: `TruckId`

### 2.4 Wire into activate_route.go

Add `trucks.Activate(creds, dbname, nic)` to `collectRouteActivations()`.

### 2.5 Register type in ui/shared.go

Add:
```go
common.RegisterType(resources, &vend.VendDeliveryTruck{}, &vend.VendDeliveryTruckList{}, "TruckId")
```

### 2.6 Build verification

```bash
cd go && go build ./...
```

---

## Phase 3: Mock Data

### 3.1 Add data arrays to `go/tests/mocks/data_vend.go`

```go
var TruckMakes = []string{"Ford", "Freightliner", "Isuzu", "Mercedes-Benz", "Ram", "Chevrolet", "GMC", "Hino"}

var TruckModels = []string{
    "E-450 Cutaway", "M2 106", "NPR-HD", "Sprinter 3500",
    "ProMaster 3500", "Express 4500", "Savana 3500", "195",
}

var TruckNames = []string{
    "Truck-NYC-01", "Truck-NYC-02", "Truck-LAX-01", "Truck-CHI-01",
    "Truck-HOU-01", "Truck-PHX-01", "Truck-SFO-01", "Truck-DAL-01",
    "Truck-ATL-01", "Truck-MIA-01",
}
```

### 3.2 Add ID slice to `go/tests/mocks/store.go`

```go
TruckIDs []string
```

### 3.3 Add generator to `go/tests/mocks/gen_fleet.go`

`func generateTrucks(store *MockDataStore) []*vend.VendDeliveryTruck` — 10 trucks with:
- Randomized makes/models/years (2018-2025)
- `milesPerGallon` between 8.0 and 22.0
- Status distribution: 60% Active, 20% En-Route, 10% Maintenance, 10% Decommissioned
- Reference `store.WarehouseIDs` for `homeDepotId`
- Random cargo capacity (400-800 cu ft), payload (3000-8000 lbs)
- Random boolean capabilities

### 3.4 Wire into `main_phases.go`

Add to `generateBusinessFoundation()` after suppliers:

```go
trucks := generateTrucks(store)
tids := extractIDs(trucks, func(v interface{}) string { return v.(*vend.VendDeliveryTruck).TruckId })
if err := runOp(client, "Trucks", "/vend/10/Truck",
    &vend.VendDeliveryTruckList{List: trucks}, tids, &store.TruckIDs); err != nil {
    return err
}
```

Add to `PrintSummary()`: `fmt.Printf("Trucks:           %d\n", len(store.TruckIDs))`

---

## Phase 4: Desktop UI

### 4.1 Update routes-config.js

Add truck service to the routes module:
```js
svc('trucks', 'Trucks', '', '/10/Truck', 'VendDeliveryTruck')
```

### 4.2 Update routes section config

Add to the services list:
```js
{ key: 'trucks', label: 'Trucks', icon: '🚛' }
```

### 4.3 Add truck enums/columns/forms to `vend-ui/routes/planning/`

Since the routes module uses `RoutePlanning` as its submodule namespace, add truck definitions to the existing `planning-enums.js`, `planning-columns.js`, and `planning-forms.js` files.

**Enums** — add to `planning-enums.js`:
- `TRUCK_STATUS`: Unspecified/Active/Maintenance/En-Route/Decommissioned
- `FUEL_TYPE`: Unspecified/Gasoline/Diesel/Electric/Hybrid
- `TRUCK_TYPE`: Unspecified/Box Truck/Cargo Van/Refrigerated/Sprinter/Pickup
- Renderers: `truckStatus` (status badge), `fuelType` (enum), `truckType` (enum)
- Add to `primaryKeys`: `VendDeliveryTruck: 'truckId'`

**Columns** — add `VendDeliveryTruck` to `planning-columns.js`:
```
truckId, name, plateNumber, type, make, model, year, status,
milesPerGallon, mileage, cargoCapacityCuFt, fuelType,
refrigerationEquipped, cashCollectionEquipped
```

**Forms** — add `VendDeliveryTruck` to `planning-forms.js`:
- Section "Identity": name, plateNumber, vin
- Section "Vehicle Specs": make, model, year, type, fuelType, cargoCapacityCuFt, maxPayloadLbs, mileage, milesPerGallon
- Section "Status": status, currentDriverId (ref VendDriver), currentRouteId (ref VendRoute), homeDepotId (ref VendWarehouse)
- Section "Maintenance": lastMaintenanceDate, nextMaintenanceDate, nextMaintenanceMileage, insuranceExpiry, registrationExpiry
- Section "Capabilities": refrigerationEquipped, cashCollectionEquipped, coinChangerEquipped

### 4.4 Update reference-registry-vend.js

Add:
```js
...ref.simple('VendDeliveryTruck', 'truckId', 'name', 'Truck'),
```

### 4.5 No changes needed to app.html, sections.js, or section HTML

The truck is added as a new service inside the existing Routes module — no new scripts or sections needed.

---

## Phase 5: Mobile UI

### 5.1 Update mobile routes enums/columns/forms

Add `VendDeliveryTruck` definitions to the mobile routes files under `m/js/routes/` (same field set as desktop, add `primary: true` on `name`, `secondary: true` on `status`).

### 5.2 Update mobile nav config

Add truck service entry to the routes module services:
```js
{ key: 'trucks', label: 'Trucks', icon: 'truck',
  endpoint: '/10/Truck', model: 'VendDeliveryTruck', idField: 'truckId' }
```

---

## Phase 6: Verification

1. `cd proto && ./make-bindings.sh` — bindings generate cleanly
2. `cd go && go build ./...` — compiles with no errors
3. Run `./run-local.sh clean` — full restart with mock data
4. Desktop: Navigate to Routes → Trucks tab → verify table shows 10 trucks with correct columns
5. Desktop: Click a truck row → verify detail popup shows all fields
6. Desktop: Add a new truck → verify form saves correctly
7. Mobile: Navigate to Routes → Trucks → verify card view with truck data
8. Mobile: Tap a truck → verify detail view

---

## Traceability Matrix

| # | Item | Phase |
|---|------|-------|
| 1 | Proto enums (TruckStatus, FuelType, TruckType) | Phase 1.1 |
| 2 | Proto message VendDeliveryTruck + List | Phase 1.2 |
| 3 | Regenerate bindings | Phase 1.3 |
| 4 | TruckService.go (ServiceName, Activate, helpers) | Phase 2.2 |
| 5 | TruckServiceCallback.go (auto-ID, validation) | Phase 2.3 |
| 6 | Wire into activate_route.go | Phase 2.4 |
| 7 | Register type in ui/shared.go | Phase 2.5 |
| 8 | Mock data arrays | Phase 3.1 |
| 9 | Store ID slice | Phase 3.2 |
| 10 | Generator function | Phase 3.3 |
| 11 | Wire into main_phases.go | Phase 3.4 |
| 12 | Desktop: routes-config.js service entry | Phase 4.1 |
| 13 | Desktop: section config service entry | Phase 4.2 |
| 14 | Desktop: enums/columns/forms definitions | Phase 4.3 |
| 15 | Desktop: reference registry entry | Phase 4.4 |
| 16 | Mobile: enums/columns/forms | Phase 5.1 |
| 17 | Mobile: nav config entry | Phase 5.2 |
| 18 | End-to-end verification | Phase 6 |
