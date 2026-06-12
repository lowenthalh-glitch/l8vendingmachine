/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var enums = MobileRoutePlanning.enums;
    var render = MobileRoutePlanning.render;
    var col = window.Layer8ColumnFactory;

    MobileRoutePlanning.columns = {
        VendRoute: [
            ...col.id('routeId'),
            { key: 'name', label: 'Name', primary: true, sortKey: 'name', filterKey: 'name' },
            { key: 'status', label: 'Status', secondary: true, sortKey: 'status', filterKey: 'status',
              enumValues: enums.ROUTE_STATUS.values,
              render: (item) => render.routeStatus(item.status) },
            ...col.col('driverId', 'Driver'),
            ...col.col('vehicleId', 'Vehicle'),
            ...col.date('plannedDate', 'Planned Date'),
            ...col.number('totalDistance', 'Total Distance')
        ],
        VendDriver: [
            ...col.id('driverId'),
            { key: 'firstName', label: 'First Name', secondary: true, sortKey: 'firstName', filterKey: 'firstName' },
            { key: 'lastName', label: 'Last Name', primary: true, sortKey: 'lastName', filterKey: 'lastName' },
            ...col.col('phone', 'Phone'),
            ...col.col('email', 'Email'),
            ...col.enum('licenseClass', 'License', null, render.licenseClass),
            ...col.date('hireDate', 'Hire Date'),
            ...col.boolean('isActive', 'Active')
        ],
        VendDeliveryTruck: [
            ...col.id('truckId'),
            { key: 'name', label: 'Name', primary: true, sortKey: 'name', filterKey: 'name' },
            { key: 'status', label: 'Status', secondary: true, sortKey: 'status', filterKey: 'status',
              enumValues: enums.TRUCK_STATUS.values,
              render: (item) => render.truckStatus(item.status) },
            ...col.col('plateNumber', 'Plate'),
            ...col.col('make', 'Make'),
            ...col.col('model', 'Model'),
            ...col.number('year', 'Year'),
            ...col.number('milesPerGallon', 'MPG'),
            ...col.number('mileage', 'Mileage'),
            ...col.enum('fuelType', 'Fuel', null, render.fuelType),
            ...col.boolean('refrigerationEquipped', 'Refrigerated')
        ]
    };
})();
