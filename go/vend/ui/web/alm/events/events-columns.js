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
// ALM Events Module - Column Definitions
// Table column configurations for Event

(function() {
    'use strict';

    window.AlmEvents = window.AlmEvents || {};

    const col = Layer8ColumnFactory;
    const render = AlmEvents.render;

    AlmEvents.columns = {
        Event: [
            ...col.id('eventId', 'Event ID'),
            ...col.enum('eventType', 'Event Type', null, render.eventType),
            ...col.status('processingState', 'Processing State', null, render.processingState),
            ...col.enum('severity', 'Severity', null, AlmAlarms.render.severity),
            ...col.col('nodeName', 'Node Name'),
            ...col.col('message', 'Message'),
            ...col.col('category', 'Category'),
            ...col.datetime('occurredAt', 'Occurred At')
        ]
    };

})();
