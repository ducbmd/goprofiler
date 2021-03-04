package main

import profiler "github.com/ducbm95/goprofiler"

func main() {
	prof := profiler.GetProfilerImpl()

	state, _ := prof.StartRecord("getFromDB")
	// the business logic for `getFromDB`
	prof.EndRecord("getFromDB", state)
}
