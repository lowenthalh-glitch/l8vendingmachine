/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var render = NayaxMachines.render;

    NayaxMachines.columns = {
        VendMachine: [
            ...col.id('machineId'),
            ...col.custom('machines', 'Machines', function(item) {
                if (!item.machines) return '0';
                var count = Object.keys(item.machines).length;
                var online = 0;
                for (var k in item.machines) {
                    if (item.machines[k].status === 'online') online++;
                }
                return '<span class="layer8d-status-badge layer8d-status-active">' + online + ' online</span> / ' + count + ' total';
            }, { sortKey: 'machineId' })
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
        ]
    };
})();
