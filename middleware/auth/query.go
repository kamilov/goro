package auth

import (
	"github.com/kamilov/goro"
	"net/http"
)

func Query(tokenName string, fn TokenAuthFunc) goro.Handler {
	return func(ctx *goro.Context) error {
		token := ctx.Request().URL.Query().Get(tokenName)
		identity, err := fn(ctx, token)

		if err != nil {
			return goro.NewError(http.StatusUnauthorized, err.Error())
		}

		ctx.Set(ContextUserKey, identity)

		return nil
	}
}
