package appsearch

import (
	"testing"
)

func TestClient(t *testing.T) {
	t.Run("Must implement APIClient", func(t *testing.T) {
		var c APIClient
		c = &client{}
		_ = c
	})
}
