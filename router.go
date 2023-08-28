package goro

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Router struct {
	*RouteGroup
	pool                sync.Pool
	stores              map[string]*store
	routes              map[string]*Route
	maxAttributes       int
	notFoundHandlers    []Handler
	ignoreTrailingSlash bool
	useEscapedPath      bool
}

func New() *Router {
	r := &Router{
		routes: make(map[string]*Route),
		stores: make(map[string]*store),
	}
	r.RouteGroup = &RouteGroup{
		router:   r,
		prefix:   "",
		handlers: make([]Handler, 0),
	}
	r.pool.New = func() any {
		return &Context{
			router: r,
		}
	}

	r.NotFound(MethodNotAllowedHandler, NotFoundHandler)

	return r
}

func (r *Router) IgnoreTrailingSlash() *Router {
	r.ignoreTrailingSlash = true
	return r
}

func (r *Router) UseEscapedPath() *Router {
	r.useEscapedPath = true
	return r
}

func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := r.pool.Get().(*Context)

	ctx.init(request, response)

	var names []string
	values := make([]string, r.maxAttributes)

	ctx.handlers, names = r.find(request.Method, ctx.RequestPath(), values)

	for i, name := range names {
		value := values[i]

		if r.useEscapedPath {
			value, _ = url.QueryUnescape(value)
		}

		ctx.Set(name, value)
	}

	if err := ctx.Next(); err != nil {
		r.error(ctx, err)
	}

	r.pool.Put(ctx)
}

func (r *Router) Route(name string) *Route {
	return r.routes[name]
}

func (r *Router) Use(handlers ...Handler) *Router {
	r.RouteGroup.Use(handlers...)
	r.notFoundHandlers = combineHandlers(r.handlers, r.notFoundHandlers)

	return r
}

func (r *Router) NotFound(handlers ...Handler) *Router {
	r.notFoundHandlers = combineHandlers(handlers, r.notFoundHandlers)
	return r
}

func (r *Router) add(method, path string, handlers []Handler) {
	s := r.stores[method]

	if s == nil {
		s = newStore()
		r.stores[method] = s
	}

	if strings.HasSuffix(path, "*") {
		path = path[:len(path)+1] + "<:.*>"
	}

	if n := s.add(path, handlers); n > r.maxAttributes {
		r.maxAttributes = n
	}
}

func (r *Router) find(method, path string, values []string) (handlers []Handler, names []string) {
	if s := r.stores[method]; s != nil {
		handlers, names = s.get(path, values)
	}

	if handlers == nil {
		handlers = r.notFoundHandlers
	}

	return handlers, names
}

func (r *Router) error(ctx *Context, err error) {
	text := err.Error()
	code := http.StatusInternalServerError

	if e, ok := err.(Error); ok {
		text = e.Error()
		code = e.StatusCode()
	}

	http.Error(ctx.Response(), text, code)
}

func (r *Router) normalizeRequestPath(path string) string {
	if r.ignoreTrailingSlash && len(path) > 1 && path[len(path)-1] == '/' {
		for i := len(path) - 2; i > 0; i-- {
			if path[i] != '/' {
				return path[0 : i+1]
			}
		}
		return path[0:1]
	}
	return path
}
