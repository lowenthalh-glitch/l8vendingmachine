/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var enums = ComplianceInspections.enums;
    var render = ComplianceInspections.render;

    ComplianceInspections.columns = {
        VendInspection: [
            ...col.id('inspectionId'),
            ...col.col('machineId', 'Machine'),
            ...col.col('locationId', 'Location'),
            ...col.enum('inspectionType', 'Type', null, render.inspectionType),
            ...col.status('status', 'Status', enums.INSPECTION_STATUS.values, render.inspectionStatus),
            ...col.date('plannedDate', 'Planned Date'),
            ...col.date('actualDate', 'Actual Date'),
            ...col.col('inspectorName', 'Inspector')
        ],
        VendInspectionFinding: [
            ...col.id('findingId'),
            ...col.col('inspectionId', 'Inspection'),
            ...col.col('title', 'Title'),
            ...col.status('severity', 'Severity', enums.FINDING_SEVERITY.values, render.findingSeverity),
            ...col.status('status', 'Status', enums.FINDING_STATUS.values, render.findingStatus),
            ...col.date('dueDate', 'Due Date'),
            ...col.col('responsibleId', 'Responsible')
        ],
        VendCertification: [
            ...col.id('certificationId'),
            ...col.col('machineId', 'Machine'),
            ...col.col('certificationType', 'Type'),
            ...col.col('certifyingBody', 'Certifying Body'),
            ...col.date('issueDate', 'Issue Date'),
            ...col.date('expiryDate', 'Expiry Date'),
            ...col.status('status', 'Status', enums.CERTIFICATION_STATUS.values, render.certificationStatus)
        ]
    };

    ComplianceInspections.primaryKeys = {
        VendInspection: 'inspectionId',
        VendInspectionFinding: 'findingId',
        VendCertification: 'certificationId'
    };
})();
