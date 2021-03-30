package appsearch

import (
	"context"
)

type DocumentAPI interface {
	// Patch a list of documents. Every document must contain "id".
	// Every document is processed separately.
	// Documents without ID will be rejected.
	// Non-existent documents will be rejected.
	PatchDocuments(ctx context.Context, engineName string, documents interface{}) (res []UpdateResponse, err error)
	// Update (replace) a list of documents
	// Every document is processed separately.
	// Documents without ID will have auto-generated ID's.
	// Non-existent documents will be automatically created.
	UpdateDocuments(ctx context.Context, engineName string, documents interface{}) (res []UpdateResponse, err error)
	// Remove a list of documents specified as []string of ID's or []interface{} of documents with "id" field
	// Every document is processed separately.
	RemoveDocuments(ctx context.Context, engineName string, documentsOrIDs interface{}) (res []DeleteResponse, err error)
	// Search documents by query
	SearchDocuments(ctx context.Context, engineName string, query Query) (response SearchResponse, err error)
}

type EngineAPI interface {
	// List an engine by name
	ListEngine(ctx context.Context, engineName string) (data EngineDescription, err error)
	// List engines with pagination
	ListEngines(ctx context.Context, page Page) (data EngineResponse, err error)
	// List all available engines
	ListAllEngines(ctx context.Context) (data []EngineDescription, err error)

	// Create engine with name
	CreateEngine(ctx context.Context, request CreateEngineRequest) (EngineDescription, error)
	// Delete engine with name
	DeleteEngine(ctx context.Context, engineName string) (err error)

	// Create engine if doesn't exist.
	// Optionally update a schema even if engine exists.
	EnsureEngine(ctx context.Context, request CreateEngineRequest, schema ...SchemaDefinition) (err error)
}

type SchemaAPI interface {
	// List a schema definition by engineName
	ListSchema(ctx context.Context, engineName string) (data SchemaDefinition, err error)

	// Update schema by engineName (create or change fields).
	// Fields cannot be deleted.
	UpdateSchema(ctx context.Context, engineName string, def SchemaDefinition) (err error)
}

// APIClient interface
type APIClient interface {
	// Engine API
	EngineAPI
	// Schema API
	SchemaAPI
	// Document API
	DocumentAPI
}
