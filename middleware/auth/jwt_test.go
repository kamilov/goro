package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/kamilov/goro"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testJWTSecretKey = "secret"

var signingMethod = jwt.SigningMethodHS256

func TestJWT(t *testing.T) {
	handler := JWT(testJWTSecretKey, JWTOptions{
		Realm:         "Need Auth",
		SigningMethod: signingMethod,
	})

	{
		token, err := NewJWT(testJWTSecretKey, signingMethod, jwt.MapClaims{
			"name": "test",
		})

		assert.Nil(t, err)

		request, _ := http.NewRequest(http.MethodGet, "/test/", nil)
		response := httptest.NewRecorder()

		request.Header.Set("Authorization", "Bearer "+token)

		ctx := goro.NewContext(request, response)
		err = handler(ctx)

		tkn := ctx.Get(ContextUserKey)

		assert.Nil(t, err)
		assert.NotNil(t, tkn)
		assert.Equal(t, "test", tkn.(*jwt.Token).Claims.(jwt.MapClaims)["name"])
	}

	{
		request, _ := http.NewRequest(http.MethodGet, "/test/", nil)
		response := httptest.NewRecorder()

		request.Header.Set("Authorization", "Bearer asd")

		ctx := goro.NewContext(request, response)
		err := handler(ctx)

		tkn := ctx.Get(ContextUserKey)

		assert.NotNil(t, err)
		assert.Nil(t, tkn)
		assert.Equal(t, `Bearer realm="Need Auth"`, response.Header().Get("WWW-Authenticate"))
	}
}
