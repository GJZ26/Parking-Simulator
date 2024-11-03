package routines

import "time"

func ExitManager(exitControl chan bool, entranceControl chan bool) {
	time.Sleep(1 * time.Second)
	entranceControl <- true
	for {
		select {
		case <-exitControl:
			entranceControl <- true
		}
	}
}
