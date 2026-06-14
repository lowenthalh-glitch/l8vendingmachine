/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var render = MobileAlmAlarms.render;

    MobileAlmAlarms.columns = {
        Alarm: [
            ...col.id('alarmId'),
            ...col.col('name', 'Name'),
            ...col.status('severity', 'Severity', null, render.severity),
            ...col.status('state', 'State', null, render.state),
            ...col.col('nodeName', 'Node'),
            ...col.date('firstOccurrence', 'First Occurrence'),
            ...col.col('occurrenceCount', 'Count')
        ],
        AlarmDefinition: [
            ...col.id('definitionId'),
            ...col.col('name', 'Name'),
            ...col.status('status', 'Status', null, render.definitionStatus),
            ...col.status('defaultSeverity', 'Default Severity', null, render.severity),
            ...col.col('eventPattern', 'Event Pattern')
        ],
        AlarmFilter: [
            ...col.id('filterId'),
            ...col.col('name', 'Name'),
            ...col.col('owner', 'Owner'),
            ...col.boolean('isShared', 'Shared'),
            ...col.boolean('isDefault', 'Default')
        ],
        Event: [
            ...col.id('eventId'),
            ...col.col('name', 'Name'),
            ...col.status('severity', 'Severity', null, render.severity),
            ...col.col('nodeId', 'Node'),
            ...col.col('sourceIdentifier', 'Source'),
            ...col.date('timestamp', 'Time')
        ]
    };

    // Add primary/secondary for card display
    MobileAlmAlarms.columns.Alarm[1].primary = true;
    MobileAlmAlarms.columns.Alarm[2].secondary = true;
    MobileAlmAlarms.columns.AlarmDefinition[1].primary = true;
    MobileAlmAlarms.columns.AlarmDefinition[2].secondary = true;
    MobileAlmAlarms.columns.AlarmFilter[1].primary = true;
    MobileAlmAlarms.columns.AlarmFilter[2].secondary = true;
    MobileAlmAlarms.columns.Event[1].primary = true;
    MobileAlmAlarms.columns.Event[2].secondary = true;
})();
