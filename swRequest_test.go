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
	account := makeswAccount("foo", "bar", "123abc")
	swr := makeswRequestHandler(account)
	swr.config.baseURI = url
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

	account := makeswAccount("foo", "bar", "123abc")
	swr := makeswRequestHandler(account)
	swr.config.baseURI = ts.URL

	params := swr.checkinParams()
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
func (swr swRequestHandler) testValidParams() map[string]string {
	ret := swr.baseParams()
	ret["serviceID"] = "flighcheckin_new"
	ret["firstName"] = swr.account.FirstName
	ret["lastName"] = swr.account.LastName
	ret["recordLocator"] = swr.account.RecordLocator
	return ret
}

func TestSWEndpointsWork(t *testing.T) {
	t.Skip("explicitly enable calls on the actual site")
	fmt.Println("SHOULDN'T HIS THIS!")
	account := makeswAccount("Hackeleen", "Fudy", "8T9HIU")
	swr := makeswRequestHandler(account)

	params := swr.checkinParams()
	paramStr := swr.paramToBody(params)
	fmt.Println("str:", paramStr)

	resp := travelInfoResponse{}

	if err := swr.fireRequest(&resp, paramStr); err != nil {
		panic(err)
	}
}
