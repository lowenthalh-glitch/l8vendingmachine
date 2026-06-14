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

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/saichler/l8alarms/go/alm/alarmdefinitions"
	"github.com/saichler/l8alarms/go/alm/alarmfilters"
	"github.com/saichler/l8alarms/go/alm/archivedalarms"
	"github.com/saichler/l8alarms/go/alm/archivedevents"
	"github.com/saichler/l8alarms/go/alm/correlationrules"
	"github.com/saichler/l8alarms/go/alm/escalationpolicies"
	"github.com/saichler/l8alarms/go/alm/events"
	"github.com/saichler/l8alarms/go/alm/maintenancewindows"
	"github.com/saichler/l8alarms/go/alm/notificationpolicies"
	"github.com/saichler/l8alarms/go/types/alm"
	"github.com/saichler/l8bus/go/overlay/vnic"
	l8common "github.com/saichler/l8common/go/common"
	evtservices "github.com/saichler/l8events/go/services"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/ipsegment"
	"github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/vend/route/optimizer"
	"github.com/saichler/l8vendingmachine/go/vend/services"
)

func main() {
	fmt.Println("[vend] Creating resources...")
	res := common.CreateResources("VendServices")
	fmt.Println("[vend] Setting network mode...")
	ifs.SetNetworkMode(ifs.NETWORK_K8s)
	fmt.Println("[vend] Creating VNIC...")
	nic := vnic.NewVirtualNetworkInterface(res, nil)
	fmt.Println("[vend] Starting VNIC...")
	nic.Start()
	fmt.Println("[vend] Waiting for VNet connection...")
	nic.WaitForConnection()
	fmt.Println("[vend] Connected to VNet!")

	// Start postgres if not local
	if len(os.Args) == 1 {
		startDb(nic)
	} else {
		// Local mode: use localhost for DB so docker-proxy routing issues are avoided
		ipsegment.MachineIP = "127.0.0.1"
		_, user, pass, _, err := nic.Resources().Security().Credential(common.DB_CREDS, common.DB_NAME, nic.Resources())
		if err == nil && user == "admin" && pass == "admin" {
			common.DB_NAME = "admin"
		}
	}

	fmt.Println("[vend] Activating all services...")
	services.ActivateAllServices(common.DB_CREDS, common.DB_NAME, nic)

	fmt.Println("[vend] Activating events service...")
	evtservices.ActivateEvents(common.DB_CREDS, common.DB_NAME, nic)

	// Activate Credentials service (for targets UI credential management)
	fmt.Println("[vend] Activating Credentials service...")
	services.ActivateCredentials(nic)

	// Activate Pollaris targets and polling configuration services
	// Activate l8alarms services individually (skip enrichment — requires l8topology)
	fmt.Println("[vend] Activating l8alarms services...")
	alarmdefinitions.Activate(common.DB_CREDS, common.DB_NAME, nic)
	// Use local alarm service without maintenance window Before hook (l8alarms bug: typed nil panic)
	l8common.ActivateService(l8common.ServiceConfig{
		ServiceName: "Alarm", ServiceArea: byte(10),
		PrimaryKey: "AlarmId",
	}, &alm.Alarm{}, &alm.AlarmList{}, common.DB_CREDS, common.DB_NAME, nic)
	events.Activate(common.DB_CREDS, common.DB_NAME, nic)
	correlationrules.Activate(common.DB_CREDS, common.DB_NAME, nic)
	notificationpolicies.Activate(common.DB_CREDS, common.DB_NAME, nic)
	escalationpolicies.Activate(common.DB_CREDS, common.DB_NAME, nic)
	maintenancewindows.Activate(common.DB_CREDS, common.DB_NAME, nic)
	alarmfilters.Activate(common.DB_CREDS, common.DB_NAME, nic)
	archivedalarms.Activate(common.DB_CREDS, common.DB_NAME, nic)
	archivedevents.Activate(common.DB_CREDS, common.DB_NAME, nic)

	fmt.Println("[vend] Activating Pollaris targets...")
	pollaris.Activate(nic)
	targets.Activate(common.DB_CREDS, common.DB_NAME, nic)
	fmt.Println("[vend] Activating route optimizer...")
	optimizer.ActivateOptimizer(nic)
	fmt.Println("[vend] All services activated!")

	common.WaitForSignal(res)
}

func startDb(nic ifs.IVNic) {
	_, user, pass, _, err := nic.Resources().Security().Credential(common.DB_CREDS, common.DB_NAME, nic.Resources())
	if err != nil {
		panic(common.DB_CREDS + " " + err.Error())
	}
	if user == "admin" && pass == "admin" {
		common.DB_NAME = "admin"
	}

	cmd := exec.Command("nohup", "/start-postgres.sh", common.DB_NAME, user, pass)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	time.Sleep(time.Second * 5)
}
