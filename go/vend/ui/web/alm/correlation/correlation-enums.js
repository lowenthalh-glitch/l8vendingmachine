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
// ALM Correlation Module - Enum Definitions
// CorrelationRuleType, CorrelationRuleStatus, TraversalDirection, ConditionOperator

(function() {
    'use strict';

    window.AlmCorrelation = window.AlmCorrelation || {};

    const factory = Layer8EnumFactory;
    const { renderEnum, createStatusRenderer } = Layer8DRenderers;

    // CorrelationRuleType: simple enum
    const CORRELATION_RULE_TYPE = factory.simple([
        'Unspecified', 'Topological', 'Temporal', 'Pattern', 'Composite'
    ]);

    // CorrelationRuleStatus: status enum with classes
    const CORRELATION_RULE_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Draft', 'draft', 'layer8d-status-pending'],
        ['Active', 'active', 'layer8d-status-active'],
        ['Disabled', 'disabled', 'layer8d-status-inactive']
    ]);

    // TraversalDirection: simple enum
    const TRAVERSAL_DIRECTION = factory.simple([
        'Unspecified', 'Upstream', 'Downstream', 'Both'
    ]);

    // ConditionOperator: simple enum
    const CONDITION_OPERATOR = factory.simple([
        'Unspecified', 'Equals', 'Not Equals', 'Contains', 'Regex',
        'Greater Than', 'Less Than', 'In'
    ]);

    // Enum exports
    AlmCorrelation.enums = {
        CORRELATION_RULE_TYPE: CORRELATION_RULE_TYPE.enum,
        CORRELATION_RULE_STATUS: CORRELATION_RULE_STATUS.enum,
        CORRELATION_RULE_STATUS_CLASSES: CORRELATION_RULE_STATUS.classes,
        TRAVERSAL_DIRECTION: TRAVERSAL_DIRECTION.enum,
        CONDITION_OPERATOR: CONDITION_OPERATOR.enum
    };

    // Renderers
    AlmCorrelation.render = {
        ruleType: (value) => renderEnum(value, CORRELATION_RULE_TYPE.enum),
        ruleStatus: createStatusRenderer(
            CORRELATION_RULE_STATUS.enum,
            CORRELATION_RULE_STATUS.classes
        ),
        traversalDirection: (value) => renderEnum(value, TRAVERSAL_DIRECTION.enum),
        conditionOperator: (value) => renderEnum(value, CONDITION_OPERATOR.enum)
    };

})();
