package blaster

import (
	"time"
)

// BlastRequest is fired by the blaster.
type BlastRequest interface {
	Send()
}

// RequestFactory builds requests
type RequestFactory interface {
	GetNext() BlastRequest
}

// RequestStatus is the status of a request on the chan
type RequestStatus struct {
	Ok   chan bool
	Err  error
	UUID uint32
}

// Handle sending a response on the chan
func (rs RequestStatus) Handle(resp bool) {
	rs.Ok <- resp
}

// BlastLock controls how the requests get fired
type BlastLock interface {
	GetChan() chan RequestStatus
	Run()
	Setup(chan RequestStatus)
	Close()
	TryClose() chan struct{}
}

// Blaster fires off the request
type Blaster interface {
	// factory, # of, starttime, request period
	Fire(RequestFactory, int, time.Time, time.Duration)
}
