package blaster

import (
	"time"
)

// BlastRequest is a netrequest to be fired by the blaster.
// Treat this like a Request factory
type BlastRequest interface {
	Run(BlastLock)
	Prep()
}

// FireBlast shoots actual logins
func FireBlast(req BlastRequest,
	count int,
	after time.Time,
	interval time.Duration,
	lock BlastLock) {

	waitDur := after.Sub(time.Now())
	<-time.After(waitDur)

	ticker := time.Tick(interval)
	for i := 0; i < count; i++ {
		req.Prep()
		go req.Run(lock)
		<-ticker
	}

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

// OnceBlastLock locks after the first task through
type OnceBlastLock struct {
	lock  chan RequestStatus
	close chan bool
}

// Run until close is called, put in its own goroutine
func (bl *OnceBlastLock) Run() {
	gotFirst := false
	for {
		select {
		case result := <-bl.lock:
			// if err false
			if result.Err != nil {
				go result.Handle(false)
				continue
			} // else no err

			// send true if first, else false
			go result.Handle(!gotFirst)
			gotFirst = true

		case <-bl.close:
			// TODO: find a way to do this more cleanly
			return
		}
	}
}

// GetChan from the bl
func (bl OnceBlastLock) GetChan() chan RequestStatus {
	return bl.lock
}

// Setup used by the bl
func (bl *OnceBlastLock) Setup(l chan RequestStatus) {
	bl.lock = l
	bl.close = make(chan bool)
}

// Close the blast loc
func (bl *OnceBlastLock) Close() {
	close(bl.close)

}
