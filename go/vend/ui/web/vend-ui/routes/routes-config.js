/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    const svc = Layer8ModuleConfigFactory.service;
    const mod = Layer8ModuleConfigFactory.module;

    Layer8ModuleConfigFactory.create({
        namespace: 'Routes',
        modules: {
            'routes': mod('Routes', '', [
                svc('routes', 'Routes', '', '/10/Route', 'VendRoute',
                    { supportedViews: ['table', 'gantt'] }),
                svc('drivers', 'Drivers', '', '/10/Driver', 'VendDriver')
            ])
        },
        submodules: ['RoutePlanning']
    });
})();
