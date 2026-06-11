/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    Layer8ModuleConfigFactory.create({
        namespace: 'Inventory',
        modules: {
            'machines': {
                label: 'Machines', icon: '🏭',
                services: [
                    { key: 'machines', label: 'Vending Machines', icon: '🏭', endpoint: '/0/VCache', model: 'VendMachine', readOnly: true }
                ]
            }
        },
        submodules: ['InventoryMachines']
    });
})();
