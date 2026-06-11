/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package boot

import "github.com/saichler/l8pollaris/go/types/l8tpollaris"

// vendCashbox deferred — per-machine endpoint only (/lynx/v1/payment/cashbox)

func createVendDevicesPoll(p *l8tpollaris.L8Pollaris) {
	p.Polling["vendDevices"] = createVendPoll(
		"vendDevices",
		"/lynx/v1/devices",
		EVERY_5_MINUTES,
		false,
		"vendmachine.paymentstatus",
		"connectionStatus:vendmachine.paymentstatus.overallstatus",
	)
}
