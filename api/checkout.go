package api

import (
	"bytes"
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

func (c Configuration) Checkout() {

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

	currentBranch := getCurrentBranch()
	fmt.Println("Current branch:", currentBranch)

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
	var replacer = strings.NewReplacer(",", "-", " ", "-", ":", "-", "(", "-", ")", "-", "_", "-", "/", "-", "--", "-", "'", "-", ".", "-", "\n", "")
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
	cmd := exec.Command("git", "branch", branchName, "master")
	cmd.Run()
	cmd = exec.Command("git", "checkout", branchName)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("git", "checkout", branchName, fmt.Sprint(err)+": "+stderr.String())
		return
	}
	result := out.String()
	if result != "" {
		fmt.Println("Result: " + out.String())
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
		if len(splited) > 0 {
			var replacer = strings.NewReplacer("* ", "")
			branchName := replacer.Replace(splited[0])
			return branchName
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

func (c Configuration) PushTmpBranch() {
	currentBranch := getCurrentBranch()
	cmd := exec.Command("git", "push", "origin", currentBranch+":"+strings.Trim(os.Args[1], "\n"))
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("git", "push", "origin", currentBranch+":"+strings.Trim(os.Args[1], "\n"), fmt.Sprint(err)+": "+stderr.String())
		return
	}
	if out.String() != "" || stderr.String() != "" {
		fmt.Println("git", "push", "origin", currentBranch+":"+strings.Trim(os.Args[1], "\n"))
		fmt.Println("Result: " + out.String() + stderr.String())
	}
}
