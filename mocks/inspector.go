package mocks

import (
	"github.com/ConnorFoody/southwest/blaster"
	"time"
)

// BlastInspector lets us inspect the blast
type BlastInspector struct {
	factory   blaster.RequestFactory
	count     int
	startTime time.Time
	period    time.Duration
}

// Fire the blast into the inspector's storage
func (b *BlastInspector) Fire(
	f blaster.RequestFactory,
	count int,
	startTime time.Time,
	period time.Duration) {

	b.factory = f
	b.count = count
	b.startTime = startTime
	b.period = period
}

// Get all the members of the blast inspector
func (b BlastInspector) Get() (blaster.RequestFactory,
	int, time.Time, time.Duration) {
	return b.factory, b.count, b.startTime, b.period
}
