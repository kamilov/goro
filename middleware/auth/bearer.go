package auth

import (
	"encoding/base64"
	"github.com/kamilov/goro"
	"net/http"
	"strings"
)

func Bearer(realm string, fn TokenAuthFunc) goro.Handler {
	return func(ctx *goro.Context) error {
		token := parseBearerCredentials(ctx.Request().Header.Get("Authorization"))
		identity, err := fn(ctx, token)

		if err == nil {
			ctx.Set(ContextUserKey, identity)
			return nil
		}

		ctx.Response().Header().Set("WWW-Authenticate", `Bearer realm="`+realm+`"`)

		return goro.NewError(http.StatusUnauthorized, err.Error())
	}
}

func parseBearerCredentials(header string) string {
	if strings.HasPrefix(header, "Bearer ") {
		if bearer, err := base64.StdEncoding.DecodeString(header[7:]); err == nil {
			return string(bearer)
		}
	}
	return ""
}
