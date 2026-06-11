/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum } = Layer8MRenderers;

    window.MobileInventoryProducts = window.MobileInventoryProducts || {};

    const PRODUCT_CATEGORY = factory.simple([
        'Unspecified', 'Cold Beverage', 'Energy Drink', 'Snack', 'Fresh Food',
        'Candy', 'Water', 'Juice', 'Coffee', 'Health'
    ]);

    const SLOT_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Stocked', 'stocked', 'active'],
        ['Low Stock', 'lowstock', 'pending'],
        ['Sold Out', 'soldout', 'terminated'],
        ['Empty', 'empty', 'inactive'],
        ['Disabled', 'disabled', 'inactive']
    ]);

    const WORK_ORDER_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'pending'],
        ['Assigned', 'assigned', 'active'],
        ['In Progress', 'inprogress', 'active'],
        ['Completed', 'completed', 'inactive'],
        ['Cancelled', 'cancelled', 'terminated']
    ]);

    MobileInventoryProducts.enums = {
        PRODUCT_CATEGORY: PRODUCT_CATEGORY,
        SLOT_STATUS: SLOT_STATUS,
        WORK_ORDER_STATUS: WORK_ORDER_STATUS
    };

    MobileInventoryProducts.render = {
        productCategory: (v) => renderEnum(v, PRODUCT_CATEGORY.enum),
        slotStatus: createStatusRenderer(SLOT_STATUS.enum, SLOT_STATUS.classes),
        workOrderStatus: createStatusRenderer(WORK_ORDER_STATUS.enum, WORK_ORDER_STATUS.classes)
    };

    MobileInventoryProducts.primaryKeys = {
        VendProduct: 'productId',
        VendPlanogram: 'planogramId',
        VendRestockOrder: 'orderId'
    };
})();
