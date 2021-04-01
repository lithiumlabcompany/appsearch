package schema

import (
	"encoding/json"
)

// Normalize input structure and marshal using json.Marshal
func Marshal(i interface{}, schema Definition) ([]byte, error) {
	normalizedMap, err := ToMap(i, schema)
	if err == nil {
		return json.Marshal(normalizedMap)
	}
	return nil, err
}

// Convert JSON-serializable type to normalized map according to Definition
func ToMap(i interface{}, schema Definition) (normalizedMap Map, err error) {
	nestedMap, err := mapFromJSON(i)
	if err != nil {
		return
	}
	return Normalize(nestedMap, schema)
}

func mapFromJSON(i interface{}) (m Map, err error) {
	err = unmarshalInto(i, &m)
	return
}
