package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tkanos/gonfig"
	"golang.org/x/crypto/ssh/terminal"
)

const BR_CONFIG_PATH = "$HOME/.br"

type Configuration struct {
	URL      string
	TEAM     string
	LOGIN    string
	PASSWORD string
}

func GetConfig() Configuration {
	configuration := Configuration{}
	gonfig.GetConf(os.ExpandEnv(BR_CONFIG_PATH+"/config.json"), &configuration)
	return configuration
}

func CreateConfig() {
	createConfigDir()

	configuration := AskUser()
	file, err := json.MarshalIndent(configuration, "", " ")
	if err != nil {
		panic(err.Error())
	}
	f, err := os.Create(os.ExpandEnv(BR_CONFIG_PATH + "/config.json"))
	if err != nil {
		panic(err.Error())
	}
	_, err = f.Write(file)
	if err != nil {
		f.Close()
		panic(err.Error())
	}
}

func AskUser() Configuration {
	configuration := Configuration{}

	fmt.Print("Enter your Jira url (e.g. https://jira.companyname.com): ")
	fmt.Scanf("%s", &configuration.URL)
	fmt.Print("Enter your Jira team (e.g. JASMINCORE): ")
	fmt.Scanf("%s", &configuration.TEAM)
	fmt.Print("Enter your Jira login: ")
	fmt.Scanf("%s", &configuration.LOGIN)
	fmt.Print("Enter your Jira password: ")
	password_, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err.Error())
	}
	configuration.PASSWORD = string(password_)
	fmt.Println("")

	return configuration
}

func createConfigDir() {
	if _, err := os.Stat(os.ExpandEnv(BR_CONFIG_PATH)); os.IsNotExist(err) {
		err := os.Mkdir(os.ExpandEnv(BR_CONFIG_PATH), 0700)
		if err != nil {
			panic(err.Error())
		}
	}
}
