package mnemonic

type Worker interface {
	Jobs() int64
	Progress() float64
}

func (m *FetcherType) Worker(e ManagerEnum) Worker {
	return m.manager[e]
}
