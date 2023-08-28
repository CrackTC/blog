package api

import (
	"encoding/json"
	"log"
	"net/http"

	"sora.zip/blog/util/git"
	"sora.zip/blog/util/redis"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h Handler) pull(w http.ResponseWriter, arguments map[string][]string) {
	if len(arguments["api_key"]) == 0 {
		json.NewEncoder(w).Encode(response{Success: false, Message: "API key not provided"})
	} else if api_key := arguments["api_key"][0]; api_key != h.apiKey {
		json.NewEncoder(w).Encode(response{Success: false, Message: "Invalid API key"})
	} else {
		log.Println("[INFO] Pulling blog")
		json.NewEncoder(w).Encode(response{Success: true, Message: "Pulling blog"})
		git.PullRepo(h.blogRoot)
		redis.RemoveKey("[recent]")
	}
}
