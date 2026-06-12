package mocks

import (
	"fmt"

	"github.com/saichler/l8vendingmachine/go/types/vend"
)

// RunAllPhases generates mock data for business-layer entities only.
// Machine inventory data comes from the Nayax simulator via the collection pipeline.
func RunAllPhases(client *VendClient, store *MockDataStore) {
	runPhase("Phase 1: Business Foundation", func() error {
		return generateBusinessFoundation(client, store)
	})
	runPhase("Phase 2: Alarm Definitions", func() error {
		return seedAlarmDefinitions(client)
	})
	runPhase("Phase 3: Historical Snapshots", func() error {
		return seedHistoricalSnapshots(client)
	})
}

func generateBusinessFoundation(client *VendClient, store *MockDataStore) error {
	locs := generateLocations()
	lids := extractIDs(locs, func(v interface{}) string { return v.(*vend.VendLocation).LocationId })
	if err := runOp(client, "Locations", "/vend/10/Location",
		&vend.VendLocationList{List: locs}, lids, &store.LocationIDs); err != nil {
		return err
	}

	groups := generateMachineGroups(nil)
	gids := extractIDs(groups, func(v interface{}) string { return v.(*vend.VendMachineGroup).GroupId })
	if err := runOp(client, "Machine Groups", "/vend/10/MachGrp",
		&vend.VendMachineGroupList{List: groups}, gids, &store.GroupIDs); err != nil {
		return err
	}

	facs := generateFacilities()
	fids := extractIDs(facs, func(v interface{}) string { return v.(*vend.VendStockingFacility).FacilityId })
	if err := runOp(client, "Facilities", "/vend/10/Facility",
		&vend.VendStockingFacilityList{List: facs}, fids, &store.FacilityIDs); err != nil {
		return err
	}

	sups := generateSuppliers()
	sids := extractIDs(sups, func(v interface{}) string { return v.(*vend.VendSupplier).SupplierId })
	if err := runOp(client, "Suppliers", "/vend/10/Supplier",
		&vend.VendSupplierList{List: sups}, sids, &store.SupplierIDs); err != nil {
		return err
	}

	trucks := generateTrucks(store)
	tids := extractIDs(trucks, func(v interface{}) string { return v.(*vend.VendDeliveryTruck).TruckId })
	if err := runOp(client, "Trucks", "/vend/10/Truck",
		&vend.VendDeliveryTruckList{List: trucks}, tids, &store.TruckIDs); err != nil {
		return err
	}

	drivers := generateDrivers(store)
	dids := extractIDs(drivers, func(v interface{}) string { return v.(*vend.VendDriver).DriverId })
	return runOp(client, "Drivers", "/vend/10/Driver",
		&vend.VendDriverList{List: drivers}, dids, &store.DriverIDs)
}

func PrintSummary(store *MockDataStore) {
	fmt.Println("\n=== Mock Data Summary ===")
	fmt.Printf("Locations:        %d\n", len(store.LocationIDs))
	fmt.Printf("Machine Groups:   %d\n", len(store.GroupIDs))
	fmt.Printf("Facilities:       %d\n", len(store.FacilityIDs))
	fmt.Printf("Suppliers:        %d\n", len(store.SupplierIDs))
	fmt.Printf("Trucks:           %d\n", len(store.TruckIDs))
	fmt.Printf("Drivers:          %d\n", len(store.DriverIDs))
	fmt.Println("========================")
	fmt.Println("Machine inventory data comes from the Nayax simulator via collection pipeline.")
}
