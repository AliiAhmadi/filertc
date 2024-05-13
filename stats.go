package main

import (
	"sync"
	"time"
)

type stats struct {
	lock      *sync.RWMutex
	nbBytes   uint64
	timeStart time.Time
	timeStop  time.Time

	timePause  time.Time
	timePaused time.Duration
}

func newStats() *stats {
	return &stats{
		lock: &sync.RWMutex{},
	}
}
