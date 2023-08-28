package auth

import "github.com/kamilov/goro"

type (
	Identity      interface{}
	TokenAuthFunc func(ctx *goro.Context, token string) (Identity, error)
)

const ContextUserKey = "User"
