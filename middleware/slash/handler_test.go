package slash

import (
	"github.com/kamilov/goro"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	remover := Handler(http.StatusMovedPermanently)
	request, _ := http.NewRequest(http.MethodGet, "/test/", nil)
	response := httptest.NewRecorder()
	ctx := goro.NewContext(request, response)
	err := remover(ctx)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusMovedPermanently, response.Code)
	assert.Equal(t, "/test", response.Header().Get("Location"))

	request, _ = http.NewRequest(http.MethodGet, "/", nil)
	response = httptest.NewRecorder()
	ctx = goro.NewContext(request, response)
	err = remover(ctx)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "", response.Header().Get("Location"))

}
