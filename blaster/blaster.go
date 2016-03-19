package blaster

import (
	"time"
)

// BlastRequest is fired by the blaster.
type BlastRequest interface {
	Send()
}

// RequestFactory builds the requests that make up a blast.
type RequestFactory interface {
	// GetNext provides the next blast request. It is up to the factory
	// to make the requests ready to fire upon return.
	GetNext() BlastRequest

	// Done checks if the requests have told the factory they are done.
	// TODO: find a cleaner way to propogate the close signal.
	Done() chan struct{}
}

// RequestStatus is the status of a request on the chan. A status message
// is sent to the master, who can reply over the Ok chan
type RequestStatus struct {
	// Ok is the send on a send/recieve struct. The process is told to
	// either exit or continue. Only the "master" (eg a lock) should write
	// to this and only the task should read from it. Whomever "owns" the
	// status closes this chan.
	Ok chan bool

	// Err is any error that the process wants to communicate back to the
	// lock.
	Err error
	// UUID of the task, should be assigned by the master.
	UUID uint32
}

// Handle sending a response on the chan. Sort of a lazy helper
func (rs RequestStatus) Handle(resp bool) {
	rs.Ok <- resp
}

// BlastLock controls how the requests get fired by acting as a dynamic
// barrier. Procs communicate their status (RequestStatus) to the lock
// and the lock responds when the task can continue. Look up "go channel
// axioms" by dave channey (last name?) to get a better idea of how this
// works.
type BlastLock interface {
	// GetChan for sending RequestStatus messages to the lock. This
	// function does not block and is threadsafe, but the channel reads may
	// block. The "client" should check that the lock has not closed.
	GetChan() chan RequestStatus

	// Run the lock until a close signal is sent. This will block.
	Run()

	// Setup the lock with the request chan it will use. Pass in the chan
	// in case we ever want to use multiple locks.
	Setup(chan RequestStatus)

	// Send a close signal to the lock. The lock will take care of cleaning
	// up anything it owns. Each task should use a select to check if the
	// close signal has been sent. I am not sure if this approach creates
	// a race condition or not. To be safe locks should put a wait between
	// sending the close signal and closing the read chan.
	Close()

	// TryClose acts as the close signal. If the lock is closed reads
	// will return immediately. Otherwise reads will block. This should be
	// used in any select that uses the main comms chan. Use anon structs
	// because they take up no mem. Mechanism is just a close(chan).
	TryClose() chan struct{}
}

// Blaster manages firing off the requests. The implementations are what
// you would expect, it is only an interface to make mock testing easier.
type Blaster interface {
	// Fire runs the request sequence. Starting at <start time> <N> tasks
	// are sent from the <factory> every <duration> period.
	Fire(RequestFactory, int, time.Time, time.Duration)
}
