/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var render = MobileInventoryProducts.render;

    MobileInventoryProducts.columns = {
        VendProduct: [
            ...col.id('productId'),
            { key: 'name', label: 'Name', primary: true, sortKey: 'name', filterKey: 'name' },
            { key: 'category', label: 'Category', secondary: true, sortKey: 'category', filterKey: 'category',
              render: (item) => render.productCategory(item.category) },
            ...col.col('upc', 'UPC'),
            ...col.money('price', 'Price'),
            ...col.number('shelfLifeDays', 'Shelf Life (Days)'),
            ...col.boolean('isActive', 'Active')
        ],
        VendPlanogram: [
            ...col.id('planogramId'),
            ...col.col('machineId', 'Machine'),
            { key: 'name', label: 'Name', primary: true, sortKey: 'name', filterKey: 'name' },
            ...col.boolean('isActive', 'Active'),
            ...col.date('effectiveDate', 'Effective Date')
        ],
        VendRestockOrder: [
            { key: 'orderId', label: 'Order ID', primary: true, sortKey: 'orderId', filterKey: 'orderId' },
            ...col.col('machineId', 'Machine'),
            { key: 'status', label: 'Status', secondary: true, sortKey: 'status', filterKey: 'status',
              render: (item) => render.workOrderStatus(item.status) },
            ...col.col('urgency', 'Urgency'),
            ...col.date('createdDate', 'Created'),
            ...col.date('dueDate', 'Due Date')
        ]
    };
})();
