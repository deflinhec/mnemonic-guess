package dispatch

import (
	"math"
	"sync"
	"sync/atomic"
)

type ManagerType struct {
	workerPool chan chan Job
	jobQueue   chan Job
	done       chan bool
	JobTotal   int64
	JobDone    int64
	wg         *sync.WaitGroup
}

func NewManager(opts ...Option) *ManagerType {
	m := &ManagerType{
		jobQueue:   make(chan Job),
		done:       make(chan bool),
		workerPool: make(chan chan Job, 8),
		wg:         &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt.apply(m)
	}
	m.Run()
	return m
}

func (m *ManagerType) Run() {
	defer m.dispatch()
	for i := 0; i < cap(m.workerPool); i++ {
		newWorker(m.workerPool, m.done).Start()
	}
	atomic.StoreInt64(&m.JobTotal, 0)
	atomic.StoreInt64(&m.JobDone, 0)
}

func (m *ManagerType) Join(work Job) {
	m.wg.Add(1)
	atomic.AddInt64(&m.JobTotal, 1)
	m.jobQueue <- work
}

func (m *ManagerType) Progress() float64 {
	if m.JobDone == 0 {
		if m.JobTotal == 0 {
			return float64(100.0)
		}
		return float64(0.0)
	}
	ratio := float64(m.JobDone) / float64(m.JobTotal)
	return math.Min(ratio*float64(100.0), float64(100.0))
}

func (m *ManagerType) Wait() {
	m.wg.Wait()
}

func (m *ManagerType) dispatch() {
	go func() {
		for {
			select {
			case job := <-m.jobQueue:
				// a job request has been received
				go func(job Job) {
					// try to obtain a worker job channel that is available.
					// this will block until a worker is idle
					jobChannel := <-m.workerPool

					// dispatch the job to the worker job channel
					jobChannel <- job
				}(job)
			case <-m.done:
				atomic.AddInt64(&m.JobDone, 1)
				m.wg.Done()
			}
		}
	}()
}
