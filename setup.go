package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"sora.zip/blog/config"
	"sora.zip/blog/util/git"
	"sora.zip/blog/util/redis"
)

func isNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

func updateRepo(path string, c <-chan time.Time) {
	for range c {
		err := git.PullRepo(path)
		if err != nil {
			log.Println("[ERROR] failed to pull repo:", err.Error())
			continue
		}

		err = git.UpdateModTime(path)
		if err != nil {
			log.Println("[ERROR] failed to update file mod time:", err.Error())
			continue
		}

		err = redis.Flush()
		if err != nil {
			log.Println("[ERROR] failed to flush redis:", err.Error())
		}
	}
}

func setModTimeZero() {
	err := filepath.Walk("web/static/blog", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("[ERROR] failed to walk file:", err.Error())
			return nil
		}
		if info.IsDir() {
			return nil
		}
		_ = os.Chtimes(path, time.Time{}, time.UnixMilli(1))
		return nil
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func setup() {
	path := "web/static/blog"
	if isNotExist(path) || isNotExist(filepath.Join(path, ".git")) {

		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatal(err.Error())
		}

		err := git.CloneRepo(config.Get().BlogRemoteURL, path)
		if err != nil {
			log.Fatal("[ERROR] failed to clone repo:", err.Error())
		}

		setModTimeZero()

		err = git.UpdateModTime(path)
		if err != nil {
			log.Println("[ERROR] failed to update file mod time:", err.Error())
		}
	}

	if duration, err := time.ParseDuration(config.Get().BlogFetchInterval); err != nil {
		log.Fatal("[ERROR] failed to parse fetch interval:", err.Error())
	} else {
		ticker := time.NewTicker(duration)
		go updateRepo(path, ticker.C)
	}
}
