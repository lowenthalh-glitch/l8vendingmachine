# Vending Machine Simulation for L8OpenSim

## Overview

Add 4 vending machine models to the l8opensim simulator, creating a new **"Vending Machines"** device category. Each machine exposes a comprehensive REST API suitable for a next-generation AI-powered vending management system. 100 simulated machines will be supported via round-robin cycling across the 4 models.

## Target Models

| # | Resource Name | Model | Type | Slots | Key Features |
|---|---------------|-------|------|-------|-------------|
| 1 | `tcn_zk_blh_64s` | TCN-ZK(22SP)+BLH-64S | Locker vending | 64 cells | Master+slave lockers, individual electronic locks, any product shape |
| 2 | `tcn_zk_blh_40s` | TCN-ZK(22SP)+BLH-40S | Locker vending | 40 cells | Same architecture, fewer cells |
| 3 | `afen_60c` | AF-60C(22SP) | Refrigerated beverage | 60 slots | Cold drinks, compressor cooling, 4-25C range |
| 4 | `afen_d900_54c` | AF-D900-54C(22SP) | Combo (snack+drink) | 54 slots | Dual-zone: ambient snacks + refrigerated drinks |

All models feature the **22SP Smart Platform**: 21.5" Android touchscreen (Rockchip RK3288), 4G/WiFi connectivity, MDB bus for payment peripherals, DEX/UCS audit data support.

---

## Protocol Selection

**Primary: HTTPS REST API** -- The 22SP machines communicate with cloud management platforms via HTTPS. The simulator will expose a REST API matching the patterns used by modern vending telemetry platforms (Cantaloupe, Nayax, TCN Yunshu). This is the same approach l8opensim uses for storage devices (Pure Storage, NetApp, Dell EMC, AWS S3).

**Secondary: SNMP** -- Basic machine health OIDs for integration with network management systems that may monitor vending infrastructure alongside network devices.

**No SSH** -- Vending machines do not expose SSH terminals.

---

## REST API Design

All endpoints are prefixed with `/api/v1`. Responses follow the envelope pattern used by existing l8opensim API resources.

### API Endpoint Catalog

#### 1. Machine Identity & Status

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/machine` | Machine identity, model, serial, firmware, uptime |
| GET | `/api/v1/machine/status` | Operational status, connectivity, last heartbeat |
| GET | `/api/v1/machine/config` | Machine configuration (locale, currency, timezone) |

#### 2. Inventory Management

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/inventory` | Full inventory: all slots with stock levels, product mapping, capacity |
| GET | `/api/v1/inventory/slots/{slotId}` | Single slot detail |
| GET | `/api/v1/inventory/alerts` | Low-stock and sold-out alerts |
| GET | `/api/v1/planogram` | Product-to-slot mapping (planogram) |

#### 3. Sales & Transactions

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/transactions` | Recent transactions (itemized: product, price, payment method, result) |
| GET | `/api/v1/transactions/summary` | Aggregated sales: total revenue, vend count, by payment method |
| GET | `/api/v1/transactions/{txnId}` | Single transaction detail |

#### 4. Payment Systems

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/payment/status` | Status of all payment peripherals (coin, bill, card, NFC, QR) |
| GET | `/api/v1/payment/cashbox` | Cash position: coins by denomination, bills, total, change available |
| GET | `/api/v1/payment/cashless` | Cashless reader status, last transaction, connectivity |

#### 5. Temperature & Refrigeration

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/temperature` | Current readings per zone, setpoints, compressor status |
| GET | `/api/v1/temperature/history` | Recent temperature history (last 24h data points) |

#### 6. Alerts & Errors

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/alerts` | Active alerts/alarms (temp, motor, payment, door, connectivity) |
| GET | `/api/v1/alerts/history` | Historical alert log |
| GET | `/api/v1/errors` | Error codes with descriptions and timestamps |

#### 7. Door & Access Events

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/access/events` | Door open/close events, duration, service visit correlation |
| GET | `/api/v1/access/locks` | Lock status per cell (locker models only) |

#### 8. Energy & Environment

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/energy` | Power draw, compressor cycles, energy-saving mode, kWh |
| GET | `/api/v1/environment` | Ambient temperature, humidity, light level |

#### 9. AI/Analytics Data

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/analytics/traffic` | Foot traffic counts by hour (from PIR/vision sensor) |
| GET | `/api/v1/analytics/performance` | Per-slot sales velocity, revenue, margin ranking |
| GET | `/api/v1/analytics/demand` | Demand forecast inputs: sales history by hour/day/week patterns |
| GET | `/api/v1/analytics/health` | Predictive maintenance signals: motor current trends, compressor health |

#### 10. DEX Audit Data

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/dex/audit` | Full DEX/UCS audit data export (structured JSON equivalent of EVA-DTS) |

### Total: 24 API endpoints per machine

---

## Response Schemas

