package api

import (
	"log"
	"net/http"
	"path/filepath"
)

type Handler struct {
	apiKey       string
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
	case "pull":
		h.pull(w, arguments)
	}
}

func NewHandler(apiKey string, blogRoot string, ignoredPaths []string) Handler {
	return Handler{apiKey: apiKey, blogRoot: blogRoot, ignoredPaths: ignoredPaths}
}
