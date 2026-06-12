/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package seeddata provides pre-configured alarm definitions, correlation rules,
// notification policies, escalation policies, maintenance windows, and alarm filters
// for the vending machine management system.
package seeddata

import (
	"github.com/saichler/l8alarms/go/types/alm"
	"github.com/saichler/l8types/go/types/l8events"
)

// GetAlarmDefinitions returns the 10 pre-configured alarm definitions for vending machines.
func GetAlarmDefinitions() []*alm.AlarmDefinition {
	return []*alm.AlarmDefinition{
		{
			DefinitionId:    "VEND-DEF-001",
			Name:            "VEND_TEMP_HIGH",
			Description:     "Cabinet temperature above safe limit for refrigerated zone",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_CRITICAL,
			EventPattern:    "VEND_TEMP_HIGH",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			AutoClearEnabled: true,
			AutoClearSeconds: 300,
			ClearEventPattern: "VEND_TEMP_NORMAL",
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
		{
			DefinitionId:    "VEND-DEF-002",
			Name:            "VEND_STOCK_LOW",
			Description:     "Machine inventory below par level threshold",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_WARNING,
			EventPattern:    "VEND_STOCK_LOW",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			AutoClearEnabled: true,
			ClearEventPattern: "VEND_STOCK_NORMAL",
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
		{
			DefinitionId:    "VEND-DEF-003",
			Name:            "VEND_SOLD_OUT",
			Description:     "Multiple slots completely empty",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_WARNING,
			EventPattern:    "VEND_SOLD_OUT",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			AutoClearEnabled: true,
			ClearEventPattern: "VEND_SOLD_OUT_CLEAR",
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
		{
			DefinitionId:    "VEND-DEF-004",
			Name:            "VEND_CASH_BOX_FULL",
			Description:     "Cash box approaching capacity (>90%)",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_WARNING,
			EventPattern:    "VEND_CASH_BOX_FULL",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			AutoClearEnabled: true,
			ClearEventPattern: "VEND_CASH_BOX_NORMAL",
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
		{
			DefinitionId:    "VEND-DEF-005",
			Name:            "VEND_EXACT_CHANGE",
			Description:     "Machine cannot make change - coin tubes empty",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_WARNING,
			EventPattern:    "VEND_EXACT_CHANGE",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			AutoClearEnabled: true,
			ClearEventPattern: "VEND_CHANGE_AVAILABLE",
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
		{
			DefinitionId:    "VEND-DEF-006",
			Name:            "VEND_PAYMENT_OFFLINE",
			Description:     "Card reader or NFC payment terminal offline",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_MAJOR,
			EventPattern:    "VEND_PAYMENT_OFFLINE",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			AutoClearEnabled: true,
			ClearEventPattern: "VEND_PAYMENT_ONLINE",
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
		{
			DefinitionId:    "VEND-DEF-007",
			Name:            "VEND_MOTOR_DEGRADED",
			Description:     "Vend motor current draw >130% of baseline (degradation)",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_WARNING,
			EventPattern:    "VEND_MOTOR_DEGRADED",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
		{
			DefinitionId:    "VEND-DEF-008",
			Name:            "VEND_COMPRESSOR_FAIL",
			Description:     "Compressor not running when temp exceeds setpoint+5C",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_CRITICAL,
			EventPattern:    "VEND_COMPRESSOR_FAIL",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			AutoClearEnabled: true,
			ClearEventPattern: "VEND_COMPRESSOR_RUNNING",
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
		{
			DefinitionId:    "VEND-DEF-009",
			Name:            "VEND_MACHINE_OFFLINE",
			Description:     "No heartbeat received for >10 minutes",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_CRITICAL,
			EventPattern:    "VEND_MACHINE_OFFLINE",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			AutoClearEnabled: true,
			ClearEventPattern: "VEND_MACHINE_ONLINE",
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
		{
			DefinitionId:    "VEND-DEF-010",
			Name:            "VEND_VEND_FAILURE_RATE",
			Description:     "Failed vend rate exceeds 5% over last hour",
			Status:          alm.AlarmDefinitionStatus_ALARM_DEFINITION_STATUS_ACTIVE,
			DefaultSeverity: l8events.Severity_SEVERITY_MAJOR,
			EventPattern:    "VEND_VEND_FAILURE_RATE",
			EventTypeFilter: alm.AlmEventType_ALM_EVENT_TYPE_THRESHOLD,
			AutoClearEnabled: true,
			ClearEventPattern: "VEND_VEND_FAILURE_NORMAL",
			DedupEnabled:      true,
			DedupKeyExpression: "sourceId+definitionId",
		},
	}
}
