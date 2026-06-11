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

package seeddata

import (
	"github.com/saichler/l8alarms/go/types/alm"
	"github.com/saichler/l8events/go/types/l8events"
)

// GetAlarmFilters returns the pre-configured saved alarm filter views.
func GetAlarmFilters() []*alm.AlarmFilter {
	return []*alm.AlarmFilter{
		{
			FilterId:          "VEND-FILTER-001",
			Name:              "Critical Only",
			Description:       "Show only CRITICAL severity alarms",
			IsShared:          true,
			IsDefault:         true,
			Severities:        []l8events.Severity{l8events.Severity_SEVERITY_CRITICAL},
			States:            []l8events.AlarmState{l8events.AlarmState_ALARM_STATE_ACTIVE},
			ExcludeSuppressed: true,
		},
		{
			FilterId:          "VEND-FILTER-002",
			Name:              "Unacknowledged",
			Description:       "All active alarms not yet acknowledged",
			IsShared:          true,
			States:            []l8events.AlarmState{l8events.AlarmState_ALARM_STATE_ACTIVE},
			ExcludeSuppressed: true,
		},
		{
			FilterId:    "VEND-FILTER-003",
			Name:        "Temperature Alarms",
			Description: "All temperature-related alarms",
			IsShared:    true,
			DefinitionIds: []string{
				"VEND-DEF-001", // VEND_TEMP_HIGH
				"VEND-DEF-008", // VEND_COMPRESSOR_FAIL
			},
			ExcludeSuppressed: true,
		},
		{
			FilterId:          "VEND-FILTER-004",
			Name:              "Inventory Alarms",
			Description:       "Stock-related alarms (low stock, sold out)",
			IsShared:          true,
			DefinitionIds:     []string{
				"VEND-DEF-002", // VEND_STOCK_LOW
				"VEND-DEF-003", // VEND_SOLD_OUT
			},
			ExcludeSuppressed: true,
		},
		{
			FilterId:       "VEND-FILTER-005",
			Name:           "Root Causes Only",
			Description:    "Show only root cause alarms (hide suppressed symptoms)",
			IsShared:       true,
			RootCauseOnly:  true,
			ExcludeSuppressed: true,
		},
	}
}
