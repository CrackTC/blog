package static

import (
	"net/http"
)

func NewHandler() http.Handler {
	return http.FileServer(http.Dir("web/static"))
}
