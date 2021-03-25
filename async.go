package iot

import (
	"sync"
	"time"
)

type AsyncResult interface {

	Wait() bool

	WaitTimeout(time.Duration) bool

	Done() <-chan struct{}

	Error() error
}

type baseAsyncResult struct {
	m        sync.RWMutex
	complete chan struct{}
	err      error
}

// Wait implements the Token Wait method.
func (b *baseAsyncResult) Wait() bool {
	<-b.complete
	return true
}

// WaitTimeout implements the Token WaitTimeout method.
func (b *baseAsyncResult) WaitTimeout(d time.Duration) bool {
	timer := time.NewTimer(d)
	select {
	case <-b.complete:
		if !timer.Stop() {
			<-timer.C
		}
		return true
	case <-timer.C:
	}

	return false
}

// Done implements the Token Done method.
func (b *baseAsyncResult) Done() <-chan struct{} {
	return b.complete
}

func (b *baseAsyncResult) flowComplete() {
	select {
	case <-b.complete:
	default:
		close(b.complete)
	}
}

func (b *baseAsyncResult) Error() error {
	b.m.RLock()
	defer b.m.RUnlock()
	return b.err
}

func (b *baseAsyncResult) setError(e error) {
	b.m.Lock()
	b.err = e
	b.flowComplete()
	b.m.Unlock()
}

type BooleanAsyncResult struct {
	baseAsyncResult

	result bool
}

func (bar *BooleanAsyncResult) Result() bool {
	bar.m.RLock()
	defer bar.m.RUnlock()

	return bar.result
}

func (bar *BooleanAsyncResult) SetResult(result bool) {
	bar.m.RLock()
	defer bar.m.RUnlock()

	bar.complete <- struct{}{}
	bar.result = result
	bar.complete <- struct{}{}
}
