package poco_test

import (
	"testing"

	"github.com/Arsfiqball/talker/poco"
)

func funcThatPanics() {
	panic("test")
}

type someProxyInterface interface {
	SomeProxyMethod()
}

type someProxyStruct struct{}

func newSomeProxyInterface() someProxyInterface {
	return &someProxyStruct{}
}

func (s *someProxyStruct) SomeProxyMethod() {
	funcThatPanics()
}

func someProxyFunc() {
	newSomeProxyInterface().SomeProxyMethod()
}

func TestRecover(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		var rp poco.RecoveredPanic

		func(e *poco.RecoveredPanic) {
			defer poco.Recover(e)

			someProxyFunc()
		}(&rp)

		if rp.Message() != "test" {
			t.Fatal("message is not 'test'")
		}

		// Test using verbose flag (-v) to print stack trace
		// for _, s := range rp.Stack() {
		// 	t.Log(s)
		// }
	})
}
