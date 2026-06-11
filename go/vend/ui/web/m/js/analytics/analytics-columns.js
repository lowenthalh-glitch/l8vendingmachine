/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;

    MobileAnalyticsData.columns = {
        VendForecast: [
            ...col.id('forecastId'),
            ...col.col('machineId', 'Machine'),
            { key: 'productId', label: 'Product', primary: true, sortKey: 'productId', filterKey: 'productId' },
            ...col.date('forecastDate', 'Forecast Date'),
            ...col.number('predictedDailyVends', 'Predicted Daily Vends'),
            ...col.col('restockUrgency', 'Restock Urgency'),
            ...col.number('confidenceScore', 'Confidence Score')
        ],
        VendSlotPerformance: [
            ...col.id('performanceId'),
            ...col.col('machineId', 'Machine'),
            ...col.col('slotId', 'Slot'),
            { key: 'productName', label: 'Product', primary: true, sortKey: 'productName', filterKey: 'productName' },
            ...col.number('vendCount', 'Vend Count'),
            ...col.money('revenue', 'Revenue'),
            ...col.number('velocity', 'Velocity'),
            ...col.number('rank', 'Rank')
        ],
        VendFleetInventory: [
            ...col.id('summaryId'),
            { key: 'productName', label: 'Product', primary: true, sortKey: 'productName', filterKey: 'productName' },
            ...col.money('unitPrice', 'Unit Price'),
            ...col.number('totalMachines', 'Total Machines'),
            ...col.number('totalUnitsInMachines', 'Units in Machines'),
            ...col.number('totalUnitsInWarehouses', 'Units in Warehouses'),
            ...col.number('totalSupplyChain', 'Total Supply Chain')
        ],
        VendInventorySnapshot: [
            ...col.id('snapshotId'),
            { key: 'machineName', label: 'Machine', primary: true, sortKey: 'machineName' },
            ...col.date('timestamp', 'Time'),
            { key: 'fillPct', label: 'Fill %', secondary: true },
            ...col.money('revenue', 'Revenue'),
            ...col.number('totalStock', 'Stock'),
            ...col.number('totalSlots', 'Slots'),
            ...col.number('emptySlots', 'Empty')
        ],
        VendRestockRecommendation: [
            ...col.id('recommendationId'),
            { key: 'machineName', label: 'Machine', primary: true },
            { key: 'priority', label: 'Priority', secondary: true },
            ...col.col('reason', 'Reason'),
            ...col.date('predictedEmptyTime', 'Predicted Empty'),
            ...col.number('currentFillPct', 'Fill %'),
            ...col.money('revenueAtRisk', 'Revenue At Risk')
        ]
    };
})();
