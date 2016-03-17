package southwest

import (
	"github.com/ConnorFoody/southwest/blaster"
)

// HoldLock "pauses" threads until it gets a confirm
type HoldLock struct {
	lock  chan blaster.RequestStatus
	pool  chan blaster.RequestStatus
	close chan bool
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
	hl.close = make(chan bool)
	hl.pool = make(chan blaster.RequestStatus)
}

// Close the blast loc
func (hl *HoldLock) Close() {
	close(hl.close)
	close(hl.pool)
}

var _ blaster.BlastLock = (*HoldLock)(nil)
