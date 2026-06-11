/*
Copyright 2024 Sharon Aicler (saichler@gmail.com)

Layer 8 Alarms is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// ALM Policies Module - Form Definitions & Primary Keys
// Form configurations for NotificationPolicy, EscalationPolicy

(function() {
    'use strict';

    window.AlmPolicies = window.AlmPolicies || {};

    const f = Layer8FormFactory;
    const enums = AlmPolicies.enums;

    // Primary keys per model
    AlmPolicies.primaryKeys = {
        NotificationPolicy: 'policyId',
        EscalationPolicy: 'policyId'
    };

    // Form definitions
    AlmPolicies.forms = {
        NotificationPolicy: f.form('Notification Policy', [
            f.section('Policy Details', [
                ...f.text('name', 'Name', true),
                ...f.textarea('description', 'Description'),
                ...f.select('status', 'Status', enums.POLICY_STATUS),
                ...f.select('minSeverity', 'Min Severity', AlmAlarms.enums.ALARM_SEVERITY),
                ...f.checkbox('notifyOnStateChange', 'Notify on State Change'),
                ...f.number('cooldownSeconds', 'Cooldown (seconds)'),
                ...f.number('maxNotificationsPerHour', 'Max Notifications/Hour')
            ]),
            f.section('Targets', [
                ...f.inlineTable('targets', 'Notification Targets', [
                    { key: 'targetId', label: 'ID', type: 'text', hidden: true },
                    { key: 'channel', label: 'Channel', type: 'select', options: enums.NOTIFICATION_CHANNEL },
                    { key: 'endpoint', label: 'Endpoint', type: 'text' },
                    { key: 'template', label: 'Template', type: 'textarea' }
                ])
            ])
        ]),

        EscalationPolicy: f.form('Escalation Policy', [
            f.section('Policy Details', [
                ...f.text('name', 'Name', true),
                ...f.textarea('description', 'Description'),
                ...f.select('status', 'Status', enums.POLICY_STATUS),
                ...f.select('minSeverity', 'Min Severity', AlmAlarms.enums.ALARM_SEVERITY)
            ]),
            f.section('Escalation Steps', [
                ...f.inlineTable('steps', 'Escalation Steps', [
                    { key: 'stepId', label: 'ID', type: 'text', hidden: true },
                    { key: 'stepOrder', label: 'Order', type: 'number' },
                    { key: 'delayMinutes', label: 'Delay (min)', type: 'number' },
                    { key: 'channel', label: 'Channel', type: 'select', options: enums.NOTIFICATION_CHANNEL },
                    { key: 'endpoint', label: 'Endpoint', type: 'text' },
                    { key: 'messageTemplate', label: 'Message Template', type: 'text' }
                ])
            ])
        ])
    };

})();
