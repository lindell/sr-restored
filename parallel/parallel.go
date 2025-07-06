package parallel

import (
	"context"
	"fmt"
	"sync"
)

type SharedGetter[T any] struct {
	fetchFunc          func(context.Context, int) (T, error)
	currentFetchesLock sync.Mutex // To protect access to currentFetches
	currentFetches     map[int][]chan valueWithError[T]
}

type valueWithError[T any] struct {
	value T
	err   error
}

func NewSharedGetter[T any](fetchFunc func(context.Context, int) (T, error)) *SharedGetter[T] {
	return &SharedGetter[T]{
		fetchFunc:      fetchFunc,
		currentFetches: make(map[int][]chan valueWithError[T]),
	}
}

func (pf *SharedGetter[T]) Fetch(ctx context.Context, id int) (T, error) {
	// If a parallel fetch is already in progress for this ID, create a new channel
	// and wait for the result.
	pf.currentFetchesLock.Lock()
	if _, exists := pf.currentFetches[id]; exists {
		newChan := make(chan valueWithError[T])
		pf.currentFetches[id] = append(pf.currentFetches[id], newChan)
		pf.currentFetchesLock.Unlock()

		val := <-newChan
		return val.value, val.err
	}

	pf.currentFetches[id] = []chan valueWithError[T]{}
	pf.currentFetchesLock.Unlock()

	result, err := pf.runWithPanicRecover(ctx, id)

	pf.currentFetchesLock.Lock()
	for _, ch := range pf.currentFetches[id] {
		ch <- valueWithError[T]{value: result, err: err}
		close(ch)
	}
	delete(pf.currentFetches, id)
	pf.currentFetchesLock.Unlock()

	return result, err
}

func (pf *SharedGetter[T]) runWithPanicRecover(ctx context.Context, id int) (res T, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("panic recover: %x", recovered)
		}
	}()
	return pf.fetchFunc(ctx, id)
}
