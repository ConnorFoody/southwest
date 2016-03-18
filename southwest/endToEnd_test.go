package southwest

import (
	"bytes"
	"fmt"
	"github.com/ConnorFoody/southwest/mocks"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func shortTimeFromNow() string {
	// 100 ms from now
	fmtStr := "Jan 2 15:04:05 -0700 MST 2006"
	return time.Now().Add(time.Duration(100 * time.Millisecond)).Format(fmtStr)
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

	rand.Seed(42)
	// load up the base json docs
	sampleCheckinData := loadSampleData("test_data/checkin.json")
	sampleBoardingData := loadSampleData("test_data/boardingpasses.json")

	// run is used to check when the status is ready to run
	var run uint32
	run = 4

	// build test server
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			canRun := rand.Intn(2) == 0
			atomic.AddUint32(&run, 1)

			var buff bytes.Buffer
			buff.ReadFrom(r.Body)

			if strings.Contains(buff.String(), "serviceID=flightcheckin_new") && run > 2 {
				fmt.Fprintln(w, sampleCheckinData)
			} else if strings.Contains(buff.String(), "serviceID=getallboardingpass") && canRun {
				fmt.Fprintln(w, sampleBoardingData)
			}

			// throw in delay
			if rand.ExpFloat64() > -1.1 {
				time.Sleep(time.Duration(15 * time.Millisecond))
			}
		}))
	defer ts.Close()

	// setup the blaster
	blastSched := BlastScheduler{}

	blastSched.SetParams(10, 10, 0)
	fmt.Println("short time from now:", shortTimeFromNow())
	blastSched.SetTime(shortTimeFromNow())

	// TODO: make a real blaster
	blastFirer := mocks.SimpleBlaster{}

	// build and run
	config := MakeConfig()
	config.BaseURI = ts.URL
	factory :=
		MakeCheckinFactory(MakeAccount("foo", "bar", "123abc"), config)

	blastSched.ScheduleBlast(&blastFirer, &factory)

	//after := time.After(3 * time.Second)
	/*
		select {
		case <-factory.lock.TryClose():
			// OK
		case <-after:
			t.Error("Expected close sooner!")

		}
	*/

	time.Sleep(500 * time.Millisecond)

}
