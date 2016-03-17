package main

import (
	"fmt"
	"github.com/ConnorFoody/southwest/blaster"
)

type swCheckinFactory struct {
	account swAccount
	lock    blaster.BlastLock
}

func (f *swCheckinFactory) GetNext() blaster.BlastRequest {
	return &swCheckinTask{account: f.account,
		lock: f.lock}
}

// swCheckinTask manages sending a southwest request
type swCheckinTask struct {
	account swAccount
	lock    blaster.BlastLock
}

// Send a sw checkin request
func (r *swCheckinTask) Send() {
	// build teh request handler
	swr := makeswRequestHandler(r.account)

	// build and send first request
	checkinParams := swr.checkinParams()
	checkinString := swr.paramToBody(checkinParams)

	checkinResp := checkinResponse{}
	err := swr.fireRequest(&checkinResp, checkinString)
	if err == nil {
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
		close(statusMsg.Ok)
		return
	}

	// if we can keep going then go get the boarding passes
	boardingParams := swr.boardingPassParams()
	boardingString := swr.paramToBody(boardingParams)

	boardingResp := boardingPassResponse{}

	err = swr.fireRequest(&boardingResp, boardingString)

	if err == nil {
	} else if !checkinResp.ok {
		err = fmt.Errorf("something wrong with status\n")
	} else if checkinResp.status != 200 {
		fmt.Println("warn: checkin OK but bad status:", checkinResp.status)
	}

	// update the status message
	statusMsg.Err = err

	r.lock.GetChan() <- statusMsg

	canContinue = <-statusMsg.Ok

	if !canContinue || err != nil {
		close(statusMsg.Ok)
		return
	}

	fmt.Println("success!")

}

// perform the actual checkin
func (r *swCheckinTask) doCheckin() {

}

var _ blaster.RequestFactory = (*swCheckinFactory)(nil)
var _ blaster.BlastRequest = (*swCheckinTask)(nil)