### GET /api/v1/machine
```json
{
  "machineId": "VM-TCN-001",
  "serialNumber": "TCN2024ZK064001",
  "model": "TCN-ZK(22SP)+BLH-64S",
  "manufacturer": "TCN Vending Technology",
  "machineType": "LOCKER",
  "firmwareVersion": "22SP-3.8.2",
  "controllerBoard": "RK3288-SP22",
  "androidVersion": "7.1.2",
  "screenSize": "21.5 inch",
  "totalSlots": 64,
  "installedDate": "2025-06-15T00:00:00Z",
  "uptime": 1847293,
  "location": {
    "name": "Building A - Lobby",
    "address": "123 Commerce Blvd, Austin, TX 78701",
    "latitude": 30.2672,
    "longitude": -97.7431,
    "locationType": "OFFICE_LOBBY",
    "timezone": "America/Chicago"
  },
  "connectivity": {
    "type": "4G_LTE",
    "carrier": "T-Mobile",
    "signalStrength": -67,
    "ipAddress": "10.0.1.101",
    "wifiBackup": true,
    "wifiSSID": "VendNet-5G"
  },
  "capabilities": ["LOCKER_CELLS", "TOUCHSCREEN", "QR_PAYMENT", "NFC", "CASH", "CAMERA", "PIR_SENSOR"]
}
```

### GET /api/v1/machine/status
```json
{
  "status": "OPERATIONAL",
  "lastHeartbeat": "2026-04-09T14:30:00Z",
  "operatingMode": "NORMAL",
  "doorStatus": "CLOSED",
  "screenStatus": "ON",
  "networkStatus": "CONNECTED",
  "paymentSystemStatus": "READY",
  "refrigerationStatus": "RUNNING",
  "vendsSinceLastService": 342,
  "daysSinceLastService": 12,
  "nextServiceDue": "2026-04-15T00:00:00Z",
  "errors": [],
  "warnings": ["LOW_STOCK_SLOTS_5"]
}
```

### GET /api/v1/inventory (example for locker model)
```json
{
  "machineId": "VM-TCN-001",
  "totalSlots": 64,
  "occupiedSlots": 47,
  "emptySlots": 17,
  "soldOutSlots": 3,
  "lowStockSlots": 5,
  "lastRestockDate": "2026-04-07T08:15:00Z",
  "slots": [
    {
      "slotId": "A01",
      "row": "A",
      "column": 1,
      "status": "STOCKED",
      "productId": "SKU-COLA-355",
      "productName": "Coca-Cola 355ml",
      "category": "COLD_BEVERAGE",
      "currentQuantity": 1,
      "capacity": 1,
      "parLevel": 1,
      "price": 175,
      "currency": "USD",
      "expirationDate": "2026-12-31",
      "lockStatus": "LOCKED",
      "lastVendTime": "2026-04-09T11:22:00Z",
      "vendCount": 23
    },
    {
      "slotId": "A02",
      "row": "A",
      "column": 2,
      "status": "SOLD_OUT",
      "productId": "SKU-WATER-500",
      "productName": "Aquafina 500ml",
      "category": "COLD_BEVERAGE",
      "currentQuantity": 0,
      "capacity": 1,
      "parLevel": 1,
      "price": 150,
      "currency": "USD",
      "expirationDate": null,
      "lockStatus": "LOCKED",
      "lastVendTime": "2026-04-09T13:45:00Z",
      "vendCount": 31,
      "soldOutSince": "2026-04-09T13:45:00Z"
    }
  ]
}
```

### GET /api/v1/inventory (example for spring/conveyor model - AF-60C)
```json
{
  "machineId": "VM-AF60-001",
  "totalSlots": 60,
  "occupiedSlots": 58,
  "emptySlots": 2,
  "soldOutSlots": 2,
  "lowStockSlots": 8,
  "lastRestockDate": "2026-04-06T07:30:00Z",
  "slots": [
    {
      "slotId": "A1",
      "row": "A",
      "column": 1,
      "status": "STOCKED",
      "productId": "SKU-PEPSI-355",
      "productName": "Pepsi 355ml Can",
      "category": "COLD_BEVERAGE",
      "currentQuantity": 8,
      "capacity": 12,
      "parLevel": 10,
      "price": 175,
      "currency": "USD",
      "expirationDate": "2027-03-15",
      "motorStatus": "OK",
      "lastVendTime": "2026-04-09T14:10:00Z",
      "vendCount": 156
    }
  ]
}
```

### GET /api/v1/transactions
```json
{
  "machineId": "VM-TCN-001",
  "period": "today",
  "totalTransactions": 34,
  "totalRevenue": 6825,
  "currency": "USD",
  "transactions": [
    {
      "transactionId": "TXN-20260409-143022-001",
      "timestamp": "2026-04-09T14:30:22Z",
      "slotId": "A01",
      "productId": "SKU-COLA-355",
      "productName": "Coca-Cola 355ml",
      "category": "COLD_BEVERAGE",
      "price": 175,
      "currency": "USD",
      "paymentMethod": "NFC_CONTACTLESS",
      "cardType": "VISA",
      "cardLastFour": "4242",
      "vendResult": "SUCCESS",
      "dispenseDuration": 2.3
    },
    {
      "transactionId": "TXN-20260409-141500-002",
      "timestamp": "2026-04-09T14:15:00Z",
      "slotId": "B05",
      "productId": "SKU-SNICKERS-52",
      "productName": "Snickers Bar 52g",
      "category": "SNACK",
      "price": 200,
      "currency": "USD",
      "paymentMethod": "QR_CODE",
      "qrProvider": "APPLE_PAY",
      "vendResult": "SUCCESS",
      "dispenseDuration": 1.8
    },
    {
      "transactionId": "TXN-20260409-140200-003",
      "timestamp": "2026-04-09T14:02:00Z",
      "slotId": "C12",
      "productId": "SKU-REDBULL-250",
      "productName": "Red Bull 250ml",
      "category": "ENERGY_DRINK",
      "price": 350,
      "currency": "USD",
      "paymentMethod": "CASH",
      "cashInserted": 500,
      "changeGiven": 150,
      "vendResult": "SUCCESS",
      "dispenseDuration": 3.1
    }
  ]
}
```

