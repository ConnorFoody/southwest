package southwest

import (
	"fmt"
	"github.com/ConnorFoody/southwest/blaster"
	"time"
)

// HoldLock blocks all but one task. If the unblocked task fails then
type HoldLock struct {
	lock  chan blaster.RequestStatus
	pool  chan blaster.RequestStatus
	close chan struct{}
}

// Run the holdlock until it is closed
func (hl *HoldLock) Run() {
	var topID uint32

	// nil when have a task
	poolTmp := hl.pool

	for {
		select {
		case result := <-hl.lock:

			// if we have a task and this tasks uuid != top task
			if poolTmp == nil && result.UUID != topID {
				go hl.spawnIntoPool(result)
				continue
			}

			if result.Err != nil {
				go result.Handle(false)
				continue
			}

			// else allow the task through
			poolTmp = nil
			topID = result.UUID
			go result.Handle(true)

		case result := <-poolTmp:
			poolTmp = nil
			topID = result.UUID
			go result.Handle(true)

		case <-hl.close:
			fmt.Println("closing hold lock!")
			return
		}
	}
}

func (hl HoldLock) spawnIntoPool(req blaster.RequestStatus) {
	if req.Err != nil {
		req.Handle(false)
		return
	}

	select {
	case hl.pool <- req:
		//done
	case <-hl.close:
		req.Handle(false)
		fmt.Println("closing id:", req.UUID)
	}
}

// GetChan from the cl
func (hl HoldLock) GetChan() chan blaster.RequestStatus {
	return hl.lock
}

// Setup used by the cl
// todo: refactor
func (hl *HoldLock) Setup(l chan blaster.RequestStatus) {
	hl.lock = l
	hl.close = make(chan struct{})
	hl.pool = make(chan blaster.RequestStatus)
}

// Close the blast loc
func (hl *HoldLock) Close() {
	close(hl.close)
	// sleep to let the close signal hit everything before closing the pool
	// this is sort of a race condition...
	time.Sleep(5 * time.Millisecond)
	close(hl.pool)
}

// TryClose checks if the lock is closed
func (hl HoldLock) TryClose() chan struct{} {
	return hl.close
}

var _ blaster.BlastLock = (*HoldLock)(nil)
