package mock

import (
	"context"
	"errors"

	"github.com/lithiumlabcompany/appsearch"
	"github.com/lithiumlabcompany/appsearch/pkg/schema"
)

func (m *mock) ListEngine(ctx context.Context, engineName string) (data appsearch.EngineDescription, err error) {
	data, ok := m.Engines[engineName]
	if !ok {
		err = appsearch.ErrEngineDoesntExist
	}
	return
}

func (m *mock) ListEngines(ctx context.Context, page appsearch.Page) (data appsearch.EngineResponse, err error) {
	return appsearch.EngineResponse{
		Results: engineValues(m.Engines),
	}, nil
}

func (m *mock) ListAllEngines(ctx context.Context) (data []appsearch.EngineDescription, err error) {
	return engineValues(m.Engines), nil
}

func (m *mock) CreateEngine(ctx context.Context, request appsearch.CreateEngineRequest) (desc appsearch.EngineDescription, err error) {
	_, ok := m.Engines[request.Name]
	if ok {
		return desc, appsearch.ErrEngineAlreadyExists
	}

	m.Engines[request.Name] = appsearch.EngineDescription{
		Name:          request.Name,
		Type:          "engine",
		Language:      &request.Language,
		DocumentCount: 0,
	}
	return m.ListEngine(ctx, request.Name)
}

func (m *mock) DeleteEngine(ctx context.Context, engineName string) (err error) {
	_, ok := m.Engines[engineName]
	if !ok {
		return appsearch.ErrEngineDoesntExist
	}
	delete(m.Engines, engineName)
	return nil
}

func (m *mock) EnsureEngine(ctx context.Context, request appsearch.CreateEngineRequest, schema ...schema.Definition) (err error) {
	_, err = m.ListEngine(ctx, request.Name)

	if errors.Is(err, appsearch.ErrEngineDoesntExist) {
		_, err = m.CreateEngine(ctx, request)
	}

	if err == nil && len(schema) > 0 {
		err = m.UpdateSchema(ctx, request.Name, schema[0])
	}

	return err
}

func engineValues(engines map[string]appsearch.EngineDescription) []appsearch.EngineDescription {
	values := make([]appsearch.EngineDescription, 0, len(engines))
	for _, engine := range engines {
		values = append(values, engine)
	}
	return values
}