### GET /api/v1/transactions/summary
```json
{
  "machineId": "VM-TCN-001",
  "period": "today",
  "totalVends": 34,
  "successfulVends": 33,
  "failedVends": 1,
  "totalRevenue": 6825,
  "currency": "USD",
  "averageTransactionValue": 207,
  "byPaymentMethod": {
    "CASH": { "count": 8, "revenue": 1550 },
    "CREDIT_CARD": { "count": 6, "revenue": 1275 },
    "NFC_CONTACTLESS": { "count": 10, "revenue": 2100 },
    "QR_CODE": { "count": 9, "revenue": 1900 },
    "MOBILE_WALLET": { "count": 0, "revenue": 0 }
  },
  "byCategory": {
    "COLD_BEVERAGE": { "count": 18, "revenue": 3150 },
    "SNACK": { "count": 9, "revenue": 1800 },
    "ENERGY_DRINK": { "count": 5, "revenue": 1750 },
    "FRESH_FOOD": { "count": 1, "revenue": 125 }
  },
  "peakHour": 12,
  "peakHourVends": 8
}
```

### GET /api/v1/payment/status
```json
{
  "machineId": "VM-TCN-001",
  "overallStatus": "READY",
  "coinAcceptor": {
    "status": "OPERATIONAL",
    "model": "MEI CF-7900",
    "protocol": "MDB",
    "acceptedDenominations": [5, 10, 25, 100],
    "lastActivity": "2026-04-09T14:02:00Z"
  },
  "billValidator": {
    "status": "OPERATIONAL",
    "model": "MEI AE-2800",
    "protocol": "MDB",
    "acceptedDenominations": [100, 500, 1000, 2000],
    "stackerStatus": "OK",
    "stackerCapacity": 600,
    "stackerCount": 145,
    "lastActivity": "2026-04-09T13:30:00Z"
  },
  "cardReader": {
    "status": "OPERATIONAL",
    "model": "Nayax VPOS Touch",
    "nfcEnabled": true,
    "chipEnabled": true,
    "contactlessEnabled": true,
    "acceptedCards": ["VISA", "MASTERCARD", "AMEX", "DISCOVER"],
    "lastTransaction": "2026-04-09T14:30:22Z"
  },
  "qrScanner": {
    "status": "OPERATIONAL",
    "supportedProviders": ["APPLE_PAY", "GOOGLE_PAY", "ALIPAY", "WECHAT_PAY"],
    "lastScan": "2026-04-09T14:15:00Z"
  }
}
```

### GET /api/v1/payment/cashbox
```json
{
  "machineId": "VM-TCN-001",
  "lastCollectionDate": "2026-04-05T09:00:00Z",
  "coinTubes": [
    { "denomination": 5, "count": 42, "capacity": 50, "value": 210 },
    { "denomination": 10, "count": 38, "capacity": 50, "value": 380 },
    { "denomination": 25, "count": 35, "capacity": 40, "value": 875 },
    { "denomination": 100, "count": 12, "capacity": 20, "value": 1200 }
  ],
  "coinBox": {
    "totalValue": 3250,
    "estimatedCount": 187,
    "status": "OK"
  },
  "billStacker": {
    "count": 145,
    "capacity": 600,
    "totalValue": 28500,
    "byDenomination": [
      { "denomination": 100, "count": 85, "value": 8500 },
      { "denomination": 500, "count": 35, "value": 17500 },
      { "denomination": 1000, "count": 20, "value": 20000 },
      { "denomination": 2000, "count": 5, "value": 10000 }
    ],
    "status": "OK"
  },
  "totalCashToCollect": 31750,
  "changeAvailable": true,
  "exactChangeRequired": false
}
```

### GET /api/v1/temperature
```json
{
  "machineId": "VM-AF60-001",
  "timestamp": "2026-04-09T14:30:00Z",
  "zones": [
    {
      "zoneId": "MAIN",
      "zoneName": "Refrigerated Cabinet",
      "currentTemp": 4.2,
      "setpoint": 4.0,
      "minTemp": 3.1,
      "maxTemp": 5.8,
      "unit": "CELSIUS",
      "status": "NORMAL",
      "compressorRunning": true,
      "compressorDutyCycle": 62.5,
      "compressorRuntime": 14400,
      "defrostStatus": "IDLE",
      "lastDefrost": "2026-04-09T04:00:00Z",
      "doorSealIntegrity": "GOOD"
    }
  ],
  "ambientTemp": 23.5,
  "humidity": 45.2,
  "glassHeaterActive": false
}
```

For the **AF-D900-54C combo** model, this includes two zones:
```json
{
  "zones": [
    { "zoneId": "UPPER", "zoneName": "Ambient (Snacks)", "currentTemp": 22.1, "setpoint": 22.0 },
    { "zoneId": "LOWER", "zoneName": "Refrigerated (Drinks)", "currentTemp": 4.3, "setpoint": 4.0, "compressorRunning": true }
  ]
}
```

For the **TCN-ZK locker** models, temperature is reported only if the optional cooling module is installed:
```json
{
  "zones": [
    { "zoneId": "CABINET", "zoneName": "Locker Cabinet", "currentTemp": 21.8, "setpoint": null, "status": "AMBIENT" }
  ]
}
```

