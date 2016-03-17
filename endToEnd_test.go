package main

import (
	"fmt"
	"github.com/ConnorFoody/southwest/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func shortTimeFromNow() string {
	// 100 ms from now
	return time.Now().Add(time.Duration(100 * time.Millisecond)).Format("Jan 2, 2013 at 3:04pm (PST)")
}

func TestShortTimeFromNow(t *testing.T) {
	tStr := shortTimeFromNow()
	tar, err := time.Parse("Jan 2, 2013 at 3:04pm (PST)", tStr)
	handleTestingError(err)

	if tar.Add(-time.Duration(100 * time.Millisecond)).After(time.Now()) {
		t.Error("expected longer time!")
	}
}

func TestEndToEnd(t *testing.T) {
	// load up the base json docs
	sampleCheckinData := loadSampleData("test_data/checkin.json")

	// build test server
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, sampleCheckinData)
		}))
	defer ts.Close()

	// setup the blaster
	blastSched := SWBlaster{}

	blastSched.SetParams(10, 100, 0)
	blastSched.SetAccount(makeswAccount("foo", "bar", "123abc"))
	blastSched.SetTime(shortTimeFromNow())

	// TODO: make a real blaster
	blastFirer := mocks.SimpleBlaster{}

	// build and run
	factory := makeCheckinFactory(makeswAccount("foo", "bar", "123abc"))
	go blastSched.ScheduleBlast(&blastFirer, &factory)
}
