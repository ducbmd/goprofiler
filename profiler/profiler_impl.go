package profiler

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// StatInfo is
type StatInfo struct {
	TotalReq    int64
	TotalTmProc int64
	PendingReq  int32
	LastTmProc  int64

	secondStats [nStatSecond]UniformStatPoint
	minuteStats [nStatMinute]UniformStatPoint
	hourStats   [nStatHour]UniformStatPoint

	timestamp int64

	mutSecond sync.Mutex
	mutMinute sync.Mutex
	mutHour   sync.Mutex
}

type ProfilerImpl struct {
	mapHistory map[string]*StatInfo
	mut        sync.Mutex
}

func NewProfilerImpl() Profiler {
	profiler := ProfilerImpl{}
	profiler.mapHistory = make(map[string]*StatInfo)

	return &profiler
}

func (statInfo *StatInfo) addHistoryStat(startTmMicrosec int64, endTmMicrosec int64) {
	tmProc := endTmMicrosec - startTmMicrosec

	timeSecs := startTmMicrosec / 1000000
	timeMins := timeSecs / 60
	timeHours := timeMins / 60

	statInfo.timestamp = startTmMicrosec

	{
		statInfo.mutSecond.Lock()
		defer statInfo.mutSecond.Unlock()

		statSecond := statInfo.secondStats[timeSecs%nStatSecond]
		statSecond.StartTime = timeSecs
		statSecond.TotalReq++
		statSecond.TotalTmProc += tmProc
		if tmProc > statSecond.PeekTimeProc {
			statSecond.PeekTimeProc = tmProc
		}
		if statInfo.PendingReq > statSecond.PeekPendingReq {
			statSecond.PeekPendingReq = statInfo.PendingReq
		}
		statInfo.secondStats[timeSecs%nStatSecond] = statSecond
	}

	{
		statInfo.mutMinute.Lock()
		defer statInfo.mutMinute.Unlock()

		statMinute := statInfo.minuteStats[timeMins%nStatMinute]
		statMinute.StartTime = timeMins
		statMinute.TotalReq++
		statMinute.TotalTmProc += tmProc
		if tmProc > statMinute.PeekTimeProc {
			statMinute.PeekTimeProc = tmProc
		}
		if statInfo.PendingReq > statMinute.PeekPendingReq {
			statMinute.PeekPendingReq = statInfo.PendingReq
		}
		statInfo.minuteStats[timeMins%nStatMinute] = statMinute
	}

	{
		statInfo.mutHour.Lock()
		defer statInfo.mutHour.Unlock()

		statHour := statInfo.hourStats[timeHours%nStatHour]
		statHour.StartTime = timeHours
		statHour.TotalReq++
		statHour.TotalTmProc += tmProc
		if tmProc > statHour.PeekTimeProc {
			statHour.PeekTimeProc = tmProc
		}
		if statInfo.PendingReq > statHour.PeekPendingReq {
			statHour.PeekPendingReq = statInfo.PendingReq
		}
		statInfo.hourStats[timeMins%nStatMinute] = statHour
	}
}

func getFuncName(api string) string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		fmt.Println("MSG: NO CALLER")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		fmt.Println("MSG CALLER WAS NIL")
	}

	return caller.Name() + "@" + api
}

// StartRecord is
func (profiler *ProfilerImpl) StartRecord(api string) (State, error) {
	funcName := getFuncName(api)
	currentTime := time.Now().UnixNano() / int64(time.Microsecond)
	state := State{
		startTime: currentTime,
	}

	profiler.mut.Lock()
	defer profiler.mut.Unlock()

	statInfo, ok := profiler.mapHistory[funcName]
	if !ok {
		statInfo = &StatInfo{}
		profiler.mapHistory[funcName] = statInfo
	}

	statInfo.TotalReq++
	statInfo.PendingReq++

	return state, nil
}

// EndRecord is
func (profiler *ProfilerImpl) EndRecord(api string, state State) error {
	funcName := getFuncName(api)
	currentTime := time.Now().UnixNano() / int64(time.Microsecond)
	timeProc := currentTime - state.startTime

	var statInfo *StatInfo
	var ok bool
	{
		profiler.mut.Lock()
		defer profiler.mut.Unlock()

		statInfo, ok = profiler.mapHistory[funcName]
		if !ok {
			return errors.New("StartRecord must be invoked before EndRecord")
		}

		statInfo.TotalTmProc += timeProc
		statInfo.PendingReq--
		statInfo.LastTmProc = timeProc
	}

	if statInfo != nil {
		statInfo.addHistoryStat(state.startTime, currentTime)
	}
	return nil
}

func (profiler *ProfilerImpl) GetRealtimeStats(fullApi string) (StatPoint, error) {
	currentTime := time.Now().UnixNano()/int64(time.Second) - 1

	statInfo := profiler.mapHistory[fullApi]

	var statPoint StatPoint
	statPoint.TotalReq = statInfo.TotalReq
	statPoint.TotalTmProc = statInfo.TotalTmProc
	statPoint.PendingReq = statInfo.PendingReq
	statPoint.LastTmProc = statInfo.LastTmProc
	statPoint.ProcRate = float64(statPoint.TotalReq / statPoint.TotalTmProc)
	statPoint.ReqRate = float64(statInfo.secondStats[currentTime%nStatSecond].TotalReq)

	return statPoint, nil
}
func (profiler *ProfilerImpl) GetHistorySecondStats(api string) ([]UniformStatPoint, error) {
	// TODO
	return nil, nil
}
func (profiler *ProfilerImpl) GetHistoryMinuteStats(api string) ([]UniformStatPoint, error) {
	// TODO
	return nil, nil
}
func (profiler *ProfilerImpl) GetHistoryHourStats(api string) ([]UniformStatPoint, error) {
	// TODO
	return nil, nil
}

// GetStats is
// func GetStats(funcName string) StatInfo {
// 	statInfo, ok := mapHistory[funcName]
// 	if ok {
// 		return *statInfo
// 	} else {
// 		fmt.Printf("%s has no stats", funcName)
// 		return StatInfo{}
// 	}
// }

// // GetAllStats is
// func GetAllStats() {
// 	table := tablewriter.NewWriter(os.Stdout)
// 	table.SetHeader([]string{"Name", "TotalReq", "PendingReq", "TotalTmProc", "LastTmProc", "ProcRate", "ReqRate"})
// 	table.SetAutoFormatHeaders(false)
// 	table.SetHeaderAlignment(tablewriter.ALIGN_RIGHT)

// 	for k, v := range mapHistory {
// 		lineData := []string{
// 			k,
// 			strconv.FormatInt(v.TotalReq, 10),
// 			strconv.FormatInt(int64(v.PendingReq), 10),
// 			strconv.FormatInt(v.TotalTmProc, 10),
// 			"0",
// 			"0",
// 			"0",
// 		}
// 		table.Append(lineData)
// 	}

// 	table.Render()
// }
