package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// UnpackSchema is used to implement custom unpacking
type UnpackSchema interface {
	UnpackSchema(normalizedMap Map) error
}

// Unpack search results into a slice
func UnpackSlice(results []Map, output interface{}) (err error) {
	sliceType := getType(output)
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
	}
	valueType := sliceType.Elem()
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	outputSlice := reflect.MakeSlice(sliceType, len(results), len(results))

	fieldIndex, tagIndex, err := buildIndex(valueType)
	if err != nil {
		return
	}

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

		outputSlice.Index(i).Set(reflect.ValueOf(newResult).Elem())
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
func Unpack(normalizedMap Map, output interface{}) (err error) {
	if unmarshal, ok := output.(UnpackSchema); ok {
		return unmarshal.UnpackSchema(normalizedMap)
	}

	outputType := reflect.TypeOf(output).Elem()
	fieldIndex, tagIndex, err := buildIndex(outputType)
	if err != nil {
		return err
	}

	nestedMap := objectify(normalizedMap, "_")
	denormalizedMap, err := denormalize(nestedMap, fieldIndex, tagIndex)
	if err == nil {
		err = unmarshalInto(denormalizedMap, output)
	}

	return err
}

// Unmarshal raw JSON
func Unmarshal(data []byte, output interface{}) (err error) {
	var raw interface{}

	err = json.Unmarshal(data, &raw)
	if err != nil {
		return
	}

	switch raw := raw.(type) {
	case []interface{}:
		return unpackInterfaceSlice(raw, output)
	case map[string]interface{}:
		return Unpack(raw, output)
	default:
		panic(fmt.Errorf("unmarshal %v", raw))
	}
}

// Return normalizedMap back to original form based on type-lookups on model interface
func Denormalize(normalizedMap Map, model interface{}) (denormalizedMap Map, err error) {
	nestedMap := objectify(normalizedMap, "_")

	fieldIndex, tagIndex, err := buildIndex(reflect.TypeOf(model))
	if err != nil {
		return nil, err
	}

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
				fieldIndex, tagIndex, err := buildIndex(innerField.Type)
				if err != nil {
					return nil, err
				}
				denormalizedMap[jsonTag], err = denormalize(innerMap, fieldIndex, tagIndex)
				if err != nil {
					return nil, err
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

		if hasInnerField && innerField.Type.Kind() == reflect.Bool {
			value = decodeBool(value)
		}

		denormalizedMap[jsonTag] = value
	}

	return
}

func unpackInterfaceSlice(raw []interface{}, output interface{}) error {
	mapSlice := make([]Map, len(raw))
	var sanity bool
	for i, raw := range raw {
		mapSlice[i], sanity = raw.(Map)
		if !sanity {
			return fmt.Errorf("cannot unmarshal slice of %T to %T", raw, output)
		}
	}
	return UnpackSlice(mapSlice, output)
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

func mapNormalizedFields(mapByField map[string]reflect.StructField) map[string]reflect.StructField {
	normalizedIndex := make(map[string]reflect.StructField)

	for k, v := range mapByField {
		normalizedIndex[NormalizeField(k)] = v
	}

	return normalizedIndex
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

func buildIndex(modelType reflect.Type) (fieldIndex map[string]reflect.StructField, tagIndex map[string]string, err error) {
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	switch modelType.Kind() {
	case reflect.Slice:
		return nil, nil, ErrCannotUnpackSlice
	case reflect.Map:
		return nil, nil, ErrCannotInferFromMap
	}

	jsonTagToField := mapFieldsByJSONTag(modelType)
	normalizedKeyToField := mapNormalizedFields(jsonTagToField)
	normalizedFieldToJSONTag := mapNormalizedToJSON(normalizedKeyToField)

	return normalizedKeyToField, normalizedFieldToJSONTag, err
}

func decodeBool(value interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case "1":
			return true
		case "0":
			return false
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
