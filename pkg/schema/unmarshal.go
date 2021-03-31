package schema

import (
	"encoding/json"
	"reflect"
	"strings"
)

// UnmarshalSchema is used to implement custom unpacking
type UnmarshalSchema interface {
	UnmarshalSchema(normalizedMap Map) error
}

// Accepts normalized Map as input and tries to unpack nested map according to struct tags
// `json:"..."` tags are used to infer original data model comparing fields via normalized schema
// By that extent its forbidden to use underscores ("_") in tags on original data model
func Unmarshal(normalizedMap Map, output interface{}) (err error) {
	if unmarshal, ok := output.(UnmarshalSchema); ok {
		return unmarshal.UnmarshalSchema(normalizedMap)
	}

	denormalizedMap, err := Denormalize(normalizedMap, output)
	if err == nil {
		err = unmarshalInto(denormalizedMap, output)
	}

	return err
}

// Return normalizedMap back to original form based on type-lookups on model interface
func Denormalize(normalizedMap Map, model interface{}) (denormalizedMap Map, err error) {
	normalizedNestedMap := objectify(normalizedMap, "_")
	denormalizedMap = make(Map)

	t := getType(model)
	jsonTagToField := mapFieldsByJSONTag(t)
	normalizedKeyToField := mapNormalizedFields(jsonTagToField)
	normalizedFieldToJSONTag := mapNormalizedToJSON(normalizedKeyToField)

	for normalizedField, value := range normalizedNestedMap {
		innerField, hasInnerField := normalizedKeyToField[normalizedField]
		jsonTag, _ := normalizedFieldToJSONTag[normalizedField]

		// Store values as-is if we can't derive original key
		if jsonTag == "" {
			jsonTag = normalizedField
		}

		if innerMap, isInnerMap := value.(Map); isInnerMap && hasInnerField {
			denormalizedMap[jsonTag], err = Denormalize(innerMap, innerField.Type)
			if err != nil {
				return
			}
		} else {
			denormalizedMap[jsonTag] = value
		}
	}

	return
}

func mapNormalizedToJSON(byNormalized map[string]reflect.StructField) map[string]string {
	normalizedToJSON := make(map[string]string)

	for normalizedKey, field := range byNormalized {
		normalizedToJSON[normalizedKey] = getJSONTagOrFieldName(field)
	}

	return normalizedToJSON
}

func mapFieldsByJSONTag(t reflect.Type) map[string]reflect.StructField {
	index := make(map[string]reflect.StructField)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		jsonField := getJSONTagOrFieldName(field)

		if jsonField == "-" {
			continue
		}

		index[jsonField] = field
	}

	return index
}

func getJSONTagOrFieldName(field reflect.StructField) string {
	// Get the field tag value as described in `json:"..."` tag
	tag := field.Tag.Get("json")

	jsonField := strings.Split(tag, ",")[0]

	if jsonField == "" {
		jsonField = field.Name
	}

	return jsonField
}

func mapNormalizedFields(mapByField map[string]reflect.StructField) map[string]reflect.StructField {
	normalizedIndex := make(map[string]reflect.StructField)

	for k, v := range mapByField {
		normalizedIndex[NormalizeField(k)] = v
	}

	return normalizedIndex
}

func getType(model interface{}) (t reflect.Type) {
	switch v := model.(type) {
	case reflect.Type:
		t = v
	default:
		t = reflect.TypeOf(model)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}
	return t
}

func unmarshalInto(input, output interface{}) error {
	b, err := json.Marshal(input)
	if err == nil {
		err = json.Unmarshal(b, output)
	}
	return err
}
