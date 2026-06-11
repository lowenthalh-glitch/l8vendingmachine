/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('analytics', {
        title: 'Analytics',
        subtitle: 'Top Performers, Forecasts, Performance, Inventory History',
        icon: '📈',
        initFn: 'initializeAnalytics',
        modules: [{
            key: 'forecasts', label: 'Analytics', icon: '📈', isDefault: true,
            services: [
                { key: 'restock', label: 'Restock', icon: '🔄', isDefault: true },
                { key: 'top-performers', label: 'Top Performers', icon: '💰' },
                { key: 'forecasts', label: 'Forecasts', icon: '📈' },
                { key: 'performance', label: 'Performance', icon: '🏆' },
                { key: 'snapshots', label: 'Inventory History', icon: '📊' }
            ]
        }]
    });
})();
