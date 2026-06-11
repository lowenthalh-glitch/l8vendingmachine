/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/warehouse/purchaseorders"
	"github.com/saichler/l8vendingmachine/go/vend/warehouse/stockmovements"
	"github.com/saichler/l8vendingmachine/go/vend/warehouse/suppliers"
	"github.com/saichler/l8vendingmachine/go/vend/warehouse/vehicleloads"
	"github.com/saichler/l8vendingmachine/go/vend/warehouse/warehouses"
	"github.com/saichler/l8vendingmachine/go/vend/warehouse/warehousestock"
)

func collectWarehouseActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { warehouses.Activate(creds, dbname, nic) },
		func() { warehousestock.Activate(creds, dbname, nic) },
		func() { suppliers.Activate(creds, dbname, nic) },
		func() { purchaseorders.Activate(creds, dbname, nic) },
		func() { stockmovements.Activate(creds, dbname, nic) },
		func() { vehicleloads.Activate(creds, dbname, nic) },
	}
}
