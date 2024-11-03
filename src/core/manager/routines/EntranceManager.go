package routines

import "time"

type EntranceManager struct{}

func NewEntranceManager() *EntranceManager {
	return new(EntranceManager)
}

func (e *EntranceManager) Run(mutex chan bool) *EntranceManager {
	for {
		select {
		case <-mutex:
		default:
			mutex <- true
		}
		time.Sleep(1500 * time.Millisecond)
	}
}
