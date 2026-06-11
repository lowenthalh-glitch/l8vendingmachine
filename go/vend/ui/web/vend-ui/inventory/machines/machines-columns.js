/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var render = InventoryMachines.render;

    InventoryMachines.columns = {
        VendMachine: [
            ...col.id('machineId'),
            ...col.col('name', 'Name'),
            ...col.custom('type', 'Type', function(item) { return render.machineType(item.type); }, { sortKey: 'type', filterKey: 'type' }),
            ...col.col('model', 'Model'),
            ...col.custom('status', 'Status', function(item) { return render.machineStatus(item.status); }, { sortKey: 'status', filterKey: 'status' }),
            ...col.col('locationCity', 'City'),
            ...col.col('locationState', 'State'),
            ...col.number('dailyTransactions', 'Daily TXN'),
            ...col.number('revenueToday', 'Revenue Today'),
            ...col.number('uptimePercent', 'Uptime %'),
            ...col.col('deviceId', 'Device ID')
        ]
    };
})();
