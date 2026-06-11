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

package common

import l8c "github.com/saichler/l8common/go/common"

// Re-export service factory types and functions from l8common.
type ServiceConfig = l8c.ServiceConfig

var ActivateService = l8c.ActivateService
var ServiceHandler = l8c.ServiceHandler
var GetEntity = l8c.GetEntity
var GetEntities = l8c.GetEntities
var PutEntity = l8c.PutEntity
var PostEntity = l8c.PostEntity
var EntityExists = l8c.EntityExists
var GetEntitiesByQuery = l8c.GetEntitiesByQuery
