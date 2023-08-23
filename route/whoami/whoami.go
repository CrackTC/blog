package whoami

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type handler struct {
	tpl *template.Template
}

type whoamiData struct {
	Title string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.tpl.ExecuteTemplate(w, "index_main", whoamiData{Title: "whoami"})
}

func NewHandler(templatePath string) handler {
	return handler{tpl: template.Must(template.ParseFiles(
		filepath.Join(templatePath, "index/whoami.html"),
		filepath.Join(templatePath, "index/index_main.html"),
	))}
}
