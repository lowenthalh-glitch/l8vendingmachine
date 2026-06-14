/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;

    MobileNayaxMachines.columns = {
        VendMachine: [
            { key: 'machineId', label: 'ID', primary: true, sortKey: 'machineId', filterKey: 'machineId' },
            ...col.custom('machines', 'Machines', function(item) {
                if (!item.machines) return '0';
                var count = Object.keys(item.machines).length;
                var online = 0;
                for (var k in item.machines) {
                    if (item.machines[k].status === 'online') online++;
                }
                return online + ' online / ' + count + ' total';
            }, { sortKey: 'machineId', secondary: true })
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

    // Add primary/secondary for card display
    MobileNayaxMachines.columns.VendMachineGroup[1].primary = true;
    MobileNayaxMachines.columns.VendMachineGroup[2].secondary = true;
    MobileNayaxMachines.columns.VendLocation[1].primary = true;
    MobileNayaxMachines.columns.VendLocation[2].secondary = true;
})();
