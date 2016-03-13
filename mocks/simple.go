package mocks

import (
	"fmt"
	"github.com/ConnorFoody/southwest/blaster"
	"time"
)

// SimpleRequest is a mock request that just prints it's runtime
type SimpleRequest struct {
	id int
}

// Send (naively) mocks a net request
func (s SimpleRequest) Send(lock blaster.BlastLock) {
	rq := blaster.RequestStatus{Ok: make(chan bool), Err: nil}
	lock.GetChan() <- rq

	canContinue := <-rq.Ok
	fmt.Println("running:", s.id, "at time:", time.Now(),
		"continue:", canContinue)
}

// SimpleFactory is used to create SimpleRequests
type SimpleFactory struct {
	id int
}

// GetNext SimpleRequest
func (factory *SimpleFactory) GetNext() blaster.BlastRequest {
	factory.id++
	return SimpleRequest{factory.id}
}

// MakeSimpleFactory creates a simple factory and lets us get
// around exporting fields
func MakeSimpleFactory() *SimpleFactory {
	return &SimpleFactory{0}
}
