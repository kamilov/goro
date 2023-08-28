package slash

import (
	"github.com/kamilov/goro"
	"net/http"
	"strings"
)

func Handler(status int) goro.Handler {
	return func(ctx *goro.Context) error {
		if ctx.Request().URL.Path != "/" && strings.HasSuffix(ctx.Request().URL.Path, "/") {
			if ctx.Request().Method != http.MethodGet {
				status = http.StatusTemporaryRedirect
			}
			ctx.Redirect(strings.TrimRight(ctx.Request().URL.Path, "/"), status)
		}
		return nil
	}
}
