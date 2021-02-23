package profiler

const (
	minInt64 = -9223372036854775808
	maxInt64 = 9223372036854775807

	nStatSecond = 120
	nStatMinute = 120
	nStatHour   = 336
)

type State struct {
	startTime int64 // `json:"epoch_millis"`
}

type StatPoint struct {
	TotalReq    int64
	TotalTmProc int64
	PendingReq  int32
	LastTmProc  int64
	ProcRate    float64
	ReqRate     float64
}

type UniformStatPoint struct {
	StartTime      int64 `json:"epoch_second or epoch_minute or epoch_hour"`
	TotalReq       int64
	TotalTmProc    int64
	PeekTimeProc   int64
	PeekPendingReq int32
}

// Profiler is ...
type Profiler interface {
	StartRecord(api string) (State, error)
	EndRecord(api string, state State) error

	GetRealtimeStats(api string) (StatPoint, error)
	GetHistorySecondStats(api string) ([]UniformStatPoint, error)
	GetHistoryMinuteStats(api string) ([]UniformStatPoint, error)
	GetHistoryHourStats(api string) ([]UniformStatPoint, error)
}
