package powercode_test

import (
	"testing"

	"errors"

	"github.com/Arsfiqball/talker/powercode"
)

func TestError(t *testing.T) {
	t.Run("one level with standard error", func(t *testing.T) {
		var err error

		stdErr := errors.New("test")
		pocoErr := powercode.NewError("TEST", "test")

		err = pocoErr.Wrap(stdErr)

		if !errors.Is(err, stdErr) {
			t.Fatal("error is not stdErr")
		}

		if !errors.Is(err, pocoErr) {
			t.Fatal("error is not pocoErr")
		}
	})

	t.Run("two level", func(t *testing.T) {
		var err error

		namedErr1 := powercode.NewError("TEST1", "test 1")
		namedErr2 := powercode.NewError("TEST2", "test 2")

		err = namedErr2.Wrap(namedErr1)

		if !errors.Is(err, namedErr1) {
			t.Fatal("error is not namedErr1")
		}

		if !errors.Is(err, namedErr2) {
			t.Fatal("error is not namedErr2")
		}
	})

	t.Run("five level with standard error", func(t *testing.T) {
		var err error

		stdErr := errors.New("test")
		namedErr1 := powercode.NewError("TEST1", "test 1")
		namedErr2 := powercode.NewError("TEST2", "test 2")
		namedErr3 := powercode.NewError("TEST3", "test 3")
		namedErr4 := powercode.NewError("TEST4", "test 4")

		err = namedErr1.Wrap(stdErr)
		err = namedErr2.Wrap(err)
		err = namedErr3.Wrap(err)
		err = namedErr4.Wrap(err)

		if !errors.Is(err, stdErr) {
			t.Fatal("error is not stdErr")
		}

		if !errors.Is(err, namedErr1) {
			t.Fatal("error is not namedErr1")
		}

		if !errors.Is(err, namedErr2) {
			t.Fatal("error is not namedErr2")
		}

		if !errors.Is(err, namedErr3) {
			t.Fatal("error is not namedErr3")
		}

		if !errors.Is(err, namedErr4) {
			t.Fatal("error is not namedErr4")
		}

		// Test using verbose flag (-v) to print stack trace
		// for _, tc := range powercode.TraceError(err) {
		// 	t.Log(tc)
		// }
	})
}

func TestErrorIsOneOf(t *testing.T) {
	t.Run("one of multiple errors", func(t *testing.T) {
		var err error

		stdErr := errors.New("test")
		namedErr1 := powercode.NewError("TEST1", "test 1")
		namedErr2 := powercode.NewError("TEST2", "test 2")
		namedErr3 := powercode.NewError("TEST3", "test 3")
		namedErr4 := powercode.NewError("TEST4", "test 4")

		err = namedErr1.Wrap(stdErr)
		err = namedErr4.Wrap(err)

		if !powercode.ErrorIsOneOf(err, namedErr1, namedErr2, namedErr3) {
			t.Fatal("error namedErr1 should be correct")
		}

		if powercode.ErrorIsOneOf(err, namedErr2, namedErr3) {
			t.Fatal("error namedErr2 should be incorrect")
		}
	})
}
