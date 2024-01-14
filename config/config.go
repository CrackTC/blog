package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	APIKey            string
	Port              int
	IgnoredPaths      []string
	BlogRemoteURL     string
	BlogFetchInterval string
	BlogsPerPage      int
	RedisURL          string
	RedisExpiration   string
}

var config *Config

func init() {
	config = &Config{
		APIKey:            "",
		Port:              8880,
		IgnoredPaths:      []string{".git", ".github", ".gitignore", ".obsidian", ".obsidian.vimrc", "img", "cedict_ts.u8"},
		BlogRemoteURL:     "",
		BlogFetchInterval: "1h",
		BlogsPerPage:      10,
		RedisURL:          "redis://localhost:6379",
		RedisExpiration:   "240h",
	}

	if (os.Getenv("API_KEY") != "") {
		config.APIKey = os.Getenv("API_KEY")
	}

	if (os.Getenv("PORT") != "") {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err == nil {
			config.Port = port
		}
	}

	if (os.Getenv("IGNORED_PATHS") != "") {
		config.IgnoredPaths = strings.Split(os.Getenv("IGNORED_PATHS"), ":")
	}

	if (os.Getenv("BLOG_REMOTE_URL") != "") {
		config.BlogRemoteURL = os.Getenv("BLOG_REMOTE_URL")
	}

	if (os.Getenv("BLOG_FETCH_INTERVAL") != "") {
		config.BlogFetchInterval = os.Getenv("BLOG_FETCH_INTERVAL")
	}

	if (os.Getenv("BLOGS_PER_PAGE") != "") {
		blogsPerPage, err := strconv.Atoi(os.Getenv("BLOGS_PER_PAGE"))
		if err == nil {
			config.BlogsPerPage = blogsPerPage
		}
	}

	if (os.Getenv("REDIS_URL") != "") {
		config.RedisURL = os.Getenv("REDIS_URL")
	}

	if (os.Getenv("REDIS_EXPIRATION") != "") {
		config.RedisExpiration = os.Getenv("REDIS_EXPIRATION")
	}
}

func Get() *Config {
	return config
}
