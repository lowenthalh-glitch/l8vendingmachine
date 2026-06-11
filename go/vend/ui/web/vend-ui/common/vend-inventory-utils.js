/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
window.VendInventoryUtils = {
    calcFillPct: function(inventory) {
        var totalStock = 0, totalCapacity = 0;
        if (inventory && inventory.length > 0) {
            inventory.forEach(function(s) {
                totalStock += s.currentStock || 0;
                totalCapacity += s.capacity || 0;
            });
        }
        if (totalCapacity === 0) return -1;
        return Math.round(totalStock / totalCapacity * 100);
    },
    fillColor: function(pct) {
        if (pct < 0) return 'var(--layer8d-text-muted)';
        if (pct > 60) return 'var(--layer8d-success)';
        if (pct > 30) return 'var(--layer8d-warning)';
        return 'var(--layer8d-error)';
    },
    productPrices: {
        'Coca-Cola 12oz': 175, 'Pepsi 12oz': 175, 'Snickers Bar': 200,
        "Lay's Classic Chips": 225, 'Monster Energy 16oz': 350, 'Bottled Water 16oz': 150,
        'KitKat Bar': 200, 'Doritos Nacho': 225, 'Red Bull 8.4oz': 325,
        'Gatorade Fruit Punch': 250, "M&M's Peanut": 200, 'Pringles Original': 275,
        'Tropicana OJ 10oz': 225, "Reese's Cups": 200, 'Smartwater 20oz': 250,
        'Nature Valley Granola': 175
    },
    getPrice: function(slot) {
        if (slot.price && slot.price > 0) return slot.price;
        return VendInventoryUtils.productPrices[slot.productName] || 200;
    },
    calcRevenue: function(inventory) {
        // Estimate revenue from items sold: sum(price * (capacity - currentStock))
        var totalCents = 0;
        if (inventory && inventory.length > 0) {
            inventory.forEach(function(s) {
                var sold = (s.capacity || 0) - (s.currentStock || 0);
                if (sold > 0) {
                    totalCents += sold * VendInventoryUtils.getPrice(s);
                }
            });
        }
        return totalCents;
    },
    formatMoney: function(cents) {
        if (!cents || cents <= 0) return '$0';
        return '$' + (cents / 100).toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 0 });
    },
    fillBar: function(pct) {
        if (pct < 0) return '<span style="color: var(--layer8d-text-muted);">-</span>';
        var color = VendInventoryUtils.fillColor(pct);
        return '<div style="display:flex;align-items:center;gap:6px;">' +
            '<div style="flex:1;height:8px;background:var(--layer8d-bg-light);border-radius:4px;overflow:hidden;min-width:60px;">' +
            '<div style="height:100%;width:' + pct + '%;background:' + color + ';border-radius:4px;"></div>' +
            '</div>' +
            '<span style="font-size:12px;font-weight:600;color:' + color + ';min-width:32px;">' + pct + '%</span>' +
            '</div>';
    }
};
