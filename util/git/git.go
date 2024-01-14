package git

import (
	"os/exec"
)

func UpdateModTime(path string) error {
	cmd := exec.Command("sh", "-c", "cd "+path+"; git log --pretty=%at --name-status --reverse | perl -ane '($x,$f)=@F;next if !$x;$t=$x,next if !defined($f)||$s{$f};$s{$f}=utime($t,$t,$f),next if $x=~/[AM]/;'")
	return cmd.Run()
}

func PullRepo(path string) error {
	cmd := exec.Command("git", "-C", path, "pull", "--force")
	return cmd.Run()
}

func CloneRepo(url string, path string) error {
	cmd := exec.Command("git", "clone", url, path)
	return cmd.Run()
}
