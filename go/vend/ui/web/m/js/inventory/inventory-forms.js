/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;
    var enums = MobileInventoryProducts.enums;

    MobileInventoryProducts.forms = {
        VendProduct: f.form('Product', [
            f.section('Product Information', [
                ...f.text('name', 'Name', true),
                ...f.select('category', 'Category', enums.PRODUCT_CATEGORY.enum),
                ...f.text('upc', 'UPC'),
                ...f.money('price', 'Price'),
                ...f.number('shelfLifeDays', 'Shelf Life (Days)'),
                ...f.checkbox('isActive', 'Active'),
                ...f.textarea('description', 'Description'),
                ...f.reference('supplierId', 'Supplier', 'VendSupplier')
            ])
        ]),
        VendPlanogram: f.form('Planogram', [
            f.section('Planogram Info', [
                ...f.text('name', 'Name'),
                ...f.reference('machineId', 'Machine', 'VendMachine'),
                ...f.checkbox('isActive', 'Active'),
                ...f.date('effectiveDate', 'Effective Date')
            ])
        ]),
        VendRestockOrder: f.form('Restock Order', [
            f.section('Order Info', [
                ...f.reference('machineId', 'Machine', 'VendMachine'),
                ...f.reference('routeId', 'Route', 'VendRoute'),
                ...f.select('status', 'Status', enums.WORK_ORDER_STATUS.enum),
                ...f.text('urgency', 'Urgency'),
                ...f.date('dueDate', 'Due Date'),
                ...f.textarea('notes', 'Notes')
            ])
        ])
    };
})();
