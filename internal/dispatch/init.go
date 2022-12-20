package dispatch

type Job interface {
	Execute()
}

type Manager interface {
	Run()

	Join(Job)

	Progress() float64

	Jobs() int64
}
