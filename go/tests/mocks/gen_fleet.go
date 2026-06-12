package mocks

import (
	"fmt"
	"math/rand"

	l8common "github.com/saichler/l8common/go/types/l8common"
	"github.com/saichler/l8vendingmachine/go/types/vend"
)

// Austin TX area coordinates for the 10 machine locations
var locationCoords = [][2]float64{
	{30.2672, -97.7431},  // Building A - Lobby (Downtown)
	{30.2950, -97.7425},  // Building C - Break Room (UT area)
	{30.2840, -97.7320},  // Gym - Main Entrance (East Austin)
	{30.3074, -97.7571},  // Hospital - Cafeteria Wing (North)
	{30.1975, -97.6664},  // Airport - Terminal B (ABIA)
	{30.2649, -97.7399},  // Hotel - Lobby (Congress Ave)
	{30.3621, -97.6987},  // Mall - Food Court (Domain)
	{30.2175, -97.7650},  // Factory - Canteen (South)
	{30.2849, -97.7341},  // University - Student Center
	{30.2632, -97.7290},  // Train Station - Platform 3
}

func generateLocations() []*vend.VendLocation {
	items := make([]*vend.VendLocation, len(LocationNames))
	for i, name := range LocationNames {
		items[i] = &vend.VendLocation{
			LocationId:   genID("loc", i),
			Name:         name,
			LocationType: LocationTypes[i],
			Timezone:     "America/Chicago",
			Coordinates:  &vend.VendGpsCoordinates{Latitude: locationCoords[i][0], Longitude: locationCoords[i][1]},
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

func generateFacilities() []*vend.VendStockingFacility {
	addresses := []struct {
		line1 string
		zip   string
		lat   float64
		lng   float64
	}{
		{"7000 Burleson Rd", "78744", 30.1974, -97.7200},
		{"13000 N Lamar Blvd", "78753", 30.3980, -97.6929},
		{"4500 S Congress Ave", "78745", 30.2120, -97.7649},
	}

	items := make([]*vend.VendStockingFacility, len(FacilityNames))
	for i, name := range FacilityNames {
		addr := addresses[i]
		items[i] = &vend.VendStockingFacility{
			FacilityId:              genID("fac", i),
			Name:                    name,
			Code:                    FacilityCodes[i],
			Address:                 &l8common.Address{Line1: addr.line1, City: "Austin", StateProvince: "TX", PostalCode: addr.zip, CountryCode: "US"},
			Coordinates:             &vend.VendGpsCoordinates{Latitude: addr.lat, Longitude: addr.lng},
			Timezone:                "America/Chicago",
			TotalStorageSqFt:        int32(8000 + i*4000),
			RefrigeratedStorageSqFt: int32(2000 + i*1000),
			LoadingDocks:            int32(3 + i),
			MaxTrucksParked:         int32(6 + i*2),
			Status:                  vend.VendFacilityStatus_VEND_FACILITY_STATUS_ACTIVE,
			OperatingHoursStart:     "05:00",
			OperatingHoursEnd:       "22:00",
			OperatingDays: []vend.VendDayOfWeek{
				vend.VendDayOfWeek_VEND_DAY_OF_WEEK_MONDAY,
				vend.VendDayOfWeek_VEND_DAY_OF_WEEK_TUESDAY,
				vend.VendDayOfWeek_VEND_DAY_OF_WEEK_WEDNESDAY,
				vend.VendDayOfWeek_VEND_DAY_OF_WEEK_THURSDAY,
				vend.VendDayOfWeek_VEND_DAY_OF_WEEK_FRIDAY,
				vend.VendDayOfWeek_VEND_DAY_OF_WEEK_SATURDAY,
			},
			Stock:        generateFacilityStock(),
			ManagerName:  fmt.Sprintf("%s %s", firstNames[i%len(firstNames)], lastNames[i%len(lastNames)]),
			ManagerPhone: randomPhone(),
			ManagerEmail: fmt.Sprintf("facility%d@vendingco.com", i+1),
			AuditInfo:    createAuditInfo(),
		}
	}
	return items
}

func generateFacilityStock() []*vend.VendFacilityStockItem {
	stock := make([]*vend.VendFacilityStockItem, len(simulatorProducts))
	for j, p := range simulatorProducts {
		maxQty := p.maxQty * 10 // facilities hold 10x truck capacity
		qty := int32(float64(maxQty) * (0.5 + rand.Float64()*0.5))
		reorder := maxQty / 4
		stock[j] = &vend.VendFacilityStockItem{
			ProductName:  p.name,
			Sku:          p.sku,
			Price:        p.price,
			Quantity:     qty,
			MaxQuantity:  maxQty,
			ReorderPoint: reorder,
		}
	}
	return stock
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
