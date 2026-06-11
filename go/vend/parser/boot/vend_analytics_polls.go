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

package boot

import "github.com/saichler/l8pollaris/go/types/l8tpollaris"

func createVendSalesByPeriodPoll(p *l8tpollaris.L8Pollaris) {
	p.Polling["vendSalesByPeriod"] = createVendPoll(
		"vendSalesByPeriod",
		"/lynx/v1/reports/sales-by-period",
		EVERY_1_HOUR,
		false,
		"vendmachine.traffic",
		"data:vendmachine.traffic",
	)
}

func createVendMachinePerformancePoll(p *l8tpollaris.L8Pollaris) {
	p.Polling["vendMachinePerformance"] = createVendPoll(
		"vendMachinePerformance",
		"/lynx/v1/reports/machine-performance",
		EVERY_1_HOUR,
		false,
		"vendmachine.health",
		"machines:vendmachine.health",
	)
}
