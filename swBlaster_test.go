package main

import (
	"fmt"
	"github.com/ConnorFoody/southwest/mocks"
	"testing"
)

func errcheck(err error) {
	if err != nil {
		panic(err)
	}
}

func TestTime(t *testing.T) {
	b := BlastScheduler{}
	errcheck(b.SetTime("mar 15, 2016 at 7:15pm (PST)"))
	b.SetParams(10, 50, 0)

	inspector := mocks.BlastInspector{}
	factory := MakeCheckinFactory(MakeAccount("foo", "bar", "123abc"), MakeConfig())
	b.ScheduleBlast(&inspector, &factory)

	_, _, runtime, _ := inspector.Get()
	fmt.Println("time is:", runtime)
}
