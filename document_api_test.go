package appsearch

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

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

		time.Sleep(time.Second)

		response, err := c.SearchDocuments(ctx, engine.Name, Query{Query: "amazing"})
		require.NoError(t, err)

		require.Len(t, response.Results, 1)
	})

	t.Run("Must list documents", func(t *testing.T) {
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

		time.Sleep(time.Second)

		response, err := c.ListDocuments(ctx, engine.Name, Page{1, 1})
		require.NoError(t, err)

		require.Len(t, response.Results, 1)
	})

	t.Run("Must decode search facets", func(t *testing.T) {
		t.Parallel()

		engine := createRandomEngine(c)
		defer deleteEngine(c, engine)

		_, err := c.UpdateDocuments(ctx, engine.Name, []m{
			{"state": "Illinois", "obscurerandomnumber": rand.Intn(12345)},
			{"state": "Illinois", "obscurerandomnumber": rand.Intn(12345)},
			{"state": "Missouri", "obscurerandomnumber": rand.Intn(12345)},
			{"state": "Missouri", "obscurerandomnumber": rand.Intn(12345)},
			{"state": "Kansas", "obscurerandomnumber": rand.Intn(12345)},
			{"state": "Arkansas", "obscurerandomnumber": rand.Intn(12345)},
			{"state": "Utah", "obscurerandomnumber": rand.Intn(12345)},
			{"state": "Utah", "obscurerandomnumber": rand.Intn(12345)},
			{"state": "Utah", "obscurerandomnumber": rand.Intn(12345)},
		})
		require.NoError(t, err)

		time.Sleep(time.Second)

		query := Query{Query: "", Facets: SearchFacets{
			"state": []Facet{{
				Type: "value",
				Name: "top_five_states",
				Sort: Sorting{"count": "desc"},
				Size: 5,
			}},
		}}
		response, err := c.SearchDocuments(ctx, engine.Name, query)
		require.NoError(t, err)

		require.Len(t, response.Facets["state"], 1)
		require.Len(t, response.Facets["state"][0].Data, 5)

		for _, facet := range response.Facets["state"][0].Data {
			require.NotEmpty(t, facet.Value)
			require.NotEmpty(t, facet.Count)
		}
	})
}
