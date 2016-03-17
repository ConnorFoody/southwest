package main

import (
	"github.com/ConnorFoody/southwest/blaster"
)

type CheckinFactory struct {
	account Account
	lock    blaster.BlastLock
	id      int
	config  Config
}

func MakeCheckinFactory(account Account,
	config Config) CheckinFactory {
	olock := &blaster.OnceBlastLock{}
	olock.Setup(make(chan blaster.RequestStatus))
	go olock.Run()

	return CheckinFactory{account: account,
		lock:   olock,
		id:     0,
		config: config,
	}
}

func (f *CheckinFactory) GetNext() blaster.BlastRequest {
	f.id++
	return &CheckinTask{account: f.account,
		lock: f.lock,
		id:   f.id,
		swr:  makeRequestHandler(f.config),
	}
}

var _ blaster.RequestFactory = (*CheckinFactory)(nil)
