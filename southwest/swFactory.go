package southwest

import (
	"github.com/ConnorFoody/southwest/blaster"
	"sync/atomic"
)

// CheckinFactory builds checkin tasks
type CheckinFactory struct {
	account Account           // acc for checkin
	lock    blaster.BlastLock // lock for handling requests
	config  Config            // network config for swr
	id      uint32            // counter to give uuids
}

// MakeCheckinFactory with account and config
func MakeCheckinFactory(account Account,
	config Config) CheckinFactory {

	// build and start the lock
	lock := &HoldLock{}
	lock.Setup(make(chan blaster.RequestStatus))
	go lock.Run()

	return CheckinFactory{account: account,
		lock:   lock,
		id:     0,
		config: config,
	}
}

// GetNext request.
// TODO: look into pre-building these
func (f *CheckinFactory) GetNext() blaster.BlastRequest {
	// atomic mostly for fun, don't think there would be unsafe access
	atomic.AddUint32(&f.id, 1)

	// TODO: make the swr come in pre build
	return &CheckinTask{account: f.account,
		lock: f.lock,
		id:   f.id,
		swr:  makeRequestHandler(f.config),
	}
}

// Done with the factory. Just falls through to the lock.
func (f *CheckinFactory) Done() chan struct{} {
	return f.lock.TryClose()
}

var _ blaster.RequestFactory = (*CheckinFactory)(nil)
