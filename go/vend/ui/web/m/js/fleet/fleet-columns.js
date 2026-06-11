/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var render = MobileFleetMachines.render;

    MobileFleetMachines.columns = {
        VendMachine: [
            ...col.id('machineId'),
            ...col.col('serialNumber', 'Serial'),
            { key: 'model', label: 'Model', primary: true, sortKey: 'model', filterKey: 'model' },
            ...col.col('manufacturer', 'Manufacturer'),
            ...col.enum('machineType', 'Type', null, render.machineType),
            { key: 'status', label: 'Status', secondary: true, sortKey: 'status', filterKey: 'status',
              render: (item) => render.machineStatus(item.status) },
            ...col.col('firmwareVersion', 'Firmware'),
            ...col.number('totalSlots', 'Slots'),
            ...col.col('locationId', 'Location'),
            ...col.date('lastHeartbeat', 'Last Heartbeat')
        ],
        VendMachineGroup: [
            ...col.id('groupId'),
            { key: 'name', label: 'Name', primary: true, sortKey: 'name', filterKey: 'name' },
            ...col.col('description', 'Description'),
            ...col.col('region', 'Region'),
            ...col.number('machineCount', 'Machines')
        ],
        VendLocation: [
            ...col.id('locationId'),
            { key: 'name', label: 'Name', primary: true, sortKey: 'name', filterKey: 'name' },
            ...col.col('locationType', 'Type'),
            ...col.col('timezone', 'Timezone'),
            ...col.col('contactName', 'Contact'),
            ...col.col('contactPhone', 'Phone')
        ],
        VendMachineProfile: [
            ...col.id('profileId'),
            { key: 'machineName', label: 'Machine', primary: true, sortKey: 'machineName' },
            { key: 'locationClass', label: 'Type', secondary: true },
            ...col.money('avgDailyRevenue', 'Revenue/Day'),
            ...col.number('avgFillPct', 'Avg Fill %'),
            ...col.number('trendMultiplier', 'Trend'),
            ...col.number('restockCount30d', 'Restocks')
        ]
    };
})();
