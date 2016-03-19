package southwest

import (
	"fmt"
	"github.com/ConnorFoody/southwest/blaster"
	"log"
)

// CheckinTask manages sending a southwest request
type CheckinTask struct {
	account Account
	lock    blaster.BlastLock
	id      uint32
	swr     requestHandler
}

// NOTE: this time there was a network brownout. We only
// hit the server at 4hz, so I don't think it was a rate issue.
// I think that either southwest got wise to what we were doing
// or something else happened. We lost the logfiles, but I think
// it was a tcp or dhcp timeout issue. Still got a19, not sure if
// the script of my own checkin got it though. I scheduled this
// run 1.5s before the checkin time based on a ping guestimate.
//
// TODO: be more sneaky:
// To be more sneaky we can mix up the user agents, versions,
// platforms etc. I think there may be more endpoints, we need to
// capture requests from something like the southwest app though.
// Loading the app into xcode is a big, big hassle. Should check
// if a real phone can be loaded in. We can also distribute it
// like we did with opinionated (may even be able to use that
// infastructure)
// NOTE: I am not sure that being sneaky is a great approach. I think
// that they may be more unhappy with some of the sneaky stuff than they
// would be with high request volume. It may be an easier flag to raise
// as well. However, to get the kind of request volume we have you would
// need several devices, so the sneaky stuff might be more natural. I
// don't think going to a bunch of different ips like we did with
// opinionated would make them happy.
//
// TODO: be more gentle with southwest:
// Each request establishes a session with the server before it
// can check in. The lifespan on these connections is long enough
// that we can prebuild them.Each connection could have a setup
// barrier that gets closed when the factory goes into use. Problem
// with this is multiple calls to send. A http.client factory (ie
// a swRequest factory) might work, or we could just go with a
// setup method.
// We may be able to guess the ping off the session init, but
// that is longer term

// Send manages a southwest checkin request.
func (r *CheckinTask) Send() {
	log.Println("starting task", r.id, "!")

	// don't run another request on a closed chan
	// TODO: look into having the factory handle this with nil nexts
	select {
	case <-r.lock.TryClose():
		return
	default:
	}

	// get the request handler
	swr := r.swr

	// establish the session with the travel info
	_, err := swr.doTravelInfo()

	// build and send first request
	checkinResp, err := swr.doCheckin(r.account)

	if err != nil {
	} else if !checkinResp.ok {
		err = fmt.Errorf("something wrong with status\n")
	} else if checkinResp.status != 200 {
		log.Println("warn: checkin OK but bad status:", checkinResp.status)
	}

	// report status and check if we can keep going
	statusMsg := blaster.RequestStatus{Ok: make(chan bool),
		Err:  err,
		UUID: r.id,
	}

	// do comms with the lock
	select {
	case r.lock.GetChan() <- statusMsg:
		// continue
	case <-r.lock.TryClose():
		log.Println("id:", r.id, "is closing")
		close(statusMsg.Ok)
		return
	}

	// check if we can continue. ok should indicate if the chan is
	// closed or not so we don't need to worry about default values
	canContinue, ok := <-statusMsg.Ok

	if !canContinue || !ok || err != nil {
		log.Println("id:", r.id, "exiting on err:", err)
		close(statusMsg.Ok)
		return
	}
	log.Println("id:", r.id, "getting boarding pass")

	// if we can keep going then go get the boarding passes
	boardingResp, err := swr.doBoardingPass(r.account)

	if err != nil {
	} else if !boardingResp.ok {
		err = fmt.Errorf("something wrong with status\n")
	} else if boardingResp.status != 200 {
		log.Println("warn: checkin OK but bad status:", checkinResp.status)
	}

	if err != nil {
		close(statusMsg.Ok)
		log.Println("id:", r.id, "exiting on err:", err)
		return
	}

	// the first task to make it all the way through should close
	// everything down
	r.lock.Close()

	log.Println("success on id:", r.id)
}

var _ blaster.BlastRequest = (*CheckinTask)(nil)
