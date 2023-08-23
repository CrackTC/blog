package index

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type handler struct {
	tpl *template.Template
}

type indexData struct {
	Title string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] Serving index")
	h.tpl.ExecuteTemplate(w, "index_main", indexData{Title: "zipped sora"})
}

func NewHandler(templatePath string) handler {
	tpl := template.Must(template.ParseFiles(filepath.Join(templatePath, "index/index.html"), filepath.Join(templatePath, "index/index_main.html")))
	handler := handler{tpl: tpl}

	return handler
}
