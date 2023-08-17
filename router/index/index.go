package index

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"sora.zip/blog/config"
	"sora.zip/blog/file"
)

var fileTree *file.FileTree

func GetFileLink(path string) string {
	c := config.Get()
	blogPath := filepath.Join(c.StaticPath, c.BlogDirectoryName)
	rel, _ := filepath.Rel(filepath.Join(blogPath), path)
	href := strings.Replace(filepath.Join("/article", rel), "%", "%25", -1)
	text := template.HTMLEscapeString(filepath.Base(path))
	return "<a href=\"" + href + "\">" + text + "</a>"
}

func GetFileListHTML(node *file.FileTree) string {
	if node == nil {
		return ""
	}

	res := "<li>"
	if node.Children != nil {
		res += template.HTMLEscapeString(filepath.Base(node.Name)) + "/"
		res += "\n<ul>\n"
		for _, child := range node.Children {
			res += GetFileListHTML(child)
		}
		res += "\n</ul>\n"
	} else {
		res += GetFileLink(node.Name)
	}
	res += "</li>\n"

	return res
}

func Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] Serving index")
	t, err := template.ParseFiles(filepath.Join(config.Get().TemplatePath, "index.html"))
	if err != nil {
		log.Println("[ERROR] Failed to parse template: index.html")
		return
	}
	t.Execute(w, template.HTML(GetFileListHTML(fileTree)))
}

func init() {
	c := config.Get()
	blogPath := filepath.Join(c.StaticPath, c.BlogDirectoryName)
	tree, err := file.GetFileTree(blogPath)
	if err != nil {
		log.Fatalln("[ERROR] Failed to get file tree")
	}
	fileTree = tree
}
