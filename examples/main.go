package main

import (
	"fmt"
	"time"

	"github.com/ducbm95/golang-profiler/profiler"
)

func main() {
	prof := profiler.GetProfilerImpl()

	for i := 0; i < 10000; i++ {
		state, _ := prof.StartRecord("getFromDB")

		apis, _ := prof.GetAllApis()
		fmt.Println(prof.GetRealtimeStats(apis[0]))

		time.Sleep(100 * time.Millisecond)
		prof.EndRecord("getFromDB", state)
	}
}
