package powercode

import (
	"errors"
	"fmt"
	"runtime"
)

type Error struct {
	code       string
	message    string
	declaredAt string
	wrappedAt  string
	data       interface{}
	parent     error
}

func NewError(code string, defaultMessage string) Error {
	var caller string

	_, file, line, ok := runtime.Caller(1)

	if ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	return Error{code: code, message: defaultMessage, declaredAt: caller}
}

func (e Error) Wrap(err error) Error {
	var caller string

	_, file, line, ok := runtime.Caller(1)

	if ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	e.parent = err
	e.wrappedAt = caller

	return e
}

func (e Error) SetInfo(message string) Error {
	e.message = message

	return e
}

func (e Error) Info() string {
	return e.message
}

func (e Error) SetData(data interface{}) Error {
	e.data = data

	return e
}

func (e Error) Data() interface{} {
	return e.data
}

func (e Error) Is(target error) bool {
	if target == nil {
		return false
	}

	pocoErr, ok := target.(Error)

	if ok && e.code == pocoErr.code {
		return true
	}

	return false
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Unwrap() error {
	return e.parent
}

func TraceError(err error) []string {
	var result []string

	for err != nil {
		pocoErr, ok := err.(Error)

		if ok {
			place := pocoErr.declaredAt

			if pocoErr.wrappedAt != "" {
				place = pocoErr.wrappedAt
			}

			result = append(result, fmt.Sprintf("%s at %s: %s", pocoErr.code, pocoErr.message, place))
			err = pocoErr.parent
		} else {
			result = append(result, fmt.Sprintf("sentinel: %s", err))
			err = nil
		}
	}

	return result
}

func ErrorIsOneOf(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}
