package static

import (
	"net/http"
	"strings"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Split(r.URL.Path, "/")[0] == "blog" {
		http.ServeFile(w, r, "web/var/"+r.URL.Path)
		return
	}
	http.ServeFile(w, r, "web/static/"+r.URL.Path)
}

func NewHandler() http.Handler {
	return &handler{}
}
