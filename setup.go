package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func copyBrToBin() {
	path, err := os.Executable()
	if err != nil {
		panic(err.Error())
	}
	stringSlice := strings.Split(os.ExpandEnv("$PATH"), ":")
	for _, v := range stringSlice {
		if strings.Contains(v, os.ExpandEnv("$HOME")) {
			if v+"/br" == path {
				break
			}
			cpCmd := exec.Command("cp", "-rf", path, v+"/br")
			err := cpCmd.Run()
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("Was copied to " + v + "/br. Now you can use the 'br' command from the console")
			break
		}
	}
}
