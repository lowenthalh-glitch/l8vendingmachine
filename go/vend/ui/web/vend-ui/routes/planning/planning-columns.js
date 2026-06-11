/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var enums = RoutePlanning.enums;
    var render = RoutePlanning.render;
    var col = window.Layer8ColumnFactory;

    RoutePlanning.columns = {
        VendRoute: [
            ...col.id('routeId'),
            ...col.col('name', 'Name'),
            ...col.status('status', 'Status', enums.ROUTE_STATUS.values, render.routeStatus),
            ...col.col('driverId', 'Driver'),
            ...col.col('vehicleId', 'Vehicle'),
            ...col.date('plannedDate', 'Planned Date'),
            ...col.number('totalDistance', 'Total Distance')
        ],
        VendDriver: [
            ...col.id('driverId'),
            ...col.col('firstName', 'First Name'),
            ...col.col('lastName', 'Last Name'),
            ...col.col('phone', 'Phone'),
            ...col.col('email', 'Email'),
            ...col.boolean('isActive', 'Active')
        ]
    };

    RoutePlanning.primaryKeys = {
        VendRoute: 'routeId',
        VendDriver: 'driverId'
    };
})();
