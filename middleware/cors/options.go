package cors

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	Wildcard = "*"
	DennyAll = "-"

	headerOrigin = "Origin"

	headerRequestedMethod  = "Access-Control-Request-Method"
	headerRequestedHeaders = "Access-Control-Request-Headers"

	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerMaxAge           = "Access-Control-Max-Age"
)

type Options struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	ExposeHeaders    string
	AllowCredentials bool
	MaxAge           time.Duration

	allowOrigin []string
	allowMethod []string
	allowHeader []string
}

var AllowAll = Options{
	AllowOrigins: Wildcard,
	AllowHeaders: Wildcard,
	AllowMethods: Wildcard,
}

func (o *Options) init() {
	o.allowHeader = buildAllow(o.AllowHeaders, false)
	o.allowMethod = buildAllow(o.AllowMethods, true)
	o.allowOrigin = buildAllow(o.AllowOrigins, true)
}

func (o *Options) isOriginAllowed(origin string) bool {
	if o.AllowOrigins == DennyAll {
		return false
	}
	return o.AllowOrigins == Wildcard || slices.Contains(o.allowOrigin, origin)
}

func (o *Options) setActualHeaders(origin string, headers http.Header) {
	if !o.isOriginAllowed(origin) {
		return
	}

	o.setOriginHeader(origin, headers)

	if o.ExposeHeaders != "" {
		headers.Set(headerExposeHeaders, o.ExposeHeaders)
	}
}

func (o *Options) setOriginHeader(origin string, headers http.Header) {
	if o.AllowOrigins == Wildcard {
		origin = Wildcard
	}

	headers.Set(headerAllowOrigin, origin)

	if o.AllowCredentials {
		headers.Set(headerAllowCredentials, "true")
	}
}

func (o *Options) setPreflightHeaders(origin, method, requestHeaders string, headers http.Header) {
	allowed, allowedHeaders := o.isPreflightAllowed(origin, method, requestHeaders)

	if !allowed {
		return
	}

	o.setOriginHeader(origin, headers)

	if o.MaxAge > time.Duration(0) {
		headers.Set(headerMaxAge, strconv.FormatInt(int64(o.MaxAge/time.Second), 10))
	}

	if o.AllowMethods == Wildcard {
		headers.Set(headerAllowMethods, method)
	} else if slices.Contains(o.allowMethod, strings.ToUpper(method)) {
		headers.Set(headerAllowMethods, o.AllowMethods)
	}

	if allowedHeaders != "" {
		headers.Set(headerAllowHeaders, allowedHeaders)
	}
}

func (o *Options) isPreflightAllowed(origin string, method string, headers string) (allowed bool, allowHeaders string) {
	if !o.isOriginAllowed(origin) {
		return
	}

	if o.AllowMethods != Wildcard && !slices.Contains(o.allowMethod, method) {
		return
	}

	if o.AllowHeaders == Wildcard || headers == "" {
		return true, headers
	}

	result := []string{}

	for _, header := range strings.Split(headers, ",") {
		header = strings.TrimSpace(header)

		if slices.Contains(o.allowHeader, strings.ToUpper(header)) {
			result = append(result, header)
		}
	}

	if len(result) > 0 {
		allowed = true
		allowHeaders = strings.Join(result, ",")
	}

	return
}

func buildAllow(allow string, caseSensitive bool) []string {
	result := []string{}

	if len(allow) > 0 {
		for _, item := range strings.Split(allow, ",") {
			item = strings.TrimSpace(item)

			if !caseSensitive {
				item = strings.ToUpper(item)
			}

			result = append(result, item)
		}
	}

	return result
}
