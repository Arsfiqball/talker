package exco

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func HttpGetCheck(url string, timeout time.Duration) Callback {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return errors.New("status code is not 200")
		}

		return nil
	}
}
