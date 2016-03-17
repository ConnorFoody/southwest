package main

import (
	"github.com/ConnorFoody/southwest/blaster"
)

type swCheckinFactory struct {
	account swAccount
	lock    blaster.BlastLock
	id      int
	config  swConfig
}

func makeCheckinFactory(account swAccount,
	config swConfig) swCheckinFactory {
	olock := &blaster.OnceBlastLock{}
	olock.Setup(make(chan blaster.RequestStatus))
	go olock.Run()

	return swCheckinFactory{account: account,
		lock:   olock,
		id:     0,
		config: config,
	}
}

func (f *swCheckinFactory) GetNext() blaster.BlastRequest {
	f.id++
	return &swCheckinTask{account: f.account,
		lock: f.lock,
		id:   f.id,
		swr:  makeswRequestHandler(f.config),
	}
}

var _ blaster.RequestFactory = (*swCheckinFactory)(nil)
