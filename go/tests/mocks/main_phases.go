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

	whs := generateWarehouses()
	wids := extractIDs(whs, func(v interface{}) string { return v.(*vend.VendWarehouse).WarehouseId })
	if err := runOp(client, "Warehouses", "/vend/10/Warehouse",
		&vend.VendWarehouseList{List: whs}, wids, &store.WarehouseIDs); err != nil {
		return err
	}

	sups := generateSuppliers()
	sids := extractIDs(sups, func(v interface{}) string { return v.(*vend.VendSupplier).SupplierId })
	return runOp(client, "Suppliers", "/vend/10/Supplier",
		&vend.VendSupplierList{List: sups}, sids, &store.SupplierIDs)
}

func PrintSummary(store *MockDataStore) {
	fmt.Println("\n=== Mock Data Summary ===")
	fmt.Printf("Locations:        %d\n", len(store.LocationIDs))
	fmt.Printf("Machine Groups:   %d\n", len(store.GroupIDs))
	fmt.Printf("Warehouses:       %d\n", len(store.WarehouseIDs))
	fmt.Printf("Suppliers:        %d\n", len(store.SupplierIDs))
	fmt.Println("========================")
	fmt.Println("Machine inventory data comes from the Nayax simulator via collection pipeline.")
}
