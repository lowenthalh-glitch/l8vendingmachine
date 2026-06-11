/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
package mocks

import (
	l8m "github.com/saichler/l8common/go/mocks"
	"net/http"
)

type VendClient = l8m.MockClient

func NewVendClient(baseURL string, httpClient *http.Client) *VendClient {
	return l8m.NewMockClient(baseURL, httpClient)
}
