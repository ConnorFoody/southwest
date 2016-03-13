package blaster_test

import (
	"fmt"
	"github.com/ConnorFoody/southwest/blaster"
	"testing"
	"time"
)

type SimpleRequest struct {
	id int
}

func (s SimpleRequest) Run(lock blaster.BlastLock) {
	rq := blaster.RequestStatus{make(chan bool), nil}
	lock.GetChan() <- rq

	canContinue := <-rq.Ok
	fmt.Println("running:", s.id, "at time:", time.Now(),
		"continue:", canContinue)
}

func (s *SimpleRequest) Prep() {
	s.id++
}

func TestSimpleRequest(t *testing.T) {
	bl := &blaster.OnceBlastLock{}
	bl.Setup(make(chan blaster.RequestStatus))
	go bl.Run()

	runT := time.Now().Add(500 * time.Millisecond)
	go blaster.FireBlast(&SimpleRequest{0},
		6,
		runT,
		time.Duration(50*time.Millisecond),
		bl)

	time.Sleep(2 * time.Second)

}

func TestSimpleLock(t *testing.T) {
	bl := blaster.OnceBlastLock{}
	(&bl).Setup(make(chan blaster.RequestStatus))
	go bl.Run()

	//rq := blaster.RequestStatus{make(chan bool), nil}
	//bl.GetChan() <- rq
	bl.Close()
}
