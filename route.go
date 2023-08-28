package goro

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Route struct {
	group    *Group
	path     string
	name     string
	handlers []Handler
}

func (r *Route) String() string {
	return r.Path()
}

func (r *Route) Name(name string) *Route {
	for n, route := range r.group.router.routes {
		if route == r {
			delete(r.group.router.routes, n)
			break
		}
	}
	r.group.router.routes[name] = r
	return r
}

func (r *Route) Path() string {
	return r.group.prefix + r.path
}

func (r *Route) Generate(values ...any) (result string) {
	path := strings.TrimRight(r.Path(), "*")
	valueIndex, openIndex, closeIndex := 0, -1, -1

	for i := 0; i < len(path); i++ {
		if path[i] == '<' && openIndex == -1 {
			openIndex = i
		} else if path[i] == '>' && openIndex >= 0 {
			value := ""

			if valueIndex < len(values) {
				value = url.QueryEscape(fmt.Sprint(values[valueIndex]))
				valueIndex++
			}

			result += path[closeIndex+1:openIndex] + value
			closeIndex = i
			openIndex = -1
		}
	}

	if closeIndex == -1 {
		result = path
	} else if closeIndex < len(path)-1 {
		result += path[closeIndex+1:]
	}

	return
}

func (r *Route) Get(handlers ...Handler) *Route {
	r.group.add(http.MethodGet, r.Path(), combineHandlers(r.handlers, handlers))
	return r
}

func (r *Route) Post(handlers ...Handler) *Route {
	r.group.add(http.MethodPost, r.Path(), combineHandlers(r.handlers, handlers))
	return r
}

func (r *Route) Put(handlers ...Handler) *Route {
	r.group.add(http.MethodPut, r.Path(), combineHandlers(r.handlers, handlers))
	return r
}

func (r *Route) Patch(handlers ...Handler) *Route {
	r.group.add(http.MethodPatch, r.Path(), combineHandlers(r.handlers, handlers))
	return r
}

func (r *Route) Delete(handlers ...Handler) *Route {
	r.group.add(http.MethodDelete, r.Path(), combineHandlers(r.handlers, handlers))
	return r
}

func (r *Route) Connect(handlers ...Handler) *Route {
	r.group.add(http.MethodConnect, r.Path(), combineHandlers(r.handlers, handlers))
	return r
}

func (r *Route) Head(handlers ...Handler) *Route {
	r.group.add(http.MethodHead, r.Path(), combineHandlers(r.handlers, handlers))
	return r
}

func (r *Route) Options(handlers ...Handler) *Route {
	r.group.add(http.MethodOptions, r.Path(), combineHandlers(r.handlers, handlers))
	return r
}

func (r *Route) Trace(handlers ...Handler) *Route {
	r.group.add(http.MethodTrace, r.Path(), combineHandlers(r.handlers, handlers))
	return r
}

func (r *Route) Any(handlers ...Handler) *Route {
	r.group.Any(r.Path(), combineHandlers(r.handlers, handlers)...)
	return r
}

func (r *Route) To(methods string, handlers ...Handler) *Route {
	r.group.To(methods, r.Path(), handlers...)
	return r
}
