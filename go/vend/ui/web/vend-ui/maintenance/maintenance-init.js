/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8DModuleFactory.create({
        namespace: 'Maintenance',
        defaultModule: 'work-orders',
        defaultService: 'work-orders',
        sectionSelector: 'work-orders',
        initializerName: 'initializeMaintenance',
        requiredNamespaces: ['MaintenanceAlerts']
    });
})();
