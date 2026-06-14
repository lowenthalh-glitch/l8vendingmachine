/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/

package topperformers

import (
	common "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	"github.com/saichler/l8types/go/ifs"
)

func newTopPerformerServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
	return common.NewValidation(&vend.VendTopPerformer{}, vnic).
		Require(func(v interface{}) string { return v.(*vend.VendTopPerformer).PerformerId }, "PerformerId").
		Build()
}
