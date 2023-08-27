package git

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func UpdateModTime(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	logs, err := repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return err
	}

	times := make(map[string]time.Time)
	err = logs.ForEach(func(commit *object.Commit) error {
		stats, err := commit.Stats()
		if err != nil {
			return err
		}
		for _, stat := range stats {
			if _, ok := times[stat.Name]; ok {
				continue
			}
			times[stat.Name] = commit.Committer.When
		}
		return nil
	})
	if err != nil {
		return err
	}

	for name, t := range times {
		_ = os.Chtimes(filepath.Join(path, name), t, t)
	}
	return nil
}

func PullRepo(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Pull(&git.PullOptions{
		Force: true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil
}

func CloneRepo(url string, path string) error {
	proxy := transport.ProxyOptions{URL: os.Getenv("http_proxy")}
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:          url,
		ProxyOptions: proxy,
	})

	return err
}
