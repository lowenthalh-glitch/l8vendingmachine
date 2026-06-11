/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/

package topperformers

import (
	common "github.com/saichler/l8vendingmachine/go/vend/common"
	vend "github.com/saichler/l8vendingmachine/go/types/vend"
	"github.com/saichler/l8types/go/ifs"
)

const (
	ServiceName = "TopPerf"
	ServiceArea = byte(10)
)

func Activate(creds, dbname string, vnic ifs.IVNic) {
	common.ActivateService(common.ServiceConfig{
		ServiceName: ServiceName, ServiceArea: ServiceArea,
		PrimaryKey: "PerformerId",
	}, &vend.VendTopPerformer{}, &vend.VendTopPerformerList{}, creds, dbname, vnic)
}
