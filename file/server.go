package file

import (
	"github.com/kamilov/goro"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type (
	Options struct {
		Root     string
		Index    string
		CatchAll string
		Allow    func(ctx *goro.Context, path string) bool
	}

	Map map[string]string
)

func Server(mapping Map, options Options) goro.Handler {
	if !filepath.IsAbs(options.Root) {
		wd, _ := os.Getwd()
		options.Root = filepath.Join(wd, options.Root)
	}

	from, to := parseMap(mapping)
	dir := http.Dir(options.Root)

	return func(ctx *goro.Context) error {
		if ctx.Request().Method != http.MethodGet && ctx.Request().Method != http.MethodHead {
			return goro.NewError(http.StatusMethodNotAllowed)
		}

		path, found := match(ctx.Request().URL.Path, from, to)

		if !found || options.Allow != nil && !options.Allow(ctx, path) {
			return goro.NewError(http.StatusForbidden)
		}

		var (
			file http.File
			stat os.FileInfo
			err  error
		)

		if file, err = dir.Open(path); err != nil {
			if options.CatchAll != "" {
				return serve(ctx, dir, options.CatchAll)
			}
			return goro.NewError(http.StatusNotFound, err.Error())
		}

		defer file.Close()

		if stat, err = file.Stat(); err != nil {
			return goro.NewError(http.StatusNotFound, err.Error())
		}

		if stat.IsDir() {
			if options.Index == "" {
				return goro.NewError(http.StatusForbidden)
			}
			return serve(ctx, dir, filepath.Join(path, options.Index))
		}

		ctx.Response().Header().Del("Content-Type")
		http.ServeContent(ctx.Response(), ctx.Request(), path, stat.ModTime(), file)
		return nil
	}
}

func Content(path string) goro.Handler {
	if !filepath.IsAbs(path) {
		wd, _ := os.Getwd()
		path = filepath.Join(wd, path)
	}

	return func(ctx *goro.Context) error {
		file, err := os.Open(path)
		if err != nil {
			return goro.NewError(http.StatusNotFound, err.Error())
		}

		defer file.Close()

		stat, err := file.Stat()

		if err != nil {
			return goro.NewError(http.StatusNotFound, err.Error())
		} else if stat.IsDir() {
			return goro.NewError(http.StatusForbidden)
		}

		ctx.Response().Header().Del("Content-Type")
		http.ServeContent(ctx.Response(), ctx.Request(), path, stat.ModTime(), file)
		return nil
	}
}

func serve(ctx *goro.Context, dir http.Dir, path string) error {
	file, err := dir.Open(path)

	if err != nil {
		return goro.NewError(http.StatusNotFound, err.Error())
	}

	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		return goro.NewError(http.StatusNotFound, err.Error())
	} else if stat.IsDir() {
		return goro.NewError(http.StatusForbidden)
	}

	ctx.Response().Header().Del("Content-Type")
	http.ServeContent(ctx.Response(), ctx.Request(), path, stat.ModTime(), file)
	return nil
}

func match(path string, from []string, to []string) (string, bool) {
	for i := len(from) - 1; i >= 0; i-- {
		prefix := from[i]

		if strings.HasPrefix(path, prefix) {
			return to[i] + path[len(prefix):], true
		}
	}

	return "", false
}

func parseMap(mapping Map) (from, to []string) {
	from, to = make([]string, len(mapping)), make([]string, len(mapping))
	offset := 0

	for key := range mapping {
		from[offset] = key
	}

	sort.Strings(from)

	for i, key := range from {
		to[i] = mapping[key]
	}

	return
}
