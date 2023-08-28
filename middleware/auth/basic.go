package auth

import (
	"encoding/base64"
	"github.com/kamilov/goro"
	"net/http"
	"strings"
)

type BasicAuthFunc func(ctx *goro.Context, username, password string) (Identity, error)

func Basic(realm string, fn BasicAuthFunc) goro.Handler {
	return func(ctx *goro.Context) error {
		username, password := parseBasicCredentials(ctx.Request().Header.Get("Authorization"))
		identity, err := fn(ctx, username, password)

		if err == nil {
			ctx.Set(ContextUserKey, identity)
			return nil
		}

		ctx.Response().Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)

		return goro.NewError(http.StatusUnauthorized, err.Error())
	}
}

func parseBasicCredentials(header string) (username, password string) {
	if strings.HasPrefix(header, "Basic ") {
		if bytes, err := base64.StdEncoding.DecodeString(header[6:]); err == nil {
			str := string(bytes)

			if i := strings.IndexByte(str, ':'); i >= 0 {
				username = str[:i]
				password = str[i+1:]
			}
		}
	}
	return
}
