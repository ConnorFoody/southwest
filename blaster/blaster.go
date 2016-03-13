package blaster

import (
	"time"
)

// BlastRequest is fired by the blaster.
type BlastRequest interface {
	Send(BlastLock)
}

// RequestFactory builds requests
type RequestFactory interface {
	GetNext() BlastRequest
}

// RequestStatus is the status of a request on the chan
type RequestStatus struct {
	Ok  chan bool
	Err error
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
}

// FireBlast shoots actual logins
// TODO: make this into a struct
func FireBlast(factory RequestFactory,
	count int,
	after time.Time,
	interval time.Duration,
	lock BlastLock) {

	waitDur := after.Sub(time.Now())
	<-time.After(waitDur)

	ticker := time.Tick(interval)
	for i := 0; i < count; i++ {
		req := factory.GetNext()
		go req.Send(lock)
		<-ticker
	}

}
