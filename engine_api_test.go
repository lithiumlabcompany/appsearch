package appsearch

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEngineAPI(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()

	c, err := Open(os.Getenv("APPSEARCH"))
	require.NoError(t, err)

	t.Run("ListEngine", func(t *testing.T) {
		t.Parallel()
		data, err := c.ListEngines(ctx, Page{Page: 1, Size: 1})
		require.NoError(t, err)
		require.NotEmpty(t, data)
	})

	t.Run("ListAllEngines", func(t *testing.T) {
		t.Parallel()
		engines, err := c.ListAllEngines(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, engines)
	})

	t.Run("DeleteEngine", func(t *testing.T) {
		t.Parallel()
		engine := createRandomEngine(c)
		require.NoError(t, err)
		err = c.DeleteEngine(ctx, engine.Name)
		require.NoError(t, err)
	})

	t.Run("CreateEngine", func(t *testing.T) {
		t.Parallel()
		t.Run("Must create engine", func(t *testing.T) {
			t.Parallel()
			engineName := fmt.Sprintf("test-%d", rand.Uint64())
			engine, err := c.CreateEngine(ctx, CreateEngineRequest{
				Name: engineName,
			})
			defer deleteEngine(c, engineName)

			require.NoError(t, err)
			require.EqualValues(t, engineName, engine.Name)
		})

		t.Run("Must return ErrEngineAlreadyExists", func(t *testing.T) {
			t.Parallel()
			engine := createRandomEngine(c)
			defer deleteEngine(c, engine)

			_, err := c.CreateEngine(ctx, CreateEngineRequest{
				Name: engine.Name,
			})
			require.ErrorIs(t, err, ErrEngineAlreadyExists)
		})
	})

	t.Run("EnsureEngine", func(t *testing.T) {
		t.Parallel()
		t.Run("Must create engine if doesn't exist", func(t *testing.T) {
			t.Parallel()
			engineName := fmt.Sprintf("test-%d", rand.Uint64())
			err := c.EnsureEngine(ctx, CreateEngineRequest{
				Name: engineName,
			})
			defer deleteEngine(c, engineName)
			require.NoError(t, err)
		})

		t.Run("Must update schema", func(t *testing.T) {
			t.Parallel()
			engineName := fmt.Sprintf("test-%d", rand.Uint64())
			schema := SchemaDefinition{
				"id":   "text",
				"foo":  "text",
				"bar":  "number",
				"baze": "date",
			}
			err := c.EnsureEngine(ctx, CreateEngineRequest{
				Name: engineName,
			}, schema)
			defer deleteEngine(c, engineName)
			require.NoError(t, err)

			def, err := c.ListSchema(ctx, engineName)
			require.NoError(t, err)
			require.EqualValues(t, schema, def)
		})
	})
}
