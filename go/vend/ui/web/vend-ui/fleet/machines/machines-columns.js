/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var render = FleetMachines.render;

    FleetMachines.columns = {
        VendFleetMachine: [
            ...col.id('machineId'),
            ...col.col('name', 'Name'),
            ...col.custom('type', 'Type', function(item) { return render.machineType(item.type); }, { sortKey: 'type', filterKey: 'type' }),
            ...col.col('model', 'Model'),
            ...col.custom('status', 'Status', function(item) { return render.machineStatus(item.status); }, { sortKey: 'status', filterKey: 'status' }),
            ...col.custom('inventoryFill', 'Inventory', function(item) {
                var pct = VendInventoryUtils.calcFillPct(item.inventory);
                return VendInventoryUtils.fillBar(pct);
            }, { sortKey: 'emptySlots' }),
            ...col.col('locationCity', 'City'),
            ...col.col('locationState', 'State'),
            ...col.number('dailyTransactions', 'Daily TXN'),
            ...col.col('deviceId', 'Device ID'),
            ...col.col('managementIp', 'Management'),
            ...col.number('totalSlots', 'Slots'),
            ...col.number('emptySlots', 'Empty'),
            ...col.number('lowStockSlots', 'Low Stock')
        ],
        VendMachineGroup: [
            ...col.id('groupId'),
            ...col.col('name', 'Name'),
            ...col.col('description', 'Description'),
            ...col.col('region', 'Region'),
            ...col.number('machineCount', 'Machines')
        ],
        VendLocation: [
            ...col.id('locationId'),
            ...col.col('name', 'Name'),
            ...col.col('locationType', 'Type'),
            ...col.col('timezone', 'Timezone'),
            ...col.col('contactName', 'Contact'),
            ...col.col('contactPhone', 'Phone')
        ],
        VendFleetInventory: [
            ...col.id('summaryId'),
            ...col.col('productName', 'Product'),
            ...col.money('unitPrice', 'Price'),
            ...col.number('totalMachines', 'Machines'),
            ...col.number('totalSlots', 'Slots'),
            ...col.number('totalUnitsInMachines', 'Units in Field'),
            ...col.number('totalCapacity', 'Capacity'),
            ...col.custom('fillPct', 'Fill %', function(item) {
                if (!item.totalCapacity || item.totalCapacity === 0) return '-';
                var pct = Math.round((item.totalUnitsInMachines / item.totalCapacity) * 100);
                return VendInventoryUtils.fillBar(pct);
            }, { sortKey: 'totalUnitsInMachines' }),
            ...col.number('fleetSoldOutCount', 'Sold Out'),
            ...col.number('fleetLowStockCount', 'Low Stock'),
            ...col.date('lastUpdated', 'Updated')
        ],
        VendMachineProfile: [
            ...col.id('profileId'),
            ...col.col('machineName', 'Machine'),
            ...col.col('locationClass', 'Location Type'),
            ...col.number('avgHourlyDepletion', 'Depletion/hr'),
            ...col.money('avgDailyRevenue', 'Avg Revenue/Day'),
            ...col.money('totalRevenue30d', '30-Day Revenue'),
            ...col.number('avgFillPct', 'Avg Fill %'),
            ...col.number('trendMultiplier', 'Trend'),
            ...col.number('restockCount30d', 'Restocks'),
            ...col.number('weekendWeekdayRatio', 'Wknd/Wkday'),
            ...col.number('cascadeThresholdPct', 'Cascade %'),
            ...col.date('lastUpdated', 'Updated')
        ]
    };
})();
