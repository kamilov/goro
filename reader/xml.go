package reader

import (
	"encoding/xml"
	"net/http"
)

type XmlReader struct{}

func (r XmlReader) Read(request *http.Request, data any) error {
	return xml.NewDecoder(request.Body).Decode(data)
}
