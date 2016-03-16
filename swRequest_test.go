package main

import (
	"fmt"
	"testing"
)

func TestUserAgent(t *testing.T) {
	config := makeswConfig()
	expected := "Southwest/2.10.1 CFNetwork/711.1.16 Darwin/14.0.0"

	if config.userAgentString != expected {
		t.Errorf("Excpected: %s, got: %s", expected, config.userAgentString)
	}
}

func TestCheckin(t *testing.T) {
	swr := makeswRequestHandler()

	checkinParams := swr.checkinParams()
	fmt.Println("params:", swr.paramToBody(checkinParams))
}
