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

package ui

import (
	"strconv"

	"github.com/saichler/l8alarms/go/types/alm"
	"github.com/saichler/l8bus/go/overlay/vnic"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
	l8events "github.com/saichler/l8types/go/types/l8events"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	"github.com/saichler/l8vendingmachine/go/vend/common"
)

func CreateVnic(vnetId uint32) ifs.IVNic {
	resources := common.CreateResources("web-" + strconv.Itoa(int(vnetId)))

	RegisterTypes(resources)

	nic := vnic.NewVirtualNetworkInterface(resources, nil)
	nic.Resources().SysConfig().KeepAliveIntervalSeconds = 60
	nic.Start()
	nic.WaitForConnection()

	return nic
}

func RegisterTypes(resources ifs.IResources) {
	// Pollaris (polling config and targets — match probler's RegisterTypes exactly)
	resources.Registry().Register(&l8tpollaris.L8Pollaris{})
	resources.Introspector().Decorators().AddPrimaryKeyDecorator(&l8tpollaris.L8PTarget{}, "TargetId")
	resources.Registry().Register(&l8tpollaris.L8PTarget{})
	resources.Registry().Register(&l8tpollaris.L8PTargetList{})
	resources.Registry().Register(&l8tpollaris.TargetAction{})
	resources.Registry().Register(&l8tpollaris.CJob{})

	// Credentials (used by targets UI for authentication config)
	common.RegisterType(resources, &l8api.L8Credentials{}, &l8api.L8CredentialsList{}, "Id")

	// Management Systems (VCache inventory)
	common.RegisterType(resources, &vend.VendMachine{}, &vend.VendMachineList{}, "MachineId")

	// Fleet (Prime Object CRUD)
	common.RegisterType(resources, &vend.VendFleetMachine{}, &vend.VendFleetMachineList{}, "MachineId")
	common.RegisterType(resources, &vend.VendMachineGroup{}, &vend.VendMachineGroupList{}, "GroupId")
	common.RegisterType(resources, &vend.VendLocation{}, &vend.VendLocationList{}, "LocationId")

	// Inventory Management — deferred to business CRUD layer
	// VendMachine inventory is now served by inv_vend (l8inventory cache at /0/VCache)

	// Sales & Transactions
	common.RegisterType(resources, &vend.VendTransaction{}, &vend.VendTransactionList{}, "TransactionId")
	common.RegisterType(resources, &vend.VendSettlement{}, &vend.VendSettlementList{}, "SettlementId")

	// Payment Systems
	common.RegisterType(resources, &vend.VendCashPosition{}, &vend.VendCashPositionList{}, "PositionId")
	common.RegisterType(resources, &vend.VendCashCollection{}, &vend.VendCashCollectionList{}, "CollectionId")

	// Temperature & Refrigeration
	common.RegisterType(resources, &vend.VendTempReading{}, &vend.VendTempReadingList{}, "ReadingId")

	// l8alarms types (Alarms & Events section)
	common.RegisterType(resources, &alm.AlarmDefinition{}, &alm.AlarmDefinitionList{}, "DefinitionId")
	common.RegisterType(resources, &alm.Alarm{}, &alm.AlarmList{}, "AlarmId")
	common.RegisterType(resources, &alm.Event{}, &alm.EventList{}, "EventId")
	common.RegisterType(resources, &alm.CorrelationRule{}, &alm.CorrelationRuleList{}, "RuleId")
	common.RegisterType(resources, &alm.NotificationPolicy{}, &alm.NotificationPolicyList{}, "PolicyId")
	common.RegisterType(resources, &alm.EscalationPolicy{}, &alm.EscalationPolicyList{}, "PolicyId")
	common.RegisterType(resources, &alm.MaintenanceWindow{}, &alm.MaintenanceWindowList{}, "WindowId")
	common.RegisterType(resources, &alm.AlarmFilter{}, &alm.AlarmFilterList{}, "FilterId")
	common.RegisterType(resources, &alm.ArchivedAlarm{}, &alm.ArchivedAlarmList{}, "AlarmId")
	common.RegisterType(resources, &alm.ArchivedEvent{}, &alm.ArchivedEventList{}, "EventId")

	// Events (l8events EventRecord for SYS module)
	common.RegisterType(resources, &l8events.EventRecord{}, &l8events.EventRecordList{}, "EventId")

	// Alerts & Maintenance (VendAlert kept for backward compat, to be removed later)
	common.RegisterType(resources, &vend.VendAlert{}, &vend.VendAlertList{}, "AlertId")
	common.RegisterType(resources, &vend.VendWorkOrder{}, &vend.VendWorkOrderList{}, "WorkOrderId")
	common.RegisterType(resources, &vend.VendServiceVisit{}, &vend.VendServiceVisitList{}, "VisitId")

	// Route Optimization
	common.RegisterType(resources, &vend.VendRoute{}, &vend.VendRouteList{}, "RouteId")
	common.RegisterType(resources, &vend.VendDriver{}, &vend.VendDriverList{}, "DriverId")
	common.RegisterType(resources, &vend.VendDeliveryTruck{}, &vend.VendDeliveryTruckList{}, "TruckId")
	resources.Registry().Register(&vend.VendRouteOptRequest{})

	// AI Analytics
	common.RegisterType(resources, &vend.VendForecast{}, &vend.VendForecastList{}, "ForecastId")
	common.RegisterType(resources, &vend.VendSlotPerformance{}, &vend.VendSlotPerformanceList{}, "PerformanceId")
	common.RegisterType(resources, &vend.VendFleetInventory{}, &vend.VendFleetInventoryList{}, "SummaryId")
	common.RegisterType(resources, &vend.VendInventorySnapshot{}, &vend.VendInventorySnapshotList{}, "SnapshotId")
	common.RegisterType(resources, &vend.VendTopPerformer{}, &vend.VendTopPerformerList{}, "PerformerId")
	common.RegisterType(resources, &vend.VendMachineProfile{}, &vend.VendMachineProfileList{}, "ProfileId")
	common.RegisterType(resources, &vend.VendRestockRecommendation{}, &vend.VendRestockRecommendationList{}, "RecommendationId")

	// Access & Security
	common.RegisterType(resources, &vend.VendAccessEvent{}, &vend.VendAccessEventList{}, "EventId")

	// DEX Audit
	common.RegisterType(resources, &vend.VendDexAudit{}, &vend.VendDexAuditList{}, "AuditId")

	// Stocking Facilities & Supply Chain
	common.RegisterType(resources, &vend.VendStockingFacility{}, &vend.VendStockingFacilityList{}, "FacilityId")
	common.RegisterType(resources, &vend.VendSupplier{}, &vend.VendSupplierList{}, "SupplierId")
	common.RegisterType(resources, &vend.VendPurchaseOrder{}, &vend.VendPurchaseOrderList{}, "OrderId")
	common.RegisterType(resources, &vend.VendStockMovement{}, &vend.VendStockMovementList{}, "MovementId")
	common.RegisterType(resources, &vend.VendVehicleLoad{}, &vend.VendVehicleLoadList{}, "LoadId")

	// Dashboard & KPIs
	common.RegisterType(resources, &vend.VendKPI{}, &vend.VendKPIList{}, "KpiId")
	common.RegisterType(resources, &vend.VendDashboard{}, &vend.VendDashboardList{}, "DashboardId")

	// Compliance & Health Inspections
	common.RegisterType(resources, &vend.VendInspection{}, &vend.VendInspectionList{}, "InspectionId")
	common.RegisterType(resources, &vend.VendInspectionFinding{}, &vend.VendInspectionFindingList{}, "FindingId")
	common.RegisterType(resources, &vend.VendCertification{}, &vend.VendCertificationList{}, "CertificationId")

	// Scheduled Reports
	common.RegisterType(resources, &vend.VendReport{}, &vend.VendReportList{}, "ReportId")

	// Data Retention (Archived)
	common.RegisterType(resources, &vend.VendArchivedTransaction{}, &vend.VendArchivedTransactionList{}, "TransactionId")
	common.RegisterType(resources, &vend.VendArchivedTempReading{}, &vend.VendArchivedTempReadingList{}, "ReadingId")
	common.RegisterType(resources, &vend.VendArchivedAccessEvent{}, &vend.VendArchivedAccessEventList{}, "EventId")
}