### GET /api/v1/alerts
```json
{
  "machineId": "VM-AF60-001",
  "activeAlerts": [
    {
      "alertId": "ALT-20260409-001",
      "timestamp": "2026-04-09T12:15:00Z",
      "severity": "WARNING",
      "category": "INVENTORY",
      "code": "LOW_STOCK",
      "description": "Slot B03 below par level (2/10, par=8)",
      "slotId": "B03",
      "currentValue": 2,
      "threshold": 8,
      "acknowledged": false
    },
    {
      "alertId": "ALT-20260409-002",
      "timestamp": "2026-04-09T13:45:00Z",
      "severity": "INFO",
      "category": "INVENTORY",
      "code": "SOLD_OUT",
      "description": "Slot A02 sold out",
      "slotId": "A02",
      "acknowledged": false
    }
  ],
  "totalActive": 2,
  "bySeverity": { "CRITICAL": 0, "WARNING": 1, "INFO": 1 }
}
```

### GET /api/v1/access/events
```json
{
  "machineId": "VM-TCN-001",
  "events": [
    {
      "eventId": "EVT-20260409-001",
      "timestamp": "2026-04-09T08:15:00Z",
      "eventType": "SERVICE_VISIT",
      "doorAction": "OPENED",
      "duration": 1800,
      "closedAt": "2026-04-09T08:45:00Z",
      "servicePersonnel": "TECH-042",
      "activities": ["RESTOCK", "CASH_COLLECTION", "CLEANING"]
    },
    {
      "eventId": "EVT-20260409-002",
      "timestamp": "2026-04-09T14:30:22Z",
      "eventType": "VEND_DISPENSE",
      "doorAction": "CELL_OPENED",
      "cellId": "A01",
      "duration": 8,
      "closedAt": "2026-04-09T14:30:30Z"
    }
  ]
}
```

### GET /api/v1/access/locks (locker models only)
```json
{
  "machineId": "VM-TCN-001",
  "totalCells": 64,
  "locks": [
    { "cellId": "A01", "status": "LOCKED", "occupied": true, "lastOpened": "2026-04-09T14:30:22Z" },
    { "cellId": "A02", "status": "LOCKED", "occupied": false, "lastOpened": "2026-04-09T13:45:00Z" },
    { "cellId": "A03", "status": "LOCKED", "occupied": true, "lastOpened": "2026-04-07T08:20:00Z" }
  ]
}
```

### GET /api/v1/energy
```json
{
  "machineId": "VM-AF60-001",
  "currentPowerDraw": 285,
  "unit": "WATTS",
  "dailyKWh": 4.8,
  "monthlyKWh": 142.5,
  "energySavingMode": "ACTIVE",
  "lightingStatus": "DIMMED",
  "compressorCycles": {
    "today": 48,
    "averageDaily": 52,
    "totalRuntime": 52200,
    "unit": "SECONDS"
  },
  "screenPowerSave": true,
  "lastWakeup": "2026-04-09T14:28:00Z"
}
```

### GET /api/v1/analytics/traffic
```json
{
  "machineId": "VM-TCN-001",
  "date": "2026-04-09",
  "totalApproaches": 187,
  "totalPurchases": 34,
  "conversionRate": 18.2,
  "hourlyTraffic": [
    { "hour": 0, "approaches": 1, "purchases": 0 },
    { "hour": 6, "approaches": 5, "purchases": 2 },
    { "hour": 7, "approaches": 18, "purchases": 5 },
    { "hour": 8, "approaches": 22, "purchases": 6 },
    { "hour": 9, "approaches": 15, "purchases": 3 },
    { "hour": 10, "approaches": 10, "purchases": 2 },
    { "hour": 11, "approaches": 14, "purchases": 3 },
    { "hour": 12, "approaches": 28, "purchases": 8 },
    { "hour": 13, "approaches": 20, "purchases": 4 },
    { "hour": 14, "approaches": 12, "purchases": 1 }
  ],
  "averageDwellTime": 12.5,
  "unit": "SECONDS"
}
```

### GET /api/v1/analytics/performance
```json
{
  "machineId": "VM-TCN-001",
  "period": "last_7_days",
  "slotPerformance": [
    {
      "slotId": "A01",
      "productName": "Coca-Cola 355ml",
      "vendCount": 23,
      "revenue": 4025,
      "velocity": 3.3,
      "rank": 1,
      "stockoutHours": 0,
      "margin": 45.2
    },
    {
      "slotId": "B05",
      "productName": "Snickers Bar 52g",
      "vendCount": 18,
      "revenue": 3600,
      "velocity": 2.6,
      "rank": 2,
      "stockoutHours": 2.5,
      "margin": 52.1
    }
  ],
  "topProducts": ["SKU-COLA-355", "SKU-SNICKERS-52", "SKU-REDBULL-250"],
  "slowMovers": ["SKU-GRANOLA-35", "SKU-TRAIL-MIX-40"],
  "recommendedSwaps": [
    {
      "removeSlot": "D08",
      "removeProduct": "SKU-GRANOLA-35",
      "reason": "Below 1 vend/day for 14 days",
      "suggestProduct": "SKU-MONSTER-473",
      "expectedLift": 340
    }
  ]
}
```

