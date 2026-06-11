/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('maintenance', {
        title: 'Maintenance',
        subtitle: 'Work Orders, Service Visits',
        icon: '🔧',
        initFn: 'initializeMaintenance',
        modules: [{
            key: 'work-orders', label: 'Work Orders', icon: '🔧', isDefault: true,
            services: [
                { key: 'work-orders', label: 'Work Orders', icon: '📝', isDefault: true },
                { key: 'service-visits', label: 'Service Visits', icon: '🔧' }
            ]
        }]
    });
})();
