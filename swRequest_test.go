package main

import (
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
	swr := makeswRequestHandler()
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

	swr := makeswRequestHandler()
	swr.config.baseURI = ts.URL

	params := swr.checkinParams()
	paramStr := swr.paramToBody(params)

	resp := checkinResponse{}

	if err := swr.fireRequest(&resp, paramStr); err != nil {
		panic(err)
	}

	if resp.ok {
		t.Errorf("did not expect to get an OK response!\n")
	}
}
