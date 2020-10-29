package main

import (
	"fmt"
	"os"

	"br/api"
	"br/config"
)

func main() {

	getVersion()

	copyBrToBin()

	configuration := config.GetConfig()
	if (configuration == config.Configuration{}) {
		config.CreateConfig()
		configuration = config.GetConfig()
	}

	api := api.Configuration{Config: configuration}

	isValid := api.CheckCredentials()
	if !isValid {
		fmt.Println("Wrong credentials")
		config.CreateConfig()
		configuration = config.GetConfig()
	}

	if len(os.Args) == 1 {
		api.GetIssues()
	} else {
		api.Checkout()
	}
}
