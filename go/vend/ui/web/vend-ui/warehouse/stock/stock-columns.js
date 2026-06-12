/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var enums = WarehouseStock.enums;
    var render = WarehouseStock.render;

    WarehouseStock.columns = {
        VendStockingFacility: [
            ...col.id('facilityId'),
            ...col.col('name', 'Name'),
            ...col.col('code', 'Code'),
            ...col.status('status', 'Status', enums.FACILITY_STATUS.values, render.facilityStatus),
            ...col.number('totalStorageSqFt', 'Total Storage (sqft)'),
            ...col.number('refrigeratedStorageSqFt', 'Refrigerated (sqft)'),
            ...col.number('loadingDocks', 'Loading Docks'),
            ...col.number('maxTrucksParked', 'Max Trucks'),
            ...col.col('managerName', 'Manager'),
            ...col.col('managerPhone', 'Manager Phone')
        ],
        VendSupplier: [
            ...col.id('supplierId'),
            ...col.col('name', 'Name'),
            ...col.col('contactName', 'Contact'),
            ...col.number('leadTimeDays', 'Lead Time (days)'),
            ...col.status('status', 'Status', enums.SUPPLIER_STATUS.values, render.supplierStatus)
        ],
        VendPurchaseOrder: [
            ...col.id('orderId'),
            ...col.col('supplierId', 'Supplier'),
            ...col.col('facilityId', 'Facility'),
            ...col.status('status', 'Status', enums.PO_STATUS.values, render.poStatus),
            ...col.date('orderDate', 'Order Date'),
            ...col.money('totalAmount', 'Total')
        ],
        VendStockMovement: [
            ...col.id('movementId'),
            ...col.col('facilityId', 'Facility'),
            ...col.col('productId', 'Product'),
            ...col.enum('movementType', 'Type', null, render.movementType),
            ...col.number('quantity', 'Quantity'),
            ...col.date('timestamp', 'Timestamp')
        ],
        VendVehicleLoad: [
            ...col.id('loadId'),
            ...col.col('routeId', 'Route'),
            ...col.col('driverId', 'Driver'),
            ...col.date('loadDate', 'Load Date'),
            ...col.status('status', 'Status', enums.VEHICLE_LOAD_STATUS.values, render.vehicleLoadStatus)
        ]
    };

    WarehouseStock.primaryKeys = {
        VendStockingFacility: 'facilityId',
        VendSupplier: 'supplierId',
        VendPurchaseOrder: 'orderId',
        VendStockMovement: 'movementId',
        VendVehicleLoad: 'loadId'
    };
})();
