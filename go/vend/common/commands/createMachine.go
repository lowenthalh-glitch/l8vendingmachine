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
package commands

import (
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
)

// CreateVendingMachine creates an L8PTarget for a vending machine REST API endpoint.
func CreateVendingMachine(ip string, port int32) *l8tpollaris.L8PTarget {
	target := &l8tpollaris.L8PTarget{}
	target.TargetId = ip
	target.LinksId = vendcommon.VendMachine_Links_ID
	target.Hosts = make(map[string]*l8tpollaris.L8PHost)
	target.InventoryType = l8tpollaris.L8PTargetType_Network_Device
	target.State = l8tpollaris.L8PTargetState_Down

	host := &l8tpollaris.L8PHost{}
	host.HostId = ip
	host.Configs = make(map[int32]*l8tpollaris.L8PHostProtocol)

	restConfig := &l8tpollaris.L8PHostProtocol{}
	restConfig.Protocol = l8tpollaris.L8PProtocol_L8PRESTAPI
	restConfig.Addr = ip
	restConfig.Port = port
	restConfig.Timeout = 30

	host.Configs[int32(restConfig.Protocol)] = restConfig
	target.Hosts[host.HostId] = host

	return target
}
