/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/temperature/tempreadings"
)

func collectTemperatureActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { tempreadings.Activate(creds, dbname, nic) },
	}
}
