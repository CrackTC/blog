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
	// return exec.Command("sh", "-c", `git ls-files | while read file; do echo $file; touch -d $(git log --date=local -1 --format="@%ct" "$file") "$file"; done`).Run()
	bytes, err := exec.Command("git", "-C", path, "ls-files").Output()
	if err != nil {
		return err
	}

	files := string(bytes)

	for _, file := range strings.Split(files, "\n") {
		file := file
		go func() {
			timeBytes, err := exec.Command("git", "-C", path, "log", "--date=local", "-1", "--format=@%ct", file).Output()
			fmt.Println(file, string(timeBytes))
			if err != nil {
				return
			}
			timeSec, err := strconv.ParseInt(strings.TrimSpace(string(timeBytes)), 10, 64)
			if err != nil {
				return
			}
			time := time.Unix(timeSec, 0)
			os.Chtimes(filepath.Join(path, file), time, time)
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
