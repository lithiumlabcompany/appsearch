package appsearch

import (
	"context"
	"net/http"

	"github.com/lithiumlabcompany/appsearch/pkg/schema"
)

// List a schema by engineName
func (c *client) ListSchema(ctx context.Context, engineName string) (data schema.Definition, err error) {
	err = c.Call(ctx, nil, &data, http.MethodGet, "engines/%s/schema", engineName)

	if data != nil {
		data["id"] = "text"
	}

	return data, err
}

// Update schema definition by engineName
func (c *client) UpdateSchema(ctx context.Context, engineName string, def schema.Definition) (err error) {
	schemaDefinition := make(schema.Definition)
	for field, fieldType := range def {
		if field != "id" {
			schemaDefinition[field] = fieldType
		}
	}

	err = c.Call(ctx, schemaDefinition, nil, http.MethodPost, "engines/%s/schema", engineName)

	return err
}
