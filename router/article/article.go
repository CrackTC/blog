package article

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"sora.zip/blog/config"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] Serving", r.URL.Path)
	t, err := template.ParseFiles(filepath.Join(config.Get().TemplatePath, "article.html"))
	if err != nil {
		log.Println("[ERROR] Failed to parse template: article.html")
		return
	}

	t.Execute(w, filepath.Base(r.URL.Path))
}
