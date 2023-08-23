package static

import (
	"log"
	"net/http"
)

type handler struct {
	staticPath string
	fs         http.Handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] Serving static:", r.URL.Path)
	h.fs.ServeHTTP(w, r)
}

func NewHandler(staticPath string) http.Handler {
	return handler{
		staticPath: staticPath, fs: http.FileServer(http.Dir(staticPath)),
	}
}
