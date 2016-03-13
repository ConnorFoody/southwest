package main

import (
	"github.com/ConnorFoody/southwest/blaster"
)

type swRequest struct {
	// TODO: flusr out the interaction with the swBlaster
}

// TODO: replace this with a "swLock"
func (r *swRequest) Send(lock blaster.BlastLock) {

}

type swRequestFactory struct {
	// TODO: flush out the interaction between this and the swBlaster
}

func (f *swRequestFactory) GetNext() *swRequest {
	return &swRequest{}
}
