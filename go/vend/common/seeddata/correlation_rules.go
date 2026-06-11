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
)

// GetCorrelationRules returns the 3 pre-configured correlation rules for vending machines.
func GetCorrelationRules() []*alm.CorrelationRule {
	return []*alm.CorrelationRule{
		{
			RuleId:              "VEND-CORR-001",
			Name:                "Compressor Failure → Temperature High",
			Description:         "When compressor fails, suppress temperature alarms as symptoms",
			RuleType:            alm.CorrelationRuleType_CORRELATION_RULE_TYPE_TEMPORAL,
			Status:              alm.CorrelationRuleStatus_CORRELATION_RULE_STATUS_ACTIVE,
			Priority:            1,
			TimeWindowSeconds:   900, // 15 minutes
			RootAlarmPattern:    ".*COMPRESSOR_FAIL.*",
			SymptomAlarmPattern: ".*TEMP_HIGH.*",
			MinSymptomCount:     1,
			AutoSuppressSymptoms: true,
		},
		{
			RuleId:              "VEND-CORR-002",
			Name:                "Payment Offline → Vend Failure Rate",
			Description:         "When payment system is offline, suppress vend failure rate alarms",
			RuleType:            alm.CorrelationRuleType_CORRELATION_RULE_TYPE_PATTERN,
			Status:              alm.CorrelationRuleStatus_CORRELATION_RULE_STATUS_ACTIVE,
			Priority:            2,
			RootAlarmPattern:    ".*PAYMENT_OFFLINE.*",
			SymptomAlarmPattern: ".*VEND_FAILURE_RATE.*",
			MinSymptomCount:     1,
			AutoSuppressSymptoms: true,
		},
		{
			RuleId:              "VEND-CORR-003",
			Name:                "Machine Offline → Suppress All",
			Description:         "When machine is offline, suppress all other alarms from that machine",
			RuleType:            alm.CorrelationRuleType_CORRELATION_RULE_TYPE_PATTERN,
			Status:              alm.CorrelationRuleStatus_CORRELATION_RULE_STATUS_ACTIVE,
			Priority:            0, // Highest priority
			RootAlarmPattern:    ".*MACHINE_OFFLINE.*",
			SymptomAlarmPattern: ".*VEND_.*",
			MinSymptomCount:     1,
			AutoSuppressSymptoms: true,
		},
	}
}
