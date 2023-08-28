package auth

import (
	"encoding/base64"
	"errors"
	"github.com/kamilov/goro"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasic(t *testing.T) {
	handler := Basic("Need Auth", func(ctx *goro.Context, username, password string) (Identity, error) {
		if username == "test" && password == "test" {
			return "test", nil
		}
		return nil, errors.New("Error")
	})

	{
		request, _ := http.NewRequest(http.MethodGet, "/test/", nil)
		response := httptest.NewRecorder()
		ctx := goro.NewContext(request, response)
		err := handler(ctx)

		assert.NotNil(t, err)
		assert.Equal(t, "Error", err.Error())
		assert.Equal(t, `Basic realm="Need Auth"`, response.Header().Get("WWW-Authenticate"))
		assert.Nil(t, ctx.Get(ContextUserKey))
	}

	{
		request, _ := http.NewRequest(http.MethodGet, "/test/", nil)
		response := httptest.NewRecorder()

		request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("test:test")))

		ctx := goro.NewContext(request, response)
		err := handler(ctx)

		assert.Nil(t, err)
		assert.Equal(t, "test", ctx.Get(ContextUserKey))
		assert.Equal(t, "", response.Header().Get("WWW-Authenticate"))
	}
}
