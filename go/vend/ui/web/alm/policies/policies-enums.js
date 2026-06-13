/*
Layer 8 Alarms - Policies Enum Definitions
*/

(function() {
    'use strict';

    window.AlmPolicies = window.AlmPolicies || {};

    const factory = Layer8EnumFactory;
    const { createStatusRenderer, renderEnum } = Layer8DRenderers;

    const POLICY_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'layer8d-status-active'],
        ['Disabled', 'disabled', 'layer8d-status-inactive']
    ]);

    const NOTIFICATION_CHANNEL = factory.simple([
        'Unspecified', 'Email', 'SMS', 'Webhook', 'Slack', 'PagerDuty'
    ]);

    AlmPolicies.enums = {
        POLICY_STATUS: POLICY_STATUS.enum,
        POLICY_STATUS_CLASSES: POLICY_STATUS.classes,
        NOTIFICATION_CHANNEL: NOTIFICATION_CHANNEL.enum
    };

    AlmPolicies.render = {
        policyStatus: createStatusRenderer(
            POLICY_STATUS.enum,
            POLICY_STATUS.classes
        ),
        notificationChannel: function(value) { return renderEnum(value, NOTIFICATION_CHANNEL.enum); }
    };

})();
