package southwest

import (
	"bytes"
	"fmt"
	"net/http"
)

// response for the travelInfo. Travel info is used to establish the
// session
type travelInfoResponse struct {
}

func (cr *travelInfoResponse) Parse(response *http.Response) {
	var buff bytes.Buffer
	buff.ReadFrom(response.Body)

	fmt.Println("resp:", buff.String())
	// TODO: maybe error checking?
}

// sends the travel info request. This request establishes the session.
// TODO: get test data for this response.
func (swr requestHandler) doTravelInfo() (travelInfoResponse, error) {
	travelInfoParams := swr.travelInfoParams()
	travelInfoString := swr.paramToBody(travelInfoParams)
	travelInfoResponse := travelInfoResponse{}
	err := swr.fireRequest(&travelInfoResponse, travelInfoString)

	return travelInfoResponse, err
}

// get the params for the travelInfo request
func (swr requestHandler) travelInfoParams() map[string]string {
	ret := swr.baseParams()
	ret["serviceID"] = "getTravelInfo"
	return ret
}
