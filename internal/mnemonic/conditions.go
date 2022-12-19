package mnemonic

type condition struct {
	maxphrases uint8
	address    string
}

type Condition interface {
	apply(s *condition)
}

type funcOption struct {
	f func(*condition)
}

func (fdo *funcOption) apply(do *condition) {
	fdo.f(do)
}

func newCondition(f func(*condition)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithMaxPhrases(max uint8) Condition {
	return newCondition(func(m *condition) {
		m.maxphrases = max
	})
}

type ConditionSequence []Condition

func (opts ConditionSequence) apply(s *condition) {
	for _, opt := range opts {
		opt.apply(s)
	}
}
