/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    const svc = Layer8ModuleConfigFactory.service;
    const mod = Layer8ModuleConfigFactory.module;

    Layer8ModuleConfigFactory.create({
        namespace: 'Maintenance',
        modules: {
            'work-orders': mod('Work Orders', '', [
                svc('work-orders', 'Work Orders', '', '/10/WorkOrder', 'VendWorkOrder',
                    { supportedViews: ['table', 'kanban', 'gantt'] }),
                svc('service-visits', 'Service Visits', '', '/10/SvcVisit', 'VendServiceVisit')
            ])
        },
        submodules: ['MaintenanceAlerts']
    });
})();
