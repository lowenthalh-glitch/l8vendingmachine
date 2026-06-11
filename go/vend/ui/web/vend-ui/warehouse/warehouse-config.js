/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    const svc = Layer8ModuleConfigFactory.service;
    const mod = Layer8ModuleConfigFactory.module;

    Layer8ModuleConfigFactory.create({
        namespace: 'Warehouse',
        modules: {
            'stock': mod('Stock', '\u{1F4E6}', [
                svc('warehouses', 'Warehouses', '\u{1F3ED}', '/10/Warehouse', 'VendWarehouse'),
                svc('warehouse-stock', 'Stock', '\u{1F4CB}', '/10/WhseStock', 'VendWarehouseStock'),
                svc('suppliers', 'Suppliers', '\u{1F4E5}', '/10/Supplier', 'VendSupplier'),
                svc('purchase-orders', 'PO', '\u{1F4C4}', '/10/PurchOrd', 'VendPurchaseOrder',
                    { supportedViews: ['table', 'kanban'] }),
                svc('movements', 'Movements', '\u{1F504}', '/10/StockMove', 'VendStockMovement',
                    { supportedViews: ['table', 'timeline'] }),
                svc('vehicle-loads', 'Loads', '\u{1F69A}', '/10/VehLoad', 'VendVehicleLoad')
            ])
        },
        submodules: ['WarehouseStock']
    });
})();
