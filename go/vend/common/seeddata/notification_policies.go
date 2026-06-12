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
	"github.com/saichler/l8types/go/types/l8events"
	"github.com/saichler/l8notify/go/types/l8notify"
)

// GetNotificationPolicies returns the pre-configured notification policies.
func GetNotificationPolicies() []*alm.NotificationPolicy {
	return []*alm.NotificationPolicy{
		{
			PolicyId:    "VEND-NOTIF-001",
			Name:        "Critical Alerts - Operations Manager",
			Description: "Email operations manager for all CRITICAL severity alarms",
			Status:      alm.AlmPolicyStatus_ALM_POLICY_STATUS_ACTIVE,
			MinSeverity: l8events.Severity_SEVERITY_CRITICAL,
			CooldownSeconds:         900, // 15 minutes
			MaxNotificationsPerHour: 10,
			Targets: []*l8notify.NotifyTarget{
				{
					TargetId: "VEND-TGT-001",
					Channel:  l8notify.NotifyChannel_NOTIFY_CHANNEL_EMAIL,
					Endpoint: "ops-manager@vendingco.com",
					Template: "[CRITICAL] {{alarm.name}} on {{alarm.nodeName}} at {{alarm.location}}: {{alarm.description}}",
				},
			},
		},
		{
			PolicyId:    "VEND-NOTIF-002",
			Name:        "Route Driver Webhook",
			Description: "Webhook notification to route driver app for WARNING+ alarms",
			Status:      alm.AlmPolicyStatus_ALM_POLICY_STATUS_ACTIVE,
			MinSeverity: l8events.Severity_SEVERITY_WARNING,
			CooldownSeconds:         900,
			MaxNotificationsPerHour: 20,
			Targets: []*l8notify.NotifyTarget{
				{
					TargetId: "VEND-TGT-002",
					Channel:  l8notify.NotifyChannel_NOTIFY_CHANNEL_WEBHOOK,
					Endpoint: "https://driver-app.vendingco.com/api/alerts",
					Template: `{"alarmId":"{{alarm.id}}","name":"{{alarm.name}}","severity":"{{alarm.severity}}","machine":"{{alarm.nodeName}}","location":"{{alarm.location}}"}`,
				},
			},
		},
		{
			PolicyId:    "VEND-NOTIF-003",
			Name:        "Maintenance Team Slack",
			Description: "Slack notification for MECHANICAL category alarms",
			Status:      alm.AlmPolicyStatus_ALM_POLICY_STATUS_ACTIVE,
			MinSeverity: l8events.Severity_SEVERITY_WARNING,
			AlarmDefinitionIds: []string{
				"VEND-DEF-007", // VEND_MOTOR_DEGRADED
				"VEND-DEF-008", // VEND_COMPRESSOR_FAIL
			},
			CooldownSeconds:         1800, // 30 minutes
			MaxNotificationsPerHour: 10,
			Targets: []*l8notify.NotifyTarget{
				{
					TargetId: "VEND-TGT-003",
					Channel:  l8notify.NotifyChannel_NOTIFY_CHANNEL_SLACK,
					Endpoint: "https://slack.example.com/placeholder-webhook",
					Template: ":warning: *{{alarm.name}}* on `{{alarm.nodeName}}`\n>Location: {{alarm.location}}\n>Severity: {{alarm.severity}}\n>{{alarm.description}}",
				},
			},
		},
		{
			PolicyId:    "VEND-NOTIF-004",
			Name:        "Machine Offline - Location Contact",
			Description: "Email location contact when machine goes offline",
			Status:      alm.AlmPolicyStatus_ALM_POLICY_STATUS_ACTIVE,
			MinSeverity: l8events.Severity_SEVERITY_CRITICAL,
			AlarmDefinitionIds: []string{
				"VEND-DEF-009", // VEND_MACHINE_OFFLINE
			},
			CooldownSeconds:         3600, // 1 hour
			MaxNotificationsPerHour: 5,
			Targets: []*l8notify.NotifyTarget{
				{
					TargetId: "VEND-TGT-004",
					Channel:  l8notify.NotifyChannel_NOTIFY_CHANNEL_EMAIL,
					Endpoint: "site-contact@vendingco.com",
					Template: "Vending machine {{alarm.nodeName}} at {{alarm.location}} is offline. Please check connectivity. Alert: {{alarm.description}}",
				},
			},
		},
	}
}

// GetEscalationPolicies returns the pre-configured escalation policies.
func GetEscalationPolicies() []*alm.EscalationPolicy {
	return []*alm.EscalationPolicy{
		{
			PolicyId:    "VEND-ESC-001",
			Name:        "Standard Escalation",
			Description: "3-step escalation: driver (0min) → supervisor (15min) → manager (60min)",
			Status:      alm.AlmPolicyStatus_ALM_POLICY_STATUS_ACTIVE,
			MinSeverity: l8events.Severity_SEVERITY_CRITICAL,
			Steps: []*l8notify.EscalationStep{
				{
					StepId:          "VEND-ESC-001-S1",
					StepOrder:       1,
					DelayMinutes:    0,
					Channel:         l8notify.NotifyChannel_NOTIFY_CHANNEL_WEBHOOK,
					Endpoint:        "https://driver-app.vendingco.com/api/escalation",
					MessageTemplate: `{"step":1,"alarmId":"{{alarm.id}}","name":"{{alarm.name}}","severity":"{{alarm.severity}}","machine":"{{alarm.nodeName}}"}`,
				},
				{
					StepId:          "VEND-ESC-001-S2",
					StepOrder:       2,
					DelayMinutes:    15,
					Channel:         l8notify.NotifyChannel_NOTIFY_CHANNEL_EMAIL,
					Endpoint:        "area-supervisor@vendingco.com",
					MessageTemplate: "[ESCALATION] Unacknowledged alert: {{alarm.name}} on {{alarm.nodeName}} at {{alarm.location}} (15 minutes). {{alarm.description}}",
				},
				{
					StepId:          "VEND-ESC-001-S3",
					StepOrder:       3,
					DelayMinutes:    60,
					Channel:         l8notify.NotifyChannel_NOTIFY_CHANNEL_EMAIL,
					Endpoint:        "ops-manager@vendingco.com",
					MessageTemplate: "[ESCALATION - FINAL] Alert unacknowledged for 60 minutes: {{alarm.name}} on {{alarm.nodeName}} at {{alarm.location}}. Immediate attention required. {{alarm.description}}",
				},
			},
		},
	}
}
