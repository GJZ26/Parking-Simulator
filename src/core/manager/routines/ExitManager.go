package routines

import "time"

func ExitManager(exitControl chan bool, entranceControl chan bool) {
	time.Sleep(1 * time.Second)
	entranceControl <- true
	for {
		select {
		case <-exitControl:
			time.Sleep(900 * time.Millisecond)
			entranceControl <- true
		}
	}
}
