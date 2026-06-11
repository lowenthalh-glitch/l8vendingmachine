/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    var f = window.Layer8FormFactory;

    FleetMachines.forms = {
        VendFleetMachine: f.form('Vending Machine', [
            f.section('Machine Information', [
                ...f.text('machineId', 'Machine ID', false),
                ...f.text('name', 'Name', false),
                ...f.text('type', 'Type', false),
                ...f.text('model', 'Model', false),
                ...f.text('status', 'Status', false),
                ...f.text('deviceId', 'Payment Device ID', false),
                ...f.text('managementIp', 'Management System', false)
            ]),
            f.section('Transactions', [
                ...f.number('dailyTransactions', 'Daily Transactions'),
                ...f.text('lastTransactionAt', 'Last Transaction', false)
            ]),
            f.section('Location', [
                ...f.text('locationAddress', 'Address', false),
                ...f.text('locationCity', 'City', false),
                ...f.text('locationState', 'State', false),
                ...f.number('locationLat', 'Latitude'),
                ...f.number('locationLng', 'Longitude')
            ]),
            f.section('Slot Inventory', [
                ...f.number('totalSlots', 'Total Slots'),
                ...f.number('emptySlots', 'Empty Slots'),
                ...f.number('lowStockSlots', 'Low Stock Slots'),
                ...f.text('inventoryLastUpdated', 'Last Updated', false),
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
        ]),
        VendMachineGroup: f.form('Machine Group', [
            f.section('Group Info', [
                ...f.text('name', 'Name', true),
                ...f.text('description', 'Description'),
                ...f.text('region', 'Region'),
                ...f.number('machineCount', 'Machine Count')
            ])
        ]),
        VendFleetInventory: f.form('Product Summary', [
            f.section('Product Details', [
                ...f.text('summaryId', 'Product', false, { readOnly: true }),
                ...f.text('productName', 'Product Name', false, { readOnly: true }),
                ...f.money('unitPrice', 'Price', false, { readOnly: true }),
                ...f.number('totalMachines', 'Machines', false, { readOnly: true }),
                ...f.number('totalSlots', 'Slots', false, { readOnly: true }),
                ...f.number('totalUnitsInMachines', 'Units in Field', false, { readOnly: true }),
                ...f.number('totalCapacity', 'Capacity', false, { readOnly: true }),
                ...f.number('fleetSoldOutCount', 'Sold Out Machines', false, { readOnly: true }),
                ...f.number('fleetLowStockCount', 'Low Stock Machines', false, { readOnly: true }),
                ...f.date('lastUpdated', 'Last Updated', false, { readOnly: true })
            ])
        ]),
        VendLocation: f.form('Location', [
            f.section('Location Info', [
                ...f.text('name', 'Name', true),
                ...f.text('locationType', 'Type'),
                ...f.text('timezone', 'Timezone')
            ]),
            f.section('Contact', [
                ...f.text('contactName', 'Contact'),
                ...f.text('contactPhone', 'Phone'),
                ...f.text('contactEmail', 'Email')
            ])
        ]),
        VendMachineProfile: f.form('Machine Profile', [
            f.section('Machine', [
                ...f.text('profileId', 'Profile ID', false, { readOnly: true }),
                ...f.text('machineName', 'Machine', false, { readOnly: true }),
                ...f.text('machineId', 'Machine ID', false, { readOnly: true }),
                ...f.text('locationClass', 'Location Type', false, { readOnly: true }),
                ...f.number('weekendWeekdayRatio', 'Weekend/Weekday Ratio', false, { readOnly: true }),
                ...f.date('lastUpdated', 'Last Updated', false, { readOnly: true })
            ]),
            f.section('Depletion & Revenue', [
                ...f.number('avgHourlyDepletion', 'Avg Depletion/hr', false, { readOnly: true }),
                ...f.money('avgDailyRevenue', 'Avg Daily Revenue', false, { readOnly: true }),
                ...f.money('totalRevenue30d', '30-Day Revenue', false, { readOnly: true }),
                ...f.number('avgFillPct', 'Avg Fill %', false, { readOnly: true }),
                ...f.number('trendMultiplier', 'Trend Multiplier', false, { readOnly: true }),
                ...f.number('cascadeThresholdPct', 'Cascade Threshold %', false, { readOnly: true })
            ]),
            f.section('Restock History', [
                ...f.number('restockCount30d', 'Restocks (30 days)', false, { readOnly: true }),
                ...f.number('avgRestockIntervalHours', 'Avg Interval (hours)', false, { readOnly: true })
            ]),
            f.section('Top Products', [
                ...f.inlineTable('topProducts', 'Product Depletion', [
                    { key: 'productName', label: 'Product', type: 'text' },
                    { key: 'depletionRatePerHour', label: 'Units/hr', type: 'number' },
                    { key: 'price', label: 'Price', type: 'money' },
                    { key: 'avgStock', label: 'Avg Stock', type: 'number' },
                    { key: 'capacity', label: 'Capacity', type: 'number' },
                    { key: 'timeToEmptyHours', label: 'Time to Empty (hrs)', type: 'number' }
                ])
            ])
        ])
    };
})();
