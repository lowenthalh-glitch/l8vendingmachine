/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    var f = window.Layer8FormFactory;

    InventoryMachines.forms = {
        VendMachine: f.form('Vending Machine', [
            f.section('Machine Information', [
                ...f.text('machineId', 'Machine ID', false),
                ...f.text('name', 'Name', false),
                ...f.text('type', 'Type', false),
                ...f.text('model', 'Model', false),
                ...f.text('status', 'Status', false),
                ...f.text('deviceId', 'Payment Device ID', false),
                ...f.number('dailyTransactions', 'Daily Transactions'),
                ...f.text('lastTransactionAt', 'Last Transaction', false)
            ]),
            f.section('Location', [
                ...f.text('locationAddress', 'Address', false),
                ...f.text('locationCity', 'City', false),
                ...f.text('locationState', 'State', false),
                ...f.number('locationLat', 'Latitude'),
                ...f.number('locationLng', 'Longitude')
            ])
        ])
    };
})();
