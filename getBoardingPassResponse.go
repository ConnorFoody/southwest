package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// response for the boardingPass
type boardingPassResponse struct {
	status   int
	ok       bool
	group    string
	position string
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

	br.group = arbJSON["boarding_group"].(string)
	br.position = arbJSON["boarding_position"].(string)

	br.ok = br.group != "" && br.position != ""

	if br.ok {
		fmt.Println("GOOD CHECKIN! group:", br.group, "pos:", br.position)
	}
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
