package goro

import (
	"github.com/kamilov/goro/reader"
	"github.com/kamilov/goro/writer"
	"net/http"
)

// Context provides contextual data and environment while processing an incoming HTTP request.
type Context struct {
	router   *Router
	request  *http.Request
	response http.ResponseWriter
	writer   writer.Writer
	data     map[string]any
	handlers []Handler
	offset   int
}

func NewContext(request *http.Request, response http.ResponseWriter, handlers ...Handler) *Context {
	ctx := &Context{
		handlers: handlers,
	}
	ctx.init(request, response)

	return ctx
}

func (c *Context) init(request *http.Request, response http.ResponseWriter) {
	c.request = request
	c.response = response
	c.data = nil
	c.offset = -1

	c.SetDataWriter(&writer.DefaultWriter{})
}

// Request returns the current request.
func (c *Context) Request() *http.Request {
	return c.request
}

// Response returns the current response.
func (c *Context) Response() http.ResponseWriter {
	return c.response
}

// RequestPath returns the current requested path
func (c *Context) RequestPath() string {
	if c.router.useEscapedPath {
		return c.request.URL.EscapedPath()
	} else {
		return c.request.URL.Path
	}
}

// SetDataWriter set the data writer.
func (c *Context) SetDataWriter(w writer.Writer) *Context {
	c.writer = w
	w.SetHeader(c.response)
	return c
}

// Write writes the transmitted data to the response in the given format.
func (c *Context) Write(data any) error {
	return c.writer.Write(c.response, data)
}

// WriteWithStatus writes the transmitted data to the response in the specified format and passes the response Status to the header.
func (c *Context) WriteWithStatus(data any, statusCode int) error {
	c.response.WriteHeader(statusCode)
	return c.writer.Write(c.response, data)
}

// Set stores data with the given name in the request context
func (c *Context) Set(name string, value any) *Context {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	c.data[name] = value
	return c
}

// Get returns the data by the specified name that was stored in the context
func (c *Context) Get(name string) any {
	return c.data[name]
}

func (c *Context) URL(route string, values ...any) string {
	if r := c.router.routes[route]; r != nil {
		return r.Generate(values...)
	}
	return ""
}

func (c *Context) Redirect(url string, status int) {
	http.Redirect(c.response, c.request, url, status)
	c.Abort()
}

func (c *Context) Next() error {
	c.offset++

	for n := len(c.handlers); c.offset < n; c.offset++ {
		if err := c.handlers[c.offset](c); err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) Abort() {
	c.offset = len(c.handlers)
}

func (c *Context) Query(name string, defaultValue ...string) string {
	if value, _ := c.request.URL.Query()[name]; len(value) > 0 {
		return value[0]
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

func (c *Context) Form(name string, defaultValue ...string) string {
	value := c.request.FormValue(name)

	if value == "" && len(defaultValue) > 0 {
		value = defaultValue[0]
	}

	return value
}

func (c *Context) PostForm(name string, defaultValue ...string) string {
	value := c.request.PostFormValue(name)

	if value == "" && len(defaultValue) > 0 {
		value = defaultValue[0]
	}

	return value
}

func (c *Context) Parse(data any) error {
	if c.request.Method != http.MethodGet {
		contentType := getContentType(c.request)

		if contentTypeReader, ok := readers[contentType]; ok {
			return contentTypeReader.Read(c.request, data)
		}
	}

	return reader.DefaultReader.Read(c.request, data)
}

var readers map[string]reader.Reader

func init() {
	readers = map[string]reader.Reader{
		MimeTypeUrlencodedForm: &reader.FormReader{},
		MimeTypeMultipartForm:  &reader.FormReader{},
		MimeTypeJson:           &reader.JsonReader{},
		MimeTypeXml:            &reader.XmlReader{},
		MimeTypeTextXml:        &reader.XmlReader{},
	}
}

func getContentType(request *http.Request) string {
	contentType := request.Header.Get("Content-Type")

	for i, char := range contentType {
		if char == ' ' || char == ';' {
			return contentType[:i]
		}
	}

	return contentType
}
