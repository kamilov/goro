package file

import (
	"github.com/kamilov/goro"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	{
		server := Server(Map{
			"/css": "/testdata/css",
		}, Options{
			Root: "./",
		})
		tests := []struct {
			name   string
			method string
			url    string
			body   string
			status int
		}{
			{"ok content", http.MethodGet, "/css/style.css", "* { color: red; }", 0},
			{"ok header", http.MethodHead, "/css/style.css", "", 0},
			{"not found", http.MethodGet, "/css/style2.css", "", http.StatusNotFound},
			{"method not allowed", http.MethodPost, "/css/style.css", "", http.StatusMethodNotAllowed},
			{"forbidden", http.MethodGet, "/css", "", http.StatusForbidden},
		}

		for _, test := range tests {
			request, _ := http.NewRequest(test.method, test.url, nil)
			response := httptest.NewRecorder()
			ctx := goro.NewContext(request, response)
			err := server(ctx)

			if test.status == 0 {
				assert.Nil(t, err, test.name+" is not error")
				assert.Equal(t, test.body, response.Body.String(), test.name+".body =")
			} else if assert.NotNil(t, err) {
				assert.Equal(t, test.status, err.(goro.Error).StatusCode(), test.name+".status =")
			}
		}
	}

	{
		server := Server(Map{
			"/css": "/testdata/css",
		}, Options{
			Root:  "./",
			Index: "index.html",
			Allow: func(ctx *goro.Context, path string) bool {
				return path != "/testdata/css/style.css"
			},
		})
		request, _ := http.NewRequest(http.MethodGet, "/css/style.css", nil)
		response := httptest.NewRecorder()
		ctx := goro.NewContext(request, response)
		err := server(ctx)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusForbidden, err.(goro.Error).StatusCode())

		request, _ = http.NewRequest(http.MethodGet, "/css", nil)
		response = httptest.NewRecorder()
		ctx = goro.NewContext(request, response)
		err = server(ctx)

		assert.Nil(t, err)
		assert.Equal(t, "css.html", response.Body.String())
	}

	{
		server := Server(Map{
			"/css": "/testdata/css",
		}, Options{
			Root:     "./",
			Index:    "index.html",
			CatchAll: "/testdata/index.html",
			Allow: func(ctx *goro.Context, path string) bool {
				return path != "/testdata/css/style.css"
			},
		})

		request, _ := http.NewRequest(http.MethodGet, "/css/style.css", nil)
		response := httptest.NewRecorder()
		ctx := goro.NewContext(request, response)
		err := server(ctx)

		assert.NotNil(t, err)

		request, _ = http.NewRequest(http.MethodGet, "/css", nil)
		response = httptest.NewRecorder()
		ctx = goro.NewContext(request, response)
		err = server(ctx)

		assert.Nil(t, err)
		assert.Equal(t, "css.html", response.Body.String())

		request, _ = http.NewRequest(http.MethodGet, "/css/main.css", nil)
		response = httptest.NewRecorder()
		ctx = goro.NewContext(request, response)
		err = server(ctx)

		assert.Nil(t, err)
		assert.Equal(t, "index.html", response.Body.String())
	}
}

func TestContent(t *testing.T) {
	h := Content("testdata/index.html")
	request, _ := http.NewRequest("GET", "/index.html", nil)
	response := httptest.NewRecorder()
	c := goro.NewContext(request, response)
	err := h(c)
	assert.Nil(t, err)
	assert.Equal(t, "index.html", response.Body.String())
}
