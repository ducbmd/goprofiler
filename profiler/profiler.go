package profiler

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Profiler is ...
type Profiler interface {
	StartRecord()
	EndRecord()
}

// StatPoint is a struct to store time-series data.
type StatPoint struct {
	StartTime   int64 `json:"epoch_millis"`
	TotalTmProc int64
}

type UniformStatPoint struct {
	StartTime      int64 `json:"epoch_second"`
	ProcRate       int64
	ReqRate        int64
	AvgTimeProc    int64
	PeekTimeProc   int64
	PeekPendingReq int64
}

type StatInfo struct {
	TotalReq    int64
	TotalTmProc int64
	PendingReq  int32

	Stats        []StatPoint
	UniformStats []UniformStatPoint

	LastReqStartTime int64 `json:"epoch_millis"`
}

func (statInfo *StatInfo) appendStat(statPoint StatPoint) {
	(*statInfo).Stats = append(statInfo.Stats, statPoint)
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

func StartRecord(api string) {
	funcName := getFuncName(api)

	currentTime := time.Now().UnixNano() / int64(time.Microsecond)
	statInfo, ok := mapHistory[funcName]
	if !ok {
		statInfo = &StatInfo{
			Stats:        make([]StatPoint, 0, 0),
			UniformStats: make([]UniformStatPoint, 0, 0),
		}
		mapHistory[funcName] = statInfo
	}

	statInfo.TotalReq++
	statInfo.PendingReq++
	statInfo.LastReqStartTime = currentTime
}

func EndRecord(api string) {
	funcName := getFuncName(api)

	currentTime := time.Now().UnixNano() / int64(time.Microsecond)
	statInfo, _ := mapHistory[funcName]

	statInfo.TotalTmProc += (currentTime - statInfo.LastReqStartTime)
	statInfo.PendingReq--

	statPoint := StatPoint{
		StartTime:   statInfo.LastReqStartTime,
		TotalTmProc: currentTime - statInfo.LastReqStartTime,
	}
	statInfo.appendStat(statPoint)
	statInfo.LastReqStartTime = 0
}

func GetStats(funcName string) StatInfo {
	statInfo, ok := mapHistory[funcName]
	if ok {
		return *statInfo
	} else {
		fmt.Printf("%s has no stats", funcName)
		return StatInfo{}
	}
}

func GetAllStats() {
	// var listStat []StatInfo = make([]StatInfo, 0)
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
