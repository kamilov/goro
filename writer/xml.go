package writer

import (
	"encoding/xml"
	"net/http"
)

type XmlWriter struct{}

func (w XmlWriter) SetHeader(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/xml; charset=UTF-8")
}

func (w XmlWriter) Write(response http.ResponseWriter, data any) (err error) {
	var bytes []byte

	if bytes, err = xml.Marshal(data); err != nil {
		return
	}

	_, err = response.Write(bytes)

	return
}
