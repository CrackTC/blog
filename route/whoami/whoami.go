package whoami

import (
	"html/template"
	"net/http"
	"path/filepath"
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

func NewHandler(templatePath string) Handler {
	return Handler{tpl: template.Must(template.ParseFiles(
		filepath.Join(templatePath, "index/whoami.html"),
		filepath.Join(templatePath, "index/index_main.html"),
	))}
}
