package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	t.Run("Must marshal value according to schema", func(t *testing.T) {
		type nestedStruct struct {
			// All trailing underscores are removed to match schema
			FooBar string `json:"_fooBar"`
			Baz    struct {
				Doo float64 `json:"_doo_"`
			} `json:"baz"`
			IncludedBasedOnFieldName string
			IgnoredJSON              string `json:"-"`
			IgnoredNotInSchema       string `json:"notInSchema"`
			BoolAsText               bool   `json:"boolAsText"`
			BoolAsNumber             bool   `json:"boolAsNumber"`
		}

		schema := Definition{
			"foobar":                   "text",
			"baz_doo":                  "number",
			"ignoredjson":              "text",
			"includedbasedonfieldname": "text",
			"boolastext":               "text",
			"boolasnumber":             "number",
		}

		item := nestedStruct{
			FooBar: "hello world",
			Baz: struct {
				Doo float64 `json:"_doo_"`
			}{Doo: 123.0},
			IncludedBasedOnFieldName: "a value",
			IgnoredJSON:              "ignored",
			IgnoredNotInSchema:       "ignored",
			BoolAsText:               true,
			BoolAsNumber:             true,
		}

		b, err := Marshal(item, schema)
		require.NoError(t, err)

		var out Map
		err = json.Unmarshal(b, &out)
		require.NoError(t, err)

		require.EqualValues(t, Map{
			"baz_doo":                  123.0,
			"foobar":                   "hello world",
			"includedbasedonfieldname": "a value",
			"boolastext":               "true",
			"boolasnumber":             1.0,
		}, out)
	})

	t.Run("Marshal slice", func(t *testing.T) {
		type model struct {
			Foo string
			Bar float64
			Baz bool
		}

		schema := Definition{
			"foo": "text",
			"bar": "number",
			"baz": "number",
		}

		input := []model{
			{Foo: "bar", Bar: 1111111, Baz: true},
			{Foo: "val", Bar: 999.123, Baz: false},
		}

		data, err := Marshal(input, schema)
		require.NoError(t, err)

		var results []model
		err = Unmarshal(data, &results)
		require.NoError(t, err)

		require.EqualValues(t, input, results)
	})
}
