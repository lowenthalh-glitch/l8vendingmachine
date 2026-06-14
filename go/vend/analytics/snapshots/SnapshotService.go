/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/

package snapshots

import (
	common "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "InvSnap"
	ServiceArea = byte(10)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService(common.ServiceConfig{
		ServiceName: ServiceName, ServiceArea: ServiceArea,
		PrimaryKey: "SnapshotId", Callback: newSnapshotServiceCallback(vnic),
	}, &vend.VendInventorySnapshot{}, &vend.VendInventorySnapshotList{}, creds, dbname, vnic)
}
