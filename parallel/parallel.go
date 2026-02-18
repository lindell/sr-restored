package parallel

import (
	"context"
	"fmt"
	"sync"
)

type SharedGetter[K fmt.Stringer, V any] struct {
	fetchFunc          func(context.Context, K) (V, error)
	currentFetchesLock sync.Mutex // To protect access to currentFetches
	currentFetches     map[string][]chan valueWithError[V]
}

type valueWithError[V any] struct {
	value V
	err   error
}

func NewSharedGetter[K fmt.Stringer, V any](fetchFunc func(context.Context, K) (V, error)) *SharedGetter[K, V] {
	return &SharedGetter[K, V]{
		fetchFunc:      fetchFunc,
		currentFetches: make(map[string][]chan valueWithError[V]),
	}
}

func (pf *SharedGetter[K, V]) Fetch(ctx context.Context, key K) (V, error) {
	strKey := key.String()

	// If a parallel fetch is already in progress for this key, create a new channel
	// and wait for the result.
	pf.currentFetchesLock.Lock()
	if _, exists := pf.currentFetches[strKey]; exists {
		newChan := make(chan valueWithError[V])
		pf.currentFetches[strKey] = append(pf.currentFetches[strKey], newChan)
		pf.currentFetchesLock.Unlock()

		val := <-newChan
		return val.value, val.err
	}

	pf.currentFetches[strKey] = []chan valueWithError[V]{}
	pf.currentFetchesLock.Unlock()

	result, err := pf.runWithPanicRecover(ctx, key)

	pf.currentFetchesLock.Lock()
	for _, ch := range pf.currentFetches[strKey] {
		ch <- valueWithError[V]{value: result, err: err}
		close(ch)
	}
	delete(pf.currentFetches, strKey)
	pf.currentFetchesLock.Unlock()

	return result, err
}

func (pf *SharedGetter[K, V]) runWithPanicRecover(ctx context.Context, key K) (res V, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("panic recover: %x", recovered)
		}
	}()
	return pf.fetchFunc(ctx, key)
}
