/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var enums = WarehouseStock.enums;
    var render = WarehouseStock.render;

    WarehouseStock.columns = {
        VendWarehouse: [
            ...col.id('warehouseId'),
            ...col.col('name', 'Name'),
            ...col.col('region', 'Region'),
            ...col.number('capacitySqft', 'Capacity (sqft)'),
            ...col.col('contactName', 'Contact'),
            ...col.boolean('isActive', 'Active')
        ],
        VendWarehouseStock: [
            ...col.id('stockId'),
            ...col.col('warehouseId', 'Warehouse'),
            ...col.col('productId', 'Product'),
            ...col.number('quantityOnHand', 'Qty On Hand'),
            ...col.number('reorderPoint', 'Reorder Point'),
            ...col.number('reorderQuantity', 'Reorder Qty')
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
            ...col.col('warehouseId', 'Warehouse'),
            ...col.status('status', 'Status', enums.PO_STATUS.values, render.poStatus),
            ...col.date('orderDate', 'Order Date'),
            ...col.money('totalAmount', 'Total')
        ],
        VendStockMovement: [
            ...col.id('movementId'),
            ...col.col('warehouseId', 'Warehouse'),
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
        VendWarehouse: 'warehouseId',
        VendWarehouseStock: 'stockId',
        VendSupplier: 'supplierId',
        VendPurchaseOrder: 'orderId',
        VendStockMovement: 'movementId',
        VendVehicleLoad: 'loadId'
    };
})();
