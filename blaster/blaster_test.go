package blaster_test

import (
	"github.com/ConnorFoody/southwest/blaster"
	"github.com/ConnorFoody/southwest/mocks"
	"testing"
	"time"
)

func TestSimpleRequest(t *testing.T) {
	bl := &blaster.OnceBlastLock{}
	bl.Setup(make(chan blaster.RequestStatus))
	go bl.Run()

	runT := time.Now().Add(50 * time.Millisecond)
	roundT := runT.Round(time.Duration(1 * time.Second))

	blast := mocks.SimpleBlaster{}
	go blast.Fire(mocks.MakeSimpleFactory(bl),
		6,
		roundT,
		time.Duration(100*time.Millisecond))

	time.Sleep(2 * time.Second)
}
