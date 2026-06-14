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
                ...f.reference('vehicleId', 'Vehicle', 'VendDeliveryTruck'),
                ...f.reference('facilityId', 'Facility', 'VendStockingFacility'),
                ...f.date('plannedDate', 'Planned Date')
            ]),
            f.section('Planned Metrics', [
                ...f.number('totalDistance', 'Distance (miles)'),
                ...f.number('totalDuration', 'Duration (min)'),
                ...f.number('estimatedFuelCost', 'Est. Fuel Cost')
            ]),
            f.section('Actual Metrics', [
                ...f.number('actualDistance', 'Actual Distance (miles)'),
                ...f.number('actualDuration', 'Actual Duration (min)'),
                ...f.number('actualFuelCost', 'Actual Fuel Cost')
            ]),
            f.section('Stops', [
                ...f.inlineTable('stops', 'Route Stops', [
                    { key: 'stopOrder', label: '#', type: 'number' },
                    { key: 'stopType', label: 'Type', type: 'text' },
                    { key: 'machineName', label: 'Name', type: 'text' },
                    { key: 'locationAddress', label: 'Address', type: 'text' },
                    { key: 'locationCity', label: 'City', type: 'text' },
                    { key: 'serviceUrgency', label: 'Urgency', type: 'text' },
                    { key: 'completionNotes', label: 'Notes', type: 'text' }
                ])
            ])
        ]),
        VendDriver: f.form('Driver', [
            f.section('Personal Info', [
                ...f.text('firstName', 'First Name', true),
                ...f.text('lastName', 'Last Name', true),
                ...f.text('phone', 'Phone'),
                ...f.text('email', 'Email'),
                ...f.date('hireDate', 'Hire Date'),
                ...f.checkbox('isActive', 'Active')
            ]),
            f.section('License & Assignment', [
                ...f.text('licenseNumber', 'License Number'),
                ...f.select('licenseClass', 'License Class', enums.LICENSE_CLASS.enum),
                ...f.reference('truckId', 'Assigned Truck', 'VendDeliveryTruck'),
                ...f.reference('homeBaseLocationId', 'Home Base', 'VendLocation')
            ]),
            f.section('Home Address', [
                ...f.address('homeAddress')
            ]),
            f.section('Weekly Schedule', [
                ...f.inlineTable('schedule', 'Schedule', [
                    { key: 'day', label: 'Day', type: 'select', options: enums.DAY_OF_WEEK.enum },
                    { key: 'startTime', label: 'Start Time', type: 'text' },
                    { key: 'startLocationId', label: 'Start Location', type: 'reference', lookupModel: 'VendLocation' }
                ])
            ])
        ]),
        VendDeliveryTruck: f.form('Delivery Truck', [
            f.section('Identity', [
                ...f.text('name', 'Name', true),
                ...f.text('plateNumber', 'Plate Number', true),
                ...f.text('vin', 'VIN')
            ]),
            f.section('Vehicle Specs', [
                ...f.text('make', 'Make'),
                ...f.text('model', 'Model'),
                ...f.number('year', 'Year'),
                ...f.select('type', 'Truck Type', enums.TRUCK_TYPE.enum),
                ...f.select('fuelType', 'Fuel Type', enums.FUEL_TYPE.enum),
                ...f.number('cargoCapacityCuFt', 'Cargo Capacity (cu ft)'),
                ...f.number('maxPayloadLbs', 'Max Payload (lbs)'),
                ...f.number('mileage', 'Mileage'),
                ...f.number('milesPerGallon', 'Miles Per Gallon')
            ]),
            f.section('Status', [
                ...f.select('status', 'Status', enums.TRUCK_STATUS.enum),
                ...f.reference('currentDriverId', 'Current Driver', 'VendDriver'),
                ...f.reference('currentRouteId', 'Current Route', 'VendRoute'),
                ...f.reference('homeDepotId', 'Home Depot', 'VendStockingFacility')
            ]),
            f.section('Maintenance', [
                ...f.date('lastMaintenanceDate', 'Last Maintenance'),
                ...f.date('nextMaintenanceDate', 'Next Maintenance'),
                ...f.number('nextMaintenanceMileage', 'Next Maintenance Mileage'),
                ...f.date('insuranceExpiry', 'Insurance Expiry'),
                ...f.date('registrationExpiry', 'Registration Expiry')
            ]),
            f.section('Capabilities', [
                ...f.checkbox('refrigerationEquipped', 'Refrigeration Equipped'),
                ...f.checkbox('cashCollectionEquipped', 'Cash Collection Equipped'),
                ...f.checkbox('coinChangerEquipped', 'Coin Changer Equipped')
            ]),
            f.section('Stock', [
                ...f.inlineTable('stock', 'Truck Stock', [
                    { key: 'productName', label: 'Product', type: 'text' },
                    { key: 'sku', label: 'SKU', type: 'text' },
                    { key: 'price', label: 'Price', type: 'money' },
                    { key: 'quantity', label: 'Qty', type: 'number' },
                    { key: 'maxQuantity', label: 'Max', type: 'number' }
                ])
            ])
        ])
    };
})();
