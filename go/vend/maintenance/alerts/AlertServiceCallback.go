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

package alerts

import (
	common "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	"github.com/saichler/l8types/go/ifs"
)

func newAlertServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
	return common.NewValidation(&vend.VendAlert{}, vnic).
		Require(func(v interface{}) string { return v.(*vend.VendAlert).AlertId }, "AlertId").
		Enum(func(v interface{}) int32 { return int32(v.(*vend.VendAlert).Severity) }, vend.VendAlertSeverity_name, "Severity").
		Enum(func(v interface{}) int32 { return int32(v.(*vend.VendAlert).Category) }, vend.VendAlertCategory_name, "Category").
		Enum(func(v interface{}) int32 { return int32(v.(*vend.VendAlert).Status) }, vend.VendAlertStatus_name, "Status").
		Build()
}
