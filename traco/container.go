package traco

func Pipe[T any](middlewares ...func(T, func() error) error) func(ctx T) error {
	return func(ctx T) error {
		var next func() error

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

func Case[FROM any, TO any](cond func(FROM) bool, typeSwitcher func(FROM) TO, pipe func(TO) error) func(FROM, func() error) error {
	return func(from FROM, next func() error) error {
		// If the condition is not met, skip this middleware.
		if !cond(from) {
			return next()
		}

		return pipe(typeSwitcher(from))
	}
}
