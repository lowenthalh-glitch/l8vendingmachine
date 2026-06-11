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
// ALM Correlation Module - Column Definitions
// Table column configurations for CorrelationRule

(function() {
    'use strict';

    window.AlmCorrelation = window.AlmCorrelation || {};

    const col = Layer8ColumnFactory;
    const render = AlmCorrelation.render;

    AlmCorrelation.columns = {
        CorrelationRule: [
            ...col.id('ruleId', 'Rule ID'),
            ...col.col('name', 'Name'),
            ...col.enum('ruleType', 'Rule Type', null, render.ruleType),
            ...col.status('status', 'Status', null, render.ruleStatus),
            ...col.col('priority', 'Priority'),
            ...col.col('timeWindowSeconds', 'Time Window (s)'),
            ...col.col('minSymptomCount', 'Min Symptoms'),
            ...col.boolean('autoSuppressSymptoms', 'Auto Suppress')
        ]
    };

})();
