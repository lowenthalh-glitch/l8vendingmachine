/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('routes', {
        title: 'Route Management',
        subtitle: 'Routes, Drivers',
        icon: '🚚',
        initFn: 'initializeRoutes',
        modules: [{
            key: 'routes', label: 'Routes', icon: '🚚', isDefault: true,
            services: [
                { key: 'routes', label: 'Routes', icon: '🗺️', isDefault: true },
                { key: 'drivers', label: 'Drivers', icon: '👤' }
            ]
        }]
    });
})();
