package main

import (
	"bytes"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type swAccount struct {
	// data needed for confirmation
	FirstName     string
	LastName      string
	RecordLocator string
}

func makeswAccount(first, last, confirm string) swAccount {
	return swAccount{FirstName: first,
		LastName:      last,
		RecordLocator: confirm}
}

// base config for all requests
// this hits the thinclient, make sure you keep the client up to
// date with the changes in the api versions and such
// as of most recent writing the settings are from here:
// https://github.com/mwynholds/southy
type swConfig struct {
	baseURI         string
	AppVersion      string
	userAgentString string
	AppID           string
	Channel         string
	Platform        string
	rcid            string
	CacheID         string
}

func makeswConfig() swConfig {

	return swConfig{
		baseURI:         "https://mobile.southwest.com/middleware/MWServlet",
		AppVersion:      "2.17.0",
		userAgentString: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.3 (KHTML, like Gecko) Version/8.0 Mobile/12A4345d Safari/600.1.4",
		AppID:           "swa",
		Channel:         "wap",
		Platform:        "thinclient",
		CacheID:         "",
		rcid:            "spaiphone",
	}
}

// wrapper on teh http response
type swResponse interface {
	Parse(*http.Response)
}

// send the actual requests
// TODO: figure out if we need to do the "create session" stuff
type swRequestHandler struct {
	config  swConfig
	account swAccount

	// client requests are sent on
	client *http.Client
}

func makeswRequestHandler(account swAccount) swRequestHandler {

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{Jar: jar}

	return swRequestHandler{config: makeswConfig(),
		account: account,
		client:  client}
}

// fire off the request
func (swr *swRequestHandler) fireRequest(
	response swResponse,
	params string) error {

	request, err := http.NewRequest("POST",
		swr.config.baseURI,
		strings.NewReader(params))

	if err != nil {
		return err
	}

	swr.buildHeader(request)

	httpResp, err := swr.client.Do(request)
	if err != nil {
		return err
	}

	response.Parse(httpResp)
	return httpResp.Body.Close()
}

func (swr swRequestHandler) buildHeader(request *http.Request) {
	request.Header.Add("User-Agent", swr.config.userAgentString)
}

// convert the params to a string of form <key>=<value>&<key>=<val>
func (swr swRequestHandler) paramToBody(params map[string]string) string {
	// use a byte buffer for effency
	var buffer bytes.Buffer

	first := true
	for key, val := range params {
		// only skip prefacing '&' for first itr
		if !first {
			buffer.WriteString("&")
		} else {
			first = false
		}

		buffer.WriteString(key)

		buffer.WriteString("=")
		buffer.WriteString(val)
	}

	return buffer.String()
}

// body params used in every call
func (swr swRequestHandler) baseParams() map[string]string {
	ret := make(map[string]string)
	params := swr.config
	ret["appID"] = params.AppID
	ret["channel"] = params.Channel
	ret["platform"] = params.Platform
	ret["cacheID"] = ""
	ret["appver"] = params.AppVersion
	ret["rcid"] = params.rcid

	return ret
}
