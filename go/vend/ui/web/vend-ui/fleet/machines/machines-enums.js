/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    window.FleetMachines = window.FleetMachines || {};

    var STATUS_MAP = {
        'online': 'Online',
        'offline': 'Offline',
        'warning': 'Warning',
        'pending': 'Pending'
    };

    var STATUS_CLASSES = {
        'online': 'layer8d-status-active',
        'offline': 'layer8d-status-inactive',
        'warning': 'layer8d-status-pending',
        'pending': 'layer8d-status-pending'
    };

    var TYPE_MAP = {
        'vending_snack': 'Snack',
        'vending_drink': 'Drink',
        'vending_combo': 'Combo',
        'coffee': 'Coffee',
        'ev_charger': 'EV Charger',
        'laundry': 'Laundry',
        'car_wash': 'Car Wash'
    };

    FleetMachines.enums = {
        STATUS_MAP: STATUS_MAP,
        TYPE_MAP: TYPE_MAP
    };

    FleetMachines.render = {
        machineStatus: function(value) {
            var label = STATUS_MAP[value] || value || 'Unknown';
            var cls = STATUS_CLASSES[value] || '';
            if (cls) {
                return '<span class="layer8d-status-badge ' + cls + '">' + label + '</span>';
            }
            return label;
        },
        machineType: function(value) {
            return TYPE_MAP[value] || value || 'Unknown';
        }
    };

    FleetMachines.primaryKeys = {
        VendFleetMachine: 'machineId',
        VendMachineGroup: 'groupId',
        VendLocation: 'locationId',
        VendFleetInventory: 'summaryId',
        VendMachineProfile: 'profileId'
    };
})();
