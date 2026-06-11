/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.RoutePlanning = window.RoutePlanning || {};

    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer } = Layer8DRenderers;

    const ROUTE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Planned', 'planned', 'layer8d-status-pending'],
        ['In Progress', 'inprogress', 'layer8d-status-active'],
        ['Completed', 'completed', 'layer8d-status-inactive'],
        ['Cancelled', 'cancelled', 'layer8d-status-terminated']
    ]);

    RoutePlanning.enums = {
        ROUTE_STATUS: ROUTE_STATUS
    };

    RoutePlanning.render = {
        routeStatus: createStatusRenderer(ROUTE_STATUS.enum, ROUTE_STATUS.classes)
    };
})();
