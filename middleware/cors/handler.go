package cors

import (
	"github.com/kamilov/goro"
	"net/http"
)

func Handler(options Options) goro.Handler {
	options.init()

	return func(ctx *goro.Context) (err error) {
		origin := ctx.Request().Header.Get(headerOrigin)

		if origin == "" {
			return
		}

		if ctx.Request().Method == http.MethodOptions {
			method := ctx.Request().Header.Get(headerRequestedMethod)

			if method == "" {
				return
			}

			headers := ctx.Request().Header.Get(headerRequestedHeaders)

			options.setPreflightHeaders(origin, method, headers, ctx.Response().Header())
			ctx.Abort()
			return
		}

		options.setActualHeaders(origin, ctx.Response().Header())
		return
	}
}
