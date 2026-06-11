/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var enums = MobileRoutePlanning.enums;
    var f = window.Layer8FormFactory;

    MobileRoutePlanning.forms = {
        VendRoute: f.form('Route', [
            f.section('Route Details', [
                ...f.text('name', 'Name', true),
                ...f.textarea('description', 'Description'),
                ...f.select('status', 'Status', enums.ROUTE_STATUS.enum),
                ...f.reference('driverId', 'Driver', 'VendDriver'),
                ...f.text('vehicleId', 'Vehicle'),
                ...f.date('plannedDate', 'Planned Date')
            ])
        ]),
        VendDriver: f.form('Driver', [
            f.section('Driver Information', [
                ...f.text('firstName', 'First Name', true),
                ...f.text('lastName', 'Last Name', true),
                ...f.text('phone', 'Phone'),
                ...f.text('email', 'Email'),
                ...f.text('licenseNumber', 'License Number'),
                ...f.text('vehicleId', 'Vehicle'),
                ...f.checkbox('isActive', 'Active')
            ])
        ])
    };
})();
