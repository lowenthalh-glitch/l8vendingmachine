# Inventory Visual Indicators — Fill Bar Column + Dashboard Widget

## Context

Currently, seeing inventory status requires navigating to the Fleet table, reading numeric columns (totalSlots, emptySlots, lowStockSlots), and mentally calculating fill percentages. There's no at-a-glance visual indicator.

This plan adds two visual elements:
1. **Fill-bar column in the Fleet table** — an inline progress bar showing total inventory fill %
2. **Dashboard inventory widget** — per-machine fill bars sorted worst-first for operational prioritization

## Data Available

`VendFleetMachine` already has the fields needed:
- `totalSlots` (int32) — total capacity across all slots
- `emptySlots` (int32) — slots with 0 stock
- `lowStockSlots` (int32) — slots below threshold
- `inventory` ([]*VendMachineSlot) — per-slot detail with `currentStock` and `capacity`

Fill % can be calculated client-side from slot-level data: `sum(currentStock) / sum(capacity) × 100`

The slot-level `inventory` array is available in the response since it's `select *`.

## Shared Helper — No Duplicate Calculation

The fill % calculation and color determination is used in both the column renderer and the dashboard widget. To avoid duplication, extract a shared helper object.

**File: `go/vend/ui/web/vend-ui/common/vend-inventory-utils.js`** (new)

```js
window.VendInventoryUtils = {
    calcFillPct: function(inventory) {
        var totalStock = 0, totalCapacity = 0;
        if (inventory && inventory.length > 0) {
            inventory.forEach(function(s) {
                totalStock += s.currentStock || 0;
                totalCapacity += s.capacity || 0;
            });
        }
        if (totalCapacity === 0) return -1; // no data
        return Math.round(totalStock / totalCapacity * 100);
    },
    fillColor: function(pct) {
        if (pct < 0) return 'var(--layer8d-text-muted)';
        if (pct > 60) return 'var(--layer8d-success)';
        if (pct > 30) return 'var(--layer8d-warning)';
        return 'var(--layer8d-error)';
    },
    fillBar: function(pct) {
        if (pct < 0) return '<span style="color: var(--layer8d-text-muted);">—</span>';
        var color = VendInventoryUtils.fillColor(pct);
        return '<div style="display:flex;align-items:center;gap:6px;">' +
            '<div style="flex:1;height:8px;background:var(--layer8d-bg-light);border-radius:4px;overflow:hidden;min-width:60px;">' +
            '<div style="height:100%;width:' + pct + '%;background:' + color + ';border-radius:4px;"></div>' +
            '</div>' +
            '<span style="font-size:12px;font-weight:600;color:' + color + ';min-width:32px;">' + pct + '%</span>' +
            '</div>';
    }
};
```

Both the column renderer and the dashboard widget call `VendInventoryUtils.calcFillPct()` and `VendInventoryUtils.fillBar()`.

## Phase 1: Shared Utility + Fill-Bar Column in Fleet Table

### Step 1: Create shared utility

**File: `go/vend/ui/web/vend-ui/common/vend-inventory-utils.js`** (new, ~25 lines)

### Step 2: Wire into app.html

Add script tag before Fleet module scripts:
```html
<script src="vend-ui/common/vend-inventory-utils.js"></script>
```

### Step 3: Add fill-bar column

**File: `go/vend/ui/web/vend-ui/fleet/machines/machines-columns.js`**

Add a custom column after `status`, before `locationCity`:
```js
...col.custom('inventoryFill', 'Inventory', function(item) {
    var pct = VendInventoryUtils.calcFillPct(item.inventory);
    return VendInventoryUtils.fillBar(pct);
}, { sortKey: 'emptySlots' }),
```

## Phase 2: Dashboard Inventory Widget

**File: `go/vend/ui/web/sections/dashboard.html`**

Add a third container for the inventory widget:
```html
<div id="dashboard-inventory" style="margin-top: 16px;"></div>
```

**File: `go/vend/ui/web/vend-ui/dashboard/dashboard-init.js`**

After the existing machine stats `.then()` handler (which already has access to the `machines` array), add a call to render the inventory widget:

