package mock

import (
	"context"

	"github.com/lithiumlabcompany/appsearch"
)

func (m *mock) PatchDocuments(ctx context.Context, engineName string, documents interface{}) (res []appsearch.UpdateResponse, err error) {
	m.impl(interfacesOf(ctx, engineName, documents), interfacesOf(&res, &err))
	return
}

func (m *mock) UpdateDocuments(ctx context.Context, engineName string, documents interface{}) (res []appsearch.UpdateResponse, err error) {
	m.impl(interfacesOf(ctx, engineName, documents), interfacesOf(&res, &err))
	return
}

func (m *mock) RemoveDocuments(ctx context.Context, engineName string, documents interface{}) (res []appsearch.DeleteResponse, err error) {
	m.impl(interfacesOf(ctx, engineName, documents), interfacesOf(&res, &err))
	return
}

func (m *mock) ListDocuments(ctx context.Context, engineName string, page appsearch.Page) (response appsearch.DocumentResponse, err error) {
	m.impl(interfacesOf(ctx, engineName, page), interfacesOf(&response, &err))
	return
}

func (m *mock) SearchDocuments(ctx context.Context, engineName string, query appsearch.Query) (response appsearch.DocumentResponse, err error) {
	m.impl(interfacesOf(ctx, engineName, query), interfacesOf(&response, &err))
	return
}
