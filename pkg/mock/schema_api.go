package mock

import (
	"context"

	"github.com/lithiumlabcompany/appsearch/pkg/schema"
)

func (m *mock) UpdateSchema(ctx context.Context, engineName string, def schema.Definition) (err error) {
	prev, ok := m.Schemas[engineName]
	if !ok {
		prev = make(schema.Definition)
		m.Schemas[engineName] = prev
	}
	for field, fieldType := range def {
		prev[field] = fieldType
	}
	return nil
}

func (m *mock) ListSchema(ctx context.Context, engineName string) (data schema.Definition, err error) {
	return m.Schemas[engineName], nil
}
