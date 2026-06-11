/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package services

import (
	"github.com/saichler/l8services/go/services/base"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
	"github.com/saichler/l8types/go/types/l8web"
	"github.com/saichler/l8utils/go/utils/web"
)

const (
	CredsServiceName = "Creds"
	CredsServiceArea = byte(75)
)

func ActivateCredentials(vnic ifs.IVNic) {
	sla := ifs.NewServiceLevelAgreement(&base.BaseService{}, CredsServiceName, CredsServiceArea, true, nil)
	sla.SetServiceGroup(ifs.SystemServiceGroup)
	sla.SetServiceItem(&l8api.L8Credentials{})
	sla.SetServiceItemList(&l8api.L8CredentialsList{})
	sla.SetVoter(true)
	sla.SetTransactional(false)
	sla.SetPrimaryKeys("Id")

	ws := web.New(CredsServiceName, CredsServiceArea, 0)
	ws.AddEndpoint(&l8api.L8Credentials{}, ifs.POST, &l8web.L8Empty{})
	ws.AddEndpoint(&l8api.L8Credentials{}, ifs.PUT, &l8web.L8Empty{})
	ws.AddEndpoint(&l8api.L8Credentials{}, ifs.PATCH, &l8web.L8Empty{})
	ws.AddEndpoint(&l8api.L8Query{}, ifs.GET, &l8api.L8CredentialsList{})
	sla.SetWebService(ws)

	base.Activate(sla, vnic)
}
