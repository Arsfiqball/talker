package powercode

import (
	"fmt"
	"runtime"
)

type RecoveredPanic struct {
	message string
	stack   []string
}

func (e RecoveredPanic) Message() string {
	return e.message
}

func (e RecoveredPanic) Stack() []string {
	return e.stack
}

func Recover(e *RecoveredPanic) {
	if e == nil {
		return
	}

	skip := 2
	depth := 10

	if r := recover(); r != nil {
		e.message = fmt.Sprintf("%v", r)

		for i := skip; i < depth-1; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}

			name := "unknown"

			fn := runtime.FuncForPC(pc)
			if fn != nil {
				name = fn.Name()
			}

			e.stack = append(e.stack, fmt.Sprintf("%s:%d %s", file, line, name))
		}
	}
}
