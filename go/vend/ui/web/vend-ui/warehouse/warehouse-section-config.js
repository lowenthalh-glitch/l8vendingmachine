/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('warehouse', {
        title: 'Warehouse & Supply Chain',
        subtitle: 'Facilities, Stock, Suppliers, Purchase Orders, Movements, Vehicle Loads',
        icon: '🏪',
        initFn: 'initializeWarehouse',
        modules: [{
            key: 'stock', label: 'Stock', icon: '🏪', isDefault: true,
            services: [
                { key: 'facilities', label: 'Facilities', icon: '🏪', isDefault: true },
                { key: 'suppliers', label: 'Suppliers', icon: '🤝' },
                { key: 'purchase-orders', label: 'Purchase Orders', icon: '📄' },
                { key: 'movements', label: 'Movements', icon: '↔️' },
                { key: 'vehicle-loads', label: 'Vehicle Loads', icon: '🚛' }
            ]
        }]
    });
})();
