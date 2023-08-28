package writer

import (
	"fmt"
	"net/http"
)

type (
	Writer interface {
		SetHeader(http.ResponseWriter)
		Write(http.ResponseWriter, any) error
	}

	DefaultWriter struct{}
)

func (w *DefaultWriter) SetHeader(response http.ResponseWriter) {}
func (w *DefaultWriter) Write(response http.ResponseWriter, data any) error {
	return Write(response, data)
}

func Write(response http.ResponseWriter, data any) error {
	var err error

	switch data.(type) {
	case []byte:
		_, err = response.Write(data.([]byte))

	case string:
		_, err = response.Write([]byte(data.(string)))

	default:
		if data != nil {
			_, err = fmt.Fprint(response, data)
		}
	}

	return err
}
