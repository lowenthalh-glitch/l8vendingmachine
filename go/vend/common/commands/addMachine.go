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
	"time"

	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8web/go/web/client"
)

// AddMachine registers a vending machine IP as an L8PTarget for REST API polling.
func AddMachine(ip string, port int32, rc *client.RestClient, resources ifs.IResources) {
	defer time.Sleep(time.Second)
	machine := CreateVendingMachine(ip, port)
	resp, err := rc.POST("0/"+targets.ServiceName, "Device",
		"", "", machine)
	if err != nil {
		resources.Logger().Error(err.Error())
		return
	}
	_, ok := resp.(*l8tpollaris.L8PTarget)
	if ok {
		resources.Logger().Info("Added vending machine ", machine.TargetId, " successfully")
	}
}
