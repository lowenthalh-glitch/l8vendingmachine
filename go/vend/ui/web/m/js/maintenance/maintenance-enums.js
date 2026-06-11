/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.MobileMaintenanceAlerts = window.MobileMaintenanceAlerts || {};

    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum } = Layer8MRenderers;

    const ALERT_SEVERITY = factory.create([
        ['Unspecified', null, ''],
        ['Info', 'info', 'active'],
        ['Warning', 'warning', 'pending'],
        ['Critical', 'critical', 'terminated']
    ]);

    const ALERT_CATEGORY = factory.simple([
        'Unspecified', 'Inventory', 'Temperature', 'Payment',
        'Mechanical', 'Connectivity', 'Security'
    ]);

    const ALERT_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'active'],
        ['Acknowledged', 'acknowledged', 'pending'],
        ['Resolved', 'resolved', 'inactive']
    ]);

    const WORK_ORDER_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'pending'],
        ['Assigned', 'assigned', 'active'],
        ['In Progress', 'inprogress', 'active'],
        ['Completed', 'completed', 'inactive'],
        ['Cancelled', 'cancelled', 'terminated']
    ]);

    MobileMaintenanceAlerts.enums = {
        ALERT_SEVERITY: ALERT_SEVERITY,
        ALERT_CATEGORY: ALERT_CATEGORY,
        ALERT_STATUS: ALERT_STATUS,
        WORK_ORDER_STATUS: WORK_ORDER_STATUS
    };

    MobileMaintenanceAlerts.render = {
        alertSeverity: createStatusRenderer(ALERT_SEVERITY.enum, ALERT_SEVERITY.classes),
        alertCategory: (value) => renderEnum(value, ALERT_CATEGORY.enum),
        alertStatus: createStatusRenderer(ALERT_STATUS.enum, ALERT_STATUS.classes),
        workOrderStatus: createStatusRenderer(WORK_ORDER_STATUS.enum, WORK_ORDER_STATUS.classes)
    };

    MobileMaintenanceAlerts.primaryKeys = {
        VendAlert: 'alertId',
        VendWorkOrder: 'workOrderId',
        VendServiceVisit: 'visitId'
    };
})();
