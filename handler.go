package goro

import (
	"net/http"
	"sort"
	"strings"
)

type Handler func(ctx *Context) error

func HTTPHandlerFunc(handler http.HandlerFunc) Handler {
	return func(ctx *Context) error {
		handler(ctx.Response(), ctx.Request())
		return nil
	}
}

func HTTPHandler(handler http.Handler) Handler {
	return func(ctx *Context) error {
		handler.ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

func NotFoundHandler(ctx *Context) error {
	return NewError(http.StatusNotFound)
}

func MethodNotAllowedHandler(ctx *Context) (err error) {
	methods := findAllowedMethods(ctx.router, ctx.RequestPath())

	if len(methods) == 0 {
		return nil
	}

	methods = append(methods, http.MethodOptions)

	sort.Strings(methods)

	ctx.Response().Header().Set("Allow", strings.Join(methods, ", "))

	if ctx.Request().Method != http.MethodOptions {
		err = NewError(http.StatusMethodNotAllowed)
	}

	ctx.Abort()

	return
}

func findAllowedMethods(router *Router, path string) (methods []string) {
	values := make([]string, router.maxAttributes)

	for method, s := range router.stores {
		if handlers, _ := s.get(path, values); handlers != nil {
			methods = append(methods, method)
		}
	}

	return
}
