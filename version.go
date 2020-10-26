package main

import (
	"fmt"
	"os"
)

func getVersion() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println("0.0.4")
		os.Exit(0)
	}
}
