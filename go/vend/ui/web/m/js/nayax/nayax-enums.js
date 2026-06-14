/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.MobileNayaxMachines = window.MobileNayaxMachines || {};

    var STATUS_MAP = {
        'online': 'Online', 'offline': 'Offline',
        'warning': 'Warning', 'pending': 'Pending'
    };
    var STATUS_CLASSES = {
        'online': 'active', 'offline': 'inactive',
        'warning': 'pending', 'pending': 'pending'
    };
    var TYPE_MAP = {
        'vending_snack': 'Snack', 'vending_drink': 'Drink', 'vending_combo': 'Combo',
        'coffee': 'Coffee', 'ev_charger': 'EV Charger', 'laundry': 'Laundry', 'car_wash': 'Car Wash'
    };

    MobileNayaxMachines.enums = { STATUS_MAP: STATUS_MAP, TYPE_MAP: TYPE_MAP };
    MobileNayaxMachines.render = {
        machineStatus: function(value) {
            var label = STATUS_MAP[value] || value || 'Unknown';
            var cls = STATUS_CLASSES[value] || '';
            return cls ? '<span class="layer8d-status-badge layer8d-status-' + cls + '">' + label + '</span>' : label;
        },
        machineType: function(value) { return TYPE_MAP[value] || value || 'Unknown'; }
    };
    MobileNayaxMachines.primaryKeys = {
        VendMachine: 'machineId', VendMachineGroup: 'groupId', VendLocation: 'locationId'
    };
})();
