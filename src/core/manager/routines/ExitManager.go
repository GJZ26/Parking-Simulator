package routines

import "time"

func ExitManager(exitControl chan bool, entranceControl chan bool) {
	entranceControl <- true
	for {
		select {
		case <-exitControl:
			time.Sleep(200 * time.Millisecond)
			entranceControl <- true
		}
	}
}
