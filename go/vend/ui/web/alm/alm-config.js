/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/

// ALM Module Configuration - Uses Layer8ModuleConfigFactory
Layer8ModuleConfigFactory.create({
    namespace: 'Alm',
    modules: {
        'alarms': {
            label: 'Alarms',
            services: [
                { key: 'alarms', label: 'Active Alarms', endpoint: '/10/Alarm', model: 'Alarm', supportedViews: ['table', 'kanban', 'chart'] },
                { key: 'alarm-definitions', label: 'Definitions', endpoint: '/10/AlmDef', model: 'AlarmDefinition' },
                { key: 'alarm-filters', label: 'Saved Filters', endpoint: '/10/AlmFilter', model: 'AlarmFilter' }
            ]
        },
        'events': {
            label: 'Events',
            services: [
                { key: 'events', label: 'Events', endpoint: '/10/Event', model: 'Event', readOnly: true }
            ]
        },
        'correlation': {
            label: 'Correlation',
            services: [
                { key: 'correlation-rules', label: 'Rules', endpoint: '/10/CorrRule', model: 'CorrelationRule' }
            ]
        },
        'policies': {
            label: 'Policies',
            services: [
                { key: 'notification-policies', label: 'Notification', endpoint: '/10/NotifPol', model: 'NotificationPolicy' },
                { key: 'escalation-policies', label: 'Escalation', endpoint: '/10/EscPolicy', model: 'EscalationPolicy' }
            ]
        },
        'maintenance': {
            label: 'Maintenance',
            services: [
                { key: 'maintenance-windows', label: 'Windows', endpoint: '/10/MaintWin', model: 'MaintenanceWindow', supportedViews: ['table', 'calendar'] }
            ]
        },
        'archive': {
            label: 'Archive',
            services: [
                { key: 'archived-alarms', label: 'Archived Alarms', endpoint: '/10/ArcAlarm', model: 'ArchivedAlarm', readOnly: true },
                { key: 'archived-events', label: 'Archived Events', endpoint: '/10/ArcEvent', model: 'ArchivedEvent', readOnly: true }
            ]
        }
    },
    submodules: ['AlmAlarms', 'AlmEvents', 'AlmCorrelation', 'AlmPolicies', 'AlmMaintenance', 'AlmArchive']
});
