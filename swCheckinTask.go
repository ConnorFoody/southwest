package main

import (
	"fmt"
	"github.com/ConnorFoody/southwest/blaster"
)

// swCheckinTask manages sending a southwest request
type swCheckinTask struct {
	account swAccount
	lock    blaster.BlastLock
	id      int
	swr     swRequestHandler
}

// Send a sw checkin request
func (r *swCheckinTask) Send() {
	// build teh request handler
	swr := r.swr

	// build and send first request
	checkinParams := swr.checkinParams(r.account)
	checkinString := swr.paramToBody(checkinParams)

	checkinResp := checkinResponse{}
	err := swr.fireRequest(&checkinResp, checkinString)
	if err != nil {
	} else if !checkinResp.ok {
		err = fmt.Errorf("something wrong with status\n")
	} else if checkinResp.status != 200 {
		fmt.Println("warn: checkin OK but bad status:", checkinResp.status)
	}

	// comm back to see if we can keep going
	statusMsg := blaster.RequestStatus{Ok: make(chan bool), Err: err}
	r.lock.GetChan() <- statusMsg

	canContinue := <-statusMsg.Ok

	if !canContinue || err != nil {
		fmt.Println("id:", r.id, "exitint on err:", err)
		close(statusMsg.Ok)
		return
	}

	// if we can keep going then go get the boarding passes
	boardingParams := swr.boardingPassParams(r.account)
	boardingString := swr.paramToBody(boardingParams)

	boardingResp := boardingPassResponse{}

	err = swr.fireRequest(&boardingResp, boardingString)

	if err != nil {
	} else if !boardingResp.ok {
		err = fmt.Errorf("something wrong with status\n")
	} else if boardingResp.status != 200 {
		fmt.Println("warn: checkin OK but bad status:", checkinResp.status)
	}

	if err != nil {
		close(statusMsg.Ok)
		fmt.Println("id:", r.id, "exitint on err:", err)
		return
	}

	fmt.Println("success on id:", r.id)

}

var _ blaster.BlastRequest = (*swCheckinTask)(nil)
