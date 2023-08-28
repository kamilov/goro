package content

import (
	"github.com/kamilov/goro"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	handler := Handler(
		goro.MimeTypeHtml,
		goro.MimeTypeJson,
		goro.MimeTypeXml,
	)

	{
		request, _ := http.NewRequest("GET", "/", nil)

		request.Header.Set("Accept", "application/json;q=1;v=1")

		response := httptest.NewRecorder()
		ctx := goro.NewContext(request, response)

		assert.Nil(t, handler(ctx))
		assert.Equal(t, `application/json`, ctx.Response().(*httptest.ResponseRecorder).Header().Get("Content-Type"))
	}

	{
		request, _ := http.NewRequest("GET", "/", nil)

		request.Header.Set("Accept", "application/xml;q=1;v=1,application/xml;q=0.9;v=0.8")

		response := httptest.NewRecorder()
		ctx := goro.NewContext(request, response)

		assert.Nil(t, handler(ctx))
		assert.Equal(t, `application/xml; charset=UTF-8`, ctx.Response().(*httptest.ResponseRecorder).Header().Get("Content-Type"))
	}

	{
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()
		ctx := goro.NewContext(request, response)

		assert.Nil(t, handler(ctx))
		assert.Equal(t, `text/html; charset=UTF-8`, ctx.Response().(*httptest.ResponseRecorder).Header().Get("Content-Type"))
	}
}
