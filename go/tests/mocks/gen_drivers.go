package mocks

import (
	"fmt"
	"time"

	l8common "github.com/saichler/l8common/go/types/l8common"
	"github.com/saichler/l8vendingmachine/go/types/vend"
)

// Austin TX area addresses and current GPS positions for the 5 drivers
var driverAddresses = []struct {
	line1    string
	city     string
	zip      string
	currLat  float64
	currLng  float64
}{
	{"1200 Barton Springs Rd", "Austin", "78704", 30.2610, -97.7550},
	{"8500 Shoal Creek Blvd", "Austin", "78757", 30.3100, -97.7350},
	{"4600 Mueller Blvd", "Austin", "78723", 30.2950, -97.7100},
	{"2100 S Lamar Blvd", "Austin", "78704", 30.2450, -97.7700},
	{"12000 Research Blvd", "Austin", "78759", 30.3750, -97.7250},
}

// Truck type → required license class
// BOX_TRUCK(1), REFRIGERATED(3) → Class B (medium truck 26,001+ lbs)
// CARGO_VAN(2), SPRINTER(4), PICKUP(5) → Class C (standard)
func licenseForTruck(truckIdx int) vend.VendLicenseClass {
	// From gen_trucks.go truckTypes array:
	// [0]=BOX, [1]=VAN, [2]=REFRIG, [3]=SPRINTER, [4]=BOX
	switch truckIdx {
	case 0, 2, 4:
		return vend.VendLicenseClass_VEND_LICENSE_CLASS_B
	default:
		return vend.VendLicenseClass_VEND_LICENSE_CLASS_C
	}
}

func generateDrivers(store *MockDataStore) []*vend.VendDriver {
	count := 5
	items := make([]*vend.VendDriver, count)

	for i := 0; i < count; i++ {
		fn := DriverFirstNames[i%len(DriverFirstNames)]
		ln := DriverLastNames[i%len(DriverLastNames)]
		addr := driverAddresses[i]
		lc := licenseForTruck(i)
		truckId := pickRef(store.TruckIDs, i)
		hireDate := time.Now().AddDate(-(1 + i), 0, 0).Unix()

		// Mon-Fri schedule, staggered start times
		// End locations: drv-0,2,4 end at facilities; drv-1,3 end at home (blank)
		startHour := 6 + (i % 3) // 06:00, 07:00, 08:00, 06:00, 07:00
		endLocId := ""
		if i%2 == 0 && len(store.FacilityIDs) > 0 {
			endLocId = pickRef(store.FacilityIDs, i/2)
		}
		// Shift durations: vary by driver (7-9 hours)
		shiftMin := int32(480 + (i%3)*30) // 480, 510, 540, 480, 510
		schedule := make([]*vend.VendDriverScheduleDay, 5)
		for d := 0; d < 5; d++ {
			schedule[d] = &vend.VendDriverScheduleDay{
				Day:                  vend.VendDayOfWeek(d + 1), // 1=Monday .. 5=Friday
				StartTime:            fmt.Sprintf("%02d:00", startHour),
				EndLocationId:        endLocId,
				ShiftDurationMinutes: shiftMin,
			}
		}

		// Driver skills — all have basic, some have specializations
		skills := []string{"restocking", "cash-handling"}
		if i == 0 || i == 2 {
			skills = append(skills, "refrigeration")
		}
		if i == 1 || i == 4 {
			skills = append(skills, "ev-charger")
		}

		items[i] = &vend.VendDriver{
			DriverId:      genID("drv", i),
			FirstName:     fn,
			LastName:      ln,
			Phone:         randomPhone(),
			Email:         fmt.Sprintf("%s.%s@vendingco.com", fn, ln),
			LicenseNumber: fmt.Sprintf("TX-%s%s-%04d", fn[:1], ln[:1], 1000+i),
			LicenseClass:  lc,
			TruckId:       truckId,
			Skills:        skills,
			IsActive:      true,
			HireDate:      hireDate,
			HomeAddress: &l8common.Address{
				Line1:         addr.line1,
				City:          addr.city,
				StateProvince: "TX",
				PostalCode:    addr.zip,
				CountryCode:   "US",
			},
			CurrentLatitude:    addr.currLat,
			CurrentLongitude:   addr.currLng,
			LastLocationUpdate: time.Now().Add(-time.Duration(5+i*3) * time.Minute).Unix(),
			Schedule:           schedule,
			AuditInfo:          createAuditInfo(),
		}
	}
	return items
}
