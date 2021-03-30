package appsearch

import (
	"context"
	"net/http"
)

// Patch a list of documents. Every document must contain "id".
// Every document is patched separately.
// Documents without ID will be rejected.
// Non-existing documents will be rejected.
func (c *client) PatchDocuments(ctx context.Context, engineName string, documents interface{}) (res []UpdateResponse, err error) {
	err = c.Call(ctx, documents, &res, http.MethodPatch, "engines/%s/documents", engineName)

	return res, err
}

// Update (replace) a list of documents
// Every document is created (or replaced) separately.
// Documents without ID will have auto-generated ID's.
// Non-existing documents will be automatically created.
func (c *client) UpdateDocuments(ctx context.Context, engineName string, documents interface{}) (res []UpdateResponse, err error) {
	err = c.Call(ctx, documents, &res, http.MethodPost, "engines/%s/documents", engineName)

	return res, err
}

// Remove a list of documents specified as string ID's or documents with "id" field
// Every document is deleted separately.
func (c *client) RemoveDocuments(ctx context.Context, engineName string, documents interface{}) (res []DeleteResponse, err error) {
	err = c.Call(ctx, documents, &res, http.MethodDelete, "engines/%s/documents", engineName)

	return res, err
}

// Search documents by query
// Refer to https://www.elastic.co/guide/en/app-search/current/search.html#search-api-request-body
func (c *client) SearchDocuments(ctx context.Context, engineName string, query Query) (response SearchResponse, err error) {
	err = c.Call(ctx, query, &response, http.MethodPost, "engines/%s/search", engineName)

	return response, err
}
