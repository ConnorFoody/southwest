package southwest

import (
	"bytes"
	"encoding/json"
	"log"
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
	log.Println("checkin json:", buff.String())
	if err != nil {
		log.Println("err decoding json:", err)
		return
	}

	// check that this has flight info
	_, cr.ok = arbJSON["output"]
}

func (swr requestHandler) doCheckin(account Account) (
	checkinResponse, error) {

	checkinParams := swr.checkinParams(account)
	checkinString := swr.paramToBody(checkinParams)

	checkinResp := checkinResponse{}
	err := swr.fireRequest(&checkinResp, checkinString)

	return checkinResp, err

}

// get the params for the checkin request
func (swr requestHandler) checkinParams(account Account) map[string]string {
	ret := swr.baseParams()
	ret["serviceID"] = "flightcheckin_new"
	ret["firstName"] = account.FirstName
	ret["lastName"] = account.LastName
	ret["recordLocator"] = account.RecordLocator
	return ret
}
