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

func createVendRestAttribute(propertyId, mapping string) *l8tpollaris.L8PAttribute {
	return createRestAttribute("vendmachine", propertyId, mapping)
}

func createRestAttribute(modelKey, propertyId, mapping string) *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = map[string]string{modelKey: propertyId}
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	rule := &l8tpollaris.L8PRule{}
	rule.Name = "RestJsonParse"
	rule.Params = make(map[string]*l8tpollaris.L8PParameter)
	rule.Params["mapping"] = &l8tpollaris.L8PParameter{Value: mapping}
	attr.Rules = append(attr.Rules, rule)
	return attr
}

func createVendArrayAttribute(propertyId, arrayPath, keyField, mapping string) *l8tpollaris.L8PAttribute {
	attr := &l8tpollaris.L8PAttribute{}
	attr.PropertyId = map[string]string{"vendmachine": propertyId}
	attr.Rules = make([]*l8tpollaris.L8PRule, 0)
	rule := &l8tpollaris.L8PRule{}
	rule.Name = "RestArrayToMap"
	rule.Params = make(map[string]*l8tpollaris.L8PParameter)
	rule.Params["array_path"] = &l8tpollaris.L8PParameter{Value: arrayPath}
	rule.Params["key_field"] = &l8tpollaris.L8PParameter{Value: keyField}
	rule.Params["mapping"] = &l8tpollaris.L8PParameter{Value: mapping}
	attr.Rules = append(attr.Rules, rule)
	return attr
}

func createVendPoll(name, endpoint string, cadence *l8tpollaris.L8PCadencePlan, always bool, propertyId, mapping string) *l8tpollaris.L8Poll {
	poll := &l8tpollaris.L8Poll{}
	poll.Name = name
	poll.What = "GET::" + endpoint + "::"
	poll.Protocol = l8tpollaris.L8PProtocol_L8PRESTAPI
	poll.Cadence = cadence
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Always = always
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createVendRestAttribute(propertyId, mapping))
	return poll
}

func createVendArrayPoll(name, endpoint string, cadence *l8tpollaris.L8PCadencePlan, always bool, propertyId, arrayPath, keyField, mapping string) *l8tpollaris.L8Poll {
	poll := &l8tpollaris.L8Poll{}
	poll.Name = name
	poll.What = "GET::" + endpoint + "::"
	poll.Protocol = l8tpollaris.L8PProtocol_L8PRESTAPI
	poll.Cadence = cadence
	poll.Timeout = DEFAULT_TIMEOUT
	poll.Always = always
	poll.Operation = l8tpollaris.L8C_Operation_L8C_Get
	poll.Attributes = make([]*l8tpollaris.L8PAttribute, 0)
	poll.Attributes = append(poll.Attributes, createVendArrayAttribute(propertyId, arrayPath, keyField, mapping))
	return poll
}
