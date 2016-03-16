package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserAgent(t *testing.T) {
	config := makeswConfig()
	expected := "Southwest/2.10.1 CFNetwork/711.1.16 Darwin/14.0.0"

	if config.userAgentString != expected {
		t.Errorf("Excpected: %s, got: %s", expected, config.userAgentString)
	}
}

func TestCheckin(t *testing.T) {
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
	fmt.Println(resp.body, paramStr)

	if resp.body != expected {
		t.Errorf("expected: %s got: %s\n", expected, resp.body)
	}

}
