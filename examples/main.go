package main

import (
	"fmt"

	"github.com/ducbm95/golang-profiler/profiler/profiler"
)

func main() {
	funcName := profiler.GetFuncName()
	profiler.StartRecord(funcName)

	profiler.EndRecord(funcName)

	fmt.Println(profiler.GetStats(funcName))
	fmt.Println(profiler.GetStats(funcName).Stats[0])
}
