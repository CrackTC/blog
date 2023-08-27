package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"sora.zip/blog/config"
	"sora.zip/blog/util/git"
)

func isNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

func updateRepo(path string, c <-chan time.Time) {
	for range c {
		err := git.PullRepo(path)
		if err != nil {
			log.Println("[ERROR] failed to update repo:", err.Error())
			continue
		}
		err = git.UpdateModTime(path)
		if err != nil {
			log.Println("[ERROR] failed to update file mod time:", err.Error())
		}
	}
}

func setModTimeZero() {
	path := filepath.Join(config.Get().StaticPath, "blog")
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
	path := filepath.Join(config.Get().StaticPath, "blog")
	if isNotExist(path) || isNotExist(filepath.Join(path, ".git")) {

		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatal(err.Error())
		}

		url := config.Get().BlogRemoteURL
		if err := git.CloneRepo(url, path); err != nil {
			log.Fatal(err.Error())
		}

		setModTimeZero()
		if err := git.UpdateModTime(path); err != nil {
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
