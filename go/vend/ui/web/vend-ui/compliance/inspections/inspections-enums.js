/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.ComplianceInspections = window.ComplianceInspections || {};

    var factory = window.Layer8EnumFactory;
    var { createStatusRenderer, renderEnum } = Layer8DRenderers;

    var INSPECTION_TYPE = factory.simple([
        'Unspecified', 'Health Department', 'Internal Audit',
        'Food Safety', 'Equipment Safety'
    ]);

    var INSPECTION_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Scheduled', 'scheduled', 'layer8d-status-pending'],
        ['In Progress', 'inprogress', 'layer8d-status-active'],
        ['Completed', 'completed', 'layer8d-status-inactive'],
        ['Cancelled', 'cancelled', 'layer8d-status-terminated']
    ]);

    var FINDING_SEVERITY = factory.create([
        ['Unspecified', null, ''],
        ['Critical', 'critical', 'layer8d-status-terminated'],
        ['High', 'high', 'layer8d-status-active'],
        ['Medium', 'medium', 'layer8d-status-pending'],
        ['Low', 'low', 'layer8d-status-inactive'],
        ['Informational', 'info', 'layer8d-status-inactive']
    ]);

    var FINDING_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'layer8d-status-active'],
        ['In Remediation', 'remediation', 'layer8d-status-pending'],
        ['Closed', 'closed', 'layer8d-status-inactive'],
        ['Deferred', 'deferred', 'layer8d-status-pending']
    ]);

    var CERTIFICATION_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'layer8d-status-active'],
        ['Pending', 'pending', 'layer8d-status-pending'],
        ['Expired', 'expired', 'layer8d-status-terminated'],
        ['Revoked', 'revoked', 'layer8d-status-terminated'],
        ['Under Renewal', 'renewal', 'layer8d-status-pending']
    ]);

    ComplianceInspections.enums = {
        INSPECTION_TYPE: INSPECTION_TYPE,
        INSPECTION_STATUS: INSPECTION_STATUS,
        FINDING_SEVERITY: FINDING_SEVERITY,
        FINDING_STATUS: FINDING_STATUS,
        CERTIFICATION_STATUS: CERTIFICATION_STATUS
    };

    ComplianceInspections.render = {
        inspectionType: function(value) { return renderEnum(value, INSPECTION_TYPE.enum); },
        inspectionStatus: createStatusRenderer(INSPECTION_STATUS.enum, INSPECTION_STATUS.classes),
        findingSeverity: createStatusRenderer(FINDING_SEVERITY.enum, FINDING_SEVERITY.classes),
        findingStatus: createStatusRenderer(FINDING_STATUS.enum, FINDING_STATUS.classes),
        certificationStatus: createStatusRenderer(CERTIFICATION_STATUS.enum, CERTIFICATION_STATUS.classes)
    };
})();
