/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    Layer8DModuleFactory.create({
        namespace: 'Nayax',
        defaultModule: 'machines',
        defaultService: 'machines',
        sectionSelector: 'machines',
        initializerName: 'initializeNayax',
        requiredNamespaces: ['NayaxMachines']
    });

    // Override row click for VendMachine to use custom detail modal
    var origInit = window.initializeNayax;
    window.initializeNayax = function() {
        if (origInit) origInit();
        if (window.Nayax && Nayax._showDetailsModal) {
            var origDetail = Nayax._showDetailsModal;
            Nayax._showDetailsModal = function(service, item, id) {
                if (service.model === 'VendMachine' && typeof showVendMachineDetail === 'function') {
                    showVendMachineDetail(item);
                } else {
                    origDetail.call(Nayax, service, item, id);
                }
            };
        }
    };
})();
