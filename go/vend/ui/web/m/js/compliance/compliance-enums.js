/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.MobileComplianceInspections = window.MobileComplianceInspections || {};

    var factory = window.Layer8EnumFactory;
    var { createStatusRenderer, renderEnum } = Layer8MRenderers;

    var INSPECTION_TYPE = factory.simple([
        'Unspecified', 'Health Department', 'Internal Audit',
        'Food Safety', 'Equipment Safety'
    ]);

    var INSPECTION_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Scheduled', 'scheduled', 'pending'],
        ['In Progress', 'inprogress', 'active'],
        ['Completed', 'completed', 'inactive'],
        ['Cancelled', 'cancelled', 'terminated']
    ]);

    var FINDING_SEVERITY = factory.create([
        ['Unspecified', null, ''],
        ['Critical', 'critical', 'terminated'],
        ['High', 'high', 'active'],
        ['Medium', 'medium', 'pending'],
        ['Low', 'low', 'inactive'],
        ['Informational', 'info', 'inactive']
    ]);

    var FINDING_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'active'],
        ['In Remediation', 'remediation', 'pending'],
        ['Closed', 'closed', 'inactive'],
        ['Deferred', 'deferred', 'pending']
    ]);

    var CERTIFICATION_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'active'],
        ['Pending', 'pending', 'pending'],
        ['Expired', 'expired', 'terminated'],
        ['Revoked', 'revoked', 'terminated'],
        ['Under Renewal', 'renewal', 'pending']
    ]);

    MobileComplianceInspections.enums = {
        INSPECTION_TYPE: INSPECTION_TYPE,
        INSPECTION_STATUS: INSPECTION_STATUS,
        FINDING_SEVERITY: FINDING_SEVERITY,
        FINDING_STATUS: FINDING_STATUS,
        CERTIFICATION_STATUS: CERTIFICATION_STATUS
    };

    MobileComplianceInspections.render = {
        inspectionType: function(value) { return renderEnum(value, INSPECTION_TYPE.enum); },
        inspectionStatus: createStatusRenderer(INSPECTION_STATUS.enum, INSPECTION_STATUS.classes),
        findingSeverity: createStatusRenderer(FINDING_SEVERITY.enum, FINDING_SEVERITY.classes),
        findingStatus: createStatusRenderer(FINDING_STATUS.enum, FINDING_STATUS.classes),
        certificationStatus: createStatusRenderer(CERTIFICATION_STATUS.enum, CERTIFICATION_STATUS.classes)
    };

    MobileComplianceInspections.primaryKeys = {
        VendInspection: 'inspectionId',
        VendInspectionFinding: 'findingId',
        VendCertification: 'certificationId'
    };
})();
