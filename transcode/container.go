package transcode

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidType = errors.New("invalid type")
)

type Ctx struct {
	locals map[string]interface{}
}

func NewCtx() Ctx {
	return Ctx{
		locals: make(map[string]interface{}),
	}
}

func GetNamed[T any](ctx Ctx, name string, out *T) error {
	v, ok := ctx.locals[name]
	if !ok {
		return fmt.Errorf("error getting %s: %w", name, ErrNotFound)
	}

	typedValue, ok := v.(T)
	if !ok {
		return fmt.Errorf("error getting %s: %w", name, ErrInvalidType)
	}

	*out = typedValue

	return nil
}

func Get[T any](ctx Ctx, out *T) error {
	name := fmt.Sprintf("%T", *out)

	return GetNamed(ctx, name, out)
}

func SetNamed[T any](ctx Ctx, name string, in T) {
	ctx.locals[name] = in
}

func Set[T any](ctx Ctx, in T) {
	name := fmt.Sprintf("%T", in)

	SetNamed(ctx, name, in)
}

type Callback func(ctx Ctx) error

type Next func() error

type Middleware func(ctx Ctx, next Next) error

func PipeCtx(middlewares ...Middleware) Callback {
	return func(ctx Ctx) error {
		var next Next

		next = func() error {
			if len(middlewares) == 0 {
				return nil
			}

			middleware := middlewares[0]
			middlewares = middlewares[1:]

			return middleware(ctx, next)
		}

		return next()
	}
}

func CaseCtx(cond func(Ctx) bool, pipe Callback) Middleware {
	return func(ctx Ctx, next Next) error {
		// If the condition is not met, skip this middleware.
		if !cond(ctx) {
			return next()
		}

		return pipe(ctx)
	}
}
