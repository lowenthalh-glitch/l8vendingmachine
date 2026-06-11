package mocks

import (
	"fmt"
	"math/rand"

	"github.com/saichler/l8vendingmachine/go/types/vend"
)

func generateLocations() []*vend.VendLocation {
	items := make([]*vend.VendLocation, len(LocationNames))
	for i, name := range LocationNames {
		items[i] = &vend.VendLocation{
			LocationId:   genID("loc", i),
			Name:         name,
			LocationType: LocationTypes[i],
			Timezone:     "America/Chicago",
			ContactName:  fmt.Sprintf("%s %s", firstNames[rand.Intn(len(firstNames))], lastNames[rand.Intn(len(lastNames))]),
			ContactPhone: randomPhone(),
			ContactEmail: fmt.Sprintf("contact%d@vendingco.com", i+1),
			AuditInfo:    createAuditInfo(),
		}
	}
	return items
}

func generateMachineGroups(machines []*vend.VendMachine) []*vend.VendMachineGroup {
	items := make([]*vend.VendMachineGroup, len(GroupNames))
	machineCount := 0
	if machines != nil {
		machineCount = len(machines) / len(GroupNames)
	}
	for i, name := range GroupNames {
		items[i] = &vend.VendMachineGroup{
			GroupId:      genID("grp", i),
			Name:         name,
			Region:       fmt.Sprintf("Region %d", i+1),
			MachineCount: int32(machineCount),
			AuditInfo:    createAuditInfo(),
		}
	}
	return items
}

func generateWarehouses() []*vend.VendWarehouse {
	items := make([]*vend.VendWarehouse, len(WarehouseNames))
	for i, name := range WarehouseNames {
		items[i] = &vend.VendWarehouse{
			WarehouseId: genID("wh", i),
			Name:        name,
			AuditInfo:   createAuditInfo(),
		}
	}
	return items
}

func generateSuppliers() []*vend.VendSupplier {
	items := make([]*vend.VendSupplier, len(SupplierNames))
	for i, name := range SupplierNames {
		items[i] = &vend.VendSupplier{
			SupplierId: genID("sup", i),
			Name:       name,
			Status:     vend.VendSupplierStatus_VEND_SUPPLIER_STATUS_ACTIVE,
			AuditInfo:  createAuditInfo(),
		}
	}
	return items
}
