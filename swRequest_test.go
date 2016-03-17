package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	testAccount = swAccount{FirstName: "foo", LastName: "bar", RecordLocator: "abc123"}
)

func handleTestingError(err error) {
	if err != nil {
		panic(err)
	}
}

// loadSampleData
func loadSampleData(name string) string {
	file, err := os.Open(name)
	handleTestingError(err)
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	handleTestingError(err)

	return string(data)
}

func buildJSONResponseServer(jsonFile string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, loadSampleData(jsonFile))
		}))

}

func buildTestRequestHandler(url string) swRequestHandler {
	config := makeswConfig()
	config.baseURI = url
	swr := makeswRequestHandler(config)
	return swr
}

func TestFireRequest(t *testing.T) {
	expected := "hello world!"
	// start the server
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, expected)
		}))
	defer ts.Close()

	swr := buildTestRequestHandler(ts.URL)

	params := swr.checkinParams(testAccount)
	paramStr := swr.paramToBody(params)

	resp := checkinResponse{}

	if err := swr.fireRequest(&resp, paramStr); err != nil {
		panic(err)
	}

	// note that we aren't expecting anything
	if resp.ok {
		t.Errorf("did not expect to get an OK response!\n")
	}
}

type testswValidResponse struct {
}

func (swv *testswValidResponse) Parse(response *http.Response) {
	var buff bytes.Buffer
	buff.ReadFrom(response.Body)

	fmt.Println("resp from SW:", buff.String())
}

// get the params for the checkin request
func (swr swRequestHandler) testValidParams(account swAccount) map[string]string {
	ret := swr.baseParams()
	ret["serviceID"] = "flighcheckin_new"
	ret["firstName"] = account.FirstName
	ret["lastName"] = account.LastName
	ret["recordLocator"] = account.RecordLocator
	return ret
}

func TestSWEndpointsWork(t *testing.T) {
	t.Skip("explicitly enable calls on the actual site")
	fmt.Println("SHOULDN'T HIS THIS!")
	account := makeswAccount("Hackeleen", "Fudy", "8T9HIU")
	swr := makeswRequestHandler(makeswConfig())

	params := swr.checkinParams(account)
	paramStr := swr.paramToBody(params)
	fmt.Println("str:", paramStr)

	resp := travelInfoResponse{}

	if err := swr.fireRequest(&resp, paramStr); err != nil {
		panic(err)
	}
}