### GET /api/v1/analytics/demand
```json
{
  "machineId": "VM-TCN-001",
  "generatedAt": "2026-04-09T00:00:00Z",
  "forecastHorizon": "7_DAYS",
  "dailyPatterns": {
    "monday": { "peakHours": [8, 12, 17], "expectedVends": 38, "expectedRevenue": 7600 },
    "tuesday": { "peakHours": [8, 12, 17], "expectedVends": 42, "expectedRevenue": 8400 },
    "wednesday": { "peakHours": [8, 12, 17], "expectedVends": 40, "expectedRevenue": 8000 },
    "thursday": { "peakHours": [8, 12, 17], "expectedVends": 41, "expectedRevenue": 8200 },
    "friday": { "peakHours": [8, 12, 15], "expectedVends": 45, "expectedRevenue": 9000 },
    "saturday": { "peakHours": [11, 14], "expectedVends": 15, "expectedRevenue": 3000 },
    "sunday": { "peakHours": [11, 14], "expectedVends": 10, "expectedRevenue": 2000 }
  },
  "predictedStockouts": [
    { "slotId": "A01", "productName": "Coca-Cola 355ml", "estimatedStockoutTime": "2026-04-10T15:00:00Z" },
    { "slotId": "C02", "productName": "Red Bull 250ml", "estimatedStockoutTime": "2026-04-11T10:00:00Z" }
  ],
  "restockUrgency": "MODERATE",
  "recommendedRestockDate": "2026-04-11T07:00:00Z"
}
```

### GET /api/v1/analytics/health
```json
{
  "machineId": "VM-AF60-001",
  "overallHealthScore": 87,
  "components": [
    {
      "component": "COMPRESSOR",
      "healthScore": 82,
      "status": "GOOD",
      "metrics": {
        "currentDrawAmps": 4.2,
        "baselineAmps": 3.8,
        "driftPercent": 10.5,
        "cyclesPerDay": 52,
        "baselineCycles": 48,
        "estimatedLifeRemaining": 18200,
        "unit": "HOURS"
      },
      "prediction": "NORMAL"
    },
    {
      "component": "VEND_MOTOR_A01",
      "healthScore": 95,
      "status": "GOOD",
      "metrics": {
        "currentDrawAmps": 1.1,
        "baselineAmps": 1.0,
        "driftPercent": 10.0,
        "jamCount30Days": 0,
        "totalVends": 1523
      },
      "prediction": "NORMAL"
    },
    {
      "component": "BILL_VALIDATOR",
      "healthScore": 71,
      "status": "ATTENTION",
      "metrics": {
        "rejectionRate": 8.5,
        "baselineRejectionRate": 3.0,
        "jamCount30Days": 2,
        "lastJam": "2026-04-03T11:22:00Z"
      },
      "prediction": "SERVICE_RECOMMENDED_14_DAYS"
    }
  ]
}
```

### GET /api/v1/dex/audit
```json
{
  "machineId": "VM-TCN-001",
  "auditTimestamp": "2026-04-09T14:30:00Z",
  "dexVersion": "UCS_06_02",
  "machineSerial": "TCN2024ZK064001",
  "machineModel": "TCN-ZK-22SP-BLH64S",
  "controllerRom": "22SP-3.8.2",
  "intervalData": {
    "vendCountCash": 8,
    "vendValueCash": 1550,
    "vendCountCashless": 25,
    "vendValueCashless": 5275,
    "vendCountFree": 1,
    "testVendCount": 0,
    "totalVendCount": 34,
    "totalVendValue": 6825
  },
  "cumulativeData": {
    "totalVendCount": 14523,
    "totalVendValue": 2904600,
    "totalCashVendCount": 4820,
    "totalCashlessVendCount": 9703
  },
  "cashAudit": {
    "coinsRecognized": 187,
    "coinsToTubes": 42,
    "coinsToCashBox": 145,
    "billsRecognized": 45,
    "totalCashIn": 5200,
    "totalChangeDispensed": 3650,
    "cashOverpay": 0
  },
  "eventLog": [
    { "code": "EVT_DOOR_OPEN", "timestamp": "2026-04-09T08:15:00Z", "count": 1 },
    { "code": "EVT_RESTOCK", "timestamp": "2026-04-07T08:15:00Z", "count": 1 },
    { "code": "EVT_CASH_COLLECT", "timestamp": "2026-04-05T09:00:00Z", "count": 1 },
    { "code": "EVT_POWER_CYCLE", "timestamp": "2026-04-01T03:22:00Z", "count": 1 }
  ],
  "selectionAudit": [
    { "slotId": "A01", "productId": "SKU-COLA-355", "price": 175, "intervalVends": 3, "cumulativeVends": 523 },
    { "slotId": "A02", "productId": "SKU-WATER-500", "price": 150, "intervalVends": 5, "cumulativeVends": 612 }
  ]
}
```

---

## SNMP OIDs

Basic machine health OIDs for network management integration. Uses a simulated enterprise OID under the TCN/Afen namespace.

### OID Schema: `1.3.6.1.4.1.99999.1.{category}.{metric}.{instance}`

| OID | Description | Type |
|-----|-------------|------|
| `1.3.6.1.2.1.1.1.0` | sysDescr - Machine model description | Static |
| `1.3.6.1.2.1.1.2.0` | sysObjectID | Static |
| `1.3.6.1.2.1.1.3.0` | sysUpTime | Static |
| `1.3.6.1.2.1.1.5.0` | sysName - Machine ID | Dynamic |
| `1.3.6.1.2.1.1.6.0` | sysLocation - Location name | Dynamic |
| `1.3.6.1.4.1.99999.1.1.1.0` | CPU utilization % | Dynamic (MetricCPUPercent) |
| `1.3.6.1.4.1.99999.1.1.2.0` | Memory used KB | Dynamic (MetricMemUsed) |
| `1.3.6.1.4.1.99999.1.1.3.0` | Memory total KB | Dynamic (MetricMemTotal) |
| `1.3.6.1.4.1.99999.1.2.1.0` | Cabinet temperature (C x10) | Dynamic (MetricTemperature) |
| `1.3.6.1.4.1.99999.1.3.1.0` | Total vend count today | Static |
| `1.3.6.1.4.1.99999.1.3.2.0` | Total revenue today (cents) | Static |
| `1.3.6.1.4.1.99999.1.4.1.0` | Machine operational status (1=OK) | Static |
| `1.3.6.1.4.1.99999.1.4.2.0` | Active alert count | Static |
| `1.3.6.1.4.1.99999.1.5.1.0` | Compressor status (1=running) | Static |
| `1.3.6.1.4.1.99999.1.5.2.0` | Door status (1=closed) | Static |

