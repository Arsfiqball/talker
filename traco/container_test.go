package traco_test

import (
	"testing"

	"github.com/Arsfiqball/talker/traco"
)

type sampleCtx struct {
	foo string
}

func (s *sampleCtx) SetFoo(foo string) {
	s.foo = foo
}

func (s *sampleCtx) Foo() string {
	return s.foo
}

func TestTypedContainer(t *testing.T) {
	t.Run("pipe set and get string", func(t *testing.T) {
		middleware1 := func(ctx *sampleCtx, next func() error) error {
			ctx.SetFoo("bar")
			return next()
		}

		middleware2 := func(ctx *sampleCtx, next func() error) error {
			if ctx.Foo() != "bar" {
				t.Fatal("ctx.Foo() != bar")
			}

			return next()
		}

		traco.Pipe(middleware1, middleware2)(&sampleCtx{})
	})

	t.Run("switch case pipe", func(t *testing.T) {
		type ctxRoute1 struct {
			baz string
		}

		type ctxRoute2 struct {
			bar string
		}

		route := func(foo string) func(*sampleCtx) bool {
			return func(ctx *sampleCtx) bool {
				return ctx.Foo() == foo
			}
		}

		toRoute1 := func(ctx *sampleCtx) *ctxRoute1 {
			return &ctxRoute1{baz: ctx.Foo()}
		}

		toRoute2 := func(ctx *sampleCtx) *ctxRoute2 {
			return &ctxRoute2{bar: ctx.Foo()}
		}

		handleRoute1 := func(ctx *ctxRoute1, next func() error) error {
			if ctx.baz != "baz" {
				t.Fatal("ctx.baz != baz")
			}

			return next()
		}

		handleRoute2 := func(ctx *ctxRoute2, next func() error) error {
			if ctx.bar != "bar" {
				t.Fatal("ctx.bar != bar")
			}

			return next()
		}

		pipe := traco.Pipe(
			traco.Case(route("baz"), toRoute1, traco.Pipe(handleRoute1)),
			traco.Case(route("bar"), toRoute2, traco.Pipe(handleRoute2)),
		)

		pipe(&sampleCtx{foo: "baz"})
		pipe(&sampleCtx{foo: "bar"})
	})
}
