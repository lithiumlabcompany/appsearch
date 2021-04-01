package schema

import (
	"encoding/json"
	"reflect"
)

// Normalize and marshal input according to schema Definition
func Marshal(input interface{}, schema Definition) (data []byte, err error) {
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Slice {
		normalized := make([]Map, value.Len())

		for i := range normalized {
			serializable := value.Index(i).Interface()

			normalized[i], err = ToMap(serializable, schema)
		}

		input = normalized
	} else {
		input, err = ToMap(input, schema)
	}
	if err == nil {
		return json.Marshal(input)
	}
	return
}

// Convert structure to normalized map according to schema Definition
func ToMap(input interface{}, schema Definition) (normalizedMap Map, err error) {
	nestedMap, err := mapFromJSON(input)
	if err != nil {
		return
	}
	return Normalize(nestedMap, schema)
}

func mapFromJSON(i interface{}) (m Map, err error) {
	err = unmarshalInto(i, &m)
	return
}
