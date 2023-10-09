package exco_test

import (
	"context"
	"testing"

	"github.com/Arsfiqball/talker/exco"
)

func TestSequential(t *testing.T) {
	t.Run("should run all callbacks", func(t *testing.T) {
		var res struct {
			a int
			b int
			c int
		}

		cb := exco.Sequential(
			func(ctx context.Context) error {
				res.a += 1
				return nil
			},
			func(ctx context.Context) error {
				res.b += 1
				return nil
			},
			func(ctx context.Context) error {
				res.c += 1
				return nil
			},
		)

		err := cb(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if res.a != 1 {
			t.Errorf("res.a should be 1, got %d", res.a)
		}

		if res.b != 1 {
			t.Errorf("res.b should be 1, got %d", res.b)
		}

		if res.c != 1 {
			t.Errorf("res.c should be 1, got %d", res.c)
		}
	})
}

func TestParallel(t *testing.T) {
	t.Run("should run all callbacks", func(t *testing.T) {
		var res struct {
			a int
			b int
			c int
		}

		cb := exco.Parallel(
			func(ctx context.Context) error {
				res.a += 1
				return nil
			},
			func(ctx context.Context) error {
				res.b += 1
				return nil
			},
			func(ctx context.Context) error {
				res.c += 1
				return nil
			},
		)

		err := cb(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if res.a != 1 {
			t.Errorf("res.a should be 1, got %d", res.a)
		}

		if res.b != 1 {
			t.Errorf("res.b should be 1, got %d", res.b)
		}

		if res.c != 1 {
			t.Errorf("res.c should be 1, got %d", res.c)
		}
	})
}

func TestTimeout(t *testing.T) {
	t.Run("should run callback with timeout", func(t *testing.T) {
		var res struct {
			a int
		}

		cb := exco.Timeout(
			func(ctx context.Context) error {
				res.a += 1
				return nil
			},
			1,
		)

		err := cb(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if res.a != 1 {
			t.Errorf("res.a should be 1, got %d", res.a)
		}
	})
}

func TestRetry(t *testing.T) {
	t.Run("should run callback with retries", func(t *testing.T) {
		var res struct {
			a int
		}

		cb := exco.Retry(
			func(ctx context.Context) error {
				res.a += 1
				return nil
			},
			3,
			0,
		)

		err := cb(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if res.a != 1 {
			t.Errorf("res.a should be 1, got %d", res.a)
		}
	})
}

func TestIgnoreError(t *testing.T) {
	t.Run("should ignore error", func(t *testing.T) {
		var res struct {
			a int
		}

		cb := exco.IgnoreError(
			func(ctx context.Context) error {
				res.a += 1
				return nil
			},
		)

		err := cb(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if res.a != 1 {
			t.Errorf("res.a should be 1, got %d", res.a)
		}
	})
}
