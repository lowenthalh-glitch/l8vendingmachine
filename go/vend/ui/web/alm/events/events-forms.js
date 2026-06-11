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
// ALM Events Module - Form Definitions & Primary Keys
// Form configurations for Event

(function() {
    'use strict';

    window.AlmEvents = window.AlmEvents || {};

    const f = Layer8FormFactory;
    const enums = AlmEvents.enums;

    // Primary keys per model
    AlmEvents.primaryKeys = {
        Event: 'eventId'
    };

    // Form definitions
    AlmEvents.forms = {
        Event: f.form('Event', [
            f.section('Event Details', [
                ...f.select('eventType', 'Event Type', enums.EVENT_TYPE, true),
                ...f.select('severity', 'Severity', AlmAlarms.enums.ALARM_SEVERITY),
                ...f.text('nodeId', 'Node ID', true),
                ...f.text('nodeName', 'Node Name'),
                ...f.text('sourceIdentifier', 'Source Identifier'),
                ...f.textarea('message', 'Message', true),
                ...f.textarea('rawData', 'Raw Data'),
                ...f.text('category', 'Category'),
                ...f.text('subcategory', 'Subcategory')
            ]),
            f.section('Processing', [
                ...f.select('processingState', 'Processing State', enums.EVENT_PROCESSING_STATE),
                ...f.reference('alarmId', 'Alarm', 'Alarm'),
                ...f.reference('definitionId', 'Alarm Definition', 'AlarmDefinition')
            ]),
            f.section('Timing', [
                ...f.datetime('occurredAt', 'Occurred At'),
                ...f.datetime('receivedAt', 'Received At'),
                ...f.datetime('processedAt', 'Processed At')
            ])
        ])
    };

})();
