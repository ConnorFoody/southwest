package southwest

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// response for the checkin. Not sure where the timing actually matters.
// it is (probably) either in the checkin or the getBoardingPasses
type checkinResponse struct {
	status int
	ok     bool
}

func (cr *checkinResponse) Parse(response *http.Response) {
	cr.status = response.StatusCode

	// set status to false before reading it
	cr.ok = false

	var buff bytes.Buffer
	buff.ReadFrom(response.Body)

	var arbJSON map[string]interface{}
	err := json.Unmarshal(buff.Bytes(), &arbJSON)

	// TODO: have this write to some temp log
	log.Println("checkin json:", buff.String())

	if err != nil {
		log.Println("err decoding json:", err)
		return
	}

	// check that this has flight info
	_, cr.ok = arbJSON["output"]
}

// wraps building the params, firing the response
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
