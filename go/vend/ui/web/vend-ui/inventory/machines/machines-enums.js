/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    window.InventoryMachines = window.InventoryMachines || {};

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

    InventoryMachines.enums = {
        STATUS_MAP: STATUS_MAP,
        STATUS_CLASSES: STATUS_CLASSES,
        TYPE_MAP: TYPE_MAP
    };

    InventoryMachines.render = {
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

    InventoryMachines.primaryKeys = {
        VendMachine: 'machineId'
    };

    // Transform: flatten the machines map from VendMachine into individual rows.
    // The server returns one VendMachine per management API with a "machines" map.
    // The UI shows each machine as a root-level row.
    InventoryMachines.transformData = function(item) {
        if (!item || !item.machines) {
            return item;
        }
        // Return the machines map values as an array
        var rows = [];
        var machines = item.machines;
        for (var key in machines) {
            if (machines.hasOwnProperty(key)) {
                var m = machines[key];
                m.machineId = m.machineId || key;
                rows.push(m);
            }
        }
        return rows;
    };
})();
