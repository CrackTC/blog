package api

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"sora.zip/blog/util/file"
	"sora.zip/blog/util/redis"
)

type data struct {
	Dest string `json:"dest"`
}

func (h handler) getWikiDestination(name string) string {
	key := "[wiki]" + name
	if val, err := redis.Get(key); err == nil {
		return val
	} else if err != redis.Nil {
		log.Printf("[ERROR] Failed to get value of key %s: %s\n", key, err.Error())
	}

	var dest string
	if strings.ContainsRune(name, '/') { // absolute path
		dest = name
	} else { // find in filesystem
		dest = file.FindFile(h.blogRoot, name, h.ignoredPaths)
	}

	if dest == "" {
		log.Printf("[WARN] Could not find destination for %s\n", name)
		return ""
	}

	dest = filepath.Join("/article", dest)
	if err := redis.Set(key, dest); err != nil {
		log.Printf("[ERROR] Failed to set key %s: %s\n", key, err.Error())
	}

	return dest
}

func (h handler) wiki(w http.ResponseWriter, arguments map[string][]string) {
	dest := data{Dest: h.getWikiDestination(arguments["name"][0])}

	if bytes, err := json.Marshal(dest); err != nil {
		log.Printf("[ERROR] Json marshal failed: %s\n", err.Error())
	} else {
		w.Write(bytes)
	}
}
