package main

import (
	"net/http"
	"path/filepath"
	"strconv"

	"sora.zip/blog/config"
	"sora.zip/blog/route/api"
	"sora.zip/blog/route/article"
	"sora.zip/blog/route/index"
	"sora.zip/blog/route/recent"
	"sora.zip/blog/route/static"
	"sora.zip/blog/route/whoami"
)

func main() {
	c := config.Get()
	blogPath := filepath.Join(c.StaticPath, "blog")
	setup()
	http.Handle("/", index.NewHandler(c.TemplatePath))
	http.Handle("/api/", http.StripPrefix("/api/", api.NewHandler(c.APIKey, blogPath, c.IgnoredPaths)))
	http.Handle("/static/", http.StripPrefix("/static/", static.NewHandler(c.StaticPath)))
	http.Handle("/article/", http.StripPrefix("/article/", article.NewHandler(blogPath, c.IgnoredPaths, c.TemplatePath)))
	http.Handle("/recent/", http.StripPrefix("/recent/", recent.NewHandler(blogPath, c.BlogsPerPage, c.TemplatePath, c.IgnoredPaths)))
	http.Handle("/whoami/", http.StripPrefix("/whoami/", whoami.NewHandler(c.TemplatePath)))
	err := http.ListenAndServe(":"+strconv.Itoa(c.Port), nil)
	if err != nil {
		panic(err)
	}
}
