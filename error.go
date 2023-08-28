package goro

import (
	"net/http"
	"strings"
)

type errorResponse struct {
	Status  int    `json:"status" xml:"status"`
	Message string `json:"message" xml:"message"`
}

type Error interface {
	error
	StatusCode() int
}

func NewError(code int, messages ...string) Error {
	if len(messages) == 0 {
		messages = append(messages, http.StatusText(code))
	}
	return &errorResponse{code, strings.Join(messages, " ")}
}

func (e *errorResponse) Error() string {
	return e.Message
}

func (e *errorResponse) StatusCode() int {
	return e.Status
}
