/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"github.com/saichler/l8parser/go/parser/boot"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	vendboot "github.com/saichler/l8vendingmachine/go/vend/parser/boot"
)

func init() {
	targets.Links = &Links{}
	boot.RegisterPollaris(vendboot.GetVendingMachinePollaris())
	boot.RegisterPollaris(vendboot.GetPerMachinePollaris())
}

const (
	Collector_Service_Name = "VColl"
	Collector_Service_Area = byte(0)

	// Management system (fleet-wide collection → VCache)
	VendMachine_Links_ID      = "Vend"
	Vend_Cache_Service_Name   = "VCache"
	Vend_Cache_Service_Area   = byte(0)
	Vend_Persist_Service_Name = "VPersist"
	Vend_Persist_Service_Area = byte(0)
	Vend_Parser_Service_Name  = "VPars"
	Vend_Parser_Service_Area  = byte(0)
	Vend_Model_Name           = "vendmachine"

	// Per-machine (per-device collection → Fleet Machine CRUD service)
	VendPerMachine_Links_ID      = "VMach"
	VMach_Cache_Service_Name     = "Machine"
	VMach_Cache_Service_Area     = byte(10)
	VMach_Parser_Service_Name    = "VMPars"
	VMach_Parser_Service_Area    = byte(10)
	VMach_Model_Name             = "vendfleetmachine"
)

type Links struct{}

func (this *Links) Collector(linkid string) (string, byte) {
	return Collector_Service_Name, Collector_Service_Area
}

func (this *Links) Parser(linkid string) (string, byte) {
	switch linkid {
	case VendPerMachine_Links_ID:
		return VMach_Parser_Service_Name, VMach_Parser_Service_Area
	}
	return Vend_Parser_Service_Name, Vend_Parser_Service_Area
}

func (this *Links) Cache(linkid string) (string, byte) {
	switch linkid {
	case VendPerMachine_Links_ID:
		return VMach_Cache_Service_Name, VMach_Cache_Service_Area
	}
	return Vend_Cache_Service_Name, Vend_Cache_Service_Area
}

func (this *Links) Persist(linkid string) (string, byte) {
	switch linkid {
	case VendPerMachine_Links_ID:
		return VMach_Cache_Service_Name, VMach_Cache_Service_Area
	}
	return Vend_Persist_Service_Name, Vend_Persist_Service_Area
}

func (this *Links) Model(linkid string) string {
	switch linkid {
	case VendPerMachine_Links_ID:
		return VMach_Model_Name
	}
	return Vend_Model_Name
}
