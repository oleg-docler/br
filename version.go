package main

import (
	"fmt"
	"os"
)

func getVersion() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println("0.0.5")
		os.Exit(0)
	}
}
