// This package from github.com/garallz/gotools/timer
// garalinluzhi@gmail.com
package timer

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	timeZoneStatus bool
	timeZoneJetLag int64
)

// SetTimeZone : verify time eg:[-08:00, +08:00]
// If want use UTC to time zone : '+00:00'
func SetTimeZone(v string) error {
	if v == "" {
		return errors.New("Time Zone string is null")
	} else if v[0] != 43 && v[0] != 45 {
		return errors.New("Time Zone string first byte not ['-' or '+']")
	} else if v == "+00:00" {
		timeZoneStatus = true
	}

	t := strings.Split(v[1:], ":")
	if len(t) != 2 {
		return errors.New("Time Zone string format wrong")
	}

	h, err := strconv.ParseUint(t[0], 10, 64)
	if err != nil {
		return errors.New("Parse Time Zone hour error: " + err.Error())
	} else if h > 23 {
		return errors.New("Time Zone hour more than 23")
	}

	m, err := strconv.ParseUint(t[1], 10, 64)
	if err != nil {
		return errors.New("Parse Time Zone minter error: " + err.Error())
	} else if m > 59 {
		return errors.New("Time Zone minter more than 59")
	}

	if v[0] == 43 {
		timeZoneJetLag = int64(h*3600 + m*60)
	} else {
		timeZoneJetLag = -int64(h*3600 + m*60)
	}
	timeZoneStatus = true

	return nil
}

// GetJetLag : get time zone differet
// Solve the time difference of time zone
// unit is second
func GetJetLag() int64 {
	if !timeZoneStatus {
		_, diff := time.Now().Zone()
		return int64(diff)
	}
	return timeZoneJetLag
}

// GetTimeNowUnix : get time unix nano now
func GetTimeNowUnix() int64 {
	return time.Now().UTC().UnixNano()
}

// GetTimeNowZone : Get current time zone nano temporary
func GetTimeNowZone() int64 {
	return time.Now().UnixNano() + GetJetLag()*SecondTimeUnit
}
