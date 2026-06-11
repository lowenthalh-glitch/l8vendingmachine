/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;
    var enums = WarehouseStock.enums;

    WarehouseStock.forms = {
        VendWarehouse: f.form('Warehouse', [
            f.section('Warehouse Information', [
                ...f.text('name', 'Name', true),
                ...f.text('region', 'Region'),
                ...f.number('capacitySqft', 'Capacity (sqft)'),
                ...f.text('contactName', 'Contact Name'),
                ...f.text('contactPhone', 'Contact Phone'),
                ...f.text('contactEmail', 'Contact Email'),
                ...f.checkbox('isActive', 'Active')
            ])
        ]),
        VendSupplier: f.form('Supplier', [
            f.section('Supplier Information', [
                ...f.text('name', 'Name', true),
                ...f.text('contactName', 'Contact Name'),
                ...f.text('contactPhone', 'Contact Phone'),
                ...f.text('contactEmail', 'Contact Email'),
                ...f.number('leadTimeDays', 'Lead Time (days)'),
                ...f.text('paymentTerms', 'Payment Terms'),
                ...f.select('status', 'Status', enums.SUPPLIER_STATUS.enum)
            ])
        ]),
        VendPurchaseOrder: f.form('Purchase Order', [
            f.section('Order Information', [
                ...f.reference('supplierId', 'Supplier', 'VendSupplier'),
                ...f.reference('warehouseId', 'Warehouse', 'VendWarehouse'),
                ...f.select('status', 'Status', enums.PO_STATUS.enum),
                ...f.date('orderDate', 'Order Date'),
                ...f.date('expectedDelivery', 'Expected Delivery'),
                ...f.textarea('notes', 'Notes')
            ])
        ]),
        VendStockMovement: f.form('Stock Movement', [
            f.section('Movement Details', [
                ...f.text('movementId', 'Movement ID', false),
                ...f.text('warehouseId', 'Warehouse', false),
                ...f.text('productId', 'Product', false),
                ...f.text('movementType', 'Type', false),
                ...f.text('quantity', 'Quantity', false),
                ...f.text('timestamp', 'Timestamp', false)
            ])
        ]),
        VendVehicleLoad: f.form('Vehicle Load', [
            f.section('Load Information', [
                ...f.reference('routeId', 'Route', 'VendRoute'),
                ...f.reference('driverId', 'Driver', 'VendDriver'),
                ...f.text('vehicleId', 'Vehicle ID'),
                ...f.date('loadDate', 'Load Date'),
                ...f.select('status', 'Status', enums.VEHICLE_LOAD_STATUS.enum)
            ])
        ]),
        VendWarehouseStock: f.form('Warehouse Stock', [
            f.section('Stock Information', [
                ...f.reference('warehouseId', 'Warehouse', 'VendWarehouse'),
                ...f.reference('productId', 'Product', 'VendProduct'),
                ...f.number('quantityOnHand', 'Quantity On Hand'),
                ...f.number('reorderPoint', 'Reorder Point'),
                ...f.number('reorderQuantity', 'Reorder Quantity')
            ])
        ])
    };
})();