---

## Device Profiles

### Vending Machine Controller Profile
```go
var profileVendingController = DeviceProfile{
    CPUBaseMin:  5,   CPUBaseMax:  20,  CPUSpike:  10,
    MemTotalKB:  2 * 1024 * 1024,  // 2 GB (RK3288 Android board)
    MemBaseMin:  55,  MemBaseMax:  75,  MemVariance: 8,
    TempBaseMin: 28,  TempBaseMax: 42,  TempSpike:  5,
}
```

Notes:
- CPU is low (5-20% base) -- the RK3288 runs Android with a touchscreen UI, not heavy compute
- Memory is 2GB (typical for RK3288 vending boards)
- Temperature represents the controller board, not the cabinet (cabinet temp is in API responses)
- The temperature MetricsCycler cycles the controller board temp; cabinet/compressor temp is static in API JSON

---

## Model-Specific Differences

The 4 models share the same API schema but differ in their response data:

| Aspect | TCN-ZK+BLH-64S | TCN-ZK+BLH-40S | AF-60C | AF-D900-54C |
|--------|----------------|----------------|--------|-------------|
| **Slots** | 64 cells | 40 cells | 60 slots | 54 slots |
| **Slot capacity** | 1 per cell | 1 per cell | 8-12 per slot | 6-10 per slot |
| **Dispense mechanism** | Electronic lock | Electronic lock | Spring motor | Spring motor (snacks) + conveyor (drinks) |
| **Temperature zones** | 1 (ambient or optional cooling) | 1 (ambient or optional cooling) | 1 (refrigerated 4C) | 2 (ambient + refrigerated) |
| **Compressor** | Optional | Optional | Yes | Yes (lower zone) |
| **Lock status endpoint** | Yes | Yes | No | No |
| **Motor status per slot** | No | No | Yes | Yes |
| **Product types** | Any (boxes, bags, bottles) | Any | Canned/bottled drinks | Snacks + drinks |
| **Manufacturer** | TCN Vending Technology | TCN Vending Technology | Afen (Hunan TCN) | Afen (Hunan TCN) |

---

## File Structure

```
go/simulator/resources/
├── tcn_zk_blh_64s/
│   ├── tcn_zk_blh_64s_snmp.json          # SNMP OIDs
│   ├── tcn_zk_blh_64s_api_machine.json    # Machine identity, status, config
│   ├── tcn_zk_blh_64s_api_inventory.json  # Inventory, planogram, alerts
│   ├── tcn_zk_blh_64s_api_sales.json      # Transactions, summary
│   ├── tcn_zk_blh_64s_api_payment.json    # Payment status, cashbox, cashless
│   ├── tcn_zk_blh_64s_api_monitor.json    # Temperature, energy, environment
│   ├── tcn_zk_blh_64s_api_access.json     # Door events, lock status
│   ├── tcn_zk_blh_64s_api_analytics.json  # Traffic, performance, demand, health
│   └── tcn_zk_blh_64s_api_dex.json        # DEX audit data
│
├── tcn_zk_blh_40s/
│   ├── tcn_zk_blh_40s_snmp.json
│   ├── tcn_zk_blh_40s_api_machine.json
│   ├── tcn_zk_blh_40s_api_inventory.json
│   ├── tcn_zk_blh_40s_api_sales.json
│   ├── tcn_zk_blh_40s_api_payment.json
│   ├── tcn_zk_blh_40s_api_monitor.json
│   ├── tcn_zk_blh_40s_api_access.json
│   ├── tcn_zk_blh_40s_api_analytics.json
│   └── tcn_zk_blh_40s_api_dex.json
│
├── afen_60c/
│   ├── afen_60c_snmp.json
│   ├── afen_60c_api_machine.json
│   ├── afen_60c_api_inventory.json
│   ├── afen_60c_api_sales.json
│   ├── afen_60c_api_payment.json
│   ├── afen_60c_api_monitor.json
│   ├── afen_60c_api_analytics.json
│   └── afen_60c_api_dex.json
│
├── afen_d900_54c/
│   ├── afen_d900_54c_snmp.json
│   ├── afen_d900_54c_api_machine.json
│   ├── afen_d900_54c_api_inventory.json
│   ├── afen_d900_54c_api_sales.json
│   ├── afen_d900_54c_api_payment.json
│   ├── afen_d900_54c_api_monitor.json
│   ├── afen_d900_54c_api_analytics.json
│   └── afen_d900_54c_api_dex.json
```

**Note:** The AF-60C and AF-D900-54C do NOT have `*_api_access.json` because they don't have per-cell locks. The access/door events are included in the monitor file for these models.

JSON files are split by domain (~9 files per model) to keep each file manageable and under 300 lines. All files within a directory are automatically merged by the resource loader.

