package goro

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute_Name(t *testing.T) {
	r := New()

	root := r.Add("/").Name("root").Get()
	test := r.Add("/test").Name("test").Get()

	assert.Equal(t, r.Route("root"), root, "r.routes[root] = ")
	assert.Equal(t, r.Route("test"), test, "r.routes[test] = ")
	assert.Equal(t, 2, len(r.routes))
}

func TestRoute_Generate(t *testing.T) {
	r := New()

	r.Add("/test/<id:\\d+>/<name:\\w+>").Name("test").Get()

	assert.Equal(t, "/test/123/name", r.Route("test").Generate(123, "name"))
	assert.Equal(t, "/test/123/", r.Route("test").Generate(123))
}
