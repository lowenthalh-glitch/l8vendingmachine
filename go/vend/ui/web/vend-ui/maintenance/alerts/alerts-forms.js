/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var enums = MaintenanceAlerts.enums;
    var f = window.Layer8FormFactory;

    MaintenanceAlerts.forms = {
        VendAlert: f.form('Alert', [
            f.section('Alert Information', [
                ...f.text('alertId', 'Alert ID', false, { readOnly: true }),
                ...f.text('machineId', 'Machine', false, { readOnly: true }),
                ...f.date('timestamp', 'Timestamp', false, { readOnly: true }),
                ...f.select('severity', 'Severity', enums.ALERT_SEVERITY.enum, false, { readOnly: true }),
                ...f.select('category', 'Category', enums.ALERT_CATEGORY, false, { readOnly: true }),
                ...f.text('code', 'Code', false, { readOnly: true }),
                ...f.textarea('description', 'Description', false, { readOnly: true }),
                ...f.select('status', 'Status', enums.ALERT_STATUS.enum, false, { readOnly: true })
            ])
        ]),
        VendWorkOrder: f.form('Work Order', [
            f.section('Work Order Details', [
                ...f.reference('machineId', 'Machine', 'VendMachine', true),
                ...f.text('workType', 'Work Type', false),
                ...f.text('priority', 'Priority', false),
                ...f.select('status', 'Status', enums.WORK_ORDER_STATUS.enum),
                ...f.reference('assignedDriverId', 'Assigned Driver', 'VendDriver'),
                ...f.textarea('description', 'Description'),
                ...f.date('scheduledDate', 'Scheduled Date')
            ])
        ]),
        VendServiceVisit: f.form('Service Visit', [
            f.section('Visit Details', [
                ...f.reference('machineId', 'Machine', 'VendMachine'),
                ...f.reference('driverId', 'Driver', 'VendDriver'),
                ...f.textarea('notes', 'Notes')
            ])
        ])
    };
})();
