package main

import (
	"github.com/ConnorFoody/southwest/blaster"
)

type swCheckinFactory struct {
	// TODO: flush out the interaction between this and the swBlaster
}

func (f *swCheckinFactory) GetNext() blaster.BlastRequest {
	return &swCheckinTask{}
}

// swCheckinTask manages sending a southwest request
type swCheckinTask struct {
	account swAccount
}

// Send a sw checkin request
func (r *swCheckinTask) Send() {

}

// perform the actual checkin
func (r *swCheckinTask) doCheckin() {

}

var _ blaster.RequestFactory = (*swCheckinFactory)(nil)
var _ blaster.BlastRequest = (*swCheckinTask)(nil)
