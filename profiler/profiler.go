package profiler

import (
	"fmt"
	"time"
)

type Profiler interface {
	StartRecord()
	EndRecord()
}

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

var mapHistory map[string]*StatInfo = make(map[string]*StatInfo)

func StartRecord(funcName string) {
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

func EndRecord(funcName string) {
	currentTime := time.Now().UnixNano() / int64(time.Microsecond)
	statInfo, _ := mapHistory[funcName]

	statInfo.TotalTmProc += (currentTime - statInfo.LastReqStartTime)
	statInfo.PendingReq--

	statPoint := StatPoint{
		StartTime:   statInfo.LastReqStartTime,
		TotalTmProc: currentTime - statInfo.LastReqStartTime,
	}
	statInfo.appendStat(statPoint)
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
