package language

import (
	"github.com/kamilov/goro"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/test/", nil)
	response := httptest.NewRecorder()
	ctx := goro.NewContext(request, response)

	request.Header.Set("Accept-Language", "ru-RU;q=1.0,ru;q=0.9,en-US;q=0.5,en;q=0.6")

	{
		handler := Handler("lang")

		assert.Nil(t, handler(ctx))
		assert.Equal(t, "en-US", ctx.Get("lang"))
	}

	{
		handler := Handler("lang", "en-US", "en", "ru-RU", "ru")

		assert.Nil(t, handler(ctx))
		assert.Equal(t, "ru-RU", ctx.Get("lang"))
	}

	{
		request.Header.Set("Accept-Language", "ru-RU;q=0")

		handler := Handler("lang", "ru", "ru-RU")

		assert.Nil(t, handler(ctx))
		assert.Equal(t, "ru", ctx.Get("lang"))
	}
}
