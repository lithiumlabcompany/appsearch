package appsearch

import (
	"strings"
	"time"

	"github.com/lithiumlabcompany/appsearch/internal/pkg/flatten"
)

// Normalize nested map into flat map as defined in schema
// Keys are stripped of trailing underscores, lowercased and flattened with underscore (_) separator
func Normalize(raw m, schema SchemaDefinition) (normalizedFlatMap m, err error) {
	flatMap, err := flatten.Flatten(raw, flatten.UnderscoreStyle)
	if err != nil {
		return
	}

	normalizedFlatMap = make(m)
	for rawKey, flatValue := range flatMap {
		normKey := strings.ToLower(strings.Trim(rawKey, "_"))
		baseKey := strings.Split(normKey, "_")[0]

		// Normalize to key (store value as is)
		if schemaType, inSchema := schema[normKey]; inSchema {
			// Make sure nil values in schema are nullified
			if flatValue == nil {
				flatValue = nullifyType(schemaType)
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
				for _, item := range values {
					stringSlice = append(stringSlice, item)
				}
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

var defaultValues = map[SchemaType]interface{}{
	"text":        "",
	"date":        time.Time{},
	"number":      0,
	"geolocation": "0.0,0.0",
}

func nullifyType(schemaType SchemaType) interface{} {
	return defaultValues[schemaType]
}
