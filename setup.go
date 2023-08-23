package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"sora.zip/blog/config"
)

func isNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

func updateModTime(path string, repo *git.Repository) error {
	logs, err := repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return err
	}

	times := make(map[string]time.Time)
	logs.ForEach(func(commit *object.Commit) error {
		stats, err := commit.Stats()
		if err != nil {
			log.Println("[ERROR] failed to get commit stats:", err.Error())
		}
		for _, stat := range stats {
			if _, ok := times[stat.Name]; ok {
				continue
			}
			times[stat.Name] = commit.Author.When
		}
		return nil
	})

	for name, time := range times {
		err := os.Chtimes(filepath.Join(path, name), time, time)
		if err != nil {
			log.Println("[ERROR] failed to get commit stats:", err.Error())
		}
	}
	return nil
}

func updateRepo(path string, c <-chan time.Time) {
	for range c {
		repo, err := git.PlainOpen(path)
		if err != nil {
			log.Fatal("[ERROR] failed to open repo:", err.Error())
		}
		wt, err := repo.Worktree()
		if err != nil {
			log.Fatal("[ERROR] failed to get repo worktree:", err.Error())
		}
		err = wt.Pull(&git.PullOptions{
			Force: true,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			log.Println("[ERROR] failed to pull repo:", err.Error())
			continue
		}
		err = updateModTime(path, repo)
		if err != nil {
			log.Println("[ERROR] failed to update file modtime:", err.Error())
			continue
		}
	}
}

func setup() {
	path := filepath.Join(config.Get().StaticPath, "blog")
	if isNotExist(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatal(err.Error())
		}
		url := config.Get().BlogRemoteURL
		proxy := transport.ProxyOptions{URL: os.Getenv("http_proxy")}
		repo, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:          url,
			ProxyOptions: proxy,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		err = updateModTime(path, repo)
		if err != nil {
			log.Println("[ERROR] failed to update file modtime:", err.Error())
		}
	}

	duration, err := time.ParseDuration(config.Get().BlogFetchInterval)
	if err != nil {
		log.Fatal("[ERROR] failed to parse fetch interval:", err.Error())
	}
	ticker := time.NewTicker(duration)
	go updateRepo(path, ticker.C)
}
