/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    const svc = Layer8ModuleConfigFactory.service;
    const mod = Layer8ModuleConfigFactory.module;

    Layer8ModuleConfigFactory.create({
        namespace: 'Reports',
        modules: {
            'reports': mod('Reports', '\u{1F4CA}', [
                svc('reports', 'Reports', '\u{1F4C4}', '/10/VendRpt', 'VendReport')
            ])
        },
        submodules: ['ReportsData']
    });
})();
