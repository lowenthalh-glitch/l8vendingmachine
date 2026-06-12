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
	"github.com/saichler/l8types/go/types/l8events"
)

// GetMaintenanceWindowTemplates returns sample maintenance window configurations.
// These are templates -- actual windows should be created with specific machine IDs
// and time ranges when scheduling service visits.
func GetMaintenanceWindowTemplates() []*l8events.MaintenanceWindow {
	return []*l8events.MaintenanceWindow{
		{
			WindowId:              "VEND-MW-TMPL-001",
			Name:                  "Scheduled Restock Window",
			Description:           "Suppress alarms during routine restocking visits",
			Status:                l8events.MaintenanceStatus_MAINTENANCE_STATUS_SCHEDULED,
			Recurrence:            l8events.RecurrenceType_RECURRENCE_TYPE_WEEKLY,
			RecurrenceInterval:    1,
			SuppressAlarms:        true,
			SuppressNotifications: true,
			CreatedBy:             "system",
		},
		{
			WindowId:              "VEND-MW-TMPL-002",
			Name:                  "Quarterly Maintenance Window",
			Description:           "Suppress alarms during quarterly preventive maintenance",
			Status:                l8events.MaintenanceStatus_MAINTENANCE_STATUS_SCHEDULED,
			Recurrence:            l8events.RecurrenceType_RECURRENCE_TYPE_MONTHLY,
			RecurrenceInterval:    3,
			SuppressAlarms:        true,
			SuppressNotifications: true,
			CreatedBy:             "system",
		},
	}
}
