/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';

    Layer8DModuleFactory.create({
        namespace: 'Alm',
        defaultModule: 'alarms',
        defaultService: 'alarms',
        sectionSelector: 'alarms',
        initializerName: 'initializeAlm',
        requiredNamespaces: ['AlmAlarms', 'AlmEvents', 'AlmCorrelation', 'AlmPolicies', 'AlmMaintenance', 'AlmArchive']
    });
})();
