# Time Ticker

** This package from github.com/garallz/gotools/timer **
** garalinluzhi@gmail.com **


### use timer eg
```go
	// init timer ticker
	// base interval is 100ms, dafualt 100ms, [100ms * interval]
	InitTicker(10)	// meat 1s interval to check

	// SetTimeZone : verify time eg:[-08:00, +08:00]
	// If want use UTC to time zone : '+00:00'
	SetTimeZone(v string) error

	// NewTimer ：make new ticker function
	// stamp -> Time timing: 15:04:05; 04:05; 05;
	// stamp -> time interval: s-m-h-d:  10s; 30m; 60h; 7d;
	// times: 	run times [-1:forever, 0:return not run]
	// run:  	defalut: running next time, if true run one times now.
	NewTimer(stamp string, times int, run bool, msg interface{}, function func(interface{})) error

	// NewRunDuration : Make a new function run
	// times: [-1 meas forever], [0 meas not run]
	NewRunDuration(duration time.Duration, times int, msg interface{}, function func(interface{}))
```

## TODO: program level