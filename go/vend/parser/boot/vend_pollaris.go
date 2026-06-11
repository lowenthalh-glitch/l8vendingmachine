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

// GetVendingMachinePollaris creates the polling configuration for vending machines.
// Polls fleet-wide Nayax Lynx endpoints from a single management API target.
func GetVendingMachinePollaris() *l8tpollaris.L8Pollaris {
	p := &l8tpollaris.L8Pollaris{}
	p.Name = "Vend"
	p.Groups = []string{"vending", "vending-machine", "nayax", "Boot_Stage_00"}
	p.Polling = make(map[string]*l8tpollaris.L8Poll)

	// Fleet-wide machine data
	createVendMachinesPoll(p)

	// Payment terminal status
	createVendDevicesPoll(p)

	// Sales & revenue
	createVendTransactionsPoll(p)
	createVendRevenuePoll(p)

	// Analytics
	createVendSalesByPeriodPoll(p)
	createVendMachinePerformancePoll(p)

	p.Order = []string{
		"vendMachines",
		"vendDevices",
		"vendTransactions", "vendRevenue",
		"vendSalesByPeriod", "vendMachinePerformance",
	}

	return p
}
