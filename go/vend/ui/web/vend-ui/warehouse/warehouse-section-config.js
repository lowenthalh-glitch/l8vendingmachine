/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('warehouse', {
        title: 'Warehouse & Supply Chain',
        subtitle: 'Warehouses, Stock, Suppliers, Purchase Orders, Movements, Vehicle Loads',
        icon: '🏪',
        initFn: 'initializeWarehouse',
        modules: [{
            key: 'stock', label: 'Stock', icon: '🏪', isDefault: true,
            services: [
                { key: 'warehouses', label: 'Warehouses', icon: '🏪', isDefault: true },
                { key: 'warehouse-stock', label: 'Stock', icon: '📦' },
                { key: 'suppliers', label: 'Suppliers', icon: '🤝' },
                { key: 'purchase-orders', label: 'Purchase Orders', icon: '📄' },
                { key: 'movements', label: 'Movements', icon: '↔️' },
                { key: 'vehicle-loads', label: 'Vehicle Loads', icon: '🚛' }
            ]
        }]
    });
})();
