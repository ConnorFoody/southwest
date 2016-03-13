package mocks

import (
	"fmt"
	"github.com/ConnorFoody/southwest/blaster"
	"time"
)

// SimpleRequest is a mock request that just prints it's runtime
type SimpleRequest struct {
	id   int
	lock blaster.BlastLock
}

// Send (naively) mocks a net request
func (s SimpleRequest) Send() {
	rq := blaster.RequestStatus{Ok: make(chan bool), Err: nil}
	s.lock.GetChan() <- rq

	canContinue := <-rq.Ok
	fmt.Println("running:", s.id, "at time:", time.Now(),
		"continue:", canContinue)
}

// SimpleFactory is used to create SimpleRequests
type SimpleFactory struct {
	id   int
	lock blaster.BlastLock
}

// GetNext SimpleRequest
func (factory *SimpleFactory) GetNext() blaster.BlastRequest {
	factory.id++
	return SimpleRequest{id: factory.id, lock: factory.lock}
}

// MakeSimpleFactory creates a simple factory and lets us get
// around exporting fields
func MakeSimpleFactory(lock blaster.BlastLock) *SimpleFactory {
	return &SimpleFactory{id: 0, lock: lock}
}

// SimpleBlaster fires off the requests
type SimpleBlaster struct {
}

// Fire off a blast
func (s SimpleBlaster) Fire(factory blaster.RequestFactory,
	count int,
	after time.Time,
	interval time.Duration) {

	waitDur := after.Sub(time.Now())
	<-time.After(waitDur)

	ticker := time.Tick(interval)
	for i := 0; i < count; i++ {
		req := factory.GetNext()
		go req.Send()
		<-ticker
	}
}
