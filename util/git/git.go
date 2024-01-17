package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func configQuotePath() error {
	return exec.Command("git", "config", "--global", "core.quotepath", "false").Run()
}

func unquotePath(path string) string {
	if path[0] == '"' && path[len(path)-1] == '"' {
		return strings.Replace(path[1:len(path)-1], `\`, "", -1)
	}
	return path
}

func UpdateModTime(path string) error {
	configQuotePath()
	bytes, err := exec.Command("git", "-C", path, "ls-files").Output()
	if err != nil {
		return err
	}

	files := strings.Split(string(bytes), "\n")

	for _, file := range files {
		if file == "" {
			continue
		}
		file := unquotePath(file)
		timeBytes, err := exec.Command("git", "-C", path, "log", "--date=local", "-1", "--format=%ct", file).Output()
		if err != nil {
			return err
		}
		timeSec, err := strconv.ParseInt(strings.TrimSpace(string(timeBytes)), 10, 64)
		if err != nil {
			return err
		}
		time := time.Unix(timeSec, 0)
		err = os.Chtimes(filepath.Join(path, file), time, time)
		if err != nil {
			return err
		}
		fmt.Println(file, time)
	}

	return nil
}

func PullRepo(path string) error {
	return exec.Command("git", "-C", path, "pull", "--force").Run()
}

func CloneRepo(url string, path string) error {
	return exec.Command("git", "clone", url, path).Run()
}
