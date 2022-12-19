package dispatch

import "sync"

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

func WithMaxJobs(limit int64) Option {
	return newOption(func(m *ManagerType) {
		m.jobQueue = make(chan Job, limit)
	})
}

func WithMaxWorker(workers uint16) Option {
	return newOption(func(m *ManagerType) {
		m.workerPool = make(chan chan Job, workers)
	})
}
