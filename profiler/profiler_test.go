package profiler_test

import (
	"strings"
	"testing"
	"time"

	"github.com/ducbm95/golang-profiler/profiler/profiler"
)

func TestLogic1(t *testing.T) {
	prof := profiler.GetProfilerImpl()

	state, err1 := prof.StartRecord("test")
	err2 := prof.EndRecord("test", state)

	if !(err1 == nil && err2 == nil) {
		t.Error("Errors must be nil")
	}

	apis, _ := prof.GetAllApis()
	if !(len(apis) == 1 && strings.HasSuffix(apis[0], "@test")) {
		t.Error("APIs must have ony 1 api with name: '@test'")
	}

	profiler.ResetProfilerImpl()
}

func TestLogic2(t *testing.T) {
	prof := profiler.GetProfilerImpl()

	state1, errStart1 := prof.StartRecord("logic1")
	errEnd1 := prof.EndRecord("logic1", state1)

	state2, errStart2 := prof.StartRecord("logic2")
	errEnd2 := prof.EndRecord("logic2", state2)

	if !(errStart1 == nil && errEnd1 == nil && errStart2 == nil && errEnd2 == nil) {
		t.Error("Errors must be nil")
	}

	apis, _ := prof.GetAllApis()
	if !(len(apis) == 2 && strings.HasSuffix(apis[0], "@logic1") && strings.HasSuffix(apis[1], "@logic2")) {
		t.Errorf("APIs must have 2 apis with name: '@logic1' and '@logic2'. Actual len: %d. Actual name: %v\n", len(apis), apis)
	}

	profiler.ResetProfilerImpl()
}

func TestLogic3(t *testing.T) {
	prof := profiler.GetProfilerImpl()
	api := "logic1"
	totalReq := 1000

	for i := 0; i < totalReq; i++ {
		state, _ := prof.StartRecord(api)
		prof.EndRecord(api, state)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < totalReq; i++ {
		state, _ := prof.StartRecord(api)
		prof.EndRecord(api, state)
	}

	apis, _ := prof.GetAllApis()
	realtimeStats, _ := prof.GetRealtimeStats(apis[0])

	if realtimeStats.TotalReq != int64(2*totalReq) {
		t.Errorf("Expect totalReq: %d. Actual: %d", 2*totalReq, realtimeStats.TotalReq)
	}

	profiler.ResetProfilerImpl()

}
