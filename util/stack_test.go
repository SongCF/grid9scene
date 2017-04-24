package util

import (
	"testing"
)

func TestRecoverPanic(t *testing.T) {
	if errFunc() != "" {
		t.Error("RecoverPanic Error")
	}
}

func errFunc() string {
	defer RecoverPanic()
	panic("test panic")
	return "error"
}
