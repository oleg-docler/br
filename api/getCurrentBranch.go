package api

import (
	"os/exec"
	"strings"
)

func getCurrentBranch() string {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err == nil {
		return strings.Trim(string(out), "\n")
	}
	return ""
}
