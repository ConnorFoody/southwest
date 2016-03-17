package southwest

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

	// try to parse out the boarding groups
	var arbJSON map[string]interface{}
	err := json.Unmarshal(buff.Bytes(), &arbJSON)
	if err != nil {
		fmt.Println("err decoding boardingpasses:", err)
		return
	}

	group, groupOK := arbJSON["boarding_group"]
	position, posOK := arbJSON["boarding_position"]

	br.ok = groupOK && posOK

	if br.ok {
		br.group = group.(string)
		br.position = position.(string)
		fmt.Println("GOOD CHECKIN! group:", br.group, "pos:", br.position)
	}
}

// get the params for the boardingPass request
func (swr requestHandler) boardingPassParams(account Account) map[string]string {
	ret := swr.baseParams()
	ret["serviceID"] = "getallboardingpass"
	ret["firstName"] = account.FirstName
	ret["lastName"] = account.LastName
	ret["recordLocator"] = account.RecordLocator
	return ret
}