---

## Code Changes

### 1. `go/simulator/types.go` -- Add to RoundRobinDeviceTypes

```go
var RoundRobinDeviceTypes = []string{
    // ... existing 28 device types ...

    // Vending Machines (4 types)
    "tcn_zk_blh_64s.json",
    "tcn_zk_blh_40s.json",
    "afen_60c.json",
    "afen_d900_54c.json",
}
```

### 2. `go/simulator/device_profiles.go` -- Add profile and map entries

```go
// Vending Machine Controller (RK3288 Android board)
var profileVendingController = DeviceProfile{
    CPUBaseMin:  5,   CPUBaseMax:  20,  CPUSpike:  10,
    MemTotalKB:  2 * 1024 * 1024,  // 2 GB
    MemBaseMin:  55,  MemBaseMax:  75,  MemVariance: 8,
    TempBaseMin: 28,  TempBaseMax: 42,  TempSpike:  5,
}

// Add to deviceProfileMap:
"tcn_zk_blh_64s.json":  profileVendingController,
"tcn_zk_blh_40s.json":  profileVendingController,
"afen_60c.json":        profileVendingController,
"afen_d900_54c.json":   profileVendingController,
```

### 3. `go/simulator/metrics_oids.go` -- Add vendor OID mappings

```go
// Vending Machine OIDs (simulated enterprise OID 99999)
"tcn_zk_blh_64s.json": {
    "1.3.6.1.4.1.99999.1.1.1.0": MetricCPUPercent,
    "1.3.6.1.4.1.99999.1.1.2.0": MetricMemUsed,
    "1.3.6.1.4.1.99999.1.1.3.0": MetricMemTotal,
    "1.3.6.1.4.1.99999.1.2.1.0": MetricTemperature,
},
"tcn_zk_blh_40s.json": {
    // same OIDs
},
"afen_60c.json": {
    // same OIDs
},
"afen_d900_54c.json": {
    // same OIDs
},
```

### 4. `go/simulator/resources.go` -- Add category recognition

In `getDeviceCategoryFromName()`:
```go
func getDeviceCategoryFromName(name string) string {
    // ... existing cases ...
    case strings.Contains(n, "tcn") || strings.Contains(n, "afen"):
        return "Vending Machines"
    // ...
}
```

In `getDeviceTypeFromName()`:
```go
func getDeviceTypeFromName(name string) string {
    // ... existing cases ...
    case strings.Contains(n, "tcn_zk"):
        return "TCN Locker Vending Machine"
    case strings.Contains(n, "afen_60c"):
        return "Afen Refrigerated Beverage Vending Machine"
    case strings.Contains(n, "afen_d900"):
        return "Afen Combo Vending Machine"
    // ...
}
```

---

## Phase Breakdown

### Phase 1: Foundation (Code Changes)
1. Add `profileVendingController` to `device_profiles.go`
2. Add 4 entries to `deviceProfileMap`
3. Add 4 entries to `vendorOIDs` in `metrics_oids.go`
4. Add 4 entries to `RoundRobinDeviceTypes` in `types.go`
5. Add vending machine category recognition to `resources.go`

### Phase 2: TCN-ZK(22SP)+BLH-64S Resource Files
Create 9 JSON resource files in `resources/tcn_zk_blh_64s/`:
- SNMP OIDs file
- 7 API domain files (machine, inventory, sales, payment, monitor, access, analytics)
- DEX audit file

This is the **reference model** -- the most feature-rich (locker cells + lock status). All response schemas are designed here first.

### Phase 3: TCN-ZK(22SP)+BLH-40S Resource Files
Create 9 JSON resource files in `resources/tcn_zk_blh_40s/`:
- Copy structure from BLH-64S
- Adjust: 40 cells instead of 64, different serial/model strings, different inventory data

### Phase 4: AF-60C(22SP) Resource Files
Create 8 JSON resource files in `resources/afen_60c/`:
- No lock status endpoint (spring motor, not lockers)
- Single refrigerated zone (compressor always running)
- Higher slot capacity (8-12 per slot vs 1 per cell)
- Motor status per slot (spring motor health)
- Different product mix (cold beverages only)

### Phase 5: AF-D900-54C(22SP) Resource Files
Create 8 JSON resource files in `resources/afen_d900_54c/`:
- No lock status endpoint
- Dual temperature zones (ambient upper + refrigerated lower)
- Mixed product categories (snacks + drinks)
- 54 slots across both zones

### Phase 6: Build & Verify
1. `cd ../l8opensim/go && go build ./...` -- verify compilation
2. `go vet ./...` -- verify no issues
3. Verify resource files load correctly (start simulator, create 4 devices one per type)
4. Verify round-robin creates all 4 types across 100 devices
5. Test API endpoints via curl for each model
6. Test SNMP OIDs via snmpget/snmpwalk for each model

---

## Traceability Matrix

| # | Item | Phase |
|---|------|-------|
| 1 | Device profile for vending controller | Phase 1 |
| 2 | Profile map entries (4 models) | Phase 1 |
| 3 | Vendor OID mappings (4 models) | Phase 1 |
| 4 | Round-robin registration (4 models) | Phase 1 |
| 5 | Category recognition in resources.go | Phase 1 |
| 6 | TCN-ZK+BLH-64S SNMP resource | Phase 2 |
| 7 | TCN-ZK+BLH-64S API resources (7 domain files) | Phase 2 |
| 8 | TCN-ZK+BLH-64S DEX audit resource | Phase 2 |
| 9 | TCN-ZK+BLH-40S all resources | Phase 3 |
| 10 | AF-60C all resources | Phase 4 |
| 11 | AF-D900-54C all resources | Phase 5 |
| 12 | Compilation verification | Phase 6 |
| 13 | Resource loading verification | Phase 6 |
| 14 | API endpoint testing | Phase 6 |
| 15 | SNMP OID testing | Phase 6 |
| 16 | Round-robin 100-device test | Phase 6 |

