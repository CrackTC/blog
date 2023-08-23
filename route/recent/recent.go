package recent

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"sora.zip/blog/util/file"
	"sora.zip/blog/util/redis"
)

type handler struct {
	blogPath     string
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

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const key string = "[recent]"

	page := getPage(r)
	start := h.blogsPerPage * (page - 1) * 3

	data := recentData{Title: "ps", HasPrev: start > 0, PrevPage: page - 1, NextPage: page + 1}
	if val, err := redis.GetList(key, int64(start), int64(h.blogsPerPage*3)); err == nil {
		data.Items = make([]itemData, len(val)/3)
		for i := 0; i < len(val); i += 3 {
			data.Items[i/3] = itemData{val[i], val[i+1], val[i+2]}
		}
		len, err := redis.Len(key)
		if err != nil {
			log.Printf("[ERROR] Failed to get length of key %s: %s\n", key, err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data.HasNext = page*h.blogsPerPage*3 < int(len)
	} else if err != redis.Nil {
		log.Printf("[ERROR] Failed to get value of key %s: %s\n", "key", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		files := file.GetFileTimesRecursive(h.blogPath, h.ignoredPaths)
		redisData := make([]any, len(files)*3)
		data.Items = make([]itemData, 0, h.blogsPerPage)
		for i, file := range files {
			// remove extension
			name := file.Name[:len(file.Name)-len(filepath.Ext(file.Name))]
			redisData[i*3] = name
			path := file.Path[len(h.blogPath)+1:]
			redisData[i*3+1] = path
			redisData[i*3+2] = file.ModTime
			if i >= start && i < start+h.blogsPerPage {
				data.Items = append(data.Items, itemData{file.Name, path, file.ModTime})
			}
		}
		if err := redis.SetList(key, redisData); err != nil {
			log.Printf("[ERROR] Failed to set value of key %s: %s\n", key, err.Error())
		}
	}
	h.tpl.ExecuteTemplate(w, "index_main", data)
}

func NewHandler(blogPath string, blogsPerPage int, templatePath string, ignoredPaths []string) handler {
	return handler{blogPath, blogsPerPage, ignoredPaths, template.Must(template.ParseFiles(
		filepath.Join(templatePath, "index/recent.html"),
		filepath.Join(templatePath, "index/index_main.html"),
	))}
}
