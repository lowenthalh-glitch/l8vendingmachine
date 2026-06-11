/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum } = Layer8MRenderers;

    window.MobileFleetMachines = window.MobileFleetMachines || {};

    const MACHINE_TYPE = factory.simple([
        'Unspecified', 'Locker', 'Refrigerated Beverage', 'Combo',
        'Snack', 'Frozen', 'Fresh Food'
    ]);

    const MACHINE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Operational', 'operational', 'active'],
        ['Out of Service', 'outofservice', 'terminated'],
        ['Maintenance', 'maintenance', 'pending'],
        ['Offline', 'offline', 'inactive'],
        ['Decommissioned', 'decommissioned', 'inactive']
    ]);

    MobileFleetMachines.enums = {
        MACHINE_TYPE: MACHINE_TYPE,
        MACHINE_STATUS: MACHINE_STATUS
    };

    MobileFleetMachines.render = {
        machineType: (v) => renderEnum(v, MACHINE_TYPE.enum),
        machineStatus: createStatusRenderer(MACHINE_STATUS.enum, MACHINE_STATUS.classes)
    };

    MobileFleetMachines.primaryKeys = {
        VendMachine: 'machineId',
        VendMachineGroup: 'groupId',
        VendLocation: 'locationId'
    };
})();
