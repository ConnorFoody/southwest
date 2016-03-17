package main

import (
	"testing"
)

func TestCheckin(t *testing.T) {
	ts := buildJSONResponseServer("test_data/checkin.json")
	defer ts.Close()

	swr := buildTestRequestHandler(ts.URL)

	params := swr.checkinParams(MakeAccount("foo", "bar", "123"))
	paramStr := swr.paramToBody(params)

	resp := checkinResponse{}

	if err := swr.fireRequest(&resp, paramStr); err != nil {
		panic(err)
	}

	if resp.status != 200 {
		t.Errorf("expected: 200 but got: %d\n", resp.status)
	}

	if !resp.ok {
		t.Error("expected ok response!")
	}
}
