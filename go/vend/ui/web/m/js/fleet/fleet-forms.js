/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;
    var enums = MobileFleetMachines.enums;

    MobileFleetMachines.forms = {
        VendMachine: f.form('Vending Machine', [
            f.section('Machine Information', [
                ...f.text('serialNumber', 'Serial Number', true),
                ...f.text('model', 'Model', true),
                ...f.text('manufacturer', 'Manufacturer'),
                ...f.select('machineType', 'Type', enums.MACHINE_TYPE.enum),
                ...f.select('status', 'Status', enums.MACHINE_STATUS.enum),
                ...f.text('firmwareVersion', 'Firmware'),
                ...f.number('totalSlots', 'Total Slots')
            ]),
            f.section('Location', [
                ...f.reference('locationId', 'Location', 'VendLocation'),
                ...f.reference('groupId', 'Group', 'VendMachineGroup'),
                ...f.reference('routeId', 'Route', 'VendRoute'),
                ...f.date('installedDate', 'Install Date')
            ])
        ]),
        VendMachineGroup: f.form('Machine Group', [
            f.section('Group Info', [
                ...f.text('name', 'Name', true),
                ...f.text('description', 'Description'),
                ...f.text('region', 'Region'),
                ...f.text('operatorId', 'Operator ID')
            ])
        ]),
        VendLocation: f.form('Location', [
            f.section('Location Info', [
                ...f.text('name', 'Name', true),
                ...f.text('locationType', 'Type'),
                ...f.text('timezone', 'Timezone')
            ]),
            f.section('Contact', [
                ...f.text('contactName', 'Contact'),
                ...f.text('contactPhone', 'Phone'),
                ...f.text('contactEmail', 'Email')
            ])
        ])
    };
})();
