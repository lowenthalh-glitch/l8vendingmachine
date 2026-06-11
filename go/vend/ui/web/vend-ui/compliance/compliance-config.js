/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    const svc = Layer8ModuleConfigFactory.service;
    const mod = Layer8ModuleConfigFactory.module;

    Layer8ModuleConfigFactory.create({
        namespace: 'Compliance',
        modules: {
            'inspections': mod('Inspections', '\u{1F50D}', [
                svc('inspections', 'Inspections', '\u{1F4CB}', '/10/Inspction', 'VendInspection',
                    { supportedViews: ['table', 'calendar'] }),
                svc('findings', 'Findings', '\u{26A0}', '/10/InspFind', 'VendInspectionFinding',
                    { supportedViews: ['table', 'kanban'] }),
                svc('certifications', 'Certs', '\u{1F4DC}', '/10/VendCert', 'VendCertification')
            ])
        },
        submodules: ['ComplianceInspections']
    });
})();
