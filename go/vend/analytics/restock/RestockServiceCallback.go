/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/

package restock

import (
	common "github.com/saichler/l8vendingmachine/go/vend/common"
	"github.com/saichler/l8vendingmachine/go/types/vend"
	"github.com/saichler/l8types/go/ifs"
)

func newRestockServiceCallback(vnic ifs.IVNic) ifs.IServiceCallback {
	return common.NewValidation(&vend.VendRestockRecommendation{}, vnic).
		Require(func(v interface{}) string { return v.(*vend.VendRestockRecommendation).RecommendationId }, "RecommendationId").
		Build()
}
