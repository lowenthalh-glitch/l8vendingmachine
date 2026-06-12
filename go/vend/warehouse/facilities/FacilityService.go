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

package facilities

import (
	common "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "Facility"
	ServiceArea = byte(10)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService(common.ServiceConfig{
		ServiceName: ServiceName, ServiceArea: ServiceArea,
		PrimaryKey: "FacilityId", Callback: newFacilityServiceCallback(vnic),
	}, &vend.VendStockingFacility{}, &vend.VendStockingFacilityList{}, creds, dbname, vnic)
}

func Facilities(vnic ifs.IVNic) (ifs.IServiceHandler, bool) {
	return common.ServiceHandler(ServiceName, ServiceArea, vnic)
}

func Facility(facilityId string, vnic ifs.IVNic) (*vend.VendStockingFacility, error) {
	result, err := common.GetEntity(ServiceName, ServiceArea, &vend.VendStockingFacility{FacilityId: facilityId}, vnic)
	if err != nil || result == nil {
		return nil, err
	}
	return result.(*vend.VendStockingFacility), nil
}
