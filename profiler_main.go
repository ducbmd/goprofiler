package goprofiler

import "sync"

var profilerInst Profiler
var syncOnce sync.Once

func GetProfilerImpl() Profiler {
	syncOnce.Do(func() {
		profilerImpl := profilerImpl{}
		profilerImpl.mapHistory = make(map[string]*statInfo)
		profilerInst = &profilerImpl

		initAPI()
	})

	return profilerInst
}

func ResetProfilerImpl() {
	profImpl := profilerImpl{}
	profImpl.mapHistory = make(map[string]*statInfo)
	profilerInst = &profImpl
}
