package schema

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lithiumlabcompany/appsearch/internal/pkg/flatten"
)

// Normalize nested map into flat map as defined in schema.
// Keys are stripped of trailing underscores, lowercased and flattened with underscore (_) separator
func Normalize(raw Map, schema Definition) (normalizedFlatMap Map, err error) {
	flatMap, err := flatten.Flatten(raw, flatten.UnderscoreStyle)
	if err != nil {
		return
	}

	normalizedFlatMap = make(Map)
	for rawKey, flatValue := range flatMap {
		normKey := NormalizeField(rawKey)
		baseKey := strings.Split(normKey, "_")[0]

		// Normalize to key (store value as is)
		if schemaType, inSchema := schema[normKey]; inSchema {
			switch value := flatValue.(type) {
			case nil:
				// Make sure nil values in schema are nullified
				flatValue = nullifyType(schemaType)
			case bool:
				// Serialize boolean to string or number
				flatValue = encodeBool(value, schemaType)
			}

			normalizedFlatMap[normKey] = flatValue
			continue
		}

		// Normalize to base key (append all to string slice)
		if _, inSchema := schema[baseKey]; inSchema {
			stringSlice := make([]string, 0)
			if prevNormalized, ok := normalizedFlatMap[baseKey].([]string); ok {
				stringSlice = prevNormalized
			}

			// Append all inner values of base key to single array
			// E.g. meanings.az.value:[1,2], meanings.ru.value:[3] -> meanings:[1,2,3]
			switch values := flatValue.(type) {
			case []string:
				stringSlice = append(stringSlice, values...)
			case []interface{}:
				for _, item := range values {
					if str, ok := item.(string); ok {
						stringSlice = append(stringSlice, str)
					}
				}
			case string:
				stringSlice = append(stringSlice, values)
			}

			normalizedFlatMap[baseKey] = stringSlice
			continue
		}
	}

	return normalizedFlatMap, nil
}

func encodeBool(value bool, schemaType Type) interface{} {
	switch schemaType {
	case TypeText:
		return fmt.Sprintf("%v", value)
	case TypeNumber:
		if value {
			return 1
		} else {
			return 0
		}
	default:
		panic(fmt.Errorf("cannot encode %v to %v", value, schemaType))
	}
}

var underscoreRe = regexp.MustCompile(`_+`)

// Get normalized variant for field name
func NormalizeField(field string) (normalized string) {
	normalized = field
	// Replace any sequence of underscores by a single underscore
	normalized = underscoreRe.ReplaceAllString(normalized, "_")
	// Trim underscores
	normalized = strings.Trim(normalized, "_")
	// Lowercase
	normalized = strings.ToLower(normalized)
	return
}

var defaultValues = map[Type]interface{}{
	"text":        "",
	"date":        time.Time{},
	"number":      0,
	"geolocation": "0.0,0.0",
}

func nullifyType(schemaType Type) interface{} {
	return defaultValues[schemaType]
}
