/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var enums = MobileWarehouseStock.enums;
    var render = MobileWarehouseStock.render;

    MobileWarehouseStock.columns = {
        VendStockingFacility: [
            ...col.id('facilityId'),
            { key: 'name', label: 'Name', primary: true, sortKey: 'name', filterKey: 'name' },
            ...col.col('code', 'Code'),
            ...col.status('status', 'Status', enums.FACILITY_STATUS.values, render.facilityStatus),
            ...col.number('totalStorageSqFt', 'Storage (sqft)'),
            ...col.number('loadingDocks', 'Loading Docks'),
            ...col.col('managerName', 'Manager')
        ],
        VendSupplier: [
            ...col.id('supplierId'),
            { key: 'name', label: 'Name', primary: true, sortKey: 'name', filterKey: 'name' },
            ...col.col('contactName', 'Contact'),
            ...col.number('leadTimeDays', 'Lead Time (days)'),
            ...col.status('status', 'Status', enums.SUPPLIER_STATUS.values, render.supplierStatus)
        ],
        VendPurchaseOrder: [
            ...col.id('orderId'),
            { key: 'supplierId', label: 'Supplier', primary: true, sortKey: 'supplierId', filterKey: 'supplierId' },
            ...col.col('facilityId', 'Facility'),
            ...col.status('status', 'Status', enums.PO_STATUS.values, render.poStatus),
            ...col.date('orderDate', 'Order Date'),
            ...col.money('totalAmount', 'Total')
        ],
        VendStockMovement: [
            ...col.id('movementId'),
            { key: 'facilityId', label: 'Facility', primary: true, sortKey: 'facilityId', filterKey: 'facilityId' },
            ...col.col('productId', 'Product'),
            ...col.enum('movementType', 'Type', null, render.movementType),
            ...col.number('quantity', 'Quantity'),
            ...col.date('timestamp', 'Timestamp')
        ],
        VendVehicleLoad: [
            ...col.id('loadId'),
            { key: 'routeId', label: 'Route', primary: true, sortKey: 'routeId', filterKey: 'routeId' },
            ...col.col('driverId', 'Driver'),
            ...col.date('loadDate', 'Load Date'),
            ...col.status('status', 'Status', enums.VEHICLE_LOAD_STATUS.values, render.vehicleLoadStatus)
        ]
    };
})();
