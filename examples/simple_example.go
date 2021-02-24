package main

import (
	"github.com/ducbm95/golang-profiler/profiler/profiler"
)

func main() {
	prof := profiler.GetProfilerImpl()

	state, _ := prof.StartRecord("getFromDB")
	// the business logic for `getFromDB`
	prof.EndRecord("getFromDB", state)
}
