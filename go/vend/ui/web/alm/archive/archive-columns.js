/*
Layer 8 Alarms - Archive Column Definitions
Table column configurations for ArchivedAlarm, ArchivedEvent
*/

(function() {
    'use strict';

    window.AlmArchive = window.AlmArchive || {};

    const col = Layer8ColumnFactory;
    const render = AlmArchive.render;

    AlmArchive.columns = {
        ArchivedAlarm: [
            ...col.id('alarmId'),
            ...col.col('name', 'Name'),
            ...col.status('severity', 'Severity', null, render.severity),
            ...col.status('state', 'State', null, render.state),
            ...col.col('nodeName', 'Node'),
            ...col.datetime('firstOccurrence', 'First Occurrence'),
            ...col.datetime('archivedAt', 'Archived At'),
            ...col.col('archivedBy', 'Archived By')
        ],

        ArchivedEvent: [
            ...col.id('eventId', 'Event ID'),
            ...col.enum('eventType', 'Event Type', null, render.eventType),
            ...col.enum('severity', 'Severity', null, render.severity),
            ...col.col('nodeName', 'Node Name'),
            ...col.col('message', 'Message'),
            ...col.datetime('occurredAt', 'Occurred At'),
            ...col.datetime('archivedAt', 'Archived At'),
            ...col.col('archivedBy', 'Archived By')
        ]
    };

})();
