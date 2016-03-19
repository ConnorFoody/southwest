package southwest

import (
	"github.com/ConnorFoody/southwest/blaster"
	"time"
)

// this lock allows a certain number of clients to go past the
// first stage then blocks all but one on the second stage
// TODO: look into making this arbitratily long.
type twoPhasePermisiveLock struct {
	// chan that requests come through
	lock chan blaster.RequestStatus

	// pools to hold the blocked requests
	// first is the starting barrier
	firstPool  chan blaster.RequestStatus
	secondPool chan blaster.RequestStatus
	close      chan struct{}
}

func (tpl twoPhasePermisiveLock) Run() {

}

// GetChan from the cl
func (tpl twoPhasePermisiveLock) GetChan() chan blaster.RequestStatus {
	return tpl.lock
}

// Setup used by the cl
func (tpl *twoPhasePermisiveLock) Setup(l chan blaster.RequestStatus) {
	tpl.lock = l
	tpl.close = make(chan struct{})
}

// Close the blast loc
func (tpl *twoPhasePermisiveLock) Close() {
	close(tpl.close)
	// sleep to let the close signal hit everything before closing the pool
	// this is sort of a race condition...
	time.Sleep(5 * time.Millisecond)
}

// TryClose checks if the lock is closed
func (tpl twoPhasePermisiveLock) TryClose() chan struct{} {
	return tpl.close
}

var _ blaster.BlastLock = (*twoPhasePermisiveLock)(nil)
