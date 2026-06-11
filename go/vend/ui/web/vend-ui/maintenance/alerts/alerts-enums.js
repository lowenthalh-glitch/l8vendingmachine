/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.MaintenanceAlerts = window.MaintenanceAlerts || {};

    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer } = Layer8DRenderers;

    const ALERT_SEVERITY = factory.create([
        ['Unspecified', null, ''],
        ['Info', 'info', 'layer8d-status-active'],
        ['Warning', 'warning', 'layer8d-status-pending'],
        ['Critical', 'critical', 'layer8d-status-terminated']
    ]);

    const ALERT_CATEGORY = factory.simple([
        'Unspecified', 'Inventory', 'Temperature', 'Payment',
        'Mechanical', 'Connectivity', 'Security'
    ]);

    const ALERT_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'layer8d-status-active'],
        ['Acknowledged', 'acknowledged', 'layer8d-status-pending'],
        ['Resolved', 'resolved', 'layer8d-status-inactive']
    ]);

    const WORK_ORDER_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'layer8d-status-pending'],
        ['Assigned', 'assigned', 'layer8d-status-active'],
        ['In Progress', 'inprogress', 'layer8d-status-active'],
        ['Completed', 'completed', 'layer8d-status-inactive'],
        ['Cancelled', 'cancelled', 'layer8d-status-terminated']
    ]);

    MaintenanceAlerts.enums = {
        ALERT_SEVERITY: ALERT_SEVERITY,
        ALERT_CATEGORY: ALERT_CATEGORY,
        ALERT_STATUS: ALERT_STATUS,
        WORK_ORDER_STATUS: WORK_ORDER_STATUS
    };

    MaintenanceAlerts.render = {
        alertSeverity: createStatusRenderer(ALERT_SEVERITY.enum, ALERT_SEVERITY.classes),
        alertCategory: (value) => Layer8DRenderers.renderEnum(value, ALERT_CATEGORY),
        alertStatus: createStatusRenderer(ALERT_STATUS.enum, ALERT_STATUS.classes),
        workOrderStatus: createStatusRenderer(WORK_ORDER_STATUS.enum, WORK_ORDER_STATUS.classes)
    };
})();
