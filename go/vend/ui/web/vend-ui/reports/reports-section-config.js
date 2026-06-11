/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('reports', {
        title: 'Reports',
        subtitle: 'Scheduled Reports',
        icon: '📋',
        initFn: 'initializeReports',
        modules: [{
            key: 'reports', label: 'Reports', icon: '📋', isDefault: true,
            services: [
                { key: 'reports', label: 'Reports', icon: '📋', isDefault: true }
            ]
        }]
    });
})();
