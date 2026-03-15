package main

import (
	"net/http"

	cache "github.com/uaysk/souin-redis/plugins/goyave"
	"goyave.dev/goyave/v4"
	"goyave.dev/goyave/v4/config"
)

func main() {
	_ = config.LoadFrom("examples/config.json")
	_ = goyave.Start(func(r *goyave.Router) {
		r.Get("/{p}", func(response *goyave.Response, r *goyave.Request) {
			_ = response.String(http.StatusOK, "Hello, World 👋!")
		}).Middleware(cache.NewHTTPCache(cache.DevDefaultConfiguration).Handle)
	})
}
