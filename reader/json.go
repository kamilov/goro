package reader

import (
	"encoding/json"
	"net/http"
)

type JsonReader struct{}

func (r JsonReader) Read(request *http.Request, data any) error {
	return json.NewDecoder(request.Body).Decode(data)
}
