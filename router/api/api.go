package api

import (
	"net/http"
	"path/filepath"
)


func Handle(w http.ResponseWriter, r *http.Request) {
	endpoint := filepath.Base(r.URL.Path)
	arguments := r.URL.Query()
	switch endpoint {
	case "wiki":
		wiki(w, arguments)
	}
}
