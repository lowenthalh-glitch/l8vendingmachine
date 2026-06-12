/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;
    var enums = WarehouseStock.enums;

    WarehouseStock.forms = {
        VendStockingFacility: f.form('Stocking Facility', [
            f.section('Identity', [
                ...f.text('name', 'Name', true),
                ...f.text('code', 'Code')
            ]),
            f.section('Location', [
                ...f.address('address'),
                ...f.text('timezone', 'Timezone')
            ]),
            f.section('Capacity', [
                ...f.number('totalStorageSqFt', 'Total Storage (sqft)'),
                ...f.number('refrigeratedStorageSqFt', 'Refrigerated Storage (sqft)'),
                ...f.number('loadingDocks', 'Loading Docks'),
                ...f.number('maxTrucksParked', 'Max Trucks Parked')
            ]),
            f.section('Operations', [
                ...f.select('status', 'Status', enums.FACILITY_STATUS.enum),
                ...f.text('operatingHoursStart', 'Operating Hours Start'),
                ...f.text('operatingHoursEnd', 'Operating Hours End')
            ]),
            f.section('Stock', [
                ...f.inlineTable('stock', 'Stock Items', [
                    { key: 'productName', label: 'Product', type: 'text' },
                    { key: 'sku', label: 'SKU', type: 'text' },
                    { key: 'price', label: 'Price', type: 'money' },
                    { key: 'quantity', label: 'Qty', type: 'number' },
                    { key: 'maxQuantity', label: 'Max', type: 'number' },
                    { key: 'reorderPoint', label: 'Reorder Point', type: 'number' }
                ])
            ]),
            f.section('Contacts', [
                ...f.text('managerName', 'Manager Name'),
                ...f.text('managerPhone', 'Manager Phone'),
                ...f.text('managerEmail', 'Manager Email')
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
                ...f.reference('facilityId', 'Facility', 'VendStockingFacility'),
                ...f.select('status', 'Status', enums.PO_STATUS.enum),
                ...f.date('orderDate', 'Order Date'),
                ...f.date('expectedDelivery', 'Expected Delivery'),
                ...f.textarea('notes', 'Notes')
            ])
        ]),
        VendStockMovement: f.form('Stock Movement', [
            f.section('Movement Details', [
                ...f.text('movementId', 'Movement ID', false),
                ...f.text('facilityId', 'Facility', false),
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
        ])
    };
})();
