package transcode_test

import (
	"errors"
	"testing"

	"github.com/Arsfiqball/talker/transcode"
)

type sampleStruct struct {
	foo string
}

func TestContainer(t *testing.T) {
	t.Run("set and get string", func(t *testing.T) {
		c := transcode.NewCtx()

		transcode.SetNamed(c, "foo", "bar")
		transcode.SetNamed(c, "num", 1)
		transcode.Set(c, sampleStruct{foo: "bar"})

		var (
			foo string
			num int
			str sampleStruct
		)

		err := errors.Join(
			transcode.GetNamed(c, "foo", &foo),
			transcode.GetNamed(c, "num", &num),
			transcode.Get(c, &str),
		)

		if err != nil {
			t.Fatal(err)
		}

		if foo != "bar" {
			t.Fatal("foo != bar")
		}

		if num != 1 {
			t.Fatal("num != 1")
		}

		if str.foo != "bar" {
			t.Fatal("str.foo != bar")
		}
	})
}

type sampleStringType string

func TestMiddleware(t *testing.T) {
	t.Run("sequential", func(t *testing.T) {
		callback := transcode.Sequential(
			func(ctx transcode.Ctx, next transcode.Next) error {
				transcode.Set(ctx, sampleStringType("bar"))
				return next()
			},
			func(ctx transcode.Ctx, next transcode.Next) error {
				transcode.SetNamed(ctx, "bar", "foo")
				return next()
			},
			func(ctx transcode.Ctx, next transcode.Next) error {
				var (
					sst sampleStringType
					bar string
				)

				err := errors.Join(
					transcode.Get(ctx, &sst),
					transcode.GetNamed(ctx, "bar", &bar),
				)

				if err != nil {
					return err
				}

				if sst != "bar" {
					t.Fatal("foo != bar")
				}

				if bar != "foo" {
					t.Fatal("bar != foo")
				}

				return next()
			},
		)

		err := callback(transcode.NewCtx())
		if err != nil {
			t.Fatal(err)
		}
	})
}
