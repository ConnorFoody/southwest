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
	go blaster.FireBlast(mocks.MakeSimpleFactory(),
		6,
		runT,
		time.Duration(1*time.Millisecond),
		bl)

	time.Sleep(2 * time.Second)
}
