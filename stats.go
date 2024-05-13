package main

import (
	"fmt"
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

func (s *stats) start() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.timeStart.IsZero() {
		s.timeStart = time.Now()
	} else if !s.timePause.IsZero() {
		s.timePaused += time.Since(s.timePause)
		s.timePause = time.Time{}
	}
}

func (s *stats) pause() {
	s.lock.RLock()

	if s.timeStart.IsZero() || !s.timeStop.IsZero() {
		s.lock.RUnlock()
		return
	}
	s.lock.RUnlock()

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.timePause.IsZero() {
		s.timePause = time.Now()
	}
}

func (s *stats) addBytes(n uint64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.nbBytes += n
}

func (s *stats) bytes() uint64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.nbBytes
}

func (s *stats) bandwidth() float64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return (float64(s.nbBytes) / (1024 * 1024)) / s.duration().Seconds()
}

func (s *stats) duration() time.Duration {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.timeStart.IsZero() {
		return 0
	} else if s.timeStop.IsZero() {
		return time.Since(s.timeStart) - s.timePaused
	}

	return s.timeStop.Sub(s.timeStart) - s.timePaused
}

func (s *stats) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return fmt.Sprintf("%v bytes | %-v | %0.4f MB/s", s.bytes(), s.duration(), s.bandwidth())
}
