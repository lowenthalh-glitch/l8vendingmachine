/*
Layer 8 Alarms - Form Definitions
Uses Layer8FormFactory for reduced boilerplate
*/

(function() {
    'use strict';

    window.AlmAlarms = window.AlmAlarms || {};

    const f = window.Layer8FormFactory;
    const enums = AlmAlarms.enums;

    // Mark fields as read-only (system-managed, not user-editable)
    function ro(fields) {
        return fields.map(function(field) { field.readOnly = true; return field; });
    }

    AlmAlarms.forms = {
        Alarm: f.form('Alarm', [
            f.section('Alarm Details', [
                // System-managed identity fields (read-only)
                ...ro(f.text('name', 'Name')),
                ...ro(f.textarea('description', 'Description')),
                ...ro(f.reference('definitionId', 'Definition', 'AlarmDefinition')),
                // Operator-editable fields
                ...f.select('severity', 'Severity', enums.ALARM_SEVERITY),
                ...f.select('state', 'State', enums.ALARM_STATE),
                // System-managed source fields (read-only)
                ...ro(f.text('nodeId', 'Node ID')),
                ...ro(f.text('nodeName', 'Node Name')),
                ...ro(f.text('linkId', 'Link ID')),
                ...ro(f.text('location', 'Location')),
                ...ro(f.text('sourceIdentifier', 'Source Identifier'))
            ]),
            f.section('Timing', [
                ...f.datetime('firstOccurrence', 'First Occurrence'),
                ...f.datetime('lastOccurrence', 'Last Occurrence'),
                ...f.datetime('acknowledgedAt', 'Acknowledged At'),
                ...f.datetime('clearedAt', 'Cleared At')
            ]),
            f.section('Notes', [
                ...f.inlineTable('notes', 'Notes', [
                    { key: 'noteId', label: 'ID', hidden: true },
                    { key: 'author', label: 'Author', type: 'text' },
                    { key: 'text', label: 'Text', type: 'text', required: true },
                    { key: 'createdAt', label: 'Created', type: 'date' }
                ])
            ])
        ]),

        AlarmDefinition: f.form('Alarm Definition', [
            f.section('Definition Details', [
                ...f.text('name', 'Name', true),
                ...f.textarea('description', 'Description'),
                ...f.select('status', 'Status', enums.ALARM_DEFINITION_STATUS),
                ...f.select('defaultSeverity', 'Default Severity', enums.ALARM_SEVERITY),
                ...f.text('eventPattern', 'Event Pattern'),
                ...f.text('nodeTypeFilter', 'Node Type Filter'),
                ...f.number('thresholdCount', 'Threshold Count'),
                ...f.number('thresholdWindowSeconds', 'Threshold Window (s)')
            ]),
            f.section('Auto-Clear Configuration', [
                ...f.checkbox('autoClearEnabled', 'Auto-Clear Enabled'),
                ...f.number('autoClearSeconds', 'Auto-Clear Seconds'),
                ...f.text('clearEventPattern', 'Clear Event Pattern')
            ]),
            f.section('Deduplication', [
                ...f.checkbox('dedupEnabled', 'Dedup Enabled'),
                ...f.text('dedupKeyExpression', 'Dedup Key Expression')
            ])
        ]),

        AlarmFilter: f.form('Alarm Filter', [
            f.section('Filter Details', [
                ...f.text('name', 'Name', true),
                ...f.text('owner', 'Owner', true),
                ...f.textarea('description', 'Description'),
                ...f.checkbox('isShared', 'Shared'),
                ...f.checkbox('isDefault', 'Default'),
                ...f.checkbox('rootCauseOnly', 'Root Cause Only'),
                ...f.checkbox('excludeSuppressed', 'Exclude Suppressed'),
                ...f.number('maxAgeHours', 'Max Age (hours)')
            ])
        ])
    };

    AlmAlarms.primaryKeys = {
        Alarm: 'alarmId',
        AlarmDefinition: 'definitionId',
        AlarmFilter: 'filterId'
    };

})();
