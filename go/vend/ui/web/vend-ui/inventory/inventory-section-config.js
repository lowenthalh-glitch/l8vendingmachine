/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('inventory', {
        title: 'Inventory',
        subtitle: 'Vending Machine Inventory from Management System',
        icon: '📦',
        initFn: 'initializeInventory',
        modules: [{
            key: 'machines', label: 'Machines', icon: '🏭', isDefault: true,
            services: [
                { key: 'machines', label: 'Vending Machines', icon: '🏭', isDefault: true }
            ]
        }]
    });
})();
