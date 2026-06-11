/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 */
package services

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8vendingmachine/go/vend/retention/archivedaccess"
	"github.com/saichler/l8vendingmachine/go/vend/retention/archivedtemps"
	"github.com/saichler/l8vendingmachine/go/vend/retention/archivedtxns"
)

func collectRetentionActivations(creds, dbname string, nic ifs.IVNic) []func() {
	return []func(){
		func() { archivedtxns.Activate(creds, dbname, nic) },
		func() { archivedtemps.Activate(creds, dbname, nic) },
		func() { archivedaccess.Activate(creds, dbname, nic) },
	}
}
