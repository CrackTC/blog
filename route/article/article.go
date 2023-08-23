package article

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"sora.zip/blog/util/redis"
	"sora.zip/blog/util/url"
)

type handler struct {
	blogRoot     string
	ignoredPaths []string
	tpl          map[string]*template.Template
}

type linkData struct {
	Name template.HTML
	URL  template.URL
}

type articleData struct {
	Title   template.HTML
	Content template.HTML
	Crumbs  []linkData
}

type dirData struct {
	Title  string
	Files  []linkData
	Crumbs []linkData
}

func GetCrumbs(path string) []linkData {
	if path == "" {
		return []linkData{{"article", "javascript:void(0)"}}
	}

	arr := strings.Split(path, "/")
	res := make([]linkData, 0, len(arr)+1)
	res = append(res, linkData{"article", "/article/"})
	for i, base := range arr {
		res = append(res, linkData{
			template.HTML(base),
			template.URL(url.Encode("/article/" + strings.Join(arr[:i+1], "/"))),
		})
	}
	res[len(res)-1].URL = "javascript:void(0)" // last one is current page
	return res
}

func (h handler) ServeArticle(w http.ResponseWriter, path string) {
	log.Println("[INFO] Serving article:", path)

	data := articleData{Title: template.HTML(filepath.Base(path)), Crumbs: GetCrumbs(path)}
	key := "[article]" + path

	if html, err := redis.Get(key); err == nil {
		data.Content = template.HTML(html)
	} else {
		if err != redis.Nil {
			log.Printf("[ERROR] Failed to get value of key %s: %s\n", key, err.Error())
		}

		path = filepath.Join(h.blogRoot, path)
		html := HtmlFromFile(path)
		data.Content = template.HTML(html)

		if err := redis.Set(key, html); err != nil {
			log.Printf("[ERROR] Failed to set key %s: %s\n", key, err.Error())
		}
	}

	h.tpl["article"].ExecuteTemplate(w, "article_main", data)
}

func (h handler) ServeDir(w http.ResponseWriter, path string) {
	log.Println("[INFO] Serving available articles:", path)

	// remove trailing slash
	if path != "" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	data := dirData{Title: filepath.Base(path), Crumbs: GetCrumbs(path)}

	// get list of files
	if files, err := os.ReadDir(filepath.Join(h.blogRoot, path)); err != nil {
		log.Printf("[ERROR] Failed to read dir %s: %s\n", path, err.Error())
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		data.Files = make([]linkData, 0, len(files))
		for _, file := range files {
			name := file.Name()
			var shouldIgnore bool
			for _, ignored := range h.ignoredPaths {
				if name == ignored {
					shouldIgnore = true
					break
				}
			}
			if shouldIgnore {
				continue
			}
			if file.IsDir() {
				name += "/"
			}
			data.Files = append(data.Files, linkData{
				template.HTML(name),
				template.URL(url.Encode(filepath.Join("/article", path, name))),
			})
		}
	}

	h.tpl["dir"].ExecuteTemplate(w, "article_main", data)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// determine whether it is a directory
	if info, err := os.Stat(filepath.Join(h.blogRoot, path)); err != nil {
		log.Printf("[ERROR] Failed to stat %s: %s\n", path, err.Error())
		http.Error(w, "404 Not Found", http.StatusNotFound)
	} else if info.IsDir() {
		h.ServeDir(w, path)
	} else {
		h.ServeArticle(w, path)
	}
}

func NewHandler(blogRoot string, ignoredPaths []string, templatePath string) handler {
	tpl := make(map[string]*template.Template)
	tpl["article"] = template.Must(template.ParseFiles(
		filepath.Join(templatePath, "article/article.html"),
		filepath.Join(templatePath, "article/article_main.html"),
	))
	tpl["dir"] = template.Must(template.ParseFiles(
		filepath.Join(templatePath, "article/dir.html"),
		filepath.Join(templatePath, "article/article_main.html"),
	))
	return handler{blogRoot: blogRoot, ignoredPaths: ignoredPaths, tpl: tpl}
}
