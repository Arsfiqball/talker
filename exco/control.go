package exco

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Callback is a function that can be used as a step in a workflow.
type Callback func(context.Context) error

// Sequential runs all callbacks sequentially.
func Sequential(callbacks ...Callback) Callback {
	return func(ctx context.Context) error {
		for _, callback := range callbacks {
			if err := callback(ctx); err != nil {
				return err
			}
		}

		return nil
	}
}

// Parallel runs all callbacks in parallel.
func Parallel(callbacks ...Callback) Callback {
	return func(ctx context.Context) error {
		var wg sync.WaitGroup

		errChan := make(chan error, len(callbacks))

		for _, callback := range callbacks {
			wg.Add(1)

			go func(w *sync.WaitGroup, callback Callback) {
				defer w.Done()
				errChan <- callback(ctx)
			}(&wg, callback)
		}

		wg.Wait()
		errs := []error{}

		for i := 0; i < len(callbacks); i++ {
			if err := <-errChan; err != nil {
				errs = append(errs, err)
			}
		}

		return errors.Join(errs...)
	}
}

// Timeout runs callback with timeout.
func Timeout(callback Callback, timeout time.Duration) Callback {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return callback(ctx)
	}
}

// Retry runs callback with retries.
func Retry(callback Callback, retries int, delay time.Duration) Callback {
	return func(ctx context.Context) error {
		var err error

		for i := 0; i < retries; i++ {
			err = callback(ctx)
			if err == nil {
				return nil
			}

			time.Sleep(delay)
		}

		return err
	}
}

// IgnoreError runs callback and ignore the error.
func IgnoreError(callback Callback) Callback {
	return func(ctx context.Context) error {
		_ = callback(ctx)
		return nil
	}
}

// Atomic runs commit and rollback in sequence.
func Atomic(commit Callback, rollback Callback) Callback {
	return func(ctx context.Context) error {
		err := commit(ctx)
		if err != nil {
			return rollback(ctx)
		}

		return nil
	}
}
