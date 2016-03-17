package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// response for the checkin
type checkinResponse struct {
	status int
	ok     bool
}

func (cr *checkinResponse) Parse(response *http.Response) {
	cr.status = response.StatusCode

	cr.ok = false
	var buff bytes.Buffer
	buff.ReadFrom(response.Body)

	var arbJSON map[string]interface{}

	err := json.Unmarshal(buff.Bytes(), &arbJSON)
	if err != nil {
		fmt.Println("err decoding json:", err)
		return
	}

	// check that this has flight info
	cr.ok = arbJSON["output"] != nil
}

// get the params for the checkin request
func (swr swRequestHandler) checkinParams(account swAccount) map[string]string {
	ret := swr.baseParams()
	ret["serviceID"] = "flightcheckin_new"
	ret["firstName"] = account.FirstName
	ret["lastName"] = account.LastName
	ret["recordLocator"] = account.RecordLocator
	return ret
}
