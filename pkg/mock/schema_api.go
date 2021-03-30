package mock

import (
	"context"

	"github.com/lithiumlabcompany/appsearch"
)

func (m *mock) UpdateSchema(ctx context.Context, engineName string, def appsearch.SchemaDefinition) (err error) {
	prev, ok := m.Schemas[engineName]
	if !ok {
		prev = make(appsearch.SchemaDefinition)
		m.Schemas[engineName] = prev
	}
	for field, fieldType := range def {
		prev[field] = fieldType
	}
	return nil
}

func (m *mock) ListSchema(ctx context.Context, engineName string) (data appsearch.SchemaDefinition, err error) {
	return m.Schemas[engineName], nil
}
