package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"br/config"

	"github.com/tidwall/gjson"
)

type Configuration struct {
	Config config.Configuration
}

func (c Configuration) getIssuesBody() []byte {
	url := c.Config.URL + "/rest/api/2/search?jql=assignee+%3D+" + c.Config.LOGIN + "+AND+status+in+%28Open%2C+%22In+Progress%22%2C+%22Under+Review%22%29"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.Config.LOGIN+":"+c.Config.PASSWORD)))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Sorry, can't get your current issues by url", url)
		fmt.Println("Config path: ", os.ExpandEnv(config.BR_CONFIG_PATH+"/config.json"))
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	return body
}

func (c Configuration) GetIssues() {

	printCurrentBranch()

	body := c.getIssuesBody()
	var myStoredVariable interface{}
	err := json.Unmarshal([]byte(string(body)), &myStoredVariable)
	if err != nil {
		println("\n\nIncorrect login/password\n")
		os.Exit(1)
	}

	issueNumExample := "2222"
	fmt.Println("Your 'open', 'in progress', 'under review' issues:")

	result := gjson.GetBytes(body, "issues")
	result.ForEach(func(key, value gjson.Result) bool {
		key = gjson.Get(value.String(), "key")
		summary := gjson.Get(value.String(), "fields.summary")
		name := gjson.Get(value.String(), "fields.issuetype.name")
		bugOrFeature := "feature"
		if name.String() == "Bug" {
			bugOrFeature = "bug"
		}

		println(key.String(), summary.String(), "("+bugOrFeature+")")

		re := regexp.MustCompile("[0-9]+")
		issueNumExample = re.FindString(key.String())

		return true // keep iterating
	})

	fmt.Println("\nUsage: br [num]. e.g. br " + issueNumExample + "\n")
}
