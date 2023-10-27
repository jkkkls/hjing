// This package from github.com/garallz/gotools/timer
// garalinluzhi@gmail.com
package timer

// parse time format
const (
	TimeFormatString = "2006-01-02 15:04:05"
	TimeFormatLength = 19
)

// time unit
const (
	NanoTimeUnit   = 1
	MicroTimeUnit  = 1000 * NanoTimeUnit
	MilliTimeUnit  = 1000 * MicroTimeUnit
	SecondTimeUnit = 1000 * MilliTimeUnit
	MinuteTimeUnit = 60 * SecondTimeUnit
	HourTimeUnit   = 60 * MinuteTimeUnit
	DayTimeUnit    = 24 * HourTimeUnit
)

// TimerFunc ticker function struct
type TimerFunc struct {
	next     int64 // Next run time
	interval int64 // 时间间隔
	times    int   // run times
	function func(interface{})
	msg      interface{}
}

// ByNext sort by next time
type ByNext []*TimerFunc

func (t ByNext) Len() int {
	return len(t)
}

func (t ByNext) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByNext) Less(i, j int) bool {
	return t[i].next < t[j].next
}

func TimeSplit(rows []*TimerFunc, timestamp int64) ([]*TimerFunc, []*TimerFunc) {
	var mins, maxs = make([]*TimerFunc, 0), make([]*TimerFunc, 0)
	for _, row := range rows {
		if row.next <= timestamp {
			mins = append(mins, row)
		} else {
			maxs = append(maxs, row)
		}
	}
	return mins, maxs
}
