package appsearch

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocumentAPI(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()

	c, err := Open(os.Getenv("APPSEARCH"))
	require.NoError(t, err)

	t.Run("Must insert document", func(t *testing.T) {
		t.Parallel()
		engine := createRandomEngine(c)
		defer deleteEngine(c, engine)

		res, err := c.UpdateDocuments(ctx, engine.Name, []m{
			{"foo": "yes-id", "id": "has-id"},
			{"none": "`none` field is forbidden by API spec"},
		})
		require.NoError(t, err)
		require.EqualValues(t, []UpdateResponse{
			{ID: "has-id", Errors: []string{}},
			{ID: "", Errors: []string{"Invalid field name: none"}},
		}, res)
	})

	t.Run("Must patch document", func(t *testing.T) {
		t.Parallel()
		engine := createRandomEngine(c)
		defer deleteEngine(c, engine)

		res, err := c.UpdateDocuments(ctx, engine.Name, []m{
			{"id": "a-document", "foo": "bar"},
		})
		require.NoError(t, err)
		require.EqualValues(t, []UpdateResponse{
			{ID: "a-document", Errors: []string{}},
		}, res)

		res, err = c.PatchDocuments(ctx, engine.Name, []m{
			{"id": "a-document", "foo": "updated"},
		})
		require.NoError(t, err)
		require.EqualValues(t, []UpdateResponse{
			{ID: "a-document", Errors: []string{}},
		}, res)
	})

	t.Run("Must remove document", func(t *testing.T) {
		t.Parallel()
		engine := createRandomEngine(c)
		defer deleteEngine(c, engine)

		res, err := c.UpdateDocuments(ctx, engine.Name, []m{
			{"id": "a-document", "foo": "bar"},
		})
		require.NoError(t, err)
		require.EqualValues(t, []UpdateResponse{
			{ID: "a-document", Errors: []string{}},
		}, res)

		// Also accepts: []m{"id": "a-document"}
		deleteRes, err := c.RemoveDocuments(ctx, engine.Name, []string{"a-document"})
		require.NoError(t, err)
		require.EqualValues(t, []DeleteResponse{
			{ID: "a-document", Deleted: true},
		}, deleteRes)
	})

	t.Run("Must find document", func(t *testing.T) {
		t.Parallel()
		engine := createRandomEngine(c)
		defer deleteEngine(c, engine)

		res, err := c.UpdateDocuments(ctx, engine.Name, []m{
			{"id": "national-parks", "title": "Amazing title"},
		})
		require.NoError(t, err)
		require.EqualValues(t, []UpdateResponse{
			{ID: "national-parks", Errors: []string{}},
		}, res)

		response, err := c.SearchDocuments(ctx, engine.Name, Query{Query: "amazing"})
		require.NoError(t, err)

		require.Len(t, response.Results, 1)
	})
}
