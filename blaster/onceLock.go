package blaster

// OnceBlastLock locks after the first task through
type OnceBlastLock struct {
	lock  chan RequestStatus
	close chan bool
}

// Run until close is called, put in its own goroutine
func (bl *OnceBlastLock) Run() {
	gotFirst := false
	for {
		select {
		case result := <-bl.lock:
			// if err false
			if result.Err != nil {
				go result.Handle(false)
				continue
			} // else no err

			// send true if first, else false
			go result.Handle(!gotFirst)
			gotFirst = true

		case <-bl.close:
			// TODO: find a way to do this more cleanly
			return
		}
	}
}

// GetChan from the bl
func (bl OnceBlastLock) GetChan() chan RequestStatus {
	return bl.lock
}

// Setup used by the bl
func (bl *OnceBlastLock) Setup(l chan RequestStatus) {
	bl.lock = l
	bl.close = make(chan bool)
}

// Close the blast loc
func (bl *OnceBlastLock) Close() {
	close(bl.close)

}
