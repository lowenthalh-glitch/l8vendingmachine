/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var enums = AnalyticsData.enums;
    var render = AnalyticsData.render;

    AnalyticsData.columns = {
        VendRestockRecommendation: [
            ...col.status('priority', 'Priority', enums.RESTOCK_PRIORITY.values, render.restockPriority),
            ...col.col('machineName', 'Machine'),
            ...col.col('reason', 'Reason'),
            ...col.enum('reasonCode', 'Type', null, render.restockReason),
            ...col.date('predictedEmptyTime', 'Predicted Empty'),
            ...col.number('currentFillPct', 'Fill %'),
            ...col.money('revenueAtRisk', 'Revenue At Risk'),
            ...col.number('confidence', 'Confidence'),
            ...col.col('locationClass', 'Location'),
            ...col.date('expiresAt', 'Expires')
        ],
        VendForecast: [
            ...col.id('forecastId'),
            ...col.col('machineId', 'Machine'),
            ...col.col('productId', 'Product'),
            ...col.date('forecastDate', 'Forecast Date'),
            ...col.number('predictedDailyVends', 'Vends/Day'),
            ...col.date('predictedStockoutTime', 'Predicted Stockout'),
            ...col.col('restockUrgency', 'Urgency'),
            ...col.number('confidenceScore', 'Confidence')
        ],
        VendSlotPerformance: [
            ...col.id('performanceId'),
            ...col.col('machineId', 'Machine'),
            ...col.col('slotId', 'Slot'),
            ...col.col('productName', 'Product'),
            ...col.number('vendCount', 'Vends'),
            ...col.number('velocity', 'Vends/Day'),
            ...col.number('rank', 'Rank'),
            ...col.number('stockoutHours', 'Stockout Hours'),
            ...col.date('periodStart', 'Period Start'),
            ...col.date('periodEnd', 'Period End')
        ],
        VendFleetInventory: [
            ...col.id('summaryId'),
            ...col.col('productName', 'Product'),
            ...col.number('totalMachines', 'Machines'),
            ...col.number('totalSlots', 'Slots'),
            ...col.number('totalUnitsInMachines', 'Units in Field'),
            ...col.number('totalCapacity', 'Capacity'),
            ...col.custom('fillPct', 'Fill %', function(item) {
                if (!item.totalCapacity || item.totalCapacity === 0) return '-';
                var pct = Math.round((item.totalUnitsInMachines / item.totalCapacity) * 100);
                return VendInventoryUtils.fillBar(pct);
            }, { sortKey: 'totalUnitsInMachines' }),
            ...col.number('fleetSoldOutCount', 'Sold Out'),
            ...col.number('fleetLowStockCount', 'Low Stock'),
            ...col.date('lastUpdated', 'Updated')
        ],
        VendInventorySnapshot: [
            ...col.id('snapshotId'),
            ...col.col('machineName', 'Machine'),
            ...col.date('timestamp', 'Time'),
            ...col.number('fillPct', 'Fill %'),
            ...col.money('revenue', 'Revenue'),
            ...col.number('totalStock', 'Stock'),
            ...col.number('totalCapacity', 'Capacity'),
            ...col.number('totalSlots', 'Slots'),
            ...col.number('emptySlots', 'Empty'),
            ...col.number('lowStockSlots', 'Low Stock')
        ],
        VendTopPerformer: [
            ...col.number('rank', 'Rank'),
            ...col.col('machineName', 'Machine'),
            ...col.money('revenue30d', '30-Day Revenue'),
            ...col.number('avgFillPct', 'Avg Fill %'),
            ...col.number('totalSnapshots', 'Data Points'),
            ...col.date('lastUpdated', 'Updated')
        ]
    };

    AnalyticsData.primaryKeys = {
        VendForecast: 'forecastId',
        VendSlotPerformance: 'performanceId',
        VendFleetInventory: 'summaryId',
        VendInventorySnapshot: 'snapshotId',
        VendTopPerformer: 'performerId',
        VendRestockRecommendation: 'recommendationId'
    };
})();
