/*
Layer 8 Alarms - Policies Enum Definitions
Uses shared L8NotifyEnums for notification channel, keeps policy status local.
*/

(function() {
    'use strict';

    window.AlmPolicies = window.AlmPolicies || {};

    const factory = Layer8EnumFactory;
    const { createStatusRenderer } = Layer8DRenderers;

    // AlmPolicyStatus: alarm-specific policy status
    const POLICY_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Active', 'active', 'layer8d-status-active'],
        ['Disabled', 'disabled', 'layer8d-status-inactive']
    ]);

    // Use shared notification channel from l8notify
    const NOTIFICATION_CHANNEL = L8NotifyEnums.NOTIFY_CHANNEL;

    // Enum exports
    AlmPolicies.enums = {
        POLICY_STATUS: POLICY_STATUS.enum,
        POLICY_STATUS_CLASSES: POLICY_STATUS.classes,
        NOTIFICATION_CHANNEL: NOTIFICATION_CHANNEL.enum
    };

    // Renderers
    AlmPolicies.render = {
        policyStatus: createStatusRenderer(
            POLICY_STATUS.enum,
            POLICY_STATUS.classes
        ),
        notificationChannel: L8NotifyEnums.render.channel
    };

})();
