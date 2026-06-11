/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;

    MobileAnalyticsData.forms = {
        VendForecast: f.form('Forecast', [
            f.section('Forecast Details', [
                ...f.text('forecastId', 'Forecast ID', false, { readOnly: true }),
                ...f.text('machineId', 'Machine', false, { readOnly: true }),
                ...f.text('productId', 'Product', false, { readOnly: true }),
                ...f.date('forecastDate', 'Forecast Date', false, { readOnly: true }),
                ...f.number('predictedDailyVends', 'Predicted Daily Vends', false, { readOnly: true }),
                ...f.text('restockUrgency', 'Restock Urgency', false, { readOnly: true }),
                ...f.number('confidenceScore', 'Confidence Score', false, { readOnly: true })
            ])
        ]),
        VendSlotPerformance: f.form('Slot Performance', [
            f.section('Performance Details', [
                ...f.text('performanceId', 'Performance ID', false, { readOnly: true }),
                ...f.text('machineId', 'Machine', false, { readOnly: true }),
                ...f.text('slotId', 'Slot', false, { readOnly: true }),
                ...f.text('productName', 'Product', false, { readOnly: true }),
                ...f.number('vendCount', 'Vend Count', false, { readOnly: true }),
                ...f.number('revenue', 'Revenue', false, { readOnly: true }),
                ...f.number('velocity', 'Velocity', false, { readOnly: true }),
                ...f.number('rank', 'Rank', false, { readOnly: true })
            ])
        ]),
        VendFleetInventory: f.form('Fleet Inventory', [
            f.section('Inventory Summary', [
                ...f.text('summaryId', 'Summary ID', false, { readOnly: true }),
                ...f.text('productName', 'Product', false, { readOnly: true }),
                ...f.number('unitPrice', 'Unit Price', false, { readOnly: true }),
                ...f.number('totalMachines', 'Total Machines', false, { readOnly: true }),
                ...f.number('totalUnitsInMachines', 'Units in Machines', false, { readOnly: true }),
                ...f.number('totalUnitsInWarehouses', 'Units in Warehouses', false, { readOnly: true }),
                ...f.number('totalSupplyChain', 'Total Supply Chain', false, { readOnly: true })
            ])
        ])
    };
})();
