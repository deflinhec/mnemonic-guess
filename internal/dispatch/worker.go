package dispatch

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	done       chan bool
	quit       chan bool
}

func newWorker(pool chan chan Job, done chan bool) Worker {
	return Worker{
		WorkerPool: pool,
		done:       done,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				job.Execute()
				w.done <- true
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
