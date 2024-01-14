package whoami

import (
	"html/template"
	"net/http"
)

type Handler struct {
	tpl *template.Template
}

type whoamiData struct {
	Title string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	err := h.tpl.ExecuteTemplate(w, "index_main", whoamiData{Title: "whoami"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewHandler() Handler {
	return Handler{tpl: template.Must(template.ParseFiles(
		"web/template/index/whoami.html",
		"web/template/index/index_main.html",
	))}
}
