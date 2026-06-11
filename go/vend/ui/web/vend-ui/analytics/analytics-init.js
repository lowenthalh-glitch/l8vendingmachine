/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8DModuleFactory.create({
        namespace: 'Analytics',
        defaultModule: 'forecasts',
        defaultService: 'restock',
        sectionSelector: 'forecasts',
        initializerName: 'initializeAnalytics',
        requiredNamespaces: ['AnalyticsData']
    });
})();
