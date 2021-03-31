package appsearch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	t.Run("resolve", func(t *testing.T) {
		t.Run("Must resolve host with token", func(t *testing.T) {
			host, token, authType, err := resolve("https://token@host")
			require.NoError(t, err)
			require.Equal(t, "https://host/api/as/v1/", host)
			require.Equal(t, "token", token)
			require.Equal(t, "Bearer", authType)
		})

		t.Run("Must resolve host with basic auth", func(t *testing.T) {
			host, token, authType, err := resolve("https://uname:pass@host")
			require.NoError(t, err)
			require.Equal(t, "https://host/api/as/v1/", host)
			require.Equal(t, "dW5hbWU6cGFzcw==", token)
			require.Equal(t, "Basic", authType)
		})

		t.Run("Must resolve host without token", func(t *testing.T) {
			host, _, _, err := resolve("https://host")
			require.NoError(t, err)
			require.Equal(t, "https://host/api/as/v1/", host)
		})

		t.Run("Must return error for invalid URL", func(t *testing.T) {
			_, _, _, err := resolve("%")
			require.Error(t, err)
		})
	})

	t.Run("Open", func(t *testing.T) {
		t.Run("Must open client from endpoint", func(t *testing.T) {
			_, err := Open("https://token@host")
			require.NoError(t, err)
		})

		t.Run("Must open client from basic auth", func(t *testing.T) {
			_, err := Open("https://uname:pass@host")
			require.NoError(t, err)
		})

		t.Run("Must open client from endpoint with token", func(t *testing.T) {
			_, err := Open("https://host", "token")
			require.NoError(t, err)
		})
	})
}
