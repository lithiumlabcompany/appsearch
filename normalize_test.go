package appsearch

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNormalize(t *testing.T) {

	t.Run("Normalize", func(t *testing.T) {
		t.Run("Should include only keys represented in schema", func(t *testing.T) {
			normalized, err := Normalize(bson.M{
				"_id":                "hello",
				"thisStuff":          1,
				"otherStuff":         2,
				"_somethingIgnored":  3,
				"justIgnored":        4,
				"_somethingIncluded": 5,
			}, SchemaDefinition{
				"id":                "text",
				"thisstuff":         "text",
				"otherstuff":        "text",
				"somethingincluded": "text",
			})
			require.NoError(t, err)

			require.EqualValues(t, bson.M{
				"id":                "hello",
				"thisstuff":         1,
				"otherstuff":        2,
				"somethingincluded": 5,
			}, normalized)
		})

		t.Run("Should include deep keys as represented in schema", func(t *testing.T) {
			normalized, err := Normalize(bson.M{
				"something": bson.M{
					"quiteNested": bson.M{
						"what": 1,
					},
				},
				"foo_ignored": 2,
				"barIgnored":  3,
				"nowhere":     4,
			}, SchemaDefinition{
				"something_quitenested_what": "text",
				"not_here":                   "text",
				"nowhere":                    "text",
			})
			require.NoError(t, err)

			require.EqualValues(t, bson.M{
				"something_quitenested_what": 1,
				"nowhere":                    4,
			}, normalized)
		})

		t.Run("Should include localized fields normalized to base key", func(t *testing.T) {
			normalized, err := Normalize(bson.M{
				"imHere": 1,
				"helloLocalized": bson.M{
					"deeply": bson.M{
						"included":  []string{"include-me"},
						"imAString": "this-is-string",
					},
					"stringSlice":    []string{"hello", "world"},
					"interfaceSlice": []interface{}{"foo", "bar"},
					"complexSlice":   []interface{}{nil, 1, 2, 3, "what?"},
				},
				"c": 3,
				"d": 4,
			}, SchemaDefinition{
				"imhere":         "text",
				"hellolocalized": "text",
			})
			require.NoError(t, err)

			localized := normalized["hellolocalized"].([]string)
			sort.Strings(localized)
			normalized["hellolocalized"] = localized

			require.EqualValues(t, bson.M{
				"imhere": 1,
				"hellolocalized": []string{
					"bar", "foo", "hello",
					"include-me", "this-is-string",
					"what?", "world"},
			}, normalized)
		})

		t.Run("Should nullify nil values to type", func(t *testing.T) {
			normalized, err := Normalize(bson.M{
				"textual":        nil,
				"numeric":        nil,
				"timestamp":      nil,
				"textualPresent": "hello",
			}, SchemaDefinition{
				"textual":           "text",
				"numeric":           "number",
				"timestamp":         "date",
				"textualpresent":    "text",
				"textualnotpresent": "text",
			})
			require.NoError(t, err)

			require.EqualValues(t, bson.M{
				"textual":        "",
				"numeric":        0,
				"timestamp":      time.Time{},
				"textualpresent": "hello",
			}, normalized)
		})
	})
}
