# Golang Profiler
A tool for profiling the performance of Golang services.

## Overview
Nowadays, in the microservices ecosystem, there are too many services run scattered on the server/cloud. Developers/Administrators will hardly to find out where is the bottleneck come from. Golang-Profiler is a solution for your business to easier to manage large amount of golang microservices.

![Alt text](resources/core-flow.png?raw=true "Golang profiler Core Flow")

## Installation
```bash
$ go get github.com/ducbm95/golang-profiler/profiler
```

## Example Usage

```golang
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
```
