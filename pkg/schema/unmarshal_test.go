package schema

import (
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

	t.Run("Unmarshal", func(t *testing.T) {
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
		err := Unmarshal(normalizedMap, &output)
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
}
