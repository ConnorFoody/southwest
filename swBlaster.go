package main

import (
	"github.com/ConnorFoody/southwest/blaster"
	"time"
)

type swBlaster struct {
	// in millieconds
	blastPeriod int

	// how many milliseconds after the start time to cover with
	// requests
	cover int

	// how early before the actual time should we send requests
	// TODO: adjust this for ping
	headstart int

	// how many requests to send, (cover + headstart)/period
	numRequests int

	startTime time.Time

	// for the user
	firstName   string
	lastName    string
	confirmCode string

	// for exiting externally
	closer chan error
}

// TODO: error checking for all of this
// SetAccount gives the blaster the form data for checkin
func (b *swBlaster) SetAccount(firstName, lastName, confirmCode string) {
	b.firstName = firstName
	b.lastName = lastName
	b.confirmCode = confirmCode
}

// SetParams for the blast. Period is the time between requests.
// cover is the ammount of time after the checkin time to send requests.
// headstart accounts for network latency.
func (b *swBlaster) SetParams(period, cover, headstart int) {
	b.blastPeriod = period
	b.cover = cover
	b.headstart = headstart
	b.numRequests = (b.headstart + b.cover) / b.blastPeriod
}

// SetTime that the blast will start (note that headstart adjusts this)
// string fmt is <month abriviation> date <time> <pm/am>
// example input date is: "jan 1 7:15 pm"
func (b *swBlaster) SetTime(timeStr string) error {
	targetTime, err := time.Parse("jan 1 7:15 pm", timeStr)
	b.startTime = targetTime
	return err
}

// ScheduleBlast at the provided param times
func (b *swBlaster) ScheduleBlast(
	factory blaster.RequestFactory,
	lock blaster.BlastLock) {

	runTime := b.startTime.Add(-time.Duration(b.headstart) * time.Millisecond)
	waitDur := runTime.Sub(time.Now())
	interval := time.Duration(b.blastPeriod) * time.Millisecond

	// wait until it is time to roll
	select {
	case <-time.After(waitDur):
		break
	case <-b.closer:
		return
	}

	// setup the intervals
	ticker := time.NewTicker(interval)
	for i := 0; i < b.cover; i++ {
		// get a request from the factory and send it out
		req := factory.GetNext()
		go req.Send(lock)

		// blocks until it is time for the next itr
		<-ticker.C
	}

	// clean up
	ticker.Stop()
}
