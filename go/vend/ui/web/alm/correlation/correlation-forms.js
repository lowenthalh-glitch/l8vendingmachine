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
// ALM Correlation Module - Form Definitions & Primary Keys
// Form configurations for CorrelationRule

(function() {
    'use strict';

    window.AlmCorrelation = window.AlmCorrelation || {};

    const f = Layer8FormFactory;
    const enums = AlmCorrelation.enums;

    // Primary keys per model
    AlmCorrelation.primaryKeys = {
        CorrelationRule: 'ruleId'
    };

    // Form definitions
    AlmCorrelation.forms = {
        CorrelationRule: f.form('Correlation Rule', [
            f.section('Rule Details', [
                ...f.text('name', 'Name', true),
                ...f.textarea('description', 'Description'),
                ...f.select('ruleType', 'Rule Type', enums.CORRELATION_RULE_TYPE, true),
                ...f.select('status', 'Status', enums.CORRELATION_RULE_STATUS),
                ...f.number('priority', 'Priority')
            ]),
            f.section('Topological', [
                ...f.select('traversalDirection', 'Traversal Direction', enums.TRAVERSAL_DIRECTION),
                ...f.number('traversalDepth', 'Traversal Depth')
            ]),
            f.section('Temporal', [
                ...f.number('timeWindowSeconds', 'Time Window (seconds)'),
                ...f.text('rootAlarmPattern', 'Root Alarm Pattern'),
                ...f.text('symptomAlarmPattern', 'Symptom Alarm Pattern')
            ]),
            f.section('Behavior', [
                ...f.number('minSymptomCount', 'Min Symptom Count'),
                ...f.checkbox('autoSuppressSymptoms', 'Auto Suppress Symptoms'),
                ...f.checkbox('autoAcknowledgeSymptoms', 'Auto Acknowledge Symptoms')
            ]),
            f.section('Conditions', [
                ...f.inlineTable('conditions', 'Conditions', [
                    { key: 'conditionId', label: 'ID', type: 'text', hidden: true },
                    { key: 'field', label: 'Field', type: 'text' },
                    { key: 'operator', label: 'Operator', type: 'select', options: enums.CONDITION_OPERATOR },
                    { key: 'value', label: 'Value', type: 'text' }
                ])
            ])
        ])
    };

})();
