/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;

    MobileNayaxMachines.forms = {
        VendMachine: f.form('Management System', [
            f.section('Management Info', [
                ...f.text('machineId', 'Management IP', false)
            ])
        ]),
        VendMachineGroup: f.form('Machine Group', [
            f.section('Group Info', [
                ...f.text('name', 'Name', true),
                ...f.text('description', 'Description'),
                ...f.text('region', 'Region'),
                ...f.number('machineCount', 'Machine Count')
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
