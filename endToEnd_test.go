package main

import (
	"bytes"
	"fmt"
	"github.com/ConnorFoody/southwest/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
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
	sampleBoardingData := loadSampleData("test_data/boardingpasses.json")

	// build test server
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var buff bytes.Buffer
			buff.ReadFrom(r.Body)

			if strings.Contains(buff.String(), "serviceID=flightcheckin_new") {
				fmt.Fprintln(w, sampleCheckinData)
			} else if strings.Contains(buff.String(),
				"serviceID=getallboardingpass") {
				fmt.Fprintln(w, sampleBoardingData)
			}
		}))
	defer ts.Close()

	// setup the blaster
	blastSched := BlastScheduler{}

	blastSched.SetParams(10, 100, 0)
	blastSched.SetTime(shortTimeFromNow())

	// TODO: make a real blaster
	blastFirer := mocks.SimpleBlaster{}

	// build and run
	config := MakeConfig()
	config.BaseURI = ts.URL
	factory :=
		MakeCheckinFactory(MakeAccount("foo", "bar", "123abc"), config)

	go blastSched.ScheduleBlast(&blastFirer, &factory)

	time.Sleep(time.Duration(3 * time.Second))
}
