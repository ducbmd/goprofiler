package profiler

const (
	minInt64 = -9223372036854775808
	maxInt64 = 9223372036854775807

	nStatSecond = 120
	nStatMinute = 120
	nStatHour   = 336
)

// State struct stores state data when StartRecord
type State struct {
	startTime int64 // start_time_microsecs
}

type StatPoint struct {
	TotalReq    int64   // total_req
	TotalTmProc int64   // total_tm_proc_microsecs
	PendingReq  int32   // pending_req
	LastTmProc  int64   // last_tm_proc_microsecs
	ProcRate    float64 // unit: req_per_second
	ReqRate     float64 // unit: req_per_second
}

type UniformStatPoint struct {
	StartTime      int64 // epoch_second_or_epoch_minute_or_epoch_hour
	TotalReq       int64 // total_req
	TotalTmProc    int64 // total_tm_proc_microsecs
	PeekTimeProc   int64 // peek_time_proc_microsecs
	PeekPendingReq int32 // peek_pending_req
}

// Profiler is ...
type Profiler interface {
	StartRecord(api string) (State, error)
	EndRecord(api string, state State) error

	// get current stats
	GetRealtimeStats(fullAPI string) (StatPoint, error)

	// Get history second stats. Newest stats will come to last of slice.
	GetHistorySecondStats(fullAPI string) ([]UniformStatPoint, error)

	// Get history minute stats. Newest stats will come to last of slice.
	GetHistoryMinuteStats(fullAPI string) ([]UniformStatPoint, error)

	// Get history hour stats. Newest stats will come to last of slice.
	GetHistoryHourStats(fullAPI string) ([]UniformStatPoint, error)

	GetAllApis() ([]string, error)
}
