package api

import (
	"log"
	"net/http"
	"path/filepath"
)

type Handler struct {
	blogRoot     string
	ignoredPaths []string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint := filepath.Base(r.URL.Path)
	arguments := r.URL.Query()
	log.Println("[INFO] Serving API:", endpoint)

	switch endpoint {
	case "wiki":
		h.wiki(w, arguments)
	}
}

func NewHandler(blogRoot string, ignoredPaths []string) Handler {
	return Handler{blogRoot: blogRoot, ignoredPaths: ignoredPaths}
}
