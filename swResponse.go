package main

import (
	"bytes"
	"net/http"
	"strings"
)

type checkinResponse struct {
	body   string
	status int
}

func (cr *checkinResponse) Parse(response *http.Response) {
	cr.status = response.StatusCode

	var buff bytes.Buffer
	buff.ReadFrom(response.Body)

	cr.body = strings.TrimSpace(buff.String())
}

// get the params for the checkin request
func (swr swRequestHandler) checkinParams() map[string]string {
	ret := swr.baseParams()
	ret["serviceID"] = "flighcheckin_new"
	ret["firstName"] = swr.account.FirstName
	ret["lastName"] = swr.account.LastName
	ret["recordLocator"] = swr.account.RecordLocator
	return ret
}
