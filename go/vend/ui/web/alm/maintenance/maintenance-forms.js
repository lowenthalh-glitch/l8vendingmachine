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
// ALM Maintenance Module - Form Definitions & Primary Keys
// Form configurations for MaintenanceWindow

(function() {
    'use strict';

    window.AlmMaintenance = window.AlmMaintenance || {};

    const f = Layer8FormFactory;
    const enums = AlmMaintenance.enums;

    // Primary keys per model
    AlmMaintenance.primaryKeys = {
        MaintenanceWindow: 'windowId'
    };

    // Form definitions
    AlmMaintenance.forms = {
        MaintenanceWindow: f.form('Maintenance Window', [
            f.section('Window Details', [
                ...f.text('name', 'Name', true),
                ...f.textarea('description', 'Description'),
                ...f.select('status', 'Status', enums.MAINTENANCE_WINDOW_STATUS),
                ...f.text('createdBy', 'Created By')
            ]),
            f.section('Schedule', [
                ...f.date('startTime', 'Start Time', true),
                ...f.date('endTime', 'End Time', true),
                ...f.select('recurrence', 'Recurrence', enums.RECURRENCE_TYPE),
                ...f.number('recurrenceInterval', 'Recurrence Interval')
            ]),
            f.section('Scope', [
                ...f.text('nodeIds', 'Node IDs'),
                ...f.text('nodeTypes', 'Node Types'),
                ...f.text('locations', 'Locations')
            ]),
            f.section('Behavior', [
                ...f.checkbox('suppressAlarms', 'Suppress Alarms'),
                ...f.checkbox('suppressNotifications', 'Suppress Notifications')
            ])
        ])
    };

})();
