package southwest

import (
	"fmt"
	"github.com/ConnorFoody/southwest/blaster"
)

// CheckinTask manages sending a southwest request
type CheckinTask struct {
	account Account
	lock    blaster.BlastLock
	id      uint32
	swr     requestHandler
}

// Send a sw checkin request
func (r *CheckinTask) Send() {
	// don't run another request on a closed chan
	select {
	case <-r.lock.TryClose():
		return
	default:
	}
	// build teh request handler
	swr := r.swr

	// get travel info
	travelInfoParams := swr.travelInfoParams()
	travelInfoString := swr.paramToBody(travelInfoParams)
	travelInfoResponse := travelInfoResponse{}
	err := swr.fireRequest(&travelInfoResponse, travelInfoString)

	// build and send first request
	checkinParams := swr.checkinParams(r.account)
	checkinString := swr.paramToBody(checkinParams)

	checkinResp := checkinResponse{}
	err = swr.fireRequest(&checkinResp, checkinString)
	if err != nil {
	} else if !checkinResp.ok {
		err = fmt.Errorf("something wrong with status\n")
	} else if checkinResp.status != 200 {
		fmt.Println("warn: checkin OK but bad status:", checkinResp.status)
	}

	// comm back to see if we can keep going
	statusMsg := blaster.RequestStatus{Ok: make(chan bool),
		Err:  err,
		UUID: r.id,
	}
	fmt.Println("id:", r.id, "trying to send", err == nil, "status")
	select {
	case r.lock.GetChan() <- statusMsg:
		// continue
	case <-r.lock.TryClose():
		fmt.Println("id:", r.id, "is closing")
		close(statusMsg.Ok)
		return
	}

	canContinue, ok := <-statusMsg.Ok

	if !canContinue || !ok || err != nil {
		fmt.Println("id:", r.id, "exiting on err:", err)
		close(statusMsg.Ok)
		return
	}
	fmt.Println("id:", r.id, "getting boarding pass")

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
		fmt.Println("id:", r.id, "exiting on err:", err)
		return
	}

	r.lock.Close()

	fmt.Println("success on id:", r.id)
}

var _ blaster.BlastRequest = (*CheckinTask)(nil)
