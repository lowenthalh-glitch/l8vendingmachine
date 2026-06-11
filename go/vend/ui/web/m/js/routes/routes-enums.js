/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.MobileRoutePlanning = window.MobileRoutePlanning || {};

    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer } = Layer8MRenderers;

    const ROUTE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Planned', 'planned', 'pending'],
        ['In Progress', 'inprogress', 'active'],
        ['Completed', 'completed', 'inactive'],
        ['Cancelled', 'cancelled', 'terminated']
    ]);

    MobileRoutePlanning.enums = {
        ROUTE_STATUS: ROUTE_STATUS
    };

    MobileRoutePlanning.render = {
        routeStatus: createStatusRenderer(ROUTE_STATUS.enum, ROUTE_STATUS.classes)
    };

    MobileRoutePlanning.primaryKeys = {
        VendRoute: 'routeId',
        VendDriver: 'driverId'
    };
})();
