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
package main

import (
	"github.com/saichler/l8bus/go/overlay/vnic"
	parserService "github.com/saichler/l8parser/go/parser/service"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	"github.com/saichler/l8vendingmachine/go/vend/parser/rules"
)

func main() {
	resources := vendcommon.CreateResources("vend-parser")
	ifs.SetNetworkMode(ifs.NETWORK_K8s)
	nic := vnic.NewVirtualNetworkInterface(resources, nil)
	nic.Start()
	nic.WaitForConnection()

	pollaris.Activate(nic)

	// Register custom parse rule for Nayax JSON array-to-map responses
	parserService.RegisterRule(&rules.RestArrayToMap{})

	// Activate management system parser (fleet-wide → VCache)
	parserService.Activate(vendcommon.VendMachine_Links_ID,
		&vend.VendMachine{}, false, nic, "MachineId")

	// Activate per-machine parser (per-device → Fleet Machine CRUD via PATCH)
	parserService.Activate(vendcommon.VendPerMachine_Links_ID,
		&vend.VendFleetMachine{}, false, nic, "MachineId")

	vendcommon.WaitForSignal(resources)
}
