package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// UnmarshalSchema is used to implement custom unpacking
type UnmarshalSchema interface {
	UnmarshalSchema(normalizedMap Map) error
}

// Return normalizedMap back to original form based on type-lookups on model interface
func Denormalize(normalizedMap Map, model interface{}) (denormalizedMap Map, err error) {
	nestedMap := objectify(normalizedMap, "_")

	modelType := getType(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	fieldIndex, tagIndex := buildIndex(modelType)
	return denormalize(nestedMap, fieldIndex, tagIndex)
}

func denormalize(nestedMap Map, fieldIndex map[string]reflect.StructField, tagIndex map[string]string) (denormalizedMap Map, err error) {
	denormalizedMap = make(Map)

	for field, value := range nestedMap {
		innerField, hasInnerField := fieldIndex[field]
		jsonTag := tagIndex[field]

		// Store values as-is if we can't derive original key
		if jsonTag == "" {
			jsonTag = field
		}

		if innerMap, isInnerMap := value.(Map); isInnerMap && hasInnerField {
			// Handle deep struct
			if innerField.Type.Kind() == reflect.Struct {
				fieldIndex, tagIndex := buildIndex(innerField.Type)
				denormalizedMap[jsonTag], err = denormalize(innerMap, fieldIndex, tagIndex)
				if err != nil {
					return
				}
				continue
			}
			// Handle unpacking of { raw } values
			rawValue, ok := innerMap["raw"]
			if !ok {
				panic(ErrRawValue)
			}
			value = rawValue
		}

		switch innerField.Type.Kind() {
		case reflect.Bool:
			denormalizedMap[jsonTag] = decodeBool(value)
		default:
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

	fieldIndex, tagIndex := buildIndex(valueType)

	for i, result := range results {
		newResult := reflect.New(valueType).Interface()

		nestedMap := objectify(result, "_")
		denormalizedMap, err := denormalize(nestedMap, fieldIndex, tagIndex)
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

func buildIndex(modelType reflect.Type) (fieldIndex map[string]reflect.StructField, tagIndex map[string]string) {
	jsonTagToField := mapFieldsByJSONTag(modelType)
	normalizedKeyToField := mapNormalizedFields(jsonTagToField)
	normalizedFieldToJSONTag := mapNormalizedToJSON(normalizedKeyToField)

	return normalizedKeyToField, normalizedFieldToJSONTag
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

func decodeBool(value interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case "true":
			return true
		case "false":
			return false
		default:
			panic(errCannotDecodeValueToBool(value))
		}
	case int:
		return decodeIntAsBool(value)
	case int32:
		return decodeIntAsBool(int(value))
	case int64:
		return decodeIntAsBool(int(value))
	case float32:
		return decodeFloatAsBool(value)
	case float64:
		return decodeFloatAsBool(float32(value))
	default:
		panic(errCannotDecodeValueToBool(value))
	}
}

func decodeFloatAsBool(value float32) bool {
	switch value {
	case 1:
		return true
	case 0:
		return false
	default:
		panic(errCannotDecodeValueToBool(value))
	}
}

func decodeIntAsBool(value int) bool {
	switch value {
	case 1:
		return true
	case 0:
		return false
	default:
		panic(errCannotDecodeValueToBool(value))
	}
}

func errCannotDecodeValueToBool(value interface{}) error {
	return fmt.Errorf("cannot decode %v to bool", value)
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
