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

package fleetinventory

import (
	common "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "FleetInv"
	ServiceArea = byte(10)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService(common.ServiceConfig{
		ServiceName: ServiceName, ServiceArea: ServiceArea,
		PrimaryKey: "SummaryId", Callback: newFleetInventoryServiceCallback(vnic),
	}, &vend.VendFleetInventory{}, &vend.VendFleetInventoryList{}, creds, dbname, vnic)
}

func FleetInventories(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func FleetInventory(id string, vnic ifs.IVNic) (*vend.VendFleetInventory, error) {
	result, err := common.GetEntity(ServiceName, ServiceArea, &vend.VendFleetInventory{SummaryId: id}, vnic)
	if err != nil || result == nil {
		return nil, err
	}
	return result.(*vend.VendFleetInventory), nil
}
