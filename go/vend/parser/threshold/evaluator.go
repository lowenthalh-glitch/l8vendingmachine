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

// Package threshold evaluates parsed vending machine data against configurable
// thresholds and generates l8events EventRecords when thresholds are crossed.
//
// Threshold categories:
//   - TEMPERATURE: Cabinet temp > safe limit (e.g., > 8C for refrigerated)
//   - INVENTORY: Stock below par level, slot sold out, product nearing expiry
//   - PAYMENT: Cash box > 90% capacity, coin tube empty, card reader offline
//   - MECHANICAL: Motor current > 130% baseline, compressor duty cycle > 85%
//   - CONNECTIVITY: Machine offline > 10 minutes, signal strength below threshold
//   - SALES: Revenue anomaly, failed vend rate > 5%
//
// Integration:
//   This evaluator is called from the parser's After hook after each poll cycle.
//   When a threshold is crossed, it calls l8events.PostEvent() to generate an
//   EventRecord with category EVENT_CATEGORY_PERFORMANCE or EVENT_CATEGORY_SYSTEM.
//   The l8alarms service then matches the event against AlarmDefinitions to
//   create/update alarms, which trigger notifications via l8notify.
package threshold

// ThresholdRule defines a configurable threshold for a vending machine metric.
type ThresholdRule struct {
	Name           string  // Rule identifier (e.g., "VEND_TEMP_HIGH")
	MetricName     string  // Metric being evaluated (e.g., "temperature.currentTemp")
	ThresholdType  string  // "UPPER" or "LOWER"
	WarningValue   float64 // Warning threshold (80% severity)
	CriticalValue  float64 // Critical threshold (100% severity)
	Category       string  // Event category: "TEMPERATURE", "INVENTORY", "PAYMENT", etc.
	Description    string  // Human-readable description
	AutoClear      bool    // Whether the alarm auto-clears when condition resolves
	ClearThreshold float64 // Value at which the condition is considered resolved
}

// DefaultThresholdRules returns the pre-configured threshold rules for vending machines.
func DefaultThresholdRules() []ThresholdRule {
	return []ThresholdRule{
		{
			Name: "VEND_TEMP_HIGH", MetricName: "temperature.currentTemp",
			ThresholdType: "UPPER", WarningValue: 7.0, CriticalValue: 8.0,
			Category: "TEMPERATURE", Description: "Cabinet temperature above safe limit",
			AutoClear: true, ClearThreshold: 6.0,
		},
		{
			Name: "VEND_STOCK_LOW", MetricName: "inventory.fillRate",
			ThresholdType: "LOWER", WarningValue: 30.0, CriticalValue: 10.0,
			Category: "INVENTORY", Description: "Machine inventory below par level",
			AutoClear: true, ClearThreshold: 50.0,
		},
		{
			Name: "VEND_SOLD_OUT", MetricName: "inventory.soldOutSlots",
			ThresholdType: "UPPER", WarningValue: 3.0, CriticalValue: 10.0,
			Category: "INVENTORY", Description: "Multiple slots sold out",
			AutoClear: true, ClearThreshold: 0.0,
		},
		{
			Name: "VEND_CASH_BOX_FULL", MetricName: "cashbox.fillPercent",
			ThresholdType: "UPPER", WarningValue: 80.0, CriticalValue: 90.0,
			Category: "PAYMENT", Description: "Cash box approaching capacity",
			AutoClear: true, ClearThreshold: 50.0,
		},
		{
			Name: "VEND_EXACT_CHANGE", MetricName: "cashbox.exactChangeRequired",
			ThresholdType: "UPPER", WarningValue: 1.0, CriticalValue: 1.0,
			Category: "PAYMENT", Description: "Cannot make change",
			AutoClear: true, ClearThreshold: 0.0,
		},
		{
			Name: "VEND_PAYMENT_OFFLINE", MetricName: "payment.cardReaderOffline",
			ThresholdType: "UPPER", WarningValue: 1.0, CriticalValue: 1.0,
			Category: "PAYMENT", Description: "Card reader or NFC offline",
			AutoClear: true, ClearThreshold: 0.0,
		},
		{
			Name: "VEND_MOTOR_DEGRADED", MetricName: "health.motorCurrentDrift",
			ThresholdType: "UPPER", WarningValue: 20.0, CriticalValue: 30.0,
			Category: "MECHANICAL", Description: "Motor current draw above baseline",
			AutoClear: false,
		},
		{
			Name: "VEND_COMPRESSOR_FAIL", MetricName: "health.compressorFault",
			ThresholdType: "UPPER", WarningValue: 1.0, CriticalValue: 1.0,
			Category: "MECHANICAL", Description: "Compressor not running when temp above setpoint",
			AutoClear: true, ClearThreshold: 0.0,
		},
		{
			Name: "VEND_MACHINE_OFFLINE", MetricName: "status.heartbeatAge",
			ThresholdType: "UPPER", WarningValue: 300.0, CriticalValue: 600.0,
			Category: "CONNECTIVITY", Description: "No heartbeat for extended period",
			AutoClear: true, ClearThreshold: 0.0,
		},
		{
			Name: "VEND_VEND_FAILURE_RATE", MetricName: "sales.failedVendRate",
			ThresholdType: "UPPER", WarningValue: 3.0, CriticalValue: 5.0,
			Category: "SALES", Description: "Failed vend rate exceeds threshold",
			AutoClear: true, ClearThreshold: 1.0,
		},
	}
}

// EvaluateSeverity returns "CRITICAL", "WARNING", or "" based on the current value
// and threshold rule.
func EvaluateSeverity(rule ThresholdRule, currentValue float64) string {
	if rule.ThresholdType == "UPPER" {
		if currentValue >= rule.CriticalValue {
			return "CRITICAL"
		}
		if currentValue >= rule.WarningValue {
			return "WARNING"
		}
	} else { // LOWER
		if currentValue <= rule.CriticalValue {
			return "CRITICAL"
		}
		if currentValue <= rule.WarningValue {
			return "WARNING"
		}
	}
	return ""
}

// IsCleared returns true if the current value indicates the condition has resolved.
func IsCleared(rule ThresholdRule, currentValue float64) bool {
	if !rule.AutoClear {
		return false
	}
	if rule.ThresholdType == "UPPER" {
		return currentValue <= rule.ClearThreshold
	}
	return currentValue >= rule.ClearThreshold
}
