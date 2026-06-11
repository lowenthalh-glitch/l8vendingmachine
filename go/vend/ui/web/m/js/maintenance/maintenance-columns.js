/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var enums = MobileMaintenanceAlerts.enums;
    var render = MobileMaintenanceAlerts.render;

    MobileMaintenanceAlerts.columns = {
        VendAlert: [
            ...col.id('alertId'),
            ...col.col('machineId', 'Machine'),
            ...col.date('timestamp', 'Timestamp'),
            { key: 'severity', label: 'Severity', secondary: true, sortKey: 'severity', filterKey: 'severity',
              enumValues: enums.ALERT_SEVERITY.values,
              render: (item) => render.alertSeverity(item.severity) },
            ...col.enum('category', 'Category', null, render.alertCategory),
            ...col.col('code', 'Code'),
            { key: 'description', label: 'Description', primary: true, sortKey: 'description', filterKey: 'description' },
            ...col.enum('status', 'Status', enums.ALERT_STATUS.values, render.alertStatus)
        ],
        VendWorkOrder: [
            { key: 'workOrderId', label: 'Work Order ID', primary: true, sortKey: 'workOrderId', filterKey: 'workOrderId' },
            ...col.col('machineId', 'Machine'),
            ...col.col('workType', 'Work Type'),
            ...col.col('priority', 'Priority'),
            { key: 'status', label: 'Status', secondary: true, sortKey: 'status', filterKey: 'status',
              enumValues: enums.WORK_ORDER_STATUS.values,
              render: (item) => render.workOrderStatus(item.status) },
            ...col.col('assignedDriverId', 'Assigned Driver'),
            ...col.date('scheduledDate', 'Scheduled Date')
        ],
        VendServiceVisit: [
            { key: 'visitId', label: 'Visit ID', primary: true, sortKey: 'visitId', filterKey: 'visitId' },
            ...col.col('machineId', 'Machine'),
            ...col.col('driverId', 'Driver'),
            ...col.date('arrivalTime', 'Arrival Time'),
            ...col.number('duration', 'Duration'),
            ...col.number('slotsRestocked', 'Slots Restocked')
        ]
    };
})();
