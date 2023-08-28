package goro

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(1, "test")

	assert.Equal(t, 1, err.StatusCode())
	assert.Equal(t, "test", err.Error())

	err = NewError(http.StatusNotFound)

	assert.Equal(t, http.StatusNotFound, err.StatusCode())
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

	str, _ := json.Marshal(err)

	assert.Equal(t, `{"status":404,"message":"Not Found"}`, string(str))
}
