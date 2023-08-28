package content

import (
	"github.com/kamilov/goro"
	"github.com/kamilov/goro/writer"
)

var writers map[string]writer.Writer

func init() {
	writers = map[string]writer.Writer{
		goro.MimeTypeHtml:    &writer.HtmlWriter{},
		goro.MimeTypeXml:     &writer.XmlWriter{},
		goro.MimeTypeTextXml: &writer.XmlWriter{},
		goro.MimeTypeJson:    &writer.JsonWriter{},
	}
}

func Handler(formats ...string) goro.Handler {
	if len(formats) == 0 {
		formats = []string{goro.MimeTypeHtml}
	}

	for _, format := range formats {
		if _, ok := writers[format]; !ok {
			panic(format + " is not supported")
		}
	}

	return func(ctx *goro.Context) error {
		format := NegotiateContentType(ctx.Request(), formats, formats[0])
		ctx.SetDataWriter(writers[format])
		return nil
	}
}
