package article

import "os/exec"

func HtmlFromFile(path string) string {
	bytes, err := exec.Command("./sharpdown", path).Output()
	if err != nil {
		return "Something went wrong."
	}
	return string(bytes)
}
