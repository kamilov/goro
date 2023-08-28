package goro

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type storeTestEntry struct {
	path            string
	attributesCount int
}

var testHandlers = []Handler{func(ctx *Context) error {
	return nil
}}

func TestStoreAdd(t *testing.T) {
	tests := []struct {
		name     string
		entries  []storeTestEntry
		expected string
	}{
		{
			"static",
			[]storeTestEntry{
				{"/test/a", 0},
				{"/test/b", 0},
				{"/test/c", 0},
			},
			`{path: , order: 0, minOrder: 0, countHandlers: 0, attributeIndex: -1, attributeNames: [], requirement: <nil>}
    {path: /test/, order: 1, minOrder: 1, countHandlers: 0, attributeIndex: -1, attributeNames: [], requirement: <nil>}
        {path: a, order: 1, minOrder: 1, countHandlers: 1, attributeIndex: -1, attributeNames: [], requirement: <nil>}
        {path: b, order: 2, minOrder: 2, countHandlers: 1, attributeIndex: -1, attributeNames: [], requirement: <nil>}
        {path: c, order: 3, minOrder: 3, countHandlers: 1, attributeIndex: -1, attributeNames: [], requirement: <nil>}
`,
		},
		{
			"attributes",
			[]storeTestEntry{
				{"/test/<id>", 1},
				{"/test/<id>/a", 1},
				{"/test/<id>/<test:\\d+>", 2},
				{`/test/<id>/<test:\d+>/b`, 2},
				{`/test/<id>/<test:\d+>/c`, 2},
			},
			`{path: , order: 0, minOrder: 0, countHandlers: 0, attributeIndex: -1, attributeNames: [], requirement: <nil>}
    {path: /test/, order: 0, minOrder: 1, countHandlers: 0, attributeIndex: -1, attributeNames: [], requirement: <nil>}
        {path: <id>, order: 1, minOrder: 1, countHandlers: 1, attributeIndex: 0, attributeNames: [id], requirement: <nil>}
            {path: /, order: 2, minOrder: 2, countHandlers: 0, attributeIndex: 0, attributeNames: [id], requirement: <nil>}
                {path: a, order: 2, minOrder: 2, countHandlers: 1, attributeIndex: 0, attributeNames: [id], requirement: <nil>}
                {path: <test:\d+>, order: 3, minOrder: 3, countHandlers: 1, attributeIndex: 1, attributeNames: [id test], requirement: ^\d+}
                    {path: /, order: 4, minOrder: 4, countHandlers: 0, attributeIndex: 1, attributeNames: [id test], requirement: <nil>}
                        {path: b, order: 4, minOrder: 4, countHandlers: 1, attributeIndex: 1, attributeNames: [id test], requirement: <nil>}
                        {path: c, order: 5, minOrder: 5, countHandlers: 1, attributeIndex: 1, attributeNames: [id test], requirement: <nil>}
`,
		},
		{
			"corners",
			[]storeTestEntry{
				{"/test/<id>/test/<name>", 2},
				{"/test/abc/<id>/<name>", 2},
			},
			`{path: , order: 0, minOrder: 0, countHandlers: 0, attributeIndex: -1, attributeNames: [], requirement: <nil>}
    {path: /test/, order: 0, minOrder: 1, countHandlers: 0, attributeIndex: -1, attributeNames: [], requirement: <nil>}
        {path: abc/, order: 0, minOrder: 2, countHandlers: 0, attributeIndex: -1, attributeNames: [], requirement: <nil>}
            {path: <id>, order: 0, minOrder: 2, countHandlers: 0, attributeIndex: 0, attributeNames: [id], requirement: <nil>}
                {path: /, order: 0, minOrder: 2, countHandlers: 0, attributeIndex: 0, attributeNames: [id], requirement: <nil>}
                    {path: <name>, order: 2, minOrder: 2, countHandlers: 1, attributeIndex: 1, attributeNames: [id name], requirement: <nil>}
        {path: <id>, order: 0, minOrder: 1, countHandlers: 0, attributeIndex: 0, attributeNames: [id], requirement: <nil>}
            {path: /test/, order: 0, minOrder: 1, countHandlers: 0, attributeIndex: 0, attributeNames: [id], requirement: <nil>}
                {path: <name>, order: 1, minOrder: 1, countHandlers: 1, attributeIndex: 1, attributeNames: [id name], requirement: <nil>}
`,
		},
	}

	for _, test := range tests {
		s := newStore()

		for _, entry := range test.entries {
			n := s.add(entry.path, testHandlers)
			assert.Equal(t, entry.attributesCount, n, test.name+" > "+entry.path+" > attributes count = ")
		}

		assert.Equal(t, test.expected, s.String(), test.name+" > store.String() = ")
	}
}

func TestStoreGet(t *testing.T) {
	s := newStore()
	m := 0
	paths := []string{
		"/test/img.png",
		"/test/test",
		"/test/<id>",
		"/test/<id>/<name:\\w+>",
		"/list/<id>/<name:\\w+>/<page:\\d+>",
	}
	tests := []struct {
		path       string
		attributes string
		handlers   []Handler
	}{
		{"/test/img.png", "", testHandlers},
		{"/test/test", "", testHandlers},
		{"/test/id", "id:id,", testHandlers},
		{"/test/test/test", "id:test,name:test,", testHandlers},
		{"/list/test/test/2", "id:test,name:test,page:2,", testHandlers},
		{"/list/test/test/a", "", nil},
	}

	for _, path := range paths {
		i := s.add(path, testHandlers)

		if i > m {
			m = i
		}
	}

	fmt.Println(s.String())

	assert.Equal(t, 3, m, "max attributes = ")

	values := make([]string, m)

	for _, test := range tests {
		handlers, names := s.get(test.path, values)

		assert.Equal(t, test.handlers, handlers, "store.Get("+test.path+") = ")

		attributes := ""

		if len(names) > 0 {
			for i, name := range names {
				attributes += fmt.Sprintf("%v:%v,", name, values[i])
			}
		}

		assert.Equal(t, test.attributes, attributes, "store.Get("+test.path+").attributes = ")
	}
}
