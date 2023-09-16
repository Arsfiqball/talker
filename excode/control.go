package excode

import (
	"context"
	"errors"
	"sync"
)

type Callback func(context.Context) error

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
