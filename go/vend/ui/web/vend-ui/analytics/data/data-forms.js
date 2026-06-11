/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;
    var enums = AnalyticsData.enums;

    AnalyticsData.forms = {
        VendRestockRecommendation: f.form('Restock Recommendation', [
            f.section('Recommendation', [
                ...f.text('recommendationId', 'ID', false, { readOnly: true }),
                ...f.text('machineName', 'Machine', false, { readOnly: true }),
                ...f.text('location', 'Location', false, { readOnly: true }),
                ...f.select('priority', 'Priority', enums.RESTOCK_PRIORITY.enum, false, { readOnly: true }),
                ...f.textarea('reason', 'Reason', false, { readOnly: true }),
                ...f.select('reasonCode', 'Reason Code', enums.RESTOCK_REASON.enum, false, { readOnly: true }),
                ...f.text('locationClass', 'Location Type', false, { readOnly: true })
            ]),
            f.section('Prediction', [
                ...f.date('predictedEmptyTime', 'Predicted Empty', false, { readOnly: true }),
                ...f.number('currentFillPct', 'Current Fill %', false, { readOnly: true }),
                ...f.number('projectedFillPct', 'Projected Fill %', false, { readOnly: true }),
                ...f.money('revenueAtRisk', 'Revenue At Risk', false, { readOnly: true }),
                ...f.number('confidence', 'Confidence', false, { readOnly: true }),
                ...f.number('revenueRank', 'Revenue Rank', false, { readOnly: true }),
                ...f.money('avgDailyRevenue', 'Avg Daily Revenue', false, { readOnly: true })
            ]),
            f.section('Suggested Products', [
                ...f.inlineTable('suggestedProducts', 'Products to Restock', [
                    { key: 'productName', label: 'Product', type: 'text' },
                    { key: 'currentStock', label: 'Current', type: 'number' },
                    { key: 'capacity', label: 'Capacity', type: 'number' },
                    { key: 'unitsToAdd', label: 'Units to Add', type: 'number' },
                    { key: 'depletionRate', label: 'Depletion/hr', type: 'number' }
                ])
            ]),
            f.section('Metadata', [
                ...f.text('routeGroupId', 'Route Group', false, { readOnly: true }),
                ...f.date('createdAt', 'Created', false, { readOnly: true }),
                ...f.date('expiresAt', 'Expires', false, { readOnly: true })
            ])
        ]),
        VendForecast: f.form('Forecast', [
            f.section('Forecast Details', [
                ...f.text('forecastId', 'Forecast ID', false, { readOnly: true }),
                ...f.text('machineId', 'Machine', false, { readOnly: true }),
                ...f.text('productId', 'Product', false, { readOnly: true }),
                ...f.date('forecastDate', 'Forecast Date', false, { readOnly: true }),
                ...f.number('horizonDays', 'Horizon (Days)', false, { readOnly: true }),
                ...f.number('predictedDailyVends', 'Predicted Vends/Day', false, { readOnly: true }),
                ...f.date('predictedStockoutTime', 'Predicted Stockout', false, { readOnly: true }),
                ...f.text('restockUrgency', 'Urgency', false, { readOnly: true }),
                ...f.number('confidenceScore', 'Confidence', false, { readOnly: true })
            ])
        ]),
        VendSlotPerformance: f.form('Slot Performance', [
            f.section('Performance Details', [
                ...f.text('performanceId', 'Performance ID', false, { readOnly: true }),
                ...f.text('machineId', 'Machine', false, { readOnly: true }),
                ...f.text('slotId', 'Slot', false, { readOnly: true }),
                ...f.text('productName', 'Product', false, { readOnly: true }),
                ...f.number('vendCount', 'Vends', false, { readOnly: true }),
                ...f.number('velocity', 'Vends/Day', false, { readOnly: true }),
                ...f.number('rank', 'Rank', false, { readOnly: true }),
                ...f.number('stockoutHours', 'Stockout Hours', false, { readOnly: true }),
                ...f.date('periodStart', 'Period Start', false, { readOnly: true }),
                ...f.date('periodEnd', 'Period End', false, { readOnly: true })
            ])
        ]),
        VendFleetInventory: f.form('Fleet Inventory', [
            f.section('Product Summary', [
                ...f.text('summaryId', 'Summary ID', false, { readOnly: true }),
                ...f.text('productName', 'Product', false, { readOnly: true }),
                ...f.number('totalMachines', 'Machines', false, { readOnly: true }),
                ...f.number('totalSlots', 'Slots', false, { readOnly: true }),
                ...f.number('totalUnitsInMachines', 'Units in Field', false, { readOnly: true }),
                ...f.number('totalCapacity', 'Capacity', false, { readOnly: true }),
                ...f.number('fleetSoldOutCount', 'Sold Out Machines', false, { readOnly: true }),
                ...f.number('fleetLowStockCount', 'Low Stock Machines', false, { readOnly: true }),
                ...f.date('lastUpdated', 'Last Updated', false, { readOnly: true })
            ])
        ]),
        VendInventorySnapshot: f.form('Inventory Snapshot', [
            f.section('Snapshot Details', [
                ...f.text('snapshotId', 'Snapshot ID', false, { readOnly: true }),
                ...f.text('machineName', 'Machine', false, { readOnly: true }),
                ...f.text('machineId', 'Machine ID', false, { readOnly: true }),
                ...f.date('timestamp', 'Timestamp', false, { readOnly: true }),
                ...f.number('fillPct', 'Fill %', false, { readOnly: true }),
                ...f.number('totalStock', 'Total Stock', false, { readOnly: true }),
                ...f.number('totalCapacity', 'Total Capacity', false, { readOnly: true }),
                ...f.number('totalSlots', 'Total Slots', false, { readOnly: true }),
                ...f.number('emptySlots', 'Empty Slots', false, { readOnly: true }),
                ...f.number('lowStockSlots', 'Low Stock Slots', false, { readOnly: true }),
                ...f.money('revenue', 'Revenue', false, { readOnly: true })
            ])
        ]),
        VendTopPerformer: f.form('Top Performer', [
            f.section('Performance Details', [
                ...f.number('rank', 'Rank', false, { readOnly: true }),
                ...f.text('machineName', 'Machine', false, { readOnly: true }),
                ...f.text('machineId', 'Machine ID', false, { readOnly: true }),
                ...f.money('revenue30d', '30-Day Revenue', false, { readOnly: true }),
                ...f.number('avgFillPct', 'Avg Fill %', false, { readOnly: true }),
                ...f.number('totalSnapshots', 'Data Points', false, { readOnly: true }),
                ...f.date('lastUpdated', 'Last Updated', false, { readOnly: true })
            ])
        ])
    };
})();
