/*
Layer 8 Alarms - Column Configurations
Table column definitions for Alarm, AlarmDefinition, AlarmFilter
*/

(function() {
    'use strict';

    window.AlmAlarms = window.AlmAlarms || {};

    const col = window.Layer8ColumnFactory;
    const render = AlmAlarms.render;

    AlmAlarms.columns = {
        Alarm: [
            ...col.id('alarmId'),
            ...col.col('name', 'Name'),
            ...col.status('severity', 'Severity', null, render.severity),
            ...col.status('state', 'State', null, render.state),
            ...col.col('nodeName', 'Node'),
            ...col.datetime('firstOccurrence', 'First Occurrence'),
            ...col.col('occurrenceCount', 'Count'),
            ...col.boolean('isRootCause', 'Root Cause'),
            ...col.col('symptomCount', 'Symptoms')
        ],

        AlarmDefinition: [
            ...col.id('definitionId'),
            ...col.col('name', 'Name'),
            ...col.status('status', 'Status', null, render.definitionStatus),
            ...col.status('defaultSeverity', 'Default Severity', null, render.severity),
            ...col.col('eventPattern', 'Event Pattern'),
            ...col.col('thresholdCount', 'Threshold')
        ],

        AlarmFilter: [
            ...col.id('filterId'),
            ...col.col('name', 'Name'),
            ...col.col('owner', 'Owner'),
            ...col.boolean('isShared', 'Shared'),
            ...col.boolean('isDefault', 'Default')
        ]
    };

})();
