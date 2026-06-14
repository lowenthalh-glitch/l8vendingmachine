/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum } = Layer8MRenderers;

    window.MobileAlmAlarms = window.MobileAlmAlarms || {};

    const ALARM_SEVERITY = factory.create([
        ['Unspecified', null, ''],
        ['Info', 'info', 'active'],
        ['Warning', 'warning', 'pending'],
        ['Minor', 'minor', 'pending'],
        ['Major', 'major', 'terminated'],
        ['Critical', 'critical', 'terminated']
    ]);

    const ALARM_STATE = factory.create([
        ['Unspecified', null, ''],
        ['Open', 'open', 'terminated'],
        ['Acknowledged', 'acknowledged', 'pending'],
        ['Cleared', 'cleared', 'active'],
        ['Closed', 'closed', 'inactive']
    ]);

    const ALARM_DEFINITION_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Draft', 'draft', 'pending'],
        ['Active', 'active', 'active'],
        ['Disabled', 'disabled', 'inactive']
    ]);

    const EVENT_TYPE = factory.simple([
        'Unspecified', 'Trap', 'Syslog', 'Threshold', 'State Change',
        'Heartbeat', 'Configuration', 'Custom'
    ]);

    MobileAlmAlarms.enums = {
        ALARM_SEVERITY: ALARM_SEVERITY,
        ALARM_STATE: ALARM_STATE,
        ALARM_DEFINITION_STATUS: ALARM_DEFINITION_STATUS,
        EVENT_TYPE: EVENT_TYPE
    };

    MobileAlmAlarms.render = {
        severity: createStatusRenderer(ALARM_SEVERITY.enum, ALARM_SEVERITY.classes),
        state: createStatusRenderer(ALARM_STATE.enum, ALARM_STATE.classes),
        definitionStatus: createStatusRenderer(ALARM_DEFINITION_STATUS.enum, ALARM_DEFINITION_STATUS.classes),
        eventType: (v) => renderEnum(v, EVENT_TYPE.enum)
    };

    MobileAlmAlarms.primaryKeys = {
        Alarm: 'alarmId',
        AlarmDefinition: 'definitionId',
        AlarmFilter: 'filterId',
        Event: 'eventId'
    };
})();
