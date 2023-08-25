package index

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Handler struct {
	tpl *template.Template
}

type indexData struct {
	Title string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	log.Println("[INFO] Serving index")
	err := h.tpl.ExecuteTemplate(w, "index_main", indexData{Title: "zipped sora"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewHandler(templatePath string) Handler {
	tpl := template.Must(template.ParseFiles(filepath.Join(templatePath, "index/index.html"), filepath.Join(templatePath, "index/index_main.html")))
	handler := Handler{tpl: tpl}

	return handler
}