---

## Estimated File Counts

| Category | Count |
|----------|-------|
| New resource directories | 4 |
| New JSON resource files | 34 (9+9+8+8) |
| Modified Go files | 4 (types.go, device_profiles.go, metrics_oids.go, resources.go) |
| New Go files | 0 |
| **Total new/modified files** | **38** |

---

## Product Catalog (Used Across All Models)

Realistic product data for populating inventory slots:

### Cold Beverages
| SKU | Product | Price (cents) | Category |
|-----|---------|--------------|----------|
| SKU-COLA-355 | Coca-Cola 355ml | 175 | COLD_BEVERAGE |
| SKU-PEPSI-355 | Pepsi 355ml Can | 175 | COLD_BEVERAGE |
| SKU-WATER-500 | Aquafina 500ml | 150 | COLD_BEVERAGE |
| SKU-SPRITE-355 | Sprite 355ml Can | 175 | COLD_BEVERAGE |
| SKU-FANTA-355 | Fanta Orange 355ml | 175 | COLD_BEVERAGE |
| SKU-DRPEPPER-355 | Dr Pepper 355ml | 175 | COLD_BEVERAGE |
| SKU-MTN-DEW-355 | Mountain Dew 355ml | 175 | COLD_BEVERAGE |
| SKU-GATORADE-591 | Gatorade Cool Blue 591ml | 225 | COLD_BEVERAGE |
| SKU-JUICE-350 | Minute Maid OJ 350ml | 200 | COLD_BEVERAGE |
| SKU-ICED-TEA-500 | Arizona Iced Tea 500ml | 175 | COLD_BEVERAGE |

### Energy Drinks
| SKU | Product | Price (cents) | Category |
|-----|---------|--------------|----------|
| SKU-REDBULL-250 | Red Bull 250ml | 350 | ENERGY_DRINK |
| SKU-MONSTER-473 | Monster Energy 473ml | 325 | ENERGY_DRINK |
| SKU-CELSIUS-355 | Celsius Sparkling 355ml | 275 | ENERGY_DRINK |
| SKU-BANG-473 | Bang Energy 473ml | 300 | ENERGY_DRINK |

### Snacks
| SKU | Product | Price (cents) | Category |
|-----|---------|--------------|----------|
| SKU-SNICKERS-52 | Snickers Bar 52g | 200 | SNACK |
| SKU-DORITOS-28 | Doritos Nacho 28g | 175 | SNACK |
| SKU-LAYS-28 | Lay's Classic 28g | 175 | SNACK |
| SKU-CHEETOS-28 | Cheetos Crunchy 28g | 175 | SNACK |
| SKU-PEANUTS-40 | Planters Peanuts 40g | 150 | SNACK |
| SKU-GRANOLA-35 | Nature Valley Granola 35g | 175 | SNACK |
| SKU-TRAIL-MIX-40 | Trail Mix 40g | 200 | SNACK |
| SKU-COOKIE-30 | Chips Ahoy Cookie 30g | 150 | SNACK |
| SKU-TWIX-50 | Twix Bar 50g | 200 | SNACK |
| SKU-KIT-KAT-42 | Kit Kat 42g | 200 | SNACK |
| SKU-MM-49 | M&M's Peanut 49g | 200 | SNACK |
| SKU-PRETZELS-28 | Rold Gold Pretzels 28g | 150 | SNACK |

### Fresh Food (locker models)
| SKU | Product | Price (cents) | Category |
|-----|---------|--------------|----------|
| SKU-SANDWICH-01 | Turkey & Swiss Sandwich | 550 | FRESH_FOOD |
| SKU-SALAD-01 | Caesar Salad Bowl | 650 | FRESH_FOOD |
| SKU-WRAP-01 | Chicken Caesar Wrap | 575 | FRESH_FOOD |
| SKU-YOGURT-170 | Greek Yogurt 170g | 225 | FRESH_FOOD |
| SKU-FRUIT-CUP | Mixed Fruit Cup | 350 | FRESH_FOOD |

---

## Notes

1. **No new Go files needed** -- all changes fit within existing files. The existing `APIServer`, resource loader, and metrics cycler handle vending machines without modification.

2. **Dynamic metrics** -- The MetricsCycler provides cycling CPU/memory/temperature for the controller board via SNMP. The API response data (cabinet temperature, inventory levels, sales) is static per the JSON resource pattern. A future enhancement could add a vending-specific metrics cycler for inventory depletion and sales accumulation, but this is beyond the current scope.

3. **100 machines** -- With 4 models in round-robin, 100 machines produces 25 of each model. Each gets a unique IP and machine ID. The existing round-robin mechanism handles this automatically.

4. **Enterprise OID 99999** -- This is a simulated/fictional OID prefix. Real TCN/Afen machines don't have registered IANA enterprise OIDs for SNMP. If real OIDs are discovered later, the JSON files can be updated without code changes.

5. **API port** -- The existing APIServer uses HTTPS on port 8443 per device. Each of the 100 vending machines gets its own IP + port 8443, matching the storage device pattern.
