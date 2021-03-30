package appsearch

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaAPI(t *testing.T) {
	t.Parallel()
	c, err := Open(os.Getenv("APPSEARCH"))
	require.NoError(t, err)

	ctx := context.TODO()

	t.Run("ListSchema", func(t *testing.T) {
		t.Parallel()
		engine := createRandomEngine(c)
		defer deleteEngine(c, engine)

		schema, err := c.ListSchema(ctx, engine.Name)
		require.NoError(t, err)
		require.NotEmpty(t, schema)
		require.EqualValues(t, "text", schema["id"])
	})

	t.Run("UpdateSchema", func(t *testing.T) {
		t.Parallel()
		engine := createRandomEngine(c)
		defer deleteEngine(c, engine)

		schema := SchemaDefinition{
			"id":  "text",
			"foo": "text",
		}

		err := c.UpdateSchema(ctx, engine.Name, schema)
		require.NoError(t, err)

		def, err := c.ListSchema(ctx, engine.Name)
		require.NoError(t, err)
		require.NotEmpty(t, def)
		require.EqualValues(t, schema, def)
	})
}
