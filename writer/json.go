package writer

import (
	"encoding/json"
	"net/http"
)

type JsonWriter struct{}

func (w JsonWriter) SetHeader(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json")
}

func (w JsonWriter) Write(response http.ResponseWriter, data any) error {
	encoder := json.NewEncoder(response)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(data)
}
