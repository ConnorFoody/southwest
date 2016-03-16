package main

import (
	"bytes"
	"fmt"
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

type swConfig struct {
	baseURI         string
	AppVersion      string
	userAgentString string
	AppID           string
	Channel         string
	Platform        string
	CacheID         string
}

func makeswConfig() swConfig {
	fmtString := "Southwest/%s CFNetwork/711.1.16 Darwin/14.0.0"
	appVersion := "2.10.1"

	return swConfig{
		baseURI:         "https://mobile.southwest.com/middleware/MWServlet",
		AppVersion:      appVersion,
		userAgentString: fmt.Sprintf(fmtString, appVersion),
		AppID:           "swa",
		Channel:         "rc",
		Platform:        "iPhone",
		CacheID:         "",
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

func makeswRequestHandler() swRequestHandler {
	account := makeswAccount("foo", "bar", "123abc")

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

	return ret
}
