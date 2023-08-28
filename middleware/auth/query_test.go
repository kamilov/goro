package auth

import (
	"errors"
	"github.com/kamilov/goro"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuery(t *testing.T) {
	handler := Query("token", func(ctx *goro.Context, token string) (Identity, error) {
		if token == "test" {
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
		assert.Nil(t, ctx.Get(ContextUserKey))
	}

	{
		request, _ := http.NewRequest(http.MethodGet, "/test/?token=test", nil)
		response := httptest.NewRecorder()
		ctx := goro.NewContext(request, response)
		err := handler(ctx)

		assert.Nil(t, err)
		assert.Equal(t, "test", ctx.Get(ContextUserKey))
	}
}
