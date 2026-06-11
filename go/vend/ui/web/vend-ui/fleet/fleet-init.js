/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    Layer8DModuleFactory.create({
        namespace: 'Fleet',
        defaultModule: 'machines',
        defaultService: 'machines',
        sectionSelector: 'machines',
        initializerName: 'initializeFleet',
        requiredNamespaces: ['FleetMachines']
    });
})();
