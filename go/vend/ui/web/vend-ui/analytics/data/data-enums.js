/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.AnalyticsData = window.AnalyticsData || {};

    var factory = window.Layer8EnumFactory;
    var { createStatusRenderer, renderEnum } = Layer8DRenderers;

    var RESTOCK_PRIORITY = factory.create([
        ['Unspecified', null, ''],
        ['Low', 'low', 'layer8d-status-inactive'],
        ['Medium', 'medium', 'layer8d-status-pending'],
        ['High', 'high', 'layer8d-status-active'],
        ['Critical', 'critical', 'layer8d-status-terminated']
    ]);

    var RESTOCK_REASON = factory.simple([
        'Unspecified', 'Weekend Demand', 'Weekday Demand', 'Rush Hour',
        'Fast Mover Empty', 'Cascade Threshold', 'Route Grouping',
        'Seasonal Trend', 'Event Driven', 'Revenue At Risk', 'Critical Prediction'
    ]);

    AnalyticsData.enums = {
        RESTOCK_PRIORITY: RESTOCK_PRIORITY,
        RESTOCK_REASON: RESTOCK_REASON
    };

    AnalyticsData.render = {
        restockPriority: createStatusRenderer(RESTOCK_PRIORITY.enum, RESTOCK_PRIORITY.classes),
        restockReason: function(value) { return renderEnum(value, RESTOCK_REASON.enum); }
    };
})();
