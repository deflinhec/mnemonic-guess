package dispatch

import (
	"sync"
	"time"
)

type Option interface {
	apply(s *ManagerType)
}

type funcOption struct {
	f func(*ManagerType)
}

func (fdo *funcOption) apply(do *ManagerType) {
	fdo.f(do)
}

func newOption(f func(*ManagerType)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithWaitGroup(wg *sync.WaitGroup) Option {
	return newOption(func(m *ManagerType) {
		m.wg = wg
	})
}

func WithQueueLimit(limit int64) Option {
	return newOption(func(m *ManagerType) {
		m.jobQueue = make(chan Job, limit)
	})
}

func WithMaxWorker(workers uint16) Option {
	return newOption(func(m *ManagerType) {
		m.workerPool = make(chan chan Job, workers)
	})
}

func WithDemand(limit int64) Option {
	return newOption(func(m *ManagerType) {
		m.join = func(j Job) {
			if m.Jobs() > limit {
				j.Execute()
				m.done <- true
			} else if len(m.jobQueue) == cap(m.jobQueue) {
				j.Execute()
				m.done <- true
			} else {
				m.jobQueue <- j
			}
		}
	})
}

func WithBlock(limit int64) Option {
	return newOption(func(m *ManagerType) {
		m.join = func(j Job) {
			for m.Jobs() > limit {
				<-time.After(time.Millisecond)
			}
			m.jobQueue <- j
		}
	})
}
