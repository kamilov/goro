package goro

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContext_Attributes(t *testing.T) {
	ctx := NewContext(nil, nil)

	ctx.Set("test1", "1")
	ctx.Set("test2", "2")

	assert.Equal(t, "1", ctx.Get("test1"))
	assert.Equal(t, "2", ctx.Get("test2"))
	assert.Equal(t, nil, ctx.Get("test3"))
}

func TestContext_init(t *testing.T) {
	ctx := NewContext(nil, nil)

	assert.Nil(t, ctx.Request())
	assert.Nil(t, ctx.Response())
	assert.Equal(t, 0, len(ctx.handlers))

	request, _ := http.NewRequest(http.MethodGet, "/test/", nil)

	ctx.init(request, httptest.NewRecorder())

	assert.NotNil(t, ctx.Request())
	assert.NotNil(t, ctx.Response())
	assert.Equal(t, -1, ctx.offset)
	assert.Nil(t, ctx.data)
}
