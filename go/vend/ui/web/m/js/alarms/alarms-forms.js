/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;
    var enums = MobileAlmAlarms.enums;

    MobileAlmAlarms.forms = {
        Alarm: f.form('Alarm', [
            f.section('Alarm Details', [
                ...f.text('name', 'Name'),
                ...f.textarea('description', 'Description'),
                ...f.select('severity', 'Severity', enums.ALARM_SEVERITY.enum),
                ...f.select('state', 'State', enums.ALARM_STATE.enum),
                ...f.text('nodeId', 'Node ID'),
                ...f.text('nodeName', 'Node Name'),
                ...f.text('sourceIdentifier', 'Source')
            ]),
            f.section('Timing', [
                ...f.datetime('firstOccurrence', 'First Occurrence'),
                ...f.datetime('lastOccurrence', 'Last Occurrence'),
                ...f.datetime('acknowledgedAt', 'Acknowledged At'),
                ...f.datetime('clearedAt', 'Cleared At')
            ])
        ]),
        AlarmDefinition: f.form('Alarm Definition', [
            f.section('Definition Details', [
                ...f.text('name', 'Name', true),
                ...f.textarea('description', 'Description'),
                ...f.select('status', 'Status', enums.ALARM_DEFINITION_STATUS.enum),
                ...f.select('defaultSeverity', 'Default Severity', enums.ALARM_SEVERITY.enum),
                ...f.text('eventPattern', 'Event Pattern'),
                ...f.number('thresholdCount', 'Threshold Count'),
                ...f.number('thresholdWindowSeconds', 'Threshold Window (s)')
            ]),
            f.section('Auto-Clear', [
                ...f.checkbox('autoClearEnabled', 'Auto-Clear Enabled'),
                ...f.number('autoClearSeconds', 'Auto-Clear Seconds')
            ])
        ]),
        AlarmFilter: f.form('Alarm Filter', [
            f.section('Filter Details', [
                ...f.text('name', 'Name', true),
                ...f.text('owner', 'Owner', true),
                ...f.textarea('description', 'Description'),
                ...f.checkbox('isShared', 'Shared'),
                ...f.checkbox('isDefault', 'Default'),
                ...f.checkbox('rootCauseOnly', 'Root Cause Only')
            ])
        ]),
        Event: f.form('Event', [
            f.section('Event Details', [
                ...f.text('name', 'Name'),
                ...f.select('severity', 'Severity', enums.ALARM_SEVERITY.enum),
                ...f.text('nodeId', 'Node'),
                ...f.text('sourceIdentifier', 'Source'),
                ...f.datetime('timestamp', 'Time'),
                ...f.textarea('description', 'Description')
            ])
        ])
    };
})();
