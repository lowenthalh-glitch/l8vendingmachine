/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.LAYER8M_NAV_CONFIG_VEND = {
        fleet: {
            subModules: [
                { key: 'machines', label: 'Machines', icon: 'fleet' }
            ],
            services: {
                'machines': [
                    { key: 'machines', label: 'Machines', icon: 'fleet', endpoint: '/10/Machine', model: 'VendMachine', idField: 'machineId' },
                    { key: 'machine-groups', label: 'Groups', icon: 'fleet', endpoint: '/10/MachGrp', model: 'VendMachineGroup', idField: 'groupId' },
                    { key: 'locations', label: 'Locations', icon: 'fleet', endpoint: '/10/Location', model: 'VendLocation', idField: 'locationId' }
                ]
            }
        },
        inventory: {
            subModules: [
                { key: 'products', label: 'Products', icon: 'inventory' }
            ],
            services: {
                'products': [
                    { key: 'products', label: 'Products', icon: 'inventory', endpoint: '/10/Product', model: 'VendProduct', idField: 'productId' },
                    { key: 'planograms', label: 'Planograms', icon: 'inventory', endpoint: '/10/Planogram', model: 'VendPlanogram', idField: 'planogramId' },
                    { key: 'restock-orders', label: 'Restock Orders', icon: 'inventory', endpoint: '/10/RstockOrd', model: 'VendRestockOrder', idField: 'orderId', supportedViews: ['table', 'kanban'] }
                ]
            }
        },
        sales: {
            subModules: [
                { key: 'transactions', label: 'Transactions', icon: 'sales' }
            ],
            services: {
                'transactions': [
                    { key: 'transactions', label: 'Transactions', icon: 'sales', endpoint: '/10/Txn', model: 'VendTransaction', idField: 'transactionId', readOnly: true },
                    { key: 'settlements', label: 'Settlements', icon: 'sales', endpoint: '/10/Settlemnt', model: 'VendSettlement', idField: 'settlementId' }
                ]
            }
        },
        maintenance: {
            subModules: [
                { key: 'alerts', label: 'Alerts', icon: 'maintenance' }
            ],
            services: {
                'alerts': [
                    { key: 'alerts', label: 'Alerts', icon: 'maintenance', endpoint: '/10/Alert', model: 'VendAlert', idField: 'alertId', supportedViews: ['table', 'kanban'] },
                    { key: 'work-orders', label: 'Work Orders', icon: 'maintenance', endpoint: '/10/WorkOrder', model: 'VendWorkOrder', idField: 'workOrderId', supportedViews: ['table', 'kanban'] },
                    { key: 'service-visits', label: 'Service Visits', icon: 'maintenance', endpoint: '/10/SvcVisit', model: 'VendServiceVisit', idField: 'visitId' }
                ]
            }
        },
        routes: {
            subModules: [
                { key: 'routes', label: 'Routes', icon: 'routes' }
            ],
            services: {
                'routes': [
                    { key: 'routes', label: 'Routes', icon: 'routes', endpoint: '/10/Route', model: 'VendRoute', idField: 'routeId' },
                    { key: 'drivers', label: 'Drivers', icon: 'routes', endpoint: '/10/Driver', model: 'VendDriver', idField: 'driverId' },
                    { key: 'trucks', label: 'Trucks', icon: 'routes', endpoint: '/10/Truck', model: 'VendDeliveryTruck', idField: 'truckId' }
                ]
            }
        },
        analytics: {
            subModules: [
                { key: 'forecasts', label: 'Forecasts', icon: 'analytics' }
            ],
            services: {
                'forecasts': [
                    { key: 'forecasts', label: 'Forecasts', icon: 'analytics', endpoint: '/10/Forecast', model: 'VendForecast', idField: 'forecastId', readOnly: true },
                    { key: 'performance', label: 'Performance', icon: 'analytics', endpoint: '/10/SlotPerf', model: 'VendSlotPerformance', idField: 'performanceId', readOnly: true },
                    { key: 'fleet-inventory', label: 'Fleet Inventory', icon: 'analytics', endpoint: '/10/FleetInv', model: 'VendFleetInventory', idField: 'summaryId', readOnly: true }
                ]
            }
        },
        warehouse: {
            subModules: [
                { key: 'stock', label: 'Stock', icon: 'warehouse' }
            ],
            services: {
                'stock': [
                    { key: 'facilities', label: 'Facilities', icon: 'warehouse', endpoint: '/10/Facility', model: 'VendStockingFacility', idField: 'facilityId' },
                    { key: 'suppliers', label: 'Suppliers', icon: 'warehouse', endpoint: '/10/Supplier', model: 'VendSupplier', idField: 'supplierId' },
                    { key: 'purchase-orders', label: 'Purchase Orders', icon: 'warehouse', endpoint: '/10/PurchOrd', model: 'VendPurchaseOrder', idField: 'orderId', supportedViews: ['table', 'kanban'] },
                    { key: 'movements', label: 'Movements', icon: 'warehouse', endpoint: '/10/StockMove', model: 'VendStockMovement', idField: 'movementId', readOnly: true },
                    { key: 'vehicle-loads', label: 'Vehicle Loads', icon: 'warehouse', endpoint: '/10/VehLoad', model: 'VendVehicleLoad', idField: 'loadId' }
                ]
            }
        },
        compliance: {
            subModules: [
                { key: 'inspections', label: 'Inspections', icon: 'compliance' }
            ],
            services: {
                'inspections': [
                    { key: 'inspections', label: 'Inspections', icon: 'compliance', endpoint: '/10/Inspction', model: 'VendInspection', idField: 'inspectionId' },
                    { key: 'findings', label: 'Findings', icon: 'compliance', endpoint: '/10/InspFind', model: 'VendInspectionFinding', idField: 'findingId', supportedViews: ['table', 'kanban'] },
                    { key: 'certifications', label: 'Certifications', icon: 'compliance', endpoint: '/10/VendCert', model: 'VendCertification', idField: 'certificationId' }
                ]
            }
        },
        reports: {
            subModules: [
                { key: 'reports', label: 'Reports', icon: 'reports' }
            ],
            services: {
                'reports': [
                    { key: 'reports', label: 'Reports', icon: 'reports', endpoint: '/10/VendRpt', model: 'VendReport', idField: 'reportId' }
                ]
            }
        },
        alarms: {
            subModules: [
                { key: 'alarms', label: 'Alarms', icon: 'alarms' }
            ],
            services: {
                'alarms': [
                    { key: 'alarms', label: 'Alarms', icon: 'alarms', endpoint: '/10/Alarm', model: 'Alarm', idField: 'alarmId' },
                    { key: 'definitions', label: 'Definitions', icon: 'alarms', endpoint: '/10/AlarmDef', model: 'AlarmDefinition', idField: 'definitionId' },
                    { key: 'filters', label: 'Filters', icon: 'alarms', endpoint: '/10/AlrmFltr', model: 'AlarmFilter', idField: 'filterId' },
                    { key: 'events', label: 'Events', icon: 'alarms', endpoint: '/10/Event', model: 'Event', idField: 'eventId', readOnly: true }
                ]
            }
        },
        nayax: {
            subModules: [
                { key: 'machines', label: 'Machines', icon: 'nayax' }
            ],
            services: {
                'machines': [
                    { key: 'machines', label: 'Machines', icon: 'nayax', endpoint: '/0/VCache', model: 'VendMachine', idField: 'machineId', readOnly: true },
                    { key: 'groups', label: 'Groups', icon: 'nayax', endpoint: '/10/MachGrp', model: 'VendMachineGroup', idField: 'groupId' },
                    { key: 'locations', label: 'Locations', icon: 'nayax', endpoint: '/10/Location', model: 'VendLocation', idField: 'locationId' }
                ]
            }
        }
    };
})();
