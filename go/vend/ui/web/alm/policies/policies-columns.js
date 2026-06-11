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
// ALM Policies Module - Column Definitions
// Table column configurations for NotificationPolicy, EscalationPolicy

(function() {
    'use strict';

    window.AlmPolicies = window.AlmPolicies || {};

    const col = Layer8ColumnFactory;
    const render = AlmPolicies.render;

    AlmPolicies.columns = {
        NotificationPolicy: [
            ...col.id('policyId', 'Policy ID'),
            ...col.col('name', 'Name'),
            ...col.status('status', 'Status', null, render.policyStatus),
            ...col.status('minSeverity', 'Min Severity', null, AlmAlarms.render.severity),
            ...col.col('cooldownSeconds', 'Cooldown (s)'),
            ...col.col('maxNotificationsPerHour', 'Max/Hour')
        ],

        EscalationPolicy: [
            ...col.id('policyId', 'Policy ID'),
            ...col.col('name', 'Name'),
            ...col.status('status', 'Status', null, render.policyStatus),
            ...col.status('minSeverity', 'Min Severity', null, AlmAlarms.render.severity)
        ]
    };

})();
