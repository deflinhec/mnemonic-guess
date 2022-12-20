package mnemonic

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/foxnut/go-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"mnemonic.deflinhec.dev/internal/dispatch"
)

type ManagerEnum int

const (
	EXPAND ManagerEnum = iota
	MATCH
)

type FetcherType struct {
	mux      sync.Mutex
	found    atomic.Bool
	Iterates atomic.Int64
	result   PhraseSequence
	wg       *sync.WaitGroup
	manager  [2]dispatch.Manager
}

func (m *FetcherType) fetch(phrases PhraseSequence, cond condition) {
	if m.found.Load() {
		return
	} else if phrases.Len() == int(cond.maxphrases) {
		m.manager[MATCH].Join(dispatch.FuncJob(func() {
			possible := phrases.String()
			if bip39.IsMnemonicValid(possible) {
				if master, err := hdwallet.NewKey(
					hdwallet.Mnemonic(possible),
				); err == nil {
					defer m.Iterates.Add(1)
					wallet, _ := master.GetWallet(
						hdwallet.CoinType(hdwallet.USDT),
					)
					addr, _ := wallet.GetAddress()
					if addr == cond.address {
						m.mux.Lock()
						defer m.mux.Unlock()
						m.result = phrases
						m.found.Store(true)
					}
				}
			}
		}))
	} else {
		m.manager[EXPAND].Join(dispatch.FuncJob(func() {
			for _, phrase := range bip39.GetWordList() {
				m.fetch(phrases.Fill(phrase), cond)
				if m.found.Load() {
					break
				}
			}
			runtime.GC()
		}))
	}
}

func (m *FetcherType) Fetch(address string,
	phrases PhraseSequence, conds ...Condition) *FetcherType {
	cond := condition{maxphrases: 12, address: address}
	ConditionSequence(conds).apply(&cond)
	defer m.fetch(phrases, cond)
	m.found.Store(false)
	m.Iterates.Store(0)
	m.result = nil
	return m
}

func (m *FetcherType) Jobs() int64 {
	return m.manager[EXPAND].Jobs()
}

func (m *FetcherType) Wait() *FetcherType {
	m.wg.Wait()
	return m
}

func (m *FetcherType) Result() PhraseSequence {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.result
}

func (m *FetcherType) Found() bool {
	return m.found.Load()
}

func Fetcher() *FetcherType {
	wg := &sync.WaitGroup{}
	return &FetcherType{
		wg: wg,
		manager: [2]dispatch.Manager{
			dispatch.NewManager(
				dispatch.WithMaxWorker(12),
				dispatch.WithWaitGroup(wg),
				dispatch.WithQueueLimit(12),
				dispatch.WithDemand(20480),
			),
			dispatch.NewManager(
				dispatch.WithMaxWorker(48),
				dispatch.WithWaitGroup(wg),
				dispatch.WithQueueLimit(1024),
				dispatch.WithBlock(20480),
			),
		},
	}
}
