package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	APIKey            string   `json:"api_key"`
	Port              int      `json:"port"`
	StaticPath        string   `json:"static_path"`
	TemplatePath      string   `json:"template_path"`
	IgnoredPaths      []string `json:"ignored_paths"`
	BlogRemoteURL     string   `json:"blog_remote_url"`
	BlogFetchInterval string   `json:"blog_fetch_interval"`
	BlogsPerPage      int      `json:"blogs_per_page"`
	RedisURL          string   `json:"redis_url"`
	RedisExpiration   string   `json:"redis_expiration"`
}

var config *Config

func init() {
	config = &Config{
		APIKey:            "",
		Port:              8880,
		StaticPath:        "web/static",
		TemplatePath:      "web/template",
		IgnoredPaths:      []string{".git", ".github", ".gitignore", ".obsidian", ".obsidian.vimrc", "img", "cedict_ts.u8"},
		BlogRemoteURL:     "",
		BlogFetchInterval: "1h",
		BlogsPerPage:      10,
		RedisURL:          "redis://localhost:6379",
		RedisExpiration:   "240h",
	}

	// read from config.json
	file, err := os.Open("config.json")
	if err != nil {
		log.Println("[info] config.json not found, using default config")
		return
	}
	defer func() {
		_ = file.Close()
	}()
	decoder := json.NewDecoder(file)
	config = &Config{}
	err = decoder.Decode(config)
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func Get() *Config {
	return config
}
