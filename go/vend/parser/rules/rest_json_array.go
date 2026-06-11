/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/

package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
)

// RestArrayToMap is a generic version of RestGpuParse that maps a JSON array
// into a protobuf map field, using any string field as the map key.
// Unlike RestGpuParse, it does not validate that the key is a PCI Bus ID.
//
// Parameters:
//   - "array_path": dot-path to the array in the JSON (e.g., "items")
//   - "key_field": JSON field name for the map key (e.g., "machineId")
//   - "mapping": comma-separated "jsonField:propertyName" pairs
type RestArrayToMap struct{}

func (this *RestArrayToMap) Name() string {
	return "RestArrayToMap"
}

func (this *RestArrayToMap) ParamNames() []string {
	return []string{"array_path", "key_field", "mapping"}
}

func (this *RestArrayToMap) Parse(resources ifs.IResources, workSpace map[string]interface{},
	params map[string]*l8tpollaris.L8PParameter, any interface{}, pollWhat string) error {

	input := workSpace["input"]
	if input == nil {
		return nil
	}

	var jsonStr string
	if cmap, ok := input.(*l8tpollaris.CMap); ok {
		jsonBytes, exists := cmap.Data["json"]
		if !exists || len(jsonBytes) == 0 {
			return errors.New("RestArrayToMap: CMap has no 'json' key")
		}
		dec := object.NewDecode(jsonBytes, 0, resources.Registry())
		val, err := dec.Get()
		if err != nil {
			return errors.New("RestArrayToMap: failed to decode json: " + err.Error())
		}
		jsonStr, _ = val.(string)
	} else if s, ok := input.(string); ok {
		jsonStr = s
	} else {
		return errors.New("RestArrayToMap: unsupported input type: " + fmt.Sprintf("%T", input))
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return errors.New("RestArrayToMap: failed to parse JSON: " + err.Error())
	}

	arrayPathParam := params["array_path"]
	mappingParam := params["mapping"]
	keyFieldParam := params["key_field"]
	if arrayPathParam == nil || mappingParam == nil || keyFieldParam == nil {
		return errors.New("RestArrayToMap: missing required parameter")
	}

	arrValue := getNestedValue(jsonData, arrayPathParam.Value)
	if arrValue == nil {
		return nil
	}
	arr, ok := arrValue.([]interface{})
	if !ok {
		return errors.New("RestArrayToMap: value at '" + arrayPathParam.Value + "' is not an array")
	}

	keyField := keyFieldParam.Value

	type fMapping struct {
		jsonField    string
		propertyName string
	}
	mappings := make([]fMapping, 0)
	for _, entry := range strings.Split(mappingParam.Value, ",") {
		parts := strings.SplitN(strings.TrimSpace(entry), ":", 2)
		if len(parts) != 2 {
			continue
		}
		mappings = append(mappings, fMapping{
			jsonField:    strings.TrimSpace(parts[0]),
			propertyName: strings.TrimSpace(parts[1]),
		})
	}

	propertyId := ""
	if pid, ok := workSpace["propertyid"]; ok {
		propertyId, _ = pid.(string)
	}

	fmt.Printf("[REST-ARRAY-TO-MAP] arrayPath=%s keyField=%s items=%d propertyId=%s\n",
		arrayPathParam.Value, keyField, len(arr), propertyId)

	for _, item := range arr {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		mapKey, ok := itemMap[keyField].(string)
		if !ok || mapKey == "" {
			continue
		}
		mapKey = strings.TrimSpace(mapKey)

		// Set the key field itself as a property
		setMapProperty(resources, propertyId, mapKey, keyField, mapKey, any)

		for _, m := range mappings {
			val := getNestedValue(itemMap, m.jsonField)
			if val == nil {
				continue
			}
			setMapProperty(resources, propertyId, mapKey, m.propertyName, val, any)
		}
	}

	fmt.Printf("[REST-ARRAY-TO-MAP-DONE] processed %d items\n", len(arr))
	return nil
}

func setMapProperty(resources ifs.IResources, propertyId, mapKey, fieldName string, val interface{}, any interface{}) {
	fullId := fmt.Sprintf("%s<{24}%s>.%s", propertyId, mapKey, fieldName)
	instance, err := properties.PropertyOf(fullId, resources)
	if err != nil || instance == nil {
		return
	}
	coerced := coerceValue(resources, val, instance, nil)
	if coerced != nil {
		instance.Set(any, coerced)
	}
}

// coerceValue converts JSON values to match property types
func coerceValue(resources ifs.IResources, value interface{}, instance *properties.Property, workSpace map[string]interface{}) interface{} {
	node := instance.Node()
	if node == nil {
		return value
	}
	typeName := node.TypeName

	switch v := value.(type) {
	case float64:
		switch typeName {
		case "int32":
			return int32(v)
		case "int64":
			return int64(v)
		case "uint32":
			return uint32(v)
		case "string":
			return fmt.Sprintf("%v", v)
		}
	case string:
		switch typeName {
		case "int32":
			return value
		case "float64":
			return value
		}
	}
	return value
}

func getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = data
	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = m[part]
		if !ok {
			return nil
		}
	}
	return current
}
