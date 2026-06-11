/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    Layer8ModuleConfigFactory.create({
        namespace: 'Nayax',
        modules: {
            'machines': {
                label: 'Machines', icon: '🏭',
                services: [
                    { key: 'machines', label: 'Vending Machines', icon: '🏭', endpoint: '/0/VCache', model: 'VendMachine', readOnly: true },
                    { key: 'machine-groups', label: 'Groups', icon: '📁', endpoint: '/10/MachGrp', model: 'VendMachineGroup' },
                    { key: 'locations', label: 'Locations', icon: '📍', endpoint: '/10/Location', model: 'VendLocation' }
                ]
            }
        },
        submodules: ['NayaxMachines']
    });
})();
