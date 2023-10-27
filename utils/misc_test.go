package utils

import (
	"testing"
	"time"
)

//go test  -bench=. -benchtime="3s"

func GetDay(reigsterTime uint64) uint64 {
	HoChiMinhLoc, _ := time.LoadLocation(HoChiMinhLocKey)
	registerTime := time.Unix(int64(reigsterTime), 0).In(HoChiMinhLoc)
	now := time.Now().In(HoChiMinhLoc)
	registerDay := time.Date(registerTime.Year(), registerTime.Month(), registerTime.Day(), 0, 0, 0, 0, HoChiMinhLoc)
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, HoChiMinhLoc)
	days := int(nowDay.Sub(registerDay).Hours() / 24)
	return uint64(days)
}

func BenchmarkGetVietnamDay(b *testing.B) {
	b.ReportAllocs()

	now := Now()
	reigsterTime := now - 7*24*3600
	for i := 0; i < b.N; i++ {
		now1 := Now()
		days := int(GetVietnamDay(now1)) - int(GetVietnamDay(reigsterTime))
		if days != 7 {
			b.Errorf("Unexpected result: %v", days)
		}
	}
}

func BenchmarkGetDay(b *testing.B) {
	b.ReportAllocs()

	now := Now()
	reigsterTime := now - 7*24*3600
	for i := 0; i < b.N; i++ {
		days := GetDay(reigsterTime)
		if days != 7 {
			b.Errorf("Unexpected result: %v", days)
		}
	}
}
