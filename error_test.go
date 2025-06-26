package ginutils_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabitee/ginutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerWithErr(t *testing.T) {
	var (
		testServer *httptest.Server
	)

	cleanupTestServer := func() {
		testServer.Close()
	}

	setupTestServer := func(t *testing.T, handler ginutils.HandlerFunc) {
		t.Cleanup(cleanupTestServer)

		engine := gin.New()
		engine.Use(ginutils.Recovery())
		engine.Any("/", ginutils.HandlerWithErr(handler))
		testServer = httptest.NewServer(engine)
	}

	t.Run("passes ErrorResponse as response", func(t *testing.T) {
		setupTestServer(t, func(context *gin.Context) error {
			return ginutils.BadRequest("some parameter is invalid")
		})

		resp, err := http.Get(testServer.URL)
		require.NoError(t, err)
		defer resp.Body.Close()
		assertErrorResponse(t, resp, ginutils.ErrorResponse{
			Status:  400,
			Message: "some parameter is invalid",
		})
	})

	t.Run("converts any error into ErrorResponse", func(t *testing.T) {
		setupTestServer(t, func(context *gin.Context) error {
			return errors.New("this is some error")
		})

		resp, err := http.Get(testServer.URL)
		require.NoError(t, err)
		defer resp.Body.Close()
		assertErrorResponse(t, resp, ginutils.ErrorResponse{
			Status:  500,
			Message: "this is some error",
		})
	})

	t.Run("responds with error object on panics (using the recover middleware)", func(t *testing.T) {
		t.Run("with string value", func(t *testing.T) {
			setupTestServer(t, func(context *gin.Context) error {
				panic("some panic")
			})

			resp, err := http.Get(testServer.URL)
			require.NoError(t, err)
			defer resp.Body.Close()
			assertErrorResponse(t, resp, ginutils.ErrorResponse{
				Status:  500,
				Message: "some panic",
			})
		})

		t.Run("with error value", func(t *testing.T) {
			setupTestServer(t, func(context *gin.Context) error {
				panic(errors.New("some error"))
			})

			resp, err := http.Get(testServer.URL)
			require.NoError(t, err)
			defer resp.Body.Close()
			assertErrorResponse(t, resp, ginutils.ErrorResponse{
				Status:  500,
				Message: "some error",
			})
		})

		t.Run("with ErrorResponse", func(t *testing.T) {
			setupTestServer(t, func(context *gin.Context) error {
				panic(ginutils.BadRequest("this is some bad request"))
			})

			resp, err := http.Get(testServer.URL)
			require.NoError(t, err)
			defer resp.Body.Close()
			assertErrorResponse(t, resp, ginutils.ErrorResponse{
				Status:  400, // The status code is preserved
				Message: "this is some bad request",
			})
		})
	})
}

func assertErrorResponse(t *testing.T, resp *http.Response, expected ginutils.ErrorResponse) {
	t.Helper()

	assert.Equal(t, expected.Status, resp.StatusCode)

	var actual ginutils.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&actual)
	require.NoError(t, err, "failed to decode response body as ErrorResponse")

	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.Message, actual.Message)
}
