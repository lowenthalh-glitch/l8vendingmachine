/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('nayax', {
        title: 'Management Systems',
        subtitle: 'Real-time Vending Machine Fleet Monitoring',
        icon: '☁️',
        initFn: 'initializeNayax',
        modules: [{
            key: 'machines', label: 'Machines', icon: '🏭', isDefault: true,
            services: [
                { key: 'machines', label: 'Vending Machines', icon: '🏭', isDefault: true },
                { key: 'machine-groups', label: 'Groups', icon: '📁' },
                { key: 'locations', label: 'Locations', icon: '📍' }
            ]
        }]
    });
})();
