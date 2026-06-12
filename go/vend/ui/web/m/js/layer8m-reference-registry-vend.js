/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var ref = window.Layer8RefFactory;
    window.Layer8MReferenceRegistryVend = {
        ...ref.simple('VendMachine', 'machineId', 'model', 'Machine'),
        ...ref.simple('VendMachineGroup', 'groupId', 'name', 'Group'),
        ...ref.simple('VendLocation', 'locationId', 'name', 'Location'),
        ...ref.simple('VendProduct', 'productId', 'name', 'Product'),
        ...ref.simple('VendDriver', 'driverId', 'lastName', 'Driver'),
        ...ref.simple('VendRoute', 'routeId', 'name', 'Route'),
        ...ref.simple('VendStockingFacility', 'facilityId', 'name', 'Facility'),
        ...ref.simple('VendSupplier', 'supplierId', 'name', 'Supplier'),
        ...ref.simple('VendPlanogram', 'planogramId', 'name', 'Planogram'),
        ...ref.simple('VendPurchaseOrder', 'orderId', 'orderId', 'PO'),
        ...ref.simple('VendInspection', 'inspectionId', 'inspectionId', 'Inspection')
    };
    Layer8MReferenceRegistry.register(window.Layer8MReferenceRegistryVend);
})();
