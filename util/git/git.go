package git

import (
	"os/exec"
)

func ConfigQuotePath() error {
	return exec.Command("git", "config", "--global", "core.quotepath", "false").Run()
}

func UpdateModTime(path string) error {
	return exec.Command("sh", "-c", "cd "+path+"; git log --pretty=%at --name-status --reverse | perl -ane '($x,$f)=@F;next if !$x;$t=$x,next if !defined($f)||$s{$f};$s{$f}=utime($t,$t,$f),next if $x=~/[AM]/;'").Run()
}

func PullRepo(path string) error {
	return exec.Command("git", "-C", path, "pull", "--force").Run()
}

func CloneRepo(url string, path string) error {
	return exec.Command("git", "clone", url, path).Run()
}
