package parallel

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func Test_newSharedGetter(t *testing.T) {
	fetchFunc, releaseLock, calls := createTester()
	sg := NewSharedGetter(fetchFunc)
	ctx := context.Background()

	go func() {
		val, err := sg.Fetch(ctx, 1)
		assert.Equal(t, val, 1, "Expected value to be 1")
		assert.NilError(t, err, "Expected no error for ID 1")
	}()
	go func() {
		val, err := sg.Fetch(ctx, 1)
		assert.Equal(t, val, 1, "Expected value to be 1")
		assert.NilError(t, err, "Expected no error for ID 1")
	}()
	go func() {
		val, err := sg.Fetch(ctx, 2)
		assert.Equal(t, val, 2, "Expected value to be 2")
		assert.NilError(t, err, "Expected no error for ID 2")
	}()
	go func() {
		_, err := sg.Fetch(ctx, 3)
		assert.Equal(t, err, testErr, "Expected error for ID 3")
	}()
	go func() {
		val, err := sg.Fetch(ctx, 2)
		assert.Equal(t, val, 2, "Expected value to be 2")
		assert.NilError(t, err, "Expected no error for ID 2")
	}()
	releaseLock()

	time.Sleep(10 * time.Millisecond) // Allow goroutines to finish

	assert.Equal(t, calls[1], 1, "Expected fetchFunc to be called once for ID 1")
	assert.Equal(t, calls[2], 1, "Expected fetchFunc to be called once for ID 2")
	assert.Equal(t, calls[3], 1, "Expected fetchFunc to be called once for ID 3")
	assert.Equal(t, len(calls), 3, "Expected fetchFunc to be called for 3 different IDs")
}

func createTester() (func(context.Context, int) (int, error), func(), map[int]int) {
	lock := &sync.Mutex{}
	calls := map[int]int{}
	lock.Lock()

	testFetchFunc := func(ctx context.Context, id int) (int, error) {
		calls[id]++
		lock.Lock()
		lock.Unlock()
		if id == 3 {
			return 0, testErr
		}
		return id, nil
	}

	releaseLock := func() {
		lock.Unlock()
	}

	return testFetchFunc, releaseLock, calls
}

var testErr = errors.New("test error")
