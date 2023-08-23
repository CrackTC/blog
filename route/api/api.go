package api

import (
	"log"
	"net/http"
	"path/filepath"
)

type handler struct {
	blogRoot     string
	ignoredPaths []string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint := filepath.Base(r.URL.Path)
	arguments := r.URL.Query()
	log.Println("[INFO] Serving API:", endpoint)

	switch endpoint {
	case "wiki":
		h.wiki(w, arguments)
	}
}

func NewHandler(blogRoot string, ignoredPaths []string) handler {
	return handler{blogRoot: blogRoot, ignoredPaths: ignoredPaths}
}
