package language

import (
	"github.com/kamilov/goro"
	"golang.org/x/text/language"
	"net/http"
)

func Handler(contextKey string, languages ...string) goro.Handler {
	if len(languages) == 0 {
		languages = []string{"en-US"}
	}

	defaultLanguage := languages[0]

	return func(ctx *goro.Context) error {
		lang := negotiateLanguage(ctx.Request(), languages, defaultLanguage)
		ctx.Set(contextKey, lang)
		return nil
	}
}

func negotiateLanguage(request *http.Request, languages []string, defaultLanguage string) string {
	bestLanguage := defaultLanguage
	bestQ := float32(-1.0)
	specs, q, _ := language.ParseAcceptLanguage(request.Header.Get("Accept-Language"))

	for _, lang := range languages {
		for i, spec := range specs {
			if q[i] > bestQ && (spec.String() == "*" || spec.String() == lang) {
				bestQ = q[i]
				bestLanguage = lang
			}
		}
	}

	if bestQ == 0 {
		bestLanguage = defaultLanguage
	}

	return bestLanguage
}