```js
renderInventoryWidget(machines);
```

New function `renderInventoryWidget(machines)`:
1. Calculate fill % for each machine using `VendInventoryUtils.calcFillPct(m.inventory)`
2. Sort by fill % ascending (worst first), skip machines with no data (pct === -1)
3. Take top 10
4. Render header with overall fleet fill % average
5. Render per-machine rows with name + fill bar (via `VendInventoryUtils.fillBar()`)
6. Add "View all in Fleet →" link at bottom

### Widget layout
```
┌─────────────────────────────────────────┐
│ Inventory Health            Fleet: 72%  │
├─────────────────────────────────────────┤
│ VM-003 Lobby Snacks    ██░░░░░░░  18%   │
│ VM-007 Break Room      ████░░░░░  35%   │
│ VM-001 Cafeteria       ██████░░░  55%   │
│ VM-005 Main Entrance   ████████░  78%   │
│ VM-002 Floor 2         █████████  92%   │
│                                         │
│              View all in Fleet →        │
└─────────────────────────────────────────┘
```

## Implementation Details

### Fill bar color thresholds
| Fill % | Color | Meaning |
|--------|-------|---------|
| > 60% | `--layer8d-success` (green) | Healthy |
| 30-60% | `--layer8d-warning` (yellow) | Low — schedule restock |
| < 30% | `--layer8d-error` (red) | Critical — immediate attention |

### No backend changes needed
All data already exists in VendFleetMachine (slot-level `inventory` array populated by the per-machine parser). This is purely a frontend rendering change.

## Mobile Parity

The mobile Fleet section currently uses model `VendMachine` (VCache type) which has a different structure — a `machines` map, not the flat `VendFleetMachine` with per-slot `inventory` array. Adding a fill-bar column on mobile would require the mobile nav config to switch to `VendFleetMachine` model and `/10/Machine` endpoint (same as desktop Fleet), which is a separate task.

The mobile dashboard does not currently exist as a custom page.

**Deferred**: Mobile fill-bar and mobile dashboard inventory widget. Flagged for next iteration when mobile fleet model alignment is addressed.

## Traceability Matrix

| # | Section | Action Item | Phase |
|---|---------|-------------|-------|
| 1 | Shared | Create VendInventoryUtils with calcFillPct, fillColor, fillBar | Phase 1 Step 1 |
| 2 | Shared | Wire vend-inventory-utils.js into app.html | Phase 1 Step 2 |
| 3 | Fleet | Add fill-bar custom column after status | Phase 1 Step 3 |
| 4 | Dashboard | Add dashboard-inventory container to dashboard.html | Phase 2 |
| 5 | Dashboard | Add renderInventoryWidget function | Phase 2 |
| 6 | Dashboard | Sort machines by fill % ascending (worst first) | Phase 2 |
| 7 | Dashboard | Limit to top 10 with "View all" link | Phase 2 |
| 8 | Mobile | Mobile fill-bar and dashboard widget | Deferred |

## Phase 3: Verification

- [ ] Navigate to **Fleet > Vending Machines** — verify fill-bar column appears with colored progress bars
- [ ] Verify sorting by emptySlots works on the Inventory column
- [ ] Navigate to **Dashboard** — verify Inventory Health widget appears below alerts
- [ ] Verify machines are sorted worst-first (lowest fill % at top)
- [ ] Verify color coding: green >60%, yellow 30-60%, red <30%
- [ ] Verify machines with no inventory data show "—" (not 0% or broken bar)
- [ ] Verify "View all in Fleet" link navigates to Fleet section

## Critical Files

| Action | File |
|--------|------|
| Create | `go/vend/ui/web/vend-ui/common/vend-inventory-utils.js` |
| Modify | `go/vend/ui/web/app.html` (add script tag) |
| Modify | `go/vend/ui/web/vend-ui/fleet/machines/machines-columns.js` |
| Modify | `go/vend/ui/web/sections/dashboard.html` |
| Modify | `go/vend/ui/web/vend-ui/dashboard/dashboard-init.js` |
