package reader

import (
	"net/http"
)

type Reader interface {
	Read(*http.Request, any) error
}

var DefaultReader Reader = &FormReader{}
