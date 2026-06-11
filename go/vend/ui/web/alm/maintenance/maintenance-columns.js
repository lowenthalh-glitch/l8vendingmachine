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
// ALM Maintenance Module - Column Definitions
// Table column configurations for MaintenanceWindow

(function() {
    'use strict';

    window.AlmMaintenance = window.AlmMaintenance || {};

    const col = Layer8ColumnFactory;
    const render = AlmMaintenance.render;

    AlmMaintenance.columns = {
        MaintenanceWindow: [
            ...col.id('windowId', 'Window ID'),
            ...col.col('name', 'Name'),
            ...col.status('status', 'Status', null, render.windowStatus),
            ...col.date('startTime', 'Start Time'),
            ...col.date('endTime', 'End Time'),
            ...col.enum('recurrence', 'Recurrence', null, render.recurrenceType),
            ...col.boolean('suppressAlarms', 'Suppress Alarms'),
            ...col.col('createdBy', 'Created By')
        ]
    };

})();
