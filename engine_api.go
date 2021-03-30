package appsearch

import (
	"context"
	"errors"
	"net/http"
)

// List all available engines
func (c *client) ListAllEngines(ctx context.Context) (engines []EngineDescription, err error) {
	page := 0
	totalPages := 1

	for page < totalPages {
		page += 1

		res, err := c.ListEngines(ctx, Page{page, 25})
		if err != nil {
			return nil, err
		}

		totalPages = res.Meta.Page.TotalPages
		engines = append(engines, res.Results...)
	}

	return engines, err
}

// List engines with pagination
func (c *client) ListEngines(ctx context.Context, page Page) (data EngineResponse, err error) {
	err = c.Call(ctx, page, &data, http.MethodGet, "engines")

	return data, err
}

// List an engine by name
func (c *client) ListEngine(ctx context.Context, engineName string) (data EngineDescription, err error) {
	err = c.Call(ctx, nil, &data, http.MethodGet, "engines/%s", engineName)

	return data, err
}

// Create an engine
func (c *client) CreateEngine(ctx context.Context, request CreateEngineRequest) (resp EngineDescription, err error) {
	err = c.Call(ctx, request, &resp, http.MethodPost, "engines")

	return
}

// Delete engine by name
func (c *client) DeleteEngine(ctx context.Context, engineName string) (err error) {
	err = c.Call(ctx, nil, nil, http.MethodDelete, "engines/%s", engineName)

	return
}

// Create engine if doesn't exist.
// Optionally update a schema even if engine exists.
func (c *client) EnsureEngine(ctx context.Context, request CreateEngineRequest, schema ...SchemaDefinition) (err error) {
	_, err = c.ListEngine(ctx, request.Name)

	if errors.Is(err, ErrEngineDoesntExist) {
		_, err = c.CreateEngine(ctx, request)
	}

	if err == nil && len(schema) > 0 {
		err = c.UpdateSchema(ctx, request.Name, schema[0])
	}

	return err
}
