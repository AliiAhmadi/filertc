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

func (s *stats) stop() {
	s.lock.RLock()

	if s.timeStart.IsZero() {
		// Can't stop if not started
		s.lock.RUnlock()
		return
	}
	s.lock.RUnlock()

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.timeStop.IsZero() {
		s.timeStop = time.Now()
	}
}
