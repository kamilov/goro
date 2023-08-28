package goro

import (
	"net/http"
	"strings"
)

type Group struct {
	router   *Router
	prefix   string
	handlers []Handler
}

func (g *Group) Group(prefix string, handlers ...Handler) *Group {
	return &Group{
		router:   g.router,
		prefix:   g.prefix + "/" + strings.TrimLeft(prefix, "/"),
		handlers: combineHandlers(handlers, g.handlers),
	}
}

func (g *Group) Use(handlers ...Handler) *Group {
	g.handlers = append(g.handlers, handlers...)
	return g
}

func (g *Group) Add(path string, handlers ...Handler) *Route {
	return &Route{
		group:    g,
		path:     path,
		handlers: handlers,
	}
}

func (g *Group) Get(path string, handlers ...Handler) *Group {
	g.add(http.MethodGet, path, handlers)
	return g
}

func (g *Group) Post(path string, handlers ...Handler) *Group {
	g.add(http.MethodPost, path, handlers)
	return g
}

func (g *Group) Put(path string, handlers ...Handler) *Group {
	g.add(http.MethodPut, path, handlers)
	return g
}

func (g *Group) Patch(path string, handlers ...Handler) *Group {
	g.add(http.MethodPatch, path, handlers)
	return g
}

func (g *Group) Delete(path string, handlers ...Handler) *Group {
	g.add(http.MethodDelete, path, handlers)
	return g
}

func (g *Group) Connect(path string, handlers ...Handler) *Group {
	g.add(http.MethodConnect, path, handlers)
	return g
}

func (g *Group) Head(path string, handlers ...Handler) *Group {
	g.add(http.MethodHead, path, handlers)
	return g
}

func (g *Group) Options(path string, handlers ...Handler) *Group {
	g.add(http.MethodOptions, path, handlers)
	return g
}

func (g *Group) Trace(path string, handlers ...Handler) *Group {
	g.add(http.MethodTrace, path, handlers)
	return g
}

func (g *Group) Any(path string, handlers ...Handler) *Group {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range methods {
		g.add(method, path, handlers)
	}

	return g
}

func (g *Group) To(methods, path string, handlers ...Handler) *Group {
	m := strings.Split(methods, ",")

	if len(methods) == 1 {
		g.add(strings.TrimSpace(methods), path, handlers)
	}

	for _, method := range m {
		g.add(strings.TrimSpace(method), path, handlers)
	}

	return g
}

func (g *Group) add(method, path string, handlers []Handler) {
	g.router.add(method, path, combineHandlers(g.handlers, handlers))
}

func combineHandlers(a []Handler, b []Handler) []Handler {
	c := make([]Handler, len(a)+len(b))

	copy(c, a)
	copy(c[len(a):], b)

	return c
}
