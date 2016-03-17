package southwest

import (
	"bytes"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// Account for the southwest checkin. Has name and record locator.
type Account struct {
	// data needed for confirmation
	FirstName     string
	LastName      string
	RecordLocator string
}

// MakeAccount using the provided params
func MakeAccount(first, last, confirm string) Account {
	return Account{FirstName: first,
		LastName:      last,
		RecordLocator: confirm}
}

// Config for the net requests
// this hits the thinclient, make sure you keep the client up to
// date with the changes in the api versions and such
// as of most recent writing the settings are from here:
// https://github.com/mwynholds/southy
type Config struct {
	BaseURI         string
	appVersion      string
	userAgentString string
	appID           string
	channel         string
	platform        string
	rcid            string
	cacheID         string
}

func MakeConfig() Config {

	return Config{
		BaseURI:         "https://mobile.southwest.com/middleware/MWServlet",
		appVersion:      "2.17.0",
		userAgentString: "Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.3 (KHTML, like Gecko) Version/8.0 Mobile/12A4345d Safari/600.1.4",
		appID:           "swa",
		channel:         "wap",
		platform:        "thinclient",
		cacheID:         "",
		rcid:            "spaiphone",
	}
}

// wrapper on teh http response
type swResponse interface {
	Parse(*http.Response)
}

// send the actual requests
// TODO: figure out if we need to do the "create session" stuff
type requestHandler struct {
	config Config

	// client requests are sent on
	client *http.Client
}

func makeRequestHandler(config Config) requestHandler {

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{Jar: jar}

	return requestHandler{config: config,
		client: client}
}

// fire off the request
func (swr *requestHandler) fireRequest(
	response swResponse, params string) error {

	request, err := http.NewRequest("POST",
		swr.config.BaseURI,
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

func (swr requestHandler) buildHeader(request *http.Request) {
	request.Header.Add("User-Agent", swr.config.userAgentString)
}

// convert the params to a string of form <key>=<value>&<key>=<val>
func (swr requestHandler) paramToBody(params map[string]string) string {
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
func (swr requestHandler) baseParams() map[string]string {
	ret := make(map[string]string)
	params := swr.config
	ret["appID"] = params.appID
	ret["channel"] = params.channel
	ret["platform"] = params.platform
	ret["cacheID"] = ""
	ret["appver"] = params.appVersion
	ret["rcid"] = params.rcid

	return ret
}
