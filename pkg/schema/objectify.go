package schema

import (
	"strings"
)

func objectify(flatMap Map, sep string) Map {
	nestedMap := make(Map)

	for key, value := range flatMap {
		deepSet(nestedMap, key, value, sep)
	}

	return nestedMap
}

func deepSet(nestedMap Map, key string, value interface{}, sep string) {
	parts := strings.Split(key, sep)

	for _, keyPart := range parts[0 : len(parts)-1] {
		innerMap, ok := nestedMap[keyPart].(Map)
		if !ok {
			nestedMap[keyPart] = make(Map)
			innerMap = nestedMap[keyPart].(Map)
		}
		nestedMap = innerMap
	}

	nestedMap[parts[len(parts)-1]] = value
}
