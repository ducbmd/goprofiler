package profiler

import "sync"

var profiler Profiler
var syncOnce sync.Once

func GetProfilerImpl() Profiler {
	syncOnce.Do(func() {
		profilerImpl := profilerImpl{}
		profilerImpl.mapHistory = make(map[string]*statInfo)
		profiler = &profilerImpl

		initAPI()
	})

	return profiler
}

func ResetProfilerImpl() {
	profImpl := profilerImpl{}
	profImpl.mapHistory = make(map[string]*statInfo)
	profiler = &profImpl
}
