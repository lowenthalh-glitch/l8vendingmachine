/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';

    Layer8ModuleConfigFactory.create({
        namespace: 'Analytics',
        modules: {
            'forecasts': {
                label: 'Analytics', icon: '📈',
                services: [
                    { key: 'restock', label: 'Restock Recommendations', icon: '🔄',
                      endpoint: '/10/Restock', model: 'VendRestockRecommendation',
                      readOnly: true, defaultSort: { column: 'priority', direction: 'desc' } },
                    { key: 'top-performers', label: 'Top Performers', icon: '💰',
                      endpoint: '/10/TopPerf', model: 'VendTopPerformer',
                      readOnly: true, supportedViews: ['chart', 'table'],
                      viewType: 'chart', defaultSort: { column: 'rank', direction: 'asc' },
                      viewConfig: { chartType: 'bar', categoryField: 'machineName', valueField: 'revenue30d',
                                    aggregation: 'sum', pageSize: 20 } },
                    { key: 'forecasts', label: 'Forecasts', icon: '📈',
                      endpoint: '/10/Forecast', model: 'VendForecast',
                      supportedViews: ['chart', 'table'], viewType: 'chart' },
                    { key: 'performance', label: 'Performance', icon: '🏆',
                      endpoint: '/10/SlotPerf', model: 'VendSlotPerformance',
                      supportedViews: ['chart', 'table'], viewType: 'chart' },
                    { key: 'snapshots', label: 'Inventory History', icon: '📊',
                      endpoint: '/10/InvSnap', model: 'VendInventorySnapshot',
                      readOnly: true, supportedViews: ['inventory-chart', 'table'],
                      viewType: 'inventory-chart' }
                ]
            }
        },
        submodules: ['AnalyticsData']
    });
})();
