/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var enums = MaintenanceAlerts.enums;
    var render = MaintenanceAlerts.render;
    var col = window.Layer8ColumnFactory;

    MaintenanceAlerts.columns = {
        VendAlert: [
            ...col.id('alertId'),
            ...col.col('machineId', 'Machine'),
            ...col.date('timestamp', 'Timestamp'),
            ...col.enum('severity', 'Severity', enums.ALERT_SEVERITY.values, render.alertSeverity),
            ...col.enum('category', 'Category', null, render.alertCategory),
            ...col.col('code', 'Code'),
            ...col.col('description', 'Description'),
            ...col.enum('status', 'Status', enums.ALERT_STATUS.values, render.alertStatus)
        ],
        VendWorkOrder: [
            ...col.id('workOrderId'),
            ...col.col('machineId', 'Machine'),
            ...col.col('workType', 'Work Type'),
            ...col.col('priority', 'Priority'),
            ...col.status('status', 'Status', enums.WORK_ORDER_STATUS.values, render.workOrderStatus),
            ...col.col('assignedDriverId', 'Assigned Driver'),
            ...col.date('scheduledDate', 'Scheduled Date')
        ],
        VendServiceVisit: [
            ...col.id('visitId'),
            ...col.col('machineId', 'Machine'),
            ...col.col('driverId', 'Driver'),
            ...col.date('arrivalTime', 'Arrival Time'),
            ...col.number('duration', 'Duration'),
            ...col.number('slotsRestocked', 'Slots Restocked')
        ]
    };

    MaintenanceAlerts.primaryKeys = {
        VendAlert: 'alertId',
        VendWorkOrder: 'workOrderId',
        VendServiceVisit: 'visitId'
    };
})();
