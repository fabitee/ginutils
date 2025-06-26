package ginutils_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabitee/ginutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPathUUID(t *testing.T) {
	t.Run("using correct path parameter", func(t *testing.T) {
		handler := gin.New()
		handler.GET("/:id", ginutils.HandlerWithErr(func(c *gin.Context) error {
			id, err := ginutils.GetPathUUID(c, "id")
			if err != nil {
				return err
			}
			c.String(http.StatusOK, id.String())
			return nil
		}))

		server := httptest.NewServer(handler)
		defer server.Close()

		t.Run("request with valid uuid", func(t *testing.T) {
			id := uuid.New()
			url := fmt.Sprintf("%s/%s", server.URL, id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err, "creating test request failed")

			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err, "sending test request failed")
			defer response.Body.Close()
			assertUUIDResponse(t, response, id)
		})

		t.Run("request with invalid uuid", func(t *testing.T) {
			id := "non-uuid-string"
			url := fmt.Sprintf("%s/%s", server.URL, id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err, "creating test request failed")

			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err, "sending test request failed")
			defer response.Body.Close()
			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	})

	t.Run("using incorrect path parameter", func(t *testing.T) {
		handler := gin.New()
		handler.GET("/:id", ginutils.HandlerWithErr(func(c *gin.Context) error {
			id, err := ginutils.GetPathUUID(c, "not_found")
			if err != nil {
				return err
			}
			c.String(http.StatusOK, id.String())
			return nil
		}))

		server := httptest.NewServer(handler)
		defer server.Close()

		t.Run("request with valid uuid will still respond with 400", func(t *testing.T) {
			id := uuid.New()
			url := fmt.Sprintf("%s/%s", server.URL, id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err, "creating test request failed")

			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err, "sending test request failed")
			defer response.Body.Close()
			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	})
}

func assertUUIDResponse(t *testing.T, response *http.Response, expectedUUID uuid.UUID) {
	t.Helper()

	assert.Equal(t, 200, response.StatusCode, "unexpected response status")
	b, err := io.ReadAll(response.Body)
	require.NoError(t, err, "reading response body failed")

	assert.Equal(t, expectedUUID.String(), string(b), "unexpected response body content")

	u, err := uuid.Parse(string(b))
	require.NoError(t, err, "parsing UUID from response body failed")
	assert.Equal(t, expectedUUID, u, "parsed UUID does not match expected UUID")
}
