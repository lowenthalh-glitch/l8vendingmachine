/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    Layer8ModuleConfigFactory.create({
        namespace: 'Fleet',
        modules: {
            'machines': {
                label: 'Machines', icon: '🏭',
                services: [
                    { key: 'machines', label: 'Vending Machines', icon: '🏭', endpoint: '/10/Machine', model: 'VendFleetMachine', readOnly: true,
                        defaultSort: { column: 'emptySlots', direction: 'desc' } },
                    { key: 'profiles', label: 'Machine Profiles', icon: '📊', endpoint: '/10/MachProf', model: 'VendMachineProfile', readOnly: true,
                        supportedViews: ['table', 'chart'],
                        viewConfig: { chartType: 'bar', categoryField: 'machineName', valueField: 'avgDailyRevenue', aggregation: 'sum' } },
                    { key: 'products', label: 'Products', icon: '📦', endpoint: '/10/FleetInv', model: 'VendFleetInventory', readOnly: true },
                    { key: 'machine-groups', label: 'Groups', icon: '📁', endpoint: '/10/MachGrp', model: 'VendMachineGroup' },
                    { key: 'locations', label: 'Locations', icon: '📍', endpoint: '/10/Location', model: 'VendLocation' }
                ]
            }
        },
        submodules: ['FleetMachines']
    });
})();
