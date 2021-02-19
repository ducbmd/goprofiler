package profiler

import (
	"fmt"
	"runtime"
)

func GetFuncName() string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(2, fpcs)
	if n == 0 {
		fmt.Println("MSG: NO CALLER")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		fmt.Println("MSG CALLER WAS NIL")
	}

	return caller.Name()
}
