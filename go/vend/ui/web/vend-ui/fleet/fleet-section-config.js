/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('fleet', {
        title: 'Fleet',
        subtitle: 'Vending Machine Fleet Management',
        icon: '🏭',
        initFn: 'initializeFleet',
        modules: [{
            key: 'machines', label: 'Machines', icon: '🏭', isDefault: true,
            services: [
                { key: 'machines', label: 'Vending Machines', icon: '🏭', isDefault: true },
                { key: 'profiles', label: 'Machine Profiles', icon: '📊' },
                { key: 'products', label: 'Products', icon: '📦' },
                { key: 'machine-groups', label: 'Groups', icon: '📁' },
                { key: 'locations', label: 'Locations', icon: '📍' }
            ]
        }]
    });
})();
