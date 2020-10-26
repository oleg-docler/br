package api

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

const BRANCH_NAME_LIMIT = 59

func (c Configuration) CreateBranch() {

	issueID := getIssueId()

	issueName := c.Config.TEAM + "-" + issueID

	body := c.getIssueDetails(issueName)

	newBranchName := generateBranchName(issueName, body)

	existingBranchName := strings.Trim(getExistingBranch(issueName), "\n ")
	if existingBranchName != "" {
		gitCheckout(existingBranchName)
	} else {
		gitCheckout(newBranchName)
	}

	printCurrentBranch()

	printAnotherBranches(issueID)
}

func (c Configuration) getIssueDetails(issueName string) []byte {
	req, err := http.NewRequest("GET", c.Config.URL+"/rest/api/2/issue/"+issueName, nil)
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.Config.LOGIN+":"+c.Config.PASSWORD)))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)

	return body
}

func generateBranchName(issueName string, body []byte) string {
	summary := gjson.GetBytes(body, "fields.summary")
	name := gjson.GetBytes(body, "fields.issuetype.name")
	bugOrFeature := "feature"
	if name.String() == "Bug" {
		bugOrFeature = "bugfix"
	}
	var replacer = strings.NewReplacer(" ", "-", ":", "-", "(", "-", ")", "-", "_", "-", "/", "-", "--", "-", "'", "-", ".", "-", "\n", "")
	branchName := replacer.Replace(issueName + "-" + summary.String())
	branchName = bugOrFeature + "/" + replacer.Replace(branchName)
	branchName = strings.Trim(branchName, "-")
	if len(branchName) > BRANCH_NAME_LIMIT {
		branchName = branchName[0:BRANCH_NAME_LIMIT]
	}
	fmt.Println(branchName)
	return branchName
}

func gitCheckout(branchName string) {
	cpCmd := exec.Command("git", "checkout", "-b", branchName, "master")
	err := cpCmd.Run()
	if err != nil {
		cpCmd := exec.Command("git", "checkout", branchName)
		cpCmd.Run()
	}
}

func getIssueId() string {
	issueId := os.Args[1]
	if _, err := strconv.Atoi(issueId); err != nil {
		fmt.Println("Wrong task Id")
		os.Exit(0)
	}
	return issueId
}

func getExistingBranch(issueId string) string {
	out, err := exec.Command("bash", "-c", "git branch | grep "+issueId).Output()
	splited := strings.Split(string(out), "\n")
	if err == nil {
		if len(splited) > 1 {
			return string(out)
		}
	}
	return ""
}

func printAnotherBranches(issueId string) {
	out, err := exec.Command("bash", "-c", "git branch | grep "+issueId).Output()
	splited := strings.Split(string(out), "\n")
	if err == nil {
		if len(splited) > 2 {
			fmt.Println("By the way, you have these branches:")
			fmt.Printf("%s", out)
		}
	}
}
