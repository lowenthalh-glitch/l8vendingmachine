/*
Layer 8 Alarms - Maintenance Module Enum Definitions
Uses shared L8EventsEnums for maintenance status and recurrence type.
*/

(function() {
    'use strict';

    window.AlmMaintenance = window.AlmMaintenance || {};

    // Use shared enums from l8events
    const MAINTENANCE_WINDOW_STATUS = L8EventsEnums.MAINTENANCE_STATUS;
    const RECURRENCE_TYPE = L8EventsEnums.RECURRENCE_TYPE;

    // Enum exports
    AlmMaintenance.enums = {
        MAINTENANCE_WINDOW_STATUS: MAINTENANCE_WINDOW_STATUS.enum,
        MAINTENANCE_WINDOW_STATUS_CLASSES: MAINTENANCE_WINDOW_STATUS.classes,
        RECURRENCE_TYPE: RECURRENCE_TYPE.enum
    };

    // Renderers
    AlmMaintenance.render = {
        windowStatus: L8EventsEnums.render.maintenanceStatus,
        recurrenceType: L8EventsEnums.render.recurrenceType
    };

})();
