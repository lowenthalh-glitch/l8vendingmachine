/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/access/accessevents"
)

func collectAccessActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { accessevents.Activate(creds, dbname, nic) },
	}
}
