package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// response for the boardingPass
type boardingPassResponse struct {
	status int
	ok     bool
}

func (br *boardingPassResponse) Parse(response *http.Response) {
	br.status = response.StatusCode

	br.ok = false
	var buff bytes.Buffer
	buff.ReadFrom(response.Body)
	fmt.Println("response is:", buff.String())

	// try to parse out the boarding groups
	var arbJSON map[string]interface{}
	err := json.Unmarshal(buff.Bytes(), &arbJSON)
	if err != nil {
		fmt.Println("err decoding boardingpasses:", err)
		return
	}

	// print the boarding pass info
	fmt.Println("group:", arbJSON["boarding_group"],
		"spot:", arbJSON["boarding_possition"])
}

// get the params for the boardingPass request
func (swr swRequestHandler) boardingPassParams() map[string]string {
	ret := swr.baseParams()
	ret["serviceID"] = "getallboardingpass"
	ret["firstName"] = swr.account.FirstName
	ret["lastName"] = swr.account.LastName
	ret["recordLocator"] = swr.account.RecordLocator
	return ret
}
