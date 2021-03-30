package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lithiumlabcompany/appsearch"
)

func TestMock(t *testing.T) {
	t.Run("Must implement APIClient", func(t *testing.T) {
		var c appsearch.APIClient
		c = &mock{}
		_ = c
	})

	t.Run("Not so dirty implementation", func(t *testing.T) {

	})

	t.Run("Dirty implementation hacks", func(t *testing.T) {
		mockRequest := appsearch.CreateEngineRequest{
			Name:     "testEngine",
			Language: "",
		}
		mockResult := appsearch.EngineDescription{
			Name:          "testEngine",
			Type:          "engine",
			Language:      &mockRequest.Language,
			DocumentCount: 0,
		}

		m := Mock(
			map[string]interface{}{
				"CreateEngine": func(ctx context.Context, request appsearch.CreateEngineRequest) (desc appsearch.EngineDescription, err error) {
					require.EqualValues(t, mockRequest, request)

					return mockResult, nil
				},
			},
		)

		res, err := m.CreateEngine(nil, mockRequest)
		require.NoError(t, err)

		require.EqualValues(t, mockResult, res)
	})
}
