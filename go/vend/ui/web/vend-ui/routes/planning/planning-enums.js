/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.RoutePlanning = window.RoutePlanning || {};

    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum } = Layer8DRenderers;

    const ROUTE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Planned', 'planned', 'layer8d-status-pending'],
        ['In Progress', 'inprogress', 'layer8d-status-active'],
        ['Completed', 'completed', 'layer8d-status-inactive'],
        ['Cancelled', 'cancelled', 'layer8d-status-terminated']
    ]);

    const LICENSE_CLASS = factory.simple([
        'Unspecified', 'Class C', 'Class B', 'Class A'
    ]);

    const DAY_OF_WEEK = factory.simple([
        'Unspecified', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'
    ]);

    const TRUCK_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'layer8d-status-active'],
        ['Maintenance', 'maintenance', 'layer8d-status-pending'],
        ['En Route', 'enroute', 'layer8d-status-active'],
        ['Decommissioned', 'decommissioned', 'layer8d-status-terminated']
    ]);

    const FUEL_TYPE = factory.simple([
        'Unspecified', 'Gasoline', 'Diesel', 'Electric', 'Hybrid'
    ]);

    const TRUCK_TYPE = factory.simple([
        'Unspecified', 'Box Truck', 'Cargo Van', 'Refrigerated', 'Sprinter', 'Pickup'
    ]);

    RoutePlanning.enums = {
        ROUTE_STATUS: ROUTE_STATUS,
        LICENSE_CLASS: LICENSE_CLASS,
        DAY_OF_WEEK: DAY_OF_WEEK,
        TRUCK_STATUS: TRUCK_STATUS,
        FUEL_TYPE: FUEL_TYPE,
        TRUCK_TYPE: TRUCK_TYPE
    };

    RoutePlanning.render = {
        routeStatus: createStatusRenderer(ROUTE_STATUS.enum, ROUTE_STATUS.classes),
        licenseClass: function(value) { return renderEnum(value, LICENSE_CLASS.enum); },
        dayOfWeek: function(value) { return renderEnum(value, DAY_OF_WEEK.enum); },
        truckStatus: createStatusRenderer(TRUCK_STATUS.enum, TRUCK_STATUS.classes),
        fuelType: function(value) { return renderEnum(value, FUEL_TYPE.enum); },
        truckType: function(value) { return renderEnum(value, TRUCK_TYPE.enum); }
    };
})();
