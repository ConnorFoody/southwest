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

	rand.Seed(42)
	// load up the base json docs
	sampleCheckinData := loadSampleData("test_data/checkin.json")
	sampleBoardingData := loadSampleData("test_data/boardingpasses.json")

	// build test server
	var run uint32
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			canRun := rand.Intn(2) == 0

			atomic.AddUint32(&run, 1)
			fmt.Println("run:", run)

			var buff bytes.Buffer
			buff.ReadFrom(r.Body)

			if strings.Contains(buff.String(), "serviceID=flightcheckin_new") && canRun {
				fmt.Fprintln(w, sampleCheckinData)
			} else if strings.Contains(buff.String(), "serviceID=getallboardingpass") && canRun && run > 2 {
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

	after := time.After(3 * time.Second)
	select {
	case <-factory.lock.TryClose():
		// OK
	case <-after:
		t.Error("Expected close sooner!")
	}

}
