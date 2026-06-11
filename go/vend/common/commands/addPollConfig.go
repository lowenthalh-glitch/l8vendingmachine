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
	"strconv"
	"time"

	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/parser/boot"
	"github.com/saichler/l8web/go/web/client"
)

// AddPollConfigs registers the vending machine pollaris configuration.
func AddPollConfigs(rc *client.RestClient, resources ifs.IResources) {
	vendPollaris := boot.GetVendingMachinePollaris()
	resp, err := rc.POST(strconv.Itoa(int(pollaris.ServiceArea))+"/"+pollaris.ServiceName,
		"Pollaris", "", "", vendPollaris)
	if err != nil {
		resources.Logger().Error(err.Error())
		return
	}
	_, ok := resp.(*l8tpollaris.L8Pollaris)
	if ok {
		resources.Logger().Info("Added ", vendPollaris.Name, " successfully")
	}
	time.Sleep(time.Second)
}
