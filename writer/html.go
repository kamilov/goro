package writer

import "net/http"

type HtmlWriter struct{}

func (w *HtmlWriter) SetHeader(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "text/html; charset=UTF-8")
}

func (w *HtmlWriter) Write(response http.ResponseWriter, data any) error {
	return Write(response, data)
}
