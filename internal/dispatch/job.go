package dispatch

type FuncJob func()

func (fn FuncJob) Execute() {
	fn()
}
