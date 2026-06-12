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
            ...col.enum('licenseClass', 'License', null, render.licenseClass),
            ...col.date('hireDate', 'Hire Date'),
            ...col.boolean('isActive', 'Active')
        ],
        VendDeliveryTruck: [
            ...col.id('truckId'),
            ...col.col('name', 'Name'),
            ...col.col('plateNumber', 'Plate'),
            ...col.enum('type', 'Type', null, render.truckType),
            ...col.col('make', 'Make'),
            ...col.col('model', 'Model'),
            ...col.number('year', 'Year'),
            ...col.status('status', 'Status', enums.TRUCK_STATUS.values, render.truckStatus),
            ...col.number('milesPerGallon', 'MPG'),
            ...col.number('mileage', 'Mileage'),
            ...col.number('cargoCapacityCuFt', 'Cargo (cu ft)'),
            ...col.enum('fuelType', 'Fuel', null, render.fuelType),
            ...col.boolean('refrigerationEquipped', 'Refrigerated'),
            ...col.boolean('cashCollectionEquipped', 'Cash Collection')
        ]
    };

    RoutePlanning.primaryKeys = {
        VendRoute: 'routeId',
        VendDriver: 'driverId',
        VendDeliveryTruck: 'truckId'
    };
})();
