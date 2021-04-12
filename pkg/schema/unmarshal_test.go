package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	t.Run("Denormalize", func(t *testing.T) {
		normalizedMap := Map{
			"helloworld":   "a value",
			"nested_stuff": "another value",
		}
		type model struct {
			HelloWorld string `json:"_helloWorld_"`
			Nested     struct {
				Stuff string `json:"___stuff___"`
			} `json:"nEsTeD"`
		}
		denormalized, err := Denormalize(normalizedMap, &model{})
		require.NoError(t, err)
		require.EqualValues(t, Map{
			"_helloWorld_": "a value",
			"nEsTeD": Map{
				"___stuff___": "another value",
			},
		}, denormalized)
	})

	t.Run("Unpack", func(t *testing.T) {
		normalizedMap := Map{
			"foo":     "1",
			"bar_baz": "2",
		}
		type model struct {
			Foo string `json:"foo"`
			Bar struct {
				Baz string `json:"baz"`
			} `json:"Bar"`
		}
		var output model
		err := Unpack(normalizedMap, &output)
		require.NoError(t, err)
		require.EqualValues(t, model{
			Foo: "1",
			Bar: struct {
				Baz string `json:"baz"`
			}{
				Baz: "2",
			},
		}, output)
	})

	t.Run("UnpackResults", func(t *testing.T) {
		type model struct {
			Foo string `json:"__foo"`
			Bar struct {
				Baz string `json:"__stUff__"`
			} `json:"dEEp_"`
			BoolAsText   bool `json:"boolAsText"`
			BoolAsNumber bool `json:"boolAsNumber"`
		}
		results := []Map{
			{
				"foo":          "1",
				"deep_stuff":   "2",
				"boolastext":   "true",
				"boolasnumber": 1,
				"nonexistent":  "hello",
				"_meta":        Map{},
			},
			{
				"foo":          "3",
				"deep_stuff":   "4",
				"boolastext":   "true",
				"boolasnumber": 1,
				"nonexistent":  "hello",
			},
			{
				"deep_stuff": Map{
					"raw": "value",
				},
				"foo": Map{
					"raw": "value",
				},
				"boolastext": Map{
					"raw": "true",
				},
				"boolasnumber": Map{
					"raw": 1,
				},
				"nonexistent": "hello",
			},
		}
		expected := []model{
			{
				Foo: "1",
				Bar: struct {
					Baz string `json:"__stUff__"`
				}{
					Baz: "2",
				},
				BoolAsText:   true,
				BoolAsNumber: true,
			},
			{
				Foo: "3",
				Bar: struct {
					Baz string `json:"__stUff__"`
				}{
					Baz: "4",
				},
				BoolAsText:   true,
				BoolAsNumber: true,
			},
			{
				Foo: "value",
				Bar: struct {
					Baz string `json:"__stUff__"`
				}{
					Baz: "value",
				},
				BoolAsText:   true,
				BoolAsNumber: true,
			},
		}
		t.Run("As []model", func(t *testing.T) {
			var output []model
			err := UnpackSlice(results, &output)
			require.NoError(t, err)
			require.EqualValues(t, expected, output)
		})
		// t.Run("As *[]model", func(t *testing.T) {
		// 	var output *[]model
		// 	err := UnpackResults(results, &output)
		// 	require.NoError(t, err)
		// 	require.EqualValues(t, expected, output)
		// })
		// t.Run("As []*model", func(t *testing.T) {
		// 	var output []*model
		// 	err := UnpackResults(results, &output)
		// 	require.NoError(t, err)
		// 	require.EqualValues(t, expected, output)
		// })
		// t.Run("As *[]*model", func(t *testing.T) {
		// 	var output *[]*model
		// 	err := UnpackResults(results, &output)
		// 	require.NoError(t, err)
		// require.EqualValues(t, expected, output)
		// })
	})

	t.Run("Unmarshal", func(t *testing.T) {
		type model struct {
			Hello string
		}
		t.Run("JSON object", func(t *testing.T) {
			var result model
			err := Unmarshal([]byte(`{"hello": "world"}`), &result)
			require.NoError(t, err)
			require.EqualValues(t, model{"world"}, result)
		})
		t.Run("JSON array", func(t *testing.T) {
			var result []model
			err := Unmarshal([]byte(`[{"hello": "world"}]`), &result)
			require.NoError(t, err)
			require.EqualValues(t, []model{{"world"}}, result)
		})
		t.Run("denormalize", func(t *testing.T) {
			t.Run("decodeValue", func(t *testing.T) {
				t.Run("decodeNumber", func(t *testing.T) {
					t.Run("float64", func(t *testing.T) {
						v := "1.123"
						n, err := decodeNumber(v, reflect.Float64)
						require.NoError(t, err)
						require.EqualValues(t, n, float64(1.123))
					})
				})
			})
		})
	})
}
