/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8DModuleFactory.create({
        namespace: 'Compliance',
        defaultModule: 'inspections',
        defaultService: 'inspections',
        sectionSelector: 'inspections',
        initializerName: 'initializeCompliance',
        requiredNamespaces: ['ComplianceInspections']
    });
})();
