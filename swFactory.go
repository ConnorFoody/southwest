package main

import (
	"github.com/ConnorFoody/southwest/blaster"
)

type swCheckinFactory struct {
	account swAccount
	lock    blaster.BlastLock
	id      int
}

func makeCheckinFactory(account swAccount) swCheckinFactory {
	olock := &blaster.OnceBlastLock{}
	olock.Setup(make(chan blaster.RequestStatus))
	return swCheckinFactory{account: account, lock: olock, id: 0}
}

func (f *swCheckinFactory) GetNext() blaster.BlastRequest {
	f.id++
	return &swCheckinTask{account: f.account,
		lock: f.lock,
		id:   f.id}
}

var _ blaster.RequestFactory = (*swCheckinFactory)(nil)
