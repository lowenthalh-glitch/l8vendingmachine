/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    var ref = window.Layer8RefFactory;

    Layer8DReferenceRegistry.register({
        ...ref.simple('VendMachine', 'machineId', 'model', 'Machine'),
        ...ref.simple('VendMachineGroup', 'groupId', 'name', 'Machine Group'),
        ...ref.simple('VendLocation', 'locationId', 'name', 'Location'),
        ...ref.simple('VendProduct', 'productId', 'name', 'Product'),
        ...ref.simple('VendDriver', 'driverId', 'lastName', 'Driver'),
        ...ref.simple('VendRoute', 'routeId', 'name', 'Route'),
        ...ref.simple('VendDeliveryTruck', 'truckId', 'name', 'Truck'),
        ...ref.simple('VendStockingFacility', 'facilityId', 'name', 'Facility'),
        ...ref.simple('VendSupplier', 'supplierId', 'name', 'Supplier'),
        ...ref.simple('VendPlanogram', 'planogramId', 'name', 'Planogram'),
        ...ref.simple('VendPurchaseOrder', 'orderId', 'orderId', 'Purchase Order'),
        ...ref.simple('VendWorkOrder', 'workOrderId', 'workOrderId', 'Work Order'),
        ...ref.simple('VendInspection', 'inspectionId', 'inspectionId', 'Inspection'),
        ...ref.simple('VendReport', 'reportId', 'name', 'Report')
    });
})();
