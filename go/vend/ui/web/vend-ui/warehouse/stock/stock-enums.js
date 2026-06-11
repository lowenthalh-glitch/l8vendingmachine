/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.WarehouseStock = window.WarehouseStock || {};

    var factory = window.Layer8EnumFactory;
    var { createStatusRenderer, renderEnum } = Layer8DRenderers;

    var PO_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Draft', 'draft', 'layer8d-status-pending'],
        ['Submitted', 'submitted', 'layer8d-status-active'],
        ['Confirmed', 'confirmed', 'layer8d-status-active'],
        ['Shipped', 'shipped', 'layer8d-status-active'],
        ['Received', 'received', 'layer8d-status-inactive'],
        ['Closed', 'closed', 'layer8d-status-inactive']
    ]);

    var MOVEMENT_TYPE = factory.simple([
        'Unspecified', 'Receive from Supplier', 'Transfer to Vehicle',
        'Return from Vehicle', 'Write Off', 'Adjustment'
    ]);

    var VEHICLE_LOAD_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Loading', 'loading', 'layer8d-status-pending'],
        ['In Transit', 'transit', 'layer8d-status-active'],
        ['Completed', 'completed', 'layer8d-status-inactive']
    ]);

    var SUPPLIER_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'layer8d-status-active'],
        ['Inactive', 'inactive', 'layer8d-status-inactive'],
        ['Suspended', 'suspended', 'layer8d-status-terminated']
    ]);

    WarehouseStock.enums = {
        PO_STATUS: PO_STATUS,
        MOVEMENT_TYPE: MOVEMENT_TYPE,
        VEHICLE_LOAD_STATUS: VEHICLE_LOAD_STATUS,
        SUPPLIER_STATUS: SUPPLIER_STATUS
    };

    WarehouseStock.render = {
        poStatus: createStatusRenderer(PO_STATUS.enum, PO_STATUS.classes),
        movementType: function(value) { return renderEnum(value, MOVEMENT_TYPE.enum); },
        vehicleLoadStatus: createStatusRenderer(VEHICLE_LOAD_STATUS.enum, VEHICLE_LOAD_STATUS.classes),
        supplierStatus: createStatusRenderer(SUPPLIER_STATUS.enum, SUPPLIER_STATUS.classes)
    };
})();
