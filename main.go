package main

import (
	"net/http"
	"strconv"

	"sora.zip/blog/config"
	"sora.zip/blog/route/api"
	"sora.zip/blog/route/article"
	"sora.zip/blog/route/index"
	"sora.zip/blog/route/recent"
	"sora.zip/blog/route/static"
	"sora.zip/blog/route/whoami"
)

func serve(c *config.Config) {
	setup()
	http.Handle("/", index.NewHandler())
	http.Handle("/api/", http.StripPrefix("/api/", api.NewHandler(c.APIKey, c.IgnoredPaths)))
	http.Handle("/static/", http.StripPrefix("/static/", static.NewHandler()))
	http.Handle("/article/", http.StripPrefix("/article/", article.NewHandler(c.IgnoredPaths)))
	http.Handle("/recent/", http.StripPrefix("/recent/", recent.NewHandler(c.BlogsPerPage, c.IgnoredPaths)))
	http.Handle("/whoami/", http.StripPrefix("/whoami/", whoami.NewHandler()))
	err := http.ListenAndServe(":"+strconv.Itoa(c.Port), nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	c := config.Get()
	serve(c)
}
