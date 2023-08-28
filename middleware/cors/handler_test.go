package cors

import (
	"github.com/kamilov/goro"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	handler := Handler(Options{
		AllowOrigins: "https://test.com, https://test.test",
		AllowMethods: "PUT, PATCH",
	})

	{
		request, _ := http.NewRequest(http.MethodOptions, "/test/", nil)
		response := httptest.NewRecorder()

		request.Header.Set(headerOrigin, "https://test.com")
		request.Header.Set(headerRequestedMethod, http.MethodPatch)

		ctx := goro.NewContext(request, response)

		assert.Nil(t, handler(ctx))
		assert.Equal(t, "https://test.com", response.Header().Get(headerAllowOrigin))
	}

	{
		request, _ := http.NewRequest(http.MethodPatch, "/test/", nil)
		response := httptest.NewRecorder()

		request.Header.Set(headerOrigin, "https://test.com")

		ctx := goro.NewContext(request, response)

		assert.Nil(t, handler(ctx))
		assert.Equal(t, "https://test.com", response.Header().Get(headerAllowOrigin))
	}

	{
		request, _ := http.NewRequest(http.MethodPatch, "/test/", nil)
		response := httptest.NewRecorder()

		ctx := goro.NewContext(request, response)

		assert.Nil(t, handler(ctx))
		assert.Equal(t, "", response.Header().Get(headerAllowOrigin))
	}
}
