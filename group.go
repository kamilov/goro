package goro

import (
	"net/http"
	"strings"
)

type RouteGroup struct {
	router   *Router
	prefix   string
	handlers []Handler
}

func (g *RouteGroup) Group(prefix string, handlers ...Handler) *RouteGroup {
	return &RouteGroup{
		router:   g.router,
		prefix:   g.prefix + "/" + strings.TrimLeft(prefix, "/"),
		handlers: combineHandlers(handlers, g.handlers),
	}
}

func (g *RouteGroup) Use(handlers ...Handler) *RouteGroup {
	g.handlers = append(g.handlers, handlers...)
	return g
}

func (g *RouteGroup) Add(path string, handlers ...Handler) *Route {
	return &Route{
		group:    g,
		path:     path,
		handlers: handlers,
	}
}

func (g *RouteGroup) Get(path string, handlers ...Handler) *RouteGroup {
	g.add(http.MethodGet, path, handlers)
	return g
}

func (g *RouteGroup) Post(path string, handlers ...Handler) *RouteGroup {
	g.add(http.MethodPost, path, handlers)
	return g
}

func (g *RouteGroup) Put(path string, handlers ...Handler) *RouteGroup {
	g.add(http.MethodPut, path, handlers)
	return g
}

func (g *RouteGroup) Patch(path string, handlers ...Handler) *RouteGroup {
	g.add(http.MethodPatch, path, handlers)
	return g
}

func (g *RouteGroup) Delete(path string, handlers ...Handler) *RouteGroup {
	g.add(http.MethodDelete, path, handlers)
	return g
}

func (g *RouteGroup) Connect(path string, handlers ...Handler) *RouteGroup {
	g.add(http.MethodConnect, path, handlers)
	return g
}

func (g *RouteGroup) Head(path string, handlers ...Handler) *RouteGroup {
	g.add(http.MethodHead, path, handlers)
	return g
}

func (g *RouteGroup) Options(path string, handlers ...Handler) *RouteGroup {
	g.add(http.MethodOptions, path, handlers)
	return g
}

func (g *RouteGroup) Trace(path string, handlers ...Handler) *RouteGroup {
	g.add(http.MethodTrace, path, handlers)
	return g
}

func (g *RouteGroup) Any(path string, handlers ...Handler) *RouteGroup {
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

func (g *RouteGroup) To(methods, path string, handlers ...Handler) *RouteGroup {
	m := strings.Split(methods, ",")

	if len(methods) == 1 {
		g.add(strings.TrimSpace(methods), path, handlers)
	}

	for _, method := range m {
		g.add(strings.TrimSpace(method), path, handlers)
	}

	return g
}

func (g *RouteGroup) add(method, path string, handlers []Handler) {
	g.router.add(method, path, combineHandlers(g.handlers, handlers))
}

func combineHandlers(a []Handler, b []Handler) []Handler {
	c := make([]Handler, len(a)+len(b))

	copy(c, a)
	copy(c[len(a):], b)

	return c
}
