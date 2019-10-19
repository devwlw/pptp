package router

import (
	"sync"
)

var isRunning bool

var lock sync.Mutex

func init() {
	lock = sync.Mutex{}
}

func GetRunningStatus() bool {
	return isRunning
}

func SetRunningStatus(status bool) {
	lock.Lock()
	isRunning = status
	lock.Unlock()
}
