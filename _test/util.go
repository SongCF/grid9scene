package _test

import (
	"runtime/debug"
	"testing"
	"runtime"
	"fmt"
)

func check(err error, dsc string) {
	if err != nil {
		panic(dsc + err.Error())
	}
}

func assert(b bool, t *testing.T) {
	if !b {
		debug.PrintStack()
		t.Fatal()
	}
}

func RecoverPanic() {
	if x := recover(); x != nil {
		fmt.Errorf("%v", x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			fmt.Errorf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}
	}
}