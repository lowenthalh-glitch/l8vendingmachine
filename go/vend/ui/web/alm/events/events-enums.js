/*
Layer 8 Alarms - Events Enum Definitions
Uses shared L8EventsEnums for event state, keeps alarm-specific event types local.
*/

(function() {
    'use strict';

    window.AlmEvents = window.AlmEvents || {};

    const factory = Layer8EnumFactory;
    const { renderEnum } = Layer8DRenderers;

    // AlmEventType: alarm-specific event types (different from l8events.EventCategory)
    const EVENT_TYPE = factory.simple([
        'Unspecified',
        'Fault',
        'Threshold',
        'StateChange',
        'ConfigChange',
        'Security',
        'Performance',
        'Syslog'
    ]);

    // Use shared event state from l8events
    const EVENT_PROCESSING_STATE = L8EventsEnums.EVENT_STATE;

    // Enum exports
    AlmEvents.enums = {
        EVENT_TYPE: EVENT_TYPE.enum,
        EVENT_PROCESSING_STATE: EVENT_PROCESSING_STATE.enum,
        EVENT_PROCESSING_STATE_CLASSES: EVENT_PROCESSING_STATE.classes
    };

    // Renderers
    AlmEvents.render = {
        eventType: (value) => renderEnum(value, EVENT_TYPE.enum),
        processingState: L8EventsEnums.render.eventState
    };

})();
