/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.MobileWarehouseStock = window.MobileWarehouseStock || {};

    var factory = window.Layer8EnumFactory;
    var { createStatusRenderer, renderEnum } = Layer8MRenderers;

    var PO_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Draft', 'draft', 'pending'],
        ['Submitted', 'submitted', 'active'],
        ['Confirmed', 'confirmed', 'active'],
        ['Shipped', 'shipped', 'active'],
        ['Received', 'received', 'inactive'],
        ['Closed', 'closed', 'inactive']
    ]);

    var MOVEMENT_TYPE = factory.simple([
        'Unspecified', 'Receive from Supplier', 'Transfer to Vehicle',
        'Return from Vehicle', 'Write Off', 'Adjustment'
    ]);

    var VEHICLE_LOAD_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Loading', 'loading', 'pending'],
        ['In Transit', 'transit', 'active'],
        ['Completed', 'completed', 'inactive']
    ]);

    var FACILITY_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'active'],
        ['Maintenance', 'maintenance', 'pending'],
        ['Closed', 'closed', 'terminated']
    ]);

    var SUPPLIER_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'active'],
        ['Inactive', 'inactive', 'inactive'],
        ['Suspended', 'suspended', 'terminated']
    ]);

    MobileWarehouseStock.enums = {
        PO_STATUS: PO_STATUS,
        MOVEMENT_TYPE: MOVEMENT_TYPE,
        VEHICLE_LOAD_STATUS: VEHICLE_LOAD_STATUS,
        FACILITY_STATUS: FACILITY_STATUS,
        SUPPLIER_STATUS: SUPPLIER_STATUS
    };

    MobileWarehouseStock.render = {
        poStatus: createStatusRenderer(PO_STATUS.enum, PO_STATUS.classes),
        movementType: function(value) { return renderEnum(value, MOVEMENT_TYPE.enum); },
        vehicleLoadStatus: createStatusRenderer(VEHICLE_LOAD_STATUS.enum, VEHICLE_LOAD_STATUS.classes),
        facilityStatus: createStatusRenderer(FACILITY_STATUS.enum, FACILITY_STATUS.classes),
        supplierStatus: createStatusRenderer(SUPPLIER_STATUS.enum, SUPPLIER_STATUS.classes)
    };

    MobileWarehouseStock.primaryKeys = {
        VendStockingFacility: 'facilityId',
        VendSupplier: 'supplierId',
        VendPurchaseOrder: 'orderId',
        VendStockMovement: 'movementId',
        VendVehicleLoad: 'loadId'
    };
})();
