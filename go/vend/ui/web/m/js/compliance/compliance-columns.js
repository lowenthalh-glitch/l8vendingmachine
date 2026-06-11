/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var enums = MobileComplianceInspections.enums;
    var render = MobileComplianceInspections.render;

    MobileComplianceInspections.columns = {
        VendInspection: [
            { key: 'inspectionId', label: 'Inspection ID', primary: true, sortKey: 'inspectionId', filterKey: 'inspectionId' },
            ...col.col('machineId', 'Machine'),
            ...col.col('locationId', 'Location'),
            ...col.enum('inspectionType', 'Type', null, render.inspectionType),
            { key: 'status', label: 'Status', secondary: true, sortKey: 'status', filterKey: 'status',
              enumValues: enums.INSPECTION_STATUS.values,
              render: (item) => render.inspectionStatus(item.status) },
            ...col.date('plannedDate', 'Planned Date'),
            ...col.date('actualDate', 'Actual Date'),
            ...col.col('inspectorName', 'Inspector')
        ],
        VendInspectionFinding: [
            ...col.id('findingId'),
            ...col.col('inspectionId', 'Inspection'),
            { key: 'title', label: 'Title', primary: true, sortKey: 'title', filterKey: 'title' },
            { key: 'severity', label: 'Severity', secondary: true, sortKey: 'severity', filterKey: 'severity',
              enumValues: enums.FINDING_SEVERITY.values,
              render: (item) => render.findingSeverity(item.severity) },
            ...col.status('status', 'Status', enums.FINDING_STATUS.values, render.findingStatus),
            ...col.date('dueDate', 'Due Date'),
            ...col.col('responsibleId', 'Responsible')
        ],
        VendCertification: [
            ...col.id('certificationId'),
            ...col.col('machineId', 'Machine'),
            { key: 'certificationType', label: 'Type', primary: true, sortKey: 'certificationType', filterKey: 'certificationType' },
            ...col.col('certifyingBody', 'Certifying Body'),
            ...col.date('issueDate', 'Issue Date'),
            ...col.date('expiryDate', 'Expiry Date'),
            ...col.status('status', 'Status', enums.CERTIFICATION_STATUS.values, render.certificationStatus)
        ]
    };
})();
