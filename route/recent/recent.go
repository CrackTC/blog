package recent

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"sora.zip/blog/util/file"
	"sora.zip/blog/util/redis"
	"sora.zip/blog/util/url"
)

type Handler struct {
	blogsPerPage int
	ignoredPaths []string
	tpl          *template.Template
}

type itemData struct {
	Title   string
	URL     string
	ModTime string
}

type recentData struct {
	Title    string
	Items    []itemData
	HasPrev  bool
	HasNext  bool
	PrevPage int
	NextPage int
}

func getPage(r *http.Request) int {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		return 1
	}

	if page, err := strconv.Atoi(pageParam); err != nil {
		log.Printf("[ERROR] Failed to parse page parameter: %s\n", err.Error())
		return 1
	} else {
		return page
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const key string = "[recent]"

	page := getPage(r)
	start := h.blogsPerPage * (page - 1)

	data := recentData{Title: "ps", HasPrev: start > 0, PrevPage: page - 1, NextPage: page + 1}
	if val, err := redis.GetList(key, int64(start*3), int64(h.blogsPerPage*3)); err == nil {
		data.Items = make([]itemData, len(val)/3)
		for i := 0; i < len(val); i += 3 {
			data.Items[i/3] = itemData{val[i], val[i+1], val[i+2]}
		}
		l, err := redis.Len(key)
		if err != nil {
			log.Printf("[ERROR] Failed to get length of key %s: %s\n", key, err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data.HasNext = page*h.blogsPerPage*3 < int(l)
	} else if err != redis.Nil {
		log.Printf("[ERROR] Failed to get value of key %s: %s\n", "key", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		files := file.GetFileTimesRecursive("web/static/blog", h.ignoredPaths)
		redisData := make([]any, len(files)*3)
		data.Items = make([]itemData, 0, h.blogsPerPage)
		data.HasNext = page*h.blogsPerPage < len(files)
		for i, f := range files {
			// remove extension
			name := f.Name[:len(f.Name)-len(filepath.Ext(f.Name))]
			path := url.Encode(f.Path[len("web/static/blog")+1:])

			redisData[i*3] = name
			redisData[i*3+1] = path
			redisData[i*3+2] = f.ModTime

			if i >= start && i < start+h.blogsPerPage {
				data.Items = append(data.Items, itemData{name, path, f.ModTime})
			}
		}
		if err := redis.SetList(key, redisData); err != nil {
			log.Printf("[ERROR] Failed to set value of key %s: %s\n", key, err.Error())
		}
	}
	err := h.tpl.ExecuteTemplate(w, "index_main", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewHandler(blogsPerPage int, ignoredPaths []string) Handler {
	return Handler{
		blogsPerPage,
		ignoredPaths,
		template.Must(template.ParseFiles(
			"web/template/index/recent.html",
			"web/template/index/index_main.html",
		))}
}
