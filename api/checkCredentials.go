package api

import (
	"encoding/json"
)

func (c Configuration) CheckCredentials() bool {

	body := c.getIssuesBody()
	var myStoredVariable interface{}
	err := json.Unmarshal([]byte(string(body)), &myStoredVariable)
	if err != nil {
		return false
	}

	return true
}
