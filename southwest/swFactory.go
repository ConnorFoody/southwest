package southwest

import (
	"github.com/ConnorFoody/southwest/blaster"
	"sync/atomic"
)

// CheckinFactory builds checkin tasks
type CheckinFactory struct {
	account Account
	lock    blaster.BlastLock
	id      uint32
	config  Config
}

// MakeCheckinFactory with account and config
func MakeCheckinFactory(account Account,
	config Config) CheckinFactory {
	lock := &HoldLock{}
	lock.Setup(make(chan blaster.RequestStatus))
	go lock.Run()

	return CheckinFactory{account: account,
		lock:   lock,
		id:     0,
		config: config,
	}
}

// GetNext request
func (f *CheckinFactory) GetNext() blaster.BlastRequest {
	atomic.AddUint32(&f.id, 1)
	return &CheckinTask{account: f.account,
		lock: f.lock,
		id:   f.id,
		swr:  makeRequestHandler(f.config),
	}
}

// Done with the factory
func (f *CheckinFactory) Done() chan struct{} {
	return f.lock.TryClose()
}

var _ blaster.RequestFactory = (*CheckinFactory)(nil)
