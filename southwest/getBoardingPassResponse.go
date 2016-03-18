package southwest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	log.Println(buff.String())

	if err != nil {
		fmt.Println("err decoding boardingpasses:", err)
		return
	}

	mbp := arbJSON["mbpPassenger"].([]interface{})
	passMap := mbp[0].(map[string]interface{})

	// see golang issue 6842 for use of tmp vars
	group, groupOK := passMap["boardingroup_text"]
	position, posOK := passMap["position1_text"]
	br.group = group.(string)
	br.position = position.(string)

	br.ok = groupOK && posOK

	if br.ok {
		log.Println("GOOD CHECKIN! group:", br.group, "pos:", br.position)
	}
}

func (swr *requestHandler) getBoardingPass(account Account) (
	boardingPassResponse, error) {

	boardingParams := swr.boardingPassParams(account)
	boardingString := swr.paramToBody(boardingParams)

	boardingResp := boardingPassResponse{}
	err := swr.fireRequest(&boardingResp, boardingString)

	return boardingResp, err
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
