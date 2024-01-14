package index

import (
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	tpl *template.Template
	fs  http.Handler
}

type indexData struct {
	Title string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.fs.ServeHTTP(w, r)
		return
	}
	log.Println("[INFO] Serving index")
	err := h.tpl.ExecuteTemplate(w, "index_main", indexData{Title: "zipped sora"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewHandler() Handler {
	tpl := template.Must(template.ParseFiles(
		"web/template/index/index.html",
		"web/template/index/index_main.html",
	))
	handler := Handler{tpl: tpl, fs: http.FileServer(http.Dir("web/static"))}
	return handler
}
