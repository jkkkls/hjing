// This package from github.com/garallz/gotools/timer
// garalinluzhi@gmail.com
package timer

import (
	"errors"
	"strconv"
	"time"
)

// NewTimer ：make new ticker function
// stamp -> Time timing: 15:04:05; 04:05; 05;
// stamp -> time interval: s-m-h-d:  10s; 30m; 60h; 7d;
// times: 	run times [-1:forever, 0:return not run]
// run:  	defalut: running next time, if true run one times now.
func NewTimer(stamp string, times int, run bool, msg interface{}, function func(interface{})) error {
	if next, interval, err := checkTime(stamp); err != nil {
		return err
	} else {
		if times < 0 {
			times = -1
		} else if times == 0 {
			return errors.New("ticker run times can not be zero")
		}

		if run {
			switch times {
			case 1:
				function(msg)
				return nil
			case -1:
				function(msg)
			default:
				times--
				function(msg)
			}
		}

		putInto(&TimerFunc{
			function: function,
			times:    times,
			next:     next,
			interval: interval,
			msg:      msg,
		})
	}
	return nil
}

// NewRunDuration : Make a new function run
// times: [-1 meas forever], [0 meas not run]
func NewRunDuration(duration time.Duration, times int, msg interface{}, function func(interface{})) {
	if times == 0 {
		return
	} else if times < 0 {
		times = -1
	}

	var data = &TimerFunc{
		next:     time.Now().Add(duration).UnixNano(),
		times:    times,
		interval: int64(duration),
		function: function,
		msg:      msg,
	}
	putInto(data)
}

// NewRunTime : Make a new function run time just one times
func NewRunTime(timestamp time.Time, msg interface{}, function func(interface{})) {
	var data = &TimerFunc{
		next:     timestamp.UnixNano(),
		times:    1,
		function: function,
		msg:      msg,
	}
	putInto(data)
}

// check timerstamp value
func checkTime(stamp string) (int64, int64, error) {
	var err error
	var temp int
	var interval int64
	var now = time.Now()
	var next = now.Unix()

	switch stamp[len(stamp)-1:] {
	case "s", "S":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * SecondTimeUnit)

	case "m", "M":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * MinuteTimeUnit)

	case "h", "H":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * HourTimeUnit)

	case "d", "D":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * DayTimeUnit)

	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		var t time.Time
		timeString := now.Format(TimeFormatString)
		switch len(stamp) {
		case 2: // second
			t, err = time.ParseInLocation(TimeFormatString, timeString[:17]+stamp, time.UTC)
			next = t.UTC().Unix() - 60
			interval = MinuteTimeUnit

		case 5: // min
			t, err = time.ParseInLocation(TimeFormatString, timeString[:14]+stamp, time.UTC)
			next = t.UTC().Unix() - GetJetLag()%3600 - 3600
			interval = HourTimeUnit

		case 8: // hour
			t, err = time.ParseInLocation(TimeFormatString, timeString[:11]+stamp, time.UTC)
			next = t.UTC().Unix() - GetJetLag() - 3600*24
			interval = DayTimeUnit

		default:
			err = errors.New("Can't parst time, please check it")
		}

	default:
		err = errors.New("Can't parst stamp value, please check it")
	}

	if err == nil && interval > 0 {
		for next <= now.Unix() {
			next += interval / SecondTimeUnit
		}
	}
	return next * SecondTimeUnit, interval, err
}
