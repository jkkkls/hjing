// This package from github.com/garallz/gotools/timer
// garalinluzhi@gmail.com
package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

const basicNumber int64 = 100 * MilliTimeUnit // 100ms

// Container : save cache data containers
type Container struct {
	cache []*TimerFunc
	count int32
	cutNo int64
	lock  sync.Mutex
}

var (
	first  = &Container{cutNo: 5 * basicNumber}
	second = &Container{cutNo: 25 * basicNumber}
	third  = &Container{cutNo: 250 * basicNumber}
)

var level int64 = 5

// InitTicker : init timer ticker
// base interval is 100ms, dafualt 100ms, [100ms * interval]
func InitTicker(interval int64) {
	if interval > 0 {
		level = interval
	}

	go func() {
		ticker := time.NewTicker(time.Nanosecond * time.Duration(basicNumber*level))

		for {
			select {
			case <-ticker.C:

				if count := atomic.AddInt32(&second.count, 1); count >= 4 {
					// run cache second check
					go checkSecondCache()
					atomic.SwapInt32(&second.count, 0)
				}

				now := time.Now().UnixNano()
				// append first arrge data
				first.lock.Lock()
				mins, maxs := TimeSplit(first.cache, now)
				first.cache = maxs
				first.lock.Unlock()

				for _, row := range mins {
					go run(row)
				}
			}
		}
	}()
}

func run(data *TimerFunc) {
	switch data.times {
	case 0:
		return
	case 1:
		data.function(data.msg)
		return
	case -1:
		data.next += data.interval
	default:
		data.times--
		data.next += data.interval
	}

	go putInto(data)
	data.function(data.msg)
}

func putInto(data *TimerFunc) {
	now := time.Now().UnixNano()

	if data.next <= (now + first.cutNo*level) {
		first.lock.Lock()
		first.cache = append(first.cache, data)
		first.lock.Unlock()
	} else if data.next > (now + second.cutNo*level) {
		third.lock.Lock()
		third.cache = append(third.cache, data)
		third.lock.Unlock()
	} else {
		second.lock.Lock()
		second.cache = append(second.cache, data)
		second.lock.Unlock()
	}
}

// When insert new data to sort
func checkSecondCache() {
	if count := atomic.AddInt32(&third.count, 1); count >= 5 {
		// run cache third check
		go checkThirdCache()
		atomic.SwapInt32(&third.count, 0)
	}

	next := time.Now().UnixNano() + first.cutNo*level

	second.lock.Lock()
	mins, maxs := TimeSplit(second.cache, next)

	first.lock.Lock()
	first.cache = append(first.cache, mins...)
	first.lock.Unlock()

	second.cache = maxs
	second.lock.Unlock()
}

// When insert new data to sort
func checkThirdCache() {
	next := time.Now().UnixNano() + second.cutNo*level

	third.lock.Lock()
	mins, maxs := TimeSplit(third.cache, next)

	second.lock.Lock()
	second.cache = append(second.cache, mins...)
	second.lock.Unlock()

	third.cache = maxs
	third.lock.Unlock()
}
