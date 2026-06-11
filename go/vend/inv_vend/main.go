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
	"time"

	"github.com/saichler/l8alarms/go/types/alm"
	"github.com/saichler/l8bus/go/overlay/vnic"
	inventory "github.com/saichler/l8inventory/go/inv/service"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/vend/fleet/machines"
	vend "github.com/saichler/l8vendingmachine/go/types/vend"
)

func main() {
	res := vendcommon.CreateResources("inv-vend")
	res.Logger().Info("Starting inv-vend")
	ifs.SetNetworkMode(ifs.NETWORK_K8s)

	nic := vnic.NewVirtualNetworkInterface(res, nil)
	nic.Start()
	nic.WaitForConnection()
	res.Logger().Info("Registering inv-vend service")

	// Register types needed for bridge and threshold communication
	nic.Resources().Registry().Register(&l8tpollaris.L8PTarget{})
	nic.Resources().Registry().Register(&l8tpollaris.L8PTargetList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&l8tpollaris.L8PTarget{}, "TargetId")
	nic.Resources().Registry().Register(&vend.VendFleetMachine{})
	nic.Resources().Registry().Register(&vend.VendFleetMachineList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendFleetMachine{}, "MachineId")
	nic.Resources().Registry().Register(&alm.Alarm{})
	nic.Resources().Registry().Register(&alm.AlarmList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&alm.Alarm{}, "AlarmId")

	// Analytics types
	nic.Resources().Registry().Register(&vend.VendFleetInventory{})
	nic.Resources().Registry().Register(&vend.VendFleetInventoryList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendFleetInventory{}, "SummaryId")
	nic.Resources().Registry().Register(&vend.VendSlotPerformance{})
	nic.Resources().Registry().Register(&vend.VendSlotPerformanceList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendSlotPerformance{}, "PerformanceId")
	nic.Resources().Registry().Register(&vend.VendForecast{})
	nic.Resources().Registry().Register(&vend.VendForecastList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendForecast{}, "ForecastId")
	nic.Resources().Registry().Register(&vend.VendInventorySnapshot{})
	nic.Resources().Registry().Register(&vend.VendInventorySnapshotList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendInventorySnapshot{}, "SnapshotId")
	nic.Resources().Registry().Register(&vend.VendTopPerformer{})
	nic.Resources().Registry().Register(&vend.VendTopPerformerList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendTopPerformer{}, "PerformerId")
	nic.Resources().Registry().Register(&vend.VendMachineProfile{})
	nic.Resources().Registry().Register(&vend.VendMachineProfileList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendMachineProfile{}, "ProfileId")
	nic.Resources().Registry().Register(&vend.VendRestockRecommendation{})
	nic.Resources().Registry().Register(&vend.VendRestockRecommendationList{})
	nic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&vend.VendRestockRecommendation{}, "RecommendationId")

	inventory.Activate(vendcommon.VendMachine_Links_ID, &vend.VendMachine{}, &vend.VendMachineList{}, nic, "MachineId")

	invCenter := inventory.Inventory(res, vendcommon.Vend_Cache_Service_Name, vendcommon.Vend_Cache_Service_Area)
	invCenter.AddMetadata("Online", Online)

	// Bridge: populate Fleet Machine service from VCache
	go bridgeVCacheToFleet(nic)

	// Threshold evaluator: check inventory levels and generate alerts
	go evaluateThresholds(nic)

	// Analytics: compute fleet inventory summaries, slot performance, and forecasts
	go computeAnalytics(nic)

	// Restock recommendations: evaluate every 30 min using profiles
	go computeRestockRecommendations(nic)

	// Data retention: monthly cleanup of old snapshots
	go cleanOldSnapshots(nic)

	vendcommon.WaitForSignal(nic.Resources())
}

func Online(any interface{}) (bool, string) {
	if any == nil {
		return false, ""
	}
	vm := any.(*vend.VendMachine)
	if vm.Machines != nil && len(vm.Machines) > 0 {
		return true, ""
	}
	return false, ""
}

func bridgeVCacheToFleet(nic ifs.IVNic) {
	// Wait for all services to be ready and VCache to have data
	time.Sleep(30 * time.Second)

	for {
		// Get all VendMachine entries from VCache via the service API
		results, err := vendcommon.GetEntities(
			vendcommon.Vend_Cache_Service_Name,
			vendcommon.Vend_Cache_Service_Area,
			&vend.VendMachine{}, nic)
		if err != nil {
			nic.Resources().Logger().Error("Bridge: failed to get VCache data: ", err.Error())
			time.Sleep(30 * time.Second)
			continue
		}
		if results != nil {
			count := 0
			for _, elem := range results {
				vm, ok := elem.(*vend.VendMachine)
				if !ok || vm.Machines == nil {
					continue
				}
				for _, info := range vm.Machines {
					fleetMachine := &vend.VendFleetMachine{
						MachineId:         info.MachineId,
						Name:              info.Name,
						Type:              info.Type,
						Model:             info.Model,
						Status:            info.Status,
						DeviceId:          info.DeviceId,
						DailyTransactions: info.DailyTransactions,
						LastTransactionAt: info.LastTransactionAt,
						ManagementIp:      vm.MachineId,
						LocationAddress:   info.LocationAddress,
						LocationCity:      info.LocationCity,
						LocationState:     info.LocationState,
						LocationLat:       info.LocationLat,
						LocationLng:       info.LocationLng,
					}
					vendcommon.PostEntity(machines.ServiceName, machines.ServiceArea, fleetMachine, nic)

					// Create per-machine target for slot inventory polling
					createPerMachineTarget(info.MachineId, vm.MachineId, nic)
					count++
				}
			}
			if count > 0 {
				fmt.Printf("[INV-VEND] Bridged %d machines from VCache to Fleet\n", count)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func createPerMachineTarget(machineId, managementIp string, nic ifs.IVNic) {
	// Check if target already exists
	existing, _ := vendcommon.GetEntity(
		targets.ServiceName, targets.ServiceArea,
		&l8tpollaris.L8PTarget{TargetId: machineId}, nic)
	if existing != nil {
		return
	}

	// Parse management IP and port
	addr := managementIp
	port := int32(8443)

	target := &l8tpollaris.L8PTarget{
		TargetId:      machineId,
		LinksId:       vendcommon.VendPerMachine_Links_ID,
		InventoryType: l8tpollaris.L8PTargetType_Vending_Machine,
		State:         l8tpollaris.L8PTargetState_Up,
		Hosts: map[string]*l8tpollaris.L8PHost{
			machineId: {
				HostId: machineId,
				Configs: map[int32]*l8tpollaris.L8PHostProtocol{
					int32(l8tpollaris.L8PProtocol_L8PRESTAPI): {
						Protocol: l8tpollaris.L8PProtocol_L8PRESTAPI,
						Addr:     addr,
						Port:     port,
						Timeout:  30,
						Ainfo:    &l8tpollaris.AuthInfo{},
					},
				},
			},
		},
	}

	// Store in targets service
	vendcommon.PostEntity(targets.ServiceName, targets.ServiceArea, target, nic)

	// Send directly to collector (PostEntity notification may prevent callback distribution)
	collectorService, collectorArea := vendcommon.Collector_Service_Name, vendcommon.Collector_Service_Area
	err := nic.RoundRobin(collectorService, collectorArea, ifs.POST, target)
	if err != nil {
		nic.Resources().Logger().Error("Failed to send target to collector: ", err.Error())
	}
	fmt.Printf("[INV-VEND] Created per-machine target: %s → %s:%d\n", machineId, addr, port)
}
