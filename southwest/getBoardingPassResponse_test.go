package southwest

import (
	"testing"
)

func TestBoardingPass(t *testing.T) {
	ts := buildJSONResponseServer("test_data/boardingpasses.json")
	defer ts.Close()

	swr := buildTestRequestHandler(ts.URL)

	params := swr.boardingPassParams(MakeAccount("foo", "bar", "123"))
	paramStr := swr.paramToBody(params)

	resp := boardingPassResponse{}

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
