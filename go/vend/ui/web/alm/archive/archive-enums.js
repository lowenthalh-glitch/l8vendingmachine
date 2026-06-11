/*
Layer 8 Alarms - Archive Enum Definitions
Reuses AlmAlarms renderers for severity/state display
*/

(function() {
    'use strict';

    window.AlmArchive = window.AlmArchive || {};

    // Archive entities reuse alarm/event renderers — no new enums needed
    AlmArchive.enums = {};

    AlmArchive.render = {
        severity: AlmAlarms.render.severity,
        state: AlmAlarms.render.state,
        eventType: AlmAlarms.render.eventType
    };

})();
