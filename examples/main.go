package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	// "github.com/ducbm95/golang-profiler/profiler/profiler"
)

func goID() int {
	var buf [12]byte
	n := runtime.Stack(buf[:], false)
	// fmt.Printf("N: %s\n", string(buf[:n]))
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func main() {
	// profiler.StartRecord("a")
	// // profiler.EndRecord("a")

	// profiler.StartRecord("b")
	// profiler.EndRecord("b")

	// // fmt.Printf("%#v\n", profiler.GetStats("main.main@a"))
	// // fmt.Printf("%#v\n", profiler.GetStats("main.main@b"))

	// profiler.GetAllStats()

	// profiler.InitUI()

	// // profiler.Profiler.EndRecord()
	x := goID()
	fmt.Println(x)

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			x := goID()
			fmt.Println(x)
			wg.Done()
		}()
	}
	wg.Wait()
}
