/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;
    var enums = MobileComplianceInspections.enums;

    MobileComplianceInspections.forms = {
        VendInspection: f.form('Inspection', [
            f.section('Inspection Information', [
                ...f.reference('machineId', 'Machine', 'VendMachine'),
                ...f.reference('locationId', 'Location', 'VendLocation'),
                ...f.select('inspectionType', 'Type', enums.INSPECTION_TYPE.enum),
                ...f.select('status', 'Status', enums.INSPECTION_STATUS.enum),
                ...f.date('plannedDate', 'Planned Date'),
                ...f.text('inspectorName', 'Inspector Name'),
                ...f.textarea('scope', 'Scope')
            ])
        ]),
        VendInspectionFinding: f.form('Finding', [
            f.section('Finding Details', [
                ...f.reference('inspectionId', 'Inspection', 'VendInspection'),
                ...f.text('title', 'Title', true),
                ...f.select('severity', 'Severity', enums.FINDING_SEVERITY.enum),
                ...f.select('status', 'Status', enums.FINDING_STATUS.enum),
                ...f.textarea('condition', 'Condition'),
                ...f.textarea('criteria', 'Criteria'),
                ...f.textarea('recommendation', 'Recommendation'),
                ...f.text('responsibleId', 'Responsible'),
                ...f.date('dueDate', 'Due Date')
            ])
        ]),
        VendCertification: f.form('Certification', [
            f.section('Certification Information', [
                ...f.reference('machineId', 'Machine', 'VendMachine'),
                ...f.reference('locationId', 'Location', 'VendLocation'),
                ...f.text('certificationType', 'Certification Type'),
                ...f.text('standard', 'Standard'),
                ...f.text('certifyingBody', 'Certifying Body'),
                ...f.text('certificateNumber', 'Certificate Number'),
                ...f.date('issueDate', 'Issue Date'),
                ...f.date('expiryDate', 'Expiry Date'),
                ...f.select('status', 'Status', enums.CERTIFICATION_STATUS.enum)
            ])
        ])
    };
})();
