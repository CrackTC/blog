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

func UpdateModTime(path string) error {
	configQuotePath()
	bytes, err := exec.Command("git", "-C", path, "ls-files").Output()
	if err != nil {
		return err
	}

	files := string(bytes)

	for _, file := range strings.Split(files, "\n") {
		file := file
		go func() {
			timeBytes, err := exec.Command("git", "-C", path, "log", "--date=local", "-1", "--format=%ct", file).Output()
			fmt.Println(file, string(timeBytes))
			if err != nil {
				fmt.Println("[ERROR] failed to get file mod time:", err.Error())
				return
			}
			timeSec, err := strconv.ParseInt(strings.TrimSpace(string(timeBytes)), 10, 64)
			if err != nil {
				fmt.Println("[ERROR] failed to parse file mod time:", err.Error())
				return
			}
			time := time.Unix(timeSec, 0)
			err = os.Chtimes(filepath.Join(path, file), time, time)
			if err != nil {
				fmt.Println("[ERROR] failed to update file mod time:", err.Error())
				return
			}
			fmt.Println(file, time)
		}()
	}

	return nil
}

func PullRepo(path string) error {
	return exec.Command("git", "-C", path, "pull", "--force").Run()
}

func CloneRepo(url string, path string) error {
	return exec.Command("git", "clone", url, path).Run()
}
