/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.MobileRoutePlanning = window.MobileRoutePlanning || {};

    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum } = Layer8MRenderers;

    const ROUTE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Planned', 'planned', 'pending'],
        ['In Progress', 'inprogress', 'active'],
        ['Completed', 'completed', 'inactive'],
        ['Cancelled', 'cancelled', 'terminated']
    ]);

    const LICENSE_CLASS = factory.simple([
        'Unspecified', 'Class C', 'Class B', 'Class A'
    ]);

    const DAY_OF_WEEK = factory.simple([
        'Unspecified', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'
    ]);

    const TRUCK_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'active'],
        ['Maintenance', 'maintenance', 'pending'],
        ['En Route', 'enroute', 'active'],
        ['Decommissioned', 'decommissioned', 'terminated']
    ]);

    const FUEL_TYPE = factory.simple([
        'Unspecified', 'Gasoline', 'Diesel', 'Electric', 'Hybrid'
    ]);

    const TRUCK_TYPE = factory.simple([
        'Unspecified', 'Box Truck', 'Cargo Van', 'Refrigerated', 'Sprinter', 'Pickup'
    ]);

    MobileRoutePlanning.enums = {
        ROUTE_STATUS: ROUTE_STATUS,
        LICENSE_CLASS: LICENSE_CLASS,
        DAY_OF_WEEK: DAY_OF_WEEK,
        TRUCK_STATUS: TRUCK_STATUS,
        FUEL_TYPE: FUEL_TYPE,
        TRUCK_TYPE: TRUCK_TYPE
    };

    MobileRoutePlanning.render = {
        routeStatus: createStatusRenderer(ROUTE_STATUS.enum, ROUTE_STATUS.classes),
        licenseClass: function(value) { return renderEnum(value, LICENSE_CLASS.enum); },
        dayOfWeek: function(value) { return renderEnum(value, DAY_OF_WEEK.enum); },
        truckStatus: createStatusRenderer(TRUCK_STATUS.enum, TRUCK_STATUS.classes),
        fuelType: function(value) { return renderEnum(value, FUEL_TYPE.enum); },
        truckType: function(value) { return renderEnum(value, TRUCK_TYPE.enum); }
    };

    MobileRoutePlanning.primaryKeys = {
        VendRoute: 'routeId',
        VendDriver: 'driverId',
        VendDeliveryTruck: 'truckId'
    };
})();
