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

var EVERY_30_SECONDS = &l8tpollaris.L8PCadencePlan{
	Enabled: true, Cadences: []int64{30},
}
var EVERY_5_MINUTES = &l8tpollaris.L8PCadencePlan{
	Enabled: true, Cadences: []int64{300},
}
var EVERY_15_MINUTES = &l8tpollaris.L8PCadencePlan{
	Enabled: true, Cadences: []int64{900},
}
var EVERY_1_HOUR = &l8tpollaris.L8PCadencePlan{
	Enabled: true, Cadences: []int64{3600},
}
var EVERY_24_HOURS = &l8tpollaris.L8PCadencePlan{
	Enabled: true, Cadences: []int64{86400},
}

const DEFAULT_TIMEOUT int64 = 30
