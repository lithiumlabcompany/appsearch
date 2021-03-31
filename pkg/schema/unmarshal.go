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

// Return normalizedMap back to original form based on type-lookups on model interface
func Denormalize(normalizedMap Map, model interface{}) (denormalizedMap Map, err error) {
	normalizedNestedMap := objectify(normalizedMap, "_")
	denormalizedMap = make(Map)

	t := getType(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	jsonTagToField := mapFieldsByJSONTag(t)
	normalizedKeyToField := mapNormalizedFields(jsonTagToField)
	normalizedFieldToJSONTag := mapNormalizedToJSON(normalizedKeyToField)

	for normalizedField, value := range normalizedNestedMap {
		innerField, hasInnerField := normalizedKeyToField[normalizedField]
		jsonTag := normalizedFieldToJSONTag[normalizedField]

		// Store values as-is if we can't derive original key
		if jsonTag == "" {
			jsonTag = normalizedField
		}

		if innerMap, isInnerMap := value.(Map); isInnerMap && hasInnerField {
			// Handle deep struct
			if innerField.Type.Kind() == reflect.Struct {
				denormalizedMap[jsonTag], err = Denormalize(innerMap, innerField.Type)
				if err != nil {
					return
				}
			} else {
				// Handle unpacking of { raw } values
				rawValue, ok := innerMap["raw"]
				if !ok {
					panic(ErrRawValue)
				}
				denormalizedMap[jsonTag] = rawValue
			}
		} else {
			denormalizedMap[jsonTag] = value
		}
	}

	return
}

// Unmarshal search results into a slice
func UnmarshalResults(results []Map, output interface{}) (err error) {
	sliceType := getType(output)
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
	}
	valueType := sliceType.Elem()
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	outputSlice := reflect.MakeSlice(sliceType, len(results), len(results))

	for i, result := range results {
		newResult := reflect.New(valueType).Interface()

		denormalizedMap, err := Denormalize(result, valueType)
		if err == nil {
			err = unmarshalInto(denormalizedMap, &newResult)
		}
		if err != nil {
			return err
		}

		if valueType.Kind() != reflect.Ptr {
			newResult = reflect.ValueOf(newResult).Elem().Interface()
		}
		outputSlice.Index(i).Set(reflect.ValueOf(newResult))
	}

	outputValue := reflect.ValueOf(output)
	if outputValue.Kind() == reflect.Ptr && !outputValue.IsZero() {
		outputValue = outputValue.Elem()
	}

	outputValue.Set(outputSlice)
	return nil
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
		return v
	default:
		return reflect.TypeOf(model)
	}
}

func unmarshalInto(input, output interface{}) error {
	b, err := json.Marshal(input)
	if err == nil {
		err = json.Unmarshal(b, output)
	}
	return err
}