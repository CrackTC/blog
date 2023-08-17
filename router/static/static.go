package static

import (
	"log"
	"net/http"
	"sora.zip/blog/config"
)

var fs http.Handler

func Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] Serving file:", r.URL.Path)
	fs.ServeHTTP(w, r)
}

func init() {
	fs = http.StripPrefix("/static/", http.FileServer(http.Dir(config.Get().StaticPath)))
}
