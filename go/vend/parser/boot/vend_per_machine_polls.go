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

// GetPerMachinePollaris creates the per-machine polling configuration.
// Name must match VendPerMachine_Links_ID ("VMach") for boot stage filtering.
// Uses $symbol substitution — replaced with the target's machineId at collection time.
func GetPerMachinePollaris() *l8tpollaris.L8Pollaris {
	p := &l8tpollaris.L8Pollaris{}
	p.Name = "VMach"
	p.Groups = []string{"vending", "vending-per-machine", "Boot_Stage_00"}
	p.Polling = make(map[string]*l8tpollaris.L8Poll)

	// Per-machine slot inventory — uses "vendfleetmachine" model key (not "vendmachine")
	poll := &l8tpollaris.L8Poll{}
	poll.Name = "vendMachineInventory"
	poll.What = "GET::/lynx/v1/machines/$symbol/inventory::"
	poll.Protocol = l8tpollaris.L8PProtocol_L8PRESTAPI
	poll.Cadence = EVERY_5_MINUTES
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Always = true
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createRestAttribute(
		"vendfleetmachine",
		"vendfleetmachine",
		"slots:vendfleetmachine.inventory,"+
			"totalSlots:vendfleetmachine.totalslots,"+
			"emptySlots:vendfleetmachine.emptyslots,"+
			"lowStockSlots:vendfleetmachine.lowstockslots,"+
			"lastUpdated:vendfleetmachine.inventorylastupdated",
	))
	p.Polling["vendMachineInventory"] = poll

	p.Order = []string{
		"vendMachineInventory",
	}

	return p
}
