/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 * Layer 8 Ecosystem - Apache 2.0
 */
package optimizer

import (
	"fmt"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	vendcommon "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/vend/route/drivers"
	"github.com/saichler/l8vendingmachine/go/vend/route/trucks"
	"github.com/saichler/l8vendingmachine/go/vend/warehouse/facilities"
)

// Assignment holds the matched truck and driver for a cluster.
type Assignment struct {
	TruckId  string
	DriverId string
	Truck    *vend.VendDeliveryTruck
	Driver   *vend.VendDriver
	StartLat float64
	StartLng float64
}

// licenseRequired returns the minimum license class for a truck type.
func licenseRequired(truckType vend.VendTruckType) vend.VendLicenseClass {
	switch truckType {
	case vend.VendTruckType_VEND_TRUCK_TYPE_BOX_TRUCK,
		vend.VendTruckType_VEND_TRUCK_TYPE_REFRIGERATED:
		return vend.VendLicenseClass_VEND_LICENSE_CLASS_B
	default:
		return vend.VendLicenseClass_VEND_LICENSE_CLASS_C
	}
}

func licenseCovers(driverClass, required vend.VendLicenseClass) bool {
	return driverClass >= required
}

// AssignResources matches a cluster to an available truck and driver.
func AssignResources(cluster *Cluster, allTrucks []*vend.VendDeliveryTruck,
	allDrivers []*vend.VendDriver, dayOfWeek vend.VendDayOfWeek,
	assignedTrucks map[string]bool, nic ifs.IVNic) (*Assignment, error) {

	for _, truck := range allTrucks {
		if truck.Status != vend.VendTruckStatus_VEND_TRUCK_STATUS_ACTIVE {
			continue
		}
		if assignedTrucks[truck.TruckId] {
			continue
		}
		reqLicense := licenseRequired(truck.Type)

		for _, driver := range allDrivers {
			if !driver.IsActive || driver.TruckId != truck.TruckId {
				continue
			}
			if !licenseCovers(driver.LicenseClass, reqLicense) {
				continue
			}
			schedEntry := findScheduleDay(driver.Schedule, dayOfWeek)
			if schedEntry == nil {
				continue
			}
			startLat, startLng := resolveStartLocation(driver, schedEntry, nic)

			assignedTrucks[truck.TruckId] = true
			return &Assignment{
				TruckId: truck.TruckId, DriverId: driver.DriverId,
				Truck: truck, Driver: driver,
				StartLat: startLat, StartLng: startLng,
			}, nil
		}
	}
	return nil, fmt.Errorf("no available truck+driver for cluster at (%.4f, %.4f)",
		cluster.CentroidLat, cluster.CentroidLng)
}

func findScheduleDay(schedule []*vend.VendDriverScheduleDay, day vend.VendDayOfWeek) *vend.VendDriverScheduleDay {
	for _, s := range schedule {
		if s.Day == day {
			return s
		}
	}
	return nil
}

func resolveStartLocation(driver *vend.VendDriver, sched *vend.VendDriverScheduleDay, nic ifs.IVNic) (float64, float64) {
	if sched.StartLocationId != "" {
		results, err := vendcommon.GetEntities("Location", byte(10), &vend.VendLocation{LocationId: sched.StartLocationId}, nic)
		if err == nil && len(results) > 0 {
			if loc, ok := results[0].(*vend.VendLocation); ok && loc.Coordinates != nil {
				return loc.Coordinates.Latitude, loc.Coordinates.Longitude
			}
		}
	}
	if driver.CurrentLatitude != 0 || driver.CurrentLongitude != 0 {
		return driver.CurrentLatitude, driver.CurrentLongitude
	}
	return 0, 0
}

// LoadAllTrucks fetches all delivery trucks.
func LoadAllTrucks(nic ifs.IVNic) []*vend.VendDeliveryTruck {
	results, err := vendcommon.GetEntities(trucks.ServiceName, trucks.ServiceArea, &vend.VendDeliveryTruck{}, nic)
	if err != nil {
		return nil
	}
	out := make([]*vend.VendDeliveryTruck, 0, len(results))
	for _, r := range results {
		if t, ok := r.(*vend.VendDeliveryTruck); ok {
			out = append(out, t)
		}
	}
	return out
}

// LoadAllDrivers fetches all drivers.
func LoadAllDrivers(nic ifs.IVNic) []*vend.VendDriver {
	results, err := vendcommon.GetEntities(drivers.ServiceName, drivers.ServiceArea, &vend.VendDriver{}, nic)
	if err != nil {
		return nil
	}
	out := make([]*vend.VendDriver, 0, len(results))
	for _, r := range results {
		if d, ok := r.(*vend.VendDriver); ok {
			out = append(out, d)
		}
	}
	return out
}

// LoadAllFacilities fetches all stocking facilities.
func LoadAllFacilities(nic ifs.IVNic) []*vend.VendStockingFacility {
	results, err := vendcommon.GetEntities(facilities.ServiceName, facilities.ServiceArea, &vend.VendStockingFacility{}, nic)
	if err != nil {
		return nil
	}
	out := make([]*vend.VendStockingFacility, 0, len(results))
	for _, r := range results {
		if f, ok := r.(*vend.VendStockingFacility); ok {
			out = append(out, f)
		}
	}
	return out
}
