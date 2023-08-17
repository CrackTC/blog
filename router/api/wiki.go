package api

import (
	"context"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"sora.zip/blog/config"
)

type data struct {
	Dest string `json:"dest"`
}

var rdb *redis.Client
var ctx context.Context

func getWikiDestination(name string) string {
	val, err := rdb.Get(ctx, name).Result()
	if err == nil {
		return val
	}

	if err != redis.Nil {
		log.Println("[ERROR] Failed to get value of key", name+":", err.Error())
		return ""
	}

	c := config.Get()
	blogPath := filepath.Join(c.StaticPath, c.BlogDirectoryName)

	if strings.ContainsRune(name, '/') { // absolute path
		return filepath.Join(blogPath, name)
	}

	var res string
	filepath.WalkDir(blogPath, func(path string, d fs.DirEntry, err error) error {
		base := d.Name()
		for _, ignorePath := range c.IgnoredPaths {
			if base == ignorePath {
				return fs.SkipDir
			}
		}
		if !d.IsDir() && base == name {
			// /static/blog/... to /article/...
			rel, _ := filepath.Rel(blogPath, path)
			res = filepath.Join("/article", rel)
			return fs.SkipAll
		}
		return nil
	})

	if res == "" {
		log.Println("[WARN] Could not find destination for", name)
		return ""
	}

	expiration, err := time.ParseDuration(c.RedisExpiration)
	if err != nil {
		log.Println("[ERROR] Invalid redis expiration:", c.RedisExpiration)
		expiration = time.Duration(0)
	}

	_, err = rdb.Set(ctx, name, res, expiration).Result()
	if err != nil {
		log.Print("[ERROR] Failed to set key", name+":", err.Error())
	}
	return res
}

func wiki(w http.ResponseWriter, arguments map[string][]string) {
	d := data{Dest: getWikiDestination(arguments["name"][0])}
	b, err := json.Marshal(d)
	if err != nil {
		log.Println("[ERROR] Json marshal failed:", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func init() {
	c := config.Get()
	opt, err := redis.ParseURL(c.RedisURL)
	if err != nil {
		log.Fatalln("[ERROR] Failed to parse redis url:", c.RedisURL)
	}

	rdb = redis.NewClient(opt)
	ctx = context.Background()
}
