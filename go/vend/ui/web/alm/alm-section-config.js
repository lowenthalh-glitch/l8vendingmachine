/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/

// ALM Section Configuration for Layer8SectionGenerator
(function() {
    'use strict';

    Layer8SectionConfigs.register('alarms', {
        title: 'Alarms & Events Management',
        subtitle: 'Monitoring, Correlation, Notification & Maintenance',
        icon: '\uD83D\uDD14',
        svgContent: Layer8SvgFactory.generate('alarms'),
        initFn: 'initializeAlm',
        modules: [
            {
                key: 'alarms', label: 'Alarms', icon: '\uD83D\uDD14', isDefault: true,
                services: [
                    { key: 'alarms', label: 'Active Alarms', icon: '\u26A0\uFE0F', isDefault: true },
                    { key: 'alarm-definitions', label: 'Definitions', icon: '\uD83D\uDCD6' },
                    { key: 'alarm-filters', label: 'Saved Filters', icon: '\uD83D\uDD0D' }
                ]
            },
            {
                key: 'events', label: 'Events', icon: '\u26A1',
                services: [
                    { key: 'events', label: 'Events', icon: '\u26A1', isDefault: true }
                ]
            },
            {
                key: 'correlation', label: 'Correlation', icon: '\uD83D\uDD17',
                services: [
                    { key: 'correlation-rules', label: 'Rules', icon: '\uD83D\uDD17', isDefault: true }
                ]
            },
            {
                key: 'policies', label: 'Policies', icon: '\uD83D\uDCDC',
                services: [
                    { key: 'notification-policies', label: 'Notification', icon: '\uD83D\uDCE8', isDefault: true },
                    { key: 'escalation-policies', label: 'Escalation', icon: '\uD83D\uDCE2' }
                ]
            },
            {
                key: 'maintenance', label: 'Maintenance', icon: '\uD83D\uDD27',
                services: [
                    { key: 'maintenance-windows', label: 'Windows', icon: '\uD83D\uDD27', isDefault: true }
                ]
            },
            {
                key: 'archive', label: 'Archive', icon: '\uD83D\uDDC4\uFE0F',
                services: [
                    { key: 'archived-alarms', label: 'Archived Alarms', icon: '\uD83D\uDD14', isDefault: true },
                    { key: 'archived-events', label: 'Archived Events', icon: '\u26A1' }
                ]
            }
        ]
    });
})();
