package profiler

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// StatInfo is
type statInfo struct {
	totalReq    int64
	totalTmProc int64
	pendingReq  int32
	lastTmProc  int64

	secondStats [nStatSecond]UniformStatPoint
	minuteStats [nStatMinute]UniformStatPoint
	hourStats   [nStatHour]UniformStatPoint

	timestamp int64

	mutSecond sync.Mutex
	mutMinute sync.Mutex
	mutHour   sync.Mutex
}

type profilerImpl struct {
	mapHistory map[string]*statInfo
	mut        sync.Mutex
}

var profiler Profiler
var syncOnce sync.Once

func GetProfilerImpl() Profiler {
	syncOnce.Do(func() {
		profilerImpl := profilerImpl{}
		profilerImpl.mapHistory = make(map[string]*statInfo)
		profiler = &profilerImpl
	})

	return profiler
}

func ResetProfilerImpl() {
	profImpl := profilerImpl{}
	profImpl.mapHistory = make(map[string]*statInfo)
	profiler = &profImpl
}

func (si *statInfo) addHistoryStat(startTmMicrosec int64, endTmMicrosec int64) {
	tmProc := endTmMicrosec - startTmMicrosec

	timeSecs := startTmMicrosec / 1000000
	timeMins := timeSecs / 60
	timeHours := timeMins / 60

	si.timestamp = startTmMicrosec

	{
		si.mutSecond.Lock()
		defer si.mutSecond.Unlock()

		statSecond := si.secondStats[timeSecs%nStatSecond]
		statSecond.StartTime = timeSecs
		statSecond.TotalReq++
		statSecond.TotalTmProc += tmProc
		if tmProc > statSecond.PeekTimeProc {
			statSecond.PeekTimeProc = tmProc
		}
		if si.pendingReq > statSecond.PeekPendingReq {
			statSecond.PeekPendingReq = si.pendingReq
		}
		si.secondStats[timeSecs%nStatSecond] = statSecond
	}

	{
		si.mutMinute.Lock()
		defer si.mutMinute.Unlock()

		statMinute := si.minuteStats[timeMins%nStatMinute]
		statMinute.StartTime = timeMins
		statMinute.TotalReq++
		statMinute.TotalTmProc += tmProc
		if tmProc > statMinute.PeekTimeProc {
			statMinute.PeekTimeProc = tmProc
		}
		if si.pendingReq > statMinute.PeekPendingReq {
			statMinute.PeekPendingReq = si.pendingReq
		}
		si.minuteStats[timeMins%nStatMinute] = statMinute
	}

	{
		si.mutHour.Lock()
		defer si.mutHour.Unlock()

		statHour := si.hourStats[timeHours%nStatHour]
		statHour.StartTime = timeHours
		statHour.TotalReq++
		statHour.TotalTmProc += tmProc
		if tmProc > statHour.PeekTimeProc {
			statHour.PeekTimeProc = tmProc
		}
		if si.pendingReq > statHour.PeekPendingReq {
			statHour.PeekPendingReq = si.pendingReq
		}
		si.hourStats[timeHours%nStatHour] = statHour
	}
}

func getFuncName(api string) string {
	return api
}

// StartRecord is
func (profiler *profilerImpl) StartRecord(api string) (State, error) {
	funcName := getFuncName(api)
	currentTime := time.Now().UnixNano() / int64(time.Microsecond)
	state := State{
		startTime: currentTime,
	}

	profiler.mut.Lock()
	defer profiler.mut.Unlock()

	si, ok := profiler.mapHistory[funcName]
	if !ok {
		si = &statInfo{}
		profiler.mapHistory[funcName] = si
	}

	si.totalReq++
	si.pendingReq++

	return state, nil
}

// EndRecord is
func (profiler *profilerImpl) EndRecord(api string, state State) error {
	funcName := getFuncName(api)
	currentTime := time.Now().UnixNano() / int64(time.Microsecond)
	timeProc := currentTime - state.startTime

	var si *statInfo
	var ok bool
	{
		profiler.mut.Lock()
		defer profiler.mut.Unlock()

		si, ok = profiler.mapHistory[funcName]
		if !ok {
			return errors.New("StartRecord must be invoked before EndRecord")
		}

		si.totalTmProc += timeProc
		si.pendingReq--
		si.lastTmProc = timeProc
	}

	if si != nil {
		si.addHistoryStat(state.startTime, currentTime)
	}
	return nil
}

func (profiler *profilerImpl) GetRealtimeStats(fullAPI string) (StatPoint, error) {
	currentTime := time.Now().UnixNano()/int64(time.Second) - 1

	si := profiler.mapHistory[fullAPI]

	var sp StatPoint
	sp.TotalReq = si.totalReq
	sp.TotalTmProc = si.totalTmProc
	sp.PendingReq = si.pendingReq
	sp.LastTmProc = si.lastTmProc
	if sp.TotalTmProc > 0 {
		sp.ProcRate = 1000000 * float64(sp.TotalReq) / float64(sp.TotalTmProc)
	}
	sp.ReqRate = float64(si.secondStats[currentTime%nStatSecond].TotalReq)

	return sp, nil
}

func (profiler *profilerImpl) GetHistorySecondStats(fullAPI string) ([]UniformStatPoint, error) {
	currentTime := time.Now().UnixNano() / int64(time.Second)
	statInfo := profiler.mapHistory[fullAPI]

	secondStats := make([]UniformStatPoint, nStatSecond)
	for i := 0; i < nStatSecond; i++ {
		idx := (currentTime + int64(i) + 1) % nStatSecond
		secondStat := statInfo.secondStats[idx]

		if secondStat.StartTime >= currentTime-int64(nStatSecond) {
			secondStats[i] = secondStat
		} else {
			secondStats[i] = UniformStatPoint{}
		}
	}

	return secondStats, nil
}

func (profiler *profilerImpl) GetHistoryMinuteStats(fullAPI string) ([]UniformStatPoint, error) {
	currentTime := time.Now().UnixNano() / int64(time.Minute)
	si := profiler.mapHistory[fullAPI]

	minuteStats := make([]UniformStatPoint, nStatMinute)
	for i := 0; i < nStatMinute; i++ {
		idx := (currentTime + int64(i) + 1) % nStatMinute
		minuteStat := si.minuteStats[idx]

		if minuteStat.StartTime >= currentTime-int64(nStatHour) {
			minuteStats[i] = minuteStat
		} else {
			minuteStats[i] = UniformStatPoint{}
		}
	}

	return minuteStats, nil
}

func (profiler *profilerImpl) GetHistoryHourStats(fullAPI string) ([]UniformStatPoint, error) {
	currentTime := time.Now().UnixNano() / int64(time.Hour)
	si := profiler.mapHistory[fullAPI]

	hourStats := make([]UniformStatPoint, nStatHour)
	for i := 0; i < nStatHour; i++ {
		idx := (currentTime + int64(i) + 1) % nStatHour
		hourStat := si.hourStats[idx]
		if hourStat.StartTime > 0 {
			fmt.Println(hourStat)
		}

		if hourStat.StartTime >= currentTime-int64(nStatHour) {
			hourStats[i] = hourStat
		} else {
			hourStats[i] = UniformStatPoint{}
		}
	}

	return hourStats, nil
}

func (profiler *profilerImpl) GetAllApis() ([]string, error) {
	apis := make([]string, 0)
	for k, _ := range profiler.mapHistory {
		apis = append(apis, k)
	}

	sort.Strings(apis)
	return apis, nil
}
