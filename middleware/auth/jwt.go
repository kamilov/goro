package auth

import (
	"github.com/kamilov/goro"
	"net/http"
	"strings"
)
import "github.com/golang-jwt/jwt"

type (
	JWTTokenHandler           func(ctx *goro.Context, token *jwt.Token) error
	JWTVerificationKeyHandler func(ctx *goro.Context) string

	JWTOptions struct {
		Realm              string
		SigningMethod      jwt.SigningMethod
		TokenHandler       JWTTokenHandler
		GetVerificationKey JWTVerificationKeyHandler
	}
)

func JWT(verificationKey string, options JWTOptions) goro.Handler {
	parser := &jwt.Parser{
		ValidMethods: []string{options.SigningMethod.Alg()},
	}

	return func(ctx *goro.Context) error {
		header := ctx.Request().Header.Get("Authorization")
		message := ""

		if options.GetVerificationKey != nil {
			verificationKey = options.GetVerificationKey(ctx)
		}

		if strings.HasPrefix(header, "Bearer ") {
			token, err := parser.Parse(header[7:], func(token *jwt.Token) (interface{}, error) {
				return []byte(verificationKey), nil
			})

			if err == nil && token.Valid {
				if options.TokenHandler == nil {
					ctx.Set(ContextUserKey, token)
				} else {
					err = options.TokenHandler(ctx, token)
				}
			}

			if err == nil {
				return nil
			}

			message = err.Error()
		}

		ctx.Response().Header().Set("WWW-Authenticate", `Bearer realm="`+options.Realm+`"`)

		if message != "" {
			return goro.NewError(http.StatusUnauthorized, message)
		}

		return goro.NewError(http.StatusUnauthorized)
	}
}

func NewJWT(signingKey string, signingMethod jwt.SigningMethod, claims jwt.MapClaims) (string, error) {
	return jwt.NewWithClaims(signingMethod, claims).SignedString([]byte(signingKey))
}
