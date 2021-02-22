package profiler

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
)

const (
	minInt64 = -9223372036854775808
	maxInt64 = 9223372036854775807
)

// Profiler is ...
type Profiler interface {
	StartRecord()
	EndRecord()
}

// StatPoint is a struct to store time-series data.
type statPoint struct {
	StartTime   int64 `json:"epoch_microsecs"`
	TotalTmProc int64 `json:"microsecs"`
}

// UniformStatPoint is
type UniformStatPoint struct {
	StartTime      int64 `json:"epoch_second"`
	TotalReq       int64
	TotalTmProc    int64
	PeekTimeProc   int64
	PeekPendingReq int64
	// StartTime      int64 `json:"epoch_second"`
	// ProcRate       float64
	// ReqRate        float64
	// AvgTimeProc    float64
	// PeekTimeProc   int64
	// PeekPendingReq int64
}

// StatInfo is
type StatInfo struct {
	TotalReq    int64
	TotalTmProc int64
	PendingReq  int32

	stats        []statPoint
	UniformStats []UniformStatPoint

	lastReqStartTime int64 // `json:"epoch_millis"`
}

func (statInfo *StatInfo) appendStat(statPoint statPoint) {
	(*statInfo).stats = append(statInfo.stats, statPoint)
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

var mapHistory map[string]*StatInfo = make(map[string]*StatInfo)
var mut sync.Mutex

// StartRecord is
func StartRecord(api string) {
	funcName := getFuncName(api)

	currentTime := time.Now().UnixNano() / int64(time.Microsecond)
	statInfo, ok := mapHistory[funcName]
	if !ok {
		statInfo = &StatInfo{
			stats:        make([]statPoint, 0, 0),
			UniformStats: make([]UniformStatPoint, 0, 0),
		}
		mapHistory[funcName] = statInfo
	}

	statInfo.TotalReq++
	statInfo.PendingReq++
	statInfo.lastReqStartTime = currentTime
}

// EndRecord is
func EndRecord(api string) {
	funcName := getFuncName(api)

	currentTime := time.Now().UnixNano() / int64(time.Microsecond)
	statInfo, _ := mapHistory[funcName]

	statInfo.TotalTmProc += (currentTime - statInfo.lastReqStartTime)
	statInfo.PendingReq--

	statPoint := statPoint{
		StartTime:   statInfo.lastReqStartTime,
		TotalTmProc: currentTime - statInfo.lastReqStartTime,
	}
	statInfo.appendStat(statPoint)
	statInfo.lastReqStartTime = 0
}

// GetStats is
func GetStats(funcName string) StatInfo {
	statInfo, ok := mapHistory[funcName]
	if ok {
		return *statInfo
	} else {
		fmt.Printf("%s has no stats", funcName)
		return StatInfo{}
	}
}

// GetAllStats is
func GetAllStats() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "TotalReq", "PendingReq", "TotalTmProc", "LastTmProc", "ProcRate", "ReqRate"})
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_RIGHT)

	for k, v := range mapHistory {
		lineData := []string{
			k,
			strconv.FormatInt(v.TotalReq, 10),
			strconv.FormatInt(int64(v.PendingReq), 10),
			strconv.FormatInt(v.TotalTmProc, 10),
			"0",
			"0",
			"0",
		}
		table.Append(lineData)
	}

	table.Render()
}

func uniformStat(statInfo *StatInfo) {
	mut.Lock()
	stats := statInfo.stats
	statInfo.stats = make([]statPoint, 0, 0)
	mut.Unlock()

	// var uniformData []UniformStatPoint

	var mapStatBySecond map[int64][]statPoint = make(map[int64][]statPoint) // map from second to list of stats
	minSecond := int64(maxInt64)
	maxSecond := int64(minInt64)
	for _, stat := range stats {
		startTimeSecond := stat.StartTime / 1000000
		if startTimeSecond > maxInt64 {
			maxSecond = startTimeSecond
		}

		if startTimeSecond < minInt64 {
			minSecond = startTimeSecond
		}

		lsStat, ok := mapStatBySecond[startTimeSecond]

		if !ok {
			lsStat = make([]statPoint, 0, 0)
			mapStatBySecond[startTimeSecond] = lsStat
		}

		lsStat = append(lsStat, stat)
	}

	for second := minSecond; second <= maxSecond; second++ {
		totalReq := int64(0)
		totalTmProc := int64(0)

		statBySecond, ok := mapStatBySecond[second]
		if ok {
			for _, stat := range statBySecond {
				totalReq++
				totalTmProc += stat.TotalTmProc
			}
		}

		var uniformStatPoint UniformStatPoint
		uniformStatPoint.StartTime = second
		uniformStatPoint.TotalReq = totalReq
		uniformStatPoint.TotalTmProc = totalTmProc

		// uniformStatPoint.ReqRate = float64(totalReq)
		// uniformStatPoint.ProcRate = float64(totalTmProc / totalReq)
		// uniformStatPoint.
	}
}
