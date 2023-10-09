package exco_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Arsfiqball/talker/exco"
)

type fakeSignal struct{}

func (fakeSignal) String() string {
	return "fake signal"
}

func (fakeSignal) Signal() {}

func httpGet(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
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

func TestProcess(t *testing.T) {
	t.Run("should start and stop process", func(t *testing.T) {
		var res struct {
			started bool
			stopped bool
		}

		proc := exco.Process{
			MonitorAddr: ":8086",
			Start: func(ctx context.Context) error {
				res.started = true
				return nil
			},
			Stop: func(ctx context.Context) error {
				res.stopped = true
				return nil
			},
		}

		sig := make(chan os.Signal, 1)

		go func() {
			time.Sleep(100 * time.Millisecond)

			if !res.started {
				t.Error("process should be started")
			}

			if res.stopped {
				t.Error("process should not be stopped yet")
			}

			time.Sleep(100 * time.Millisecond)

			err := httpGet("http://localhost:8086/live")
			if err != nil {
				t.Error(fmt.Errorf("live check not ok: %w", err))
			}

			err = httpGet("http://localhost:8086/ready")
			if err != nil {
				t.Error(fmt.Errorf("ready check not ok: %w", err))
			}

			time.Sleep(100 * time.Millisecond)

			sig <- fakeSignal{} // stop process
		}()

		exco.Serve(proc, sig)

		time.Sleep(100 * time.Millisecond)

		if !res.stopped {
			t.Error("process should be stopped")
		}
	})
}
