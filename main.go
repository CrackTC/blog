package main

import (
	"net/http"
	"strconv"

	"sora.zip/blog/config"
	"sora.zip/blog/router/api"
	"sora.zip/blog/router/article"
	"sora.zip/blog/router/index"
	"sora.zip/blog/router/static"
)

func main() {
	http.HandleFunc("/", index.Handle)
	http.HandleFunc("/api/", api.Handle)
	http.HandleFunc("/article/", article.Handle)
	http.HandleFunc("/static/", static.Handle)
	err := http.ListenAndServe(":"+strconv.Itoa(config.Get().Port), nil)
	if err != nil {
		panic(err)
	}
}
