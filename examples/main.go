package main

import (
	"github.com/ducbm95/golang-profiler/profiler/profiler"
)

func main() {
	profiler.StartRecord("a")
	// profiler.EndRecord("a")

	profiler.StartRecord("b")
	profiler.EndRecord("b")

	// fmt.Printf("%#v\n", profiler.GetStats("main.main@a"))
	// fmt.Printf("%#v\n", profiler.GetStats("main.main@b"))

	profiler.GetAllStats()
}
