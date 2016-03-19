package main

import (
	"bufio"
	"fmt"
	"github.com/ConnorFoody/southwest/mocks"
	"github.com/ConnorFoody/southwest/southwest"
	"os"
	"strings"
)

func main() {
	// TODO:make the time more intuitive...
	fmtStr := "Jan 2 15:04:05 -0700 PDT 2006"
	fmt.Println("welcome to the SW checkin bomber")
	fmt.Println("FORMAT:", fmtStr)
	fmt.Print("enter time: ")

	// read from the cmd line
	cmdLine := bufio.NewReader(os.Stdin)

	// get the input line
	timeStr, err := cmdLine.ReadString('\n')
	timeStr = strings.Trim(timeStr, "\n")
	if err != nil {
		err = fmt.Errorf("error parsing text \"%s\", error: %s\n", timeStr, err)
		panic(err)
	}

	// build sched
	sched := southwest.BlastScheduler{}
	if err := sched.SetTime(timeStr); err != nil {
		panic(err)
	}

	// for now just use a hardcoded set
	sched.SetParams(250, 250, 000)

	// build account
	account := southwest.Account{}

	fmt.Print("enter full name: ")
	name, err := cmdLine.ReadString('\n')
	name = strings.Trim(name, "\n")
	if err != nil {
		panic(err)
	}

	names := strings.Split(name, " ")
	if len(names) != 2 {
		panic("expected only first and last names, but got more")
	}
	account.FirstName = names[0]
	account.LastName = names[1]

	fmt.Print("enter record locator: ")
	code, err := cmdLine.ReadString('\n')
	account.RecordLocator = code

	fmt.Printf("first: %s last: %s, confirm: %s\n", account.FirstName, account.LastName, account.RecordLocator)

	// build and submit the blast
	factory := southwest.MakeCheckinFactory(account, southwest.MakeConfig())
	blaster := mocks.SimpleBlaster{}

	sched.ScheduleBlast(&blaster, &factory)

	// wait until the factory is all done
	<-factory.Done()
}
