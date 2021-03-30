package appsearch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	t.Run("resolve", func(t *testing.T) {
		t.Run("Must resolve host with key", func(t *testing.T) {
			host, key, err := resolve("https://key@host")
			require.NoError(t, err)
			require.Equal(t, "https://host/api/as/v1/", host)
			require.Equal(t, "key", key)
		})

		t.Run("Must resolve host without key", func(t *testing.T) {
			host, key, err := resolve("https://host")
			require.NoError(t, err)
			require.Equal(t, "https://host/api/as/v1/", host)
			require.Equal(t, "", key)
		})

		t.Run("Must return error for invalid URL", func(t *testing.T) {
			_, _, err := resolve("%")
			require.Error(t, err)
		})
	})

	t.Run("getHostURL", func(t *testing.T) {
		t.Run("Must resolve endpoint with auth", func(t *testing.T) {
			host, key, err := getHostURL([]string{"https://key@host"})
			require.NoError(t, err)
			require.Equal(t, "https://host/api/as/v1/", host)
			require.Equal(t, "key", key)
		})

		t.Run("Must resolve endpoint with key", func(t *testing.T) {
			host, key, err := getHostURL([]string{"https://host", "key"})
			require.NoError(t, err)
			require.Equal(t, "https://host/api/as/v1/", host)
			require.Equal(t, "key", key)
		})

		t.Run("Must return error for invalid URL", func(t *testing.T) {
			_, _, err := getHostURL([]string{"%"})
			require.Error(t, err)
		})

		t.Run("Must return error for invalid usage", func(t *testing.T) {
			_, _, err := getHostURL([]string{})
			require.ErrorIs(t, err, ErrInvalidParams)
		})
	})

	t.Run("Open", func(t *testing.T) {
		t.Run("Must open client from endpoint", func(t *testing.T) {
			_, err := Open("https://key@host")
			require.NoError(t, err)
		})

		t.Run("Must open client from endpoint with key", func(t *testing.T) {
			_, err := Open("https://host", "key")
			require.NoError(t, err)
		})
	})
}
