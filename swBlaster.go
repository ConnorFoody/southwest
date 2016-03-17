package main

import (
	"github.com/ConnorFoody/southwest/blaster"
	"time"
)

// SWBlaster builds a blast for a southwest checkin
type SWBlaster struct {
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
	account swAccount

	// for exiting externally
	closer chan error
}

// SetAccount gives the blaster the form data for checkin
func (b *SWBlaster) SetAccount(account swAccount) {
	b.account = account
}

// SetParams for the blast. Period is the time between requests.
// cover is the ammount of time after the checkin time to send requests.
// headstart accounts for network latency.
func (b *SWBlaster) SetParams(period, cover, headstart int) {
	b.blastPeriod = period
	b.cover = cover
	b.headstart = headstart
	b.numRequests = (b.headstart + b.cover) / b.blastPeriod
}

// SetTime that the blast will start (note that headstart adjusts this)
// string fmt is <month abriviation> date <time> <pm/am>
// example input date is: "jan 1 7:15 pm"
func (b *SWBlaster) SetTime(timeStr string) error {
	targetTime, err := time.Parse("Jan 2, 2006 at 3:04pm (PST)", timeStr)
	b.startTime = targetTime
	return err
}

// ScheduleBlast at the provided param times
// NOTE: this doesn't make a lick of sense, neeed to wrap the runner
// in an interface, the swBlaster is actually the blastBuilder and
// it helps set up the request factory (maybe it is the request factory?)
// this
func (b *SWBlaster) ScheduleBlast(blast blaster.Blaster,
	factory blaster.RequestFactory) {

	runTime := b.startTime.Add(-time.Duration(b.headstart) * time.Millisecond)
	interval := time.Duration(b.blastPeriod) * time.Millisecond

	blast.Fire(factory, b.numRequests, runTime, interval)
}
