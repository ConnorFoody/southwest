package main

import (
	"bufio"
	"fmt"
	"github.com/ConnorFoody/southwest/mocks"
	"github.com/ConnorFoody/southwest/southwest"
	"os"
	"strings"
	"time"
)

func shortTimeFromNow() string {
	// 100 ms from now
	fmtStr := "Jan 2 15:04:05 -0700 MST 2006"
	return time.Now().Add(time.Duration(100 * time.Millisecond)).Format(fmtStr)
}

func main() {
	layoutStr := "Jan 1, 2016 at 1:00pm (PST)"
	fmt.Println("welcome to the SW checkin bomber")
	fmt.Println("FORMAT:", layoutStr)
	//fmt.Print("enter time: ")

	// get the input line
	cmdLine := bufio.NewReader(os.Stdin)
	/*timeStr, err := cmdLine.ReadString('\n')
	timeStr = strings.Trim(timeStr, "\n")
	if err != nil {
		err = fmt.Errorf("error parsing text \"%s\", error: %s\n", timeStr, err)
		panic(err)
	}
	*/
	timeStr := shortTimeFromNow()
	fmt.Println("tstr:", timeStr)

	// build sched
	sched := southwest.BlastScheduler{}
	if err := sched.SetTime(timeStr); err != nil {
		panic(err)
	}

	// for now just use a hardcoded set
	sched.SetParams(10, 10, 0)

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

	// build the actual date
	factory := southwest.MakeCheckinFactory(account, southwest.MakeConfig())
	blaster := mocks.SimpleBlaster{}

	sched.ScheduleBlast(&blaster, &factory)

	<-factory.Done()
}
