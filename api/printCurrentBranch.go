package api

import (
	"fmt"
	"os/exec"
)

func printCurrentBranch() {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err == nil {
		fmt.Println("Current branch:", string(out))
	}
}
