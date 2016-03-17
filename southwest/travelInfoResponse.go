package southwest

import (
	"bytes"
	"fmt"
	"net/http"
)

// response for the travelInfo
type travelInfoResponse struct {
}

func (cr *travelInfoResponse) Parse(response *http.Response) {
	var buff bytes.Buffer
	buff.ReadFrom(response.Body)

	fmt.Println("resp:", buff.String())
}

// get the params for the travelInfo request
func (swr requestHandler) travelInfoParams() map[string]string {
	ret := swr.baseParams()
	ret["serviceID"] = "getTravelInfo"
	return ret
}
