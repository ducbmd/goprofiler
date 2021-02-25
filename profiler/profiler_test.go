package profiler_test

import (
	"strings"
	"sync"
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
	if !(len(apis) == 1 && strings.HasSuffix(apis[0], "test")) {
		t.Error("APIs must have ony 1 api with name: 'test'")
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
	if !(len(apis) == 2 && strings.HasSuffix(apis[0], "logic1") && strings.HasSuffix(apis[1], "logic2")) {
		t.Errorf("APIs must have 2 apis with name: 'logic1' and 'logic2'. Actual len: %d. Actual name: %v\n", len(apis), apis)
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

func TestLogic4(t *testing.T) {
	prof := profiler.GetProfilerImpl()
	api := "logic4"

	procTmList := []int{
		100, 200, 300, 400,
		500, 200, 300,
		700, 100, 200,
		300, 400, 300,
	}

	now := time.Now().UnixNano()
	delayTime := now/int64(time.Second)*1000_000_000 - now + 1000_000_000
	t.Log(delayTime)
	time.Sleep(time.Duration(delayTime * int64(time.Nanosecond)))

	for _, val := range procTmList {
		state, _ := prof.StartRecord(api)
		time.Sleep(time.Duration(val * int(time.Millisecond)))
		prof.EndRecord(api, state)
	}

	statPoint, _ := prof.GetRealtimeStats(api)
	t.Log(statPoint)
	if !(int(statPoint.TotalReq) == len(procTmList)) {
		t.Errorf("Expect TotalReq: %d. Actual: %d", len(procTmList), statPoint.TotalReq)
	}
	if !(statPoint.TotalTmProc/1_000_000 == 4) {
		t.Errorf("Expect TotalTmProc: %d. Actual: %d", 4, statPoint.TotalTmProc/1_000_000)
	}

	secondStatPoints, _ := prof.GetHistorySecondStats(api)
	t.Log(secondStatPoints)
	if !(secondStatPoints[profiler.NStatSecond-5].TotalReq == 4 &&
		secondStatPoints[profiler.NStatSecond-4].TotalReq == 3 &&
		secondStatPoints[profiler.NStatSecond-3].TotalReq == 3 &&
		secondStatPoints[profiler.NStatSecond-2].TotalReq == 3) {
		t.Errorf("History second not correct")
	}
	if !(secondStatPoints[profiler.NStatSecond-5].TotalTmProc/1_000_000 == 1 &&
		secondStatPoints[profiler.NStatSecond-4].TotalTmProc/1_000_000 == 1 &&
		secondStatPoints[profiler.NStatSecond-3].TotalTmProc/1_000_000 == 1 &&
		secondStatPoints[profiler.NStatSecond-2].TotalTmProc/1_000_000 == 1) {
		t.Errorf("History second not correct")
	}

	minuteStatPoints, _ := prof.GetHistoryMinuteStats(api)
	t.Log(minuteStatPoints)
	if !(minuteStatPoints[profiler.NStatMinute-1].TotalReq == 13 &&
		minuteStatPoints[profiler.NStatMinute-1].TotalTmProc/1_000_000 == 4) {
		t.Errorf("History minute not correct")
	}

	hourStatPoints, _ := prof.GetHistoryHourStats(api)
	t.Log(hourStatPoints)
	if !(hourStatPoints[profiler.NStatHour-1].TotalReq == 13 &&
		hourStatPoints[profiler.NStatHour-1].TotalTmProc/1_000_000 == 4) {
		t.Errorf("History hour not correct")
	}

	profiler.ResetProfilerImpl()
}

func aFunc(wg *sync.WaitGroup) {
	prof := profiler.GetProfilerImpl()
	api := "logic1"

	state, _ := prof.StartRecord(api)
	prof.EndRecord(api, state)

	wg.Done()
}

func TestConcurrent1(t *testing.T) {
	prof := profiler.GetProfilerImpl()
	totalReq := 1000

	var wg sync.WaitGroup
	for i := 0; i < totalReq; i++ {
		wg.Add(1)
		go aFunc(&wg)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < totalReq; i++ {
		wg.Add(1)
		go aFunc(&wg)
	}
	wg.Wait()

	apis, _ := prof.GetAllApis()
	realtimeStats, _ := prof.GetRealtimeStats(apis[0])

	if realtimeStats.TotalReq != int64(2*totalReq) {
		t.Errorf("Expect totalReq: %d. Actual: %d", 2*totalReq, realtimeStats.TotalReq)
	}

	profiler.ResetProfilerImpl()
}

func TestConcurrent2(t *testing.T) {
	prof := profiler.GetProfilerImpl()
	totalReq := 1000
	api := "logic2"

	var wg sync.WaitGroup
	for i := 0; i < totalReq; i++ {
		wg.Add(1)
		go func() {
			state, _ := prof.StartRecord(api)
			prof.EndRecord(api, state)

			wg.Done()
		}()
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < totalReq; i++ {
		wg.Add(1)
		go func() {
			state, _ := prof.StartRecord(api)
			prof.EndRecord(api, state)

			wg.Done()
		}()
	}
	wg.Wait()

	apis, _ := prof.GetAllApis()
	realtimeStats, _ := prof.GetRealtimeStats(apis[0])

	if realtimeStats.TotalReq != int64(2*totalReq) {
		t.Errorf("Expect totalReq: %d. Actual: %d", 2*totalReq, realtimeStats.TotalReq)
	}

	profiler.ResetProfilerImpl()
}

func TestConcurrent3(t *testing.T) {
	prof := profiler.GetProfilerImpl()

	totalReq1 := 100
	api1 := "logic1"

	totalReq2 := 200
	api2 := "logic2"

	totalReq3 := 150
	api3 := "logic3"

	delayMillis := 100
	delayMillisDuration := time.Duration(delayMillis)

	var wg sync.WaitGroup
	for i := 0; i < totalReq1; i++ {
		wg.Add(1)
		go func() {
			state, _ := prof.StartRecord(api1)
			time.Sleep(delayMillisDuration * time.Millisecond)
			prof.EndRecord(api1, state)

			wg.Done()
		}()
	}

	for i := 0; i < totalReq2; i++ {
		wg.Add(1)
		go func() {
			state, _ := prof.StartRecord(api2)
			time.Sleep(delayMillisDuration * time.Millisecond)
			prof.EndRecord(api2, state)

			wg.Done()
		}()
	}

	for i := 0; i < totalReq3; i++ {
		wg.Add(1)
		go func() {
			state, _ := prof.StartRecord(api3)
			time.Sleep(delayMillisDuration * time.Millisecond)
			prof.EndRecord(api3, state)

			wg.Done()
		}()
	}
	wg.Wait()

	apis, _ := prof.GetAllApis()

	for _, api := range apis {
		realtimeStats, _ := prof.GetRealtimeStats(api)
		// t.Log(realtimeStats)
		switch api {
		case api1:
			if realtimeStats.TotalReq != int64(totalReq1) {
				t.Errorf("Expect totalReq: %d. Actual: %d", totalReq1, realtimeStats.TotalReq)
			}
			t.Logf("Expect TotalTmProc: %d. Actual: %d", delayMillis*totalReq1*1000, realtimeStats.TotalTmProc)
		case api2:
			if realtimeStats.TotalReq != int64(totalReq2) {
				t.Errorf("Expect totalReq: %d. Actual: %d", totalReq2, realtimeStats.TotalReq)
			}
			t.Logf("Expect TotalTmProc: %d. Actual: %d", delayMillis*totalReq2*1000, realtimeStats.TotalTmProc)
		case api3:
			if realtimeStats.TotalReq != int64(totalReq3) {
				t.Errorf("Expect totalReq: %d. Actual: %d", totalReq3, realtimeStats.TotalReq)
			}
			t.Logf("Expect TotalTmProc: %d. Actual: %d", delayMillis*totalReq3*1000, realtimeStats.TotalTmProc)
		}
	}

	profiler.ResetProfilerImpl()
}
