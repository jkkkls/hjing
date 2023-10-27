package utils

// Copyright 2017 guangbo. All rights reserved.

//常用接口

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// SendHttpRequest 发送http请求
func SendHttpRequest(method string, url string, body string) ([]byte, error) {
	client := &http.Client{}
	req, err1 := http.NewRequest(method, url, strings.NewReader(body))
	if err1 != nil {
		return nil, err1
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, err2
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// PathExists 目录是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// CreateDir 创建一个目录
func CreateDir(path string) {
	os.MkdirAll(path, os.ModePerm)
}

// ExportServiceFunction 从服务中获取可以通过cmd调用的函数
func ExportServiceFunction(u interface{}) map[uint16]string {
	funcs := make(map[uint16]string)
	t := reflect.TypeOf(u)

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		l := len(m.Name)
		if l <= 5 {
			continue
		}

		s := string([]byte(m.Name)[l-4 : l])
		cmd, err := strconv.ParseInt(s, 16, 32)
		if err != nil {
			continue
		}

		funcs[uint16(cmd)] = m.Name
	}

	return funcs
}

func InIntArray(array []int, key int) bool {
	for i := 0; i < len(array); i++ {
		if (array)[i] == key {
			return true
		}
	}
	return false
}

func InArray(array *[]string, key string) bool {
	for i := 0; i < len(*array); i++ {
		if (*array)[i] == key {
			return true
		}
	}
	return false
}

// Now 返回当前时间戳
func Now() uint64 {
	return uint64(time.Now().Unix())
}

// NowNano 返回当前时间戳
func NowNano() uint64 {
	return uint64(time.Now().UnixNano())
}

// ITS interface -> string
func ITS(i interface{}) string {
	if i == nil {
		return ""
	}
	return i.(string)
}

// STU32 string -> uint32
func STU32(str string) uint32 {
	n, _ := strconv.Atoi(str)
	return uint32(n)
}

// GetLocatione ...
func GetLocatione(timeZone string) *time.Location {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return loc
	}
	loc, _ = time.LoadLocation("Local")
	return loc
}

// GetNextDuration 获取指定时区指定时间差
func GetNextDuration(timeZone, date string) time.Duration {
	now := time.Now()
	loc := GetLocatione(timeZone)
	arr := strings.Split(date, ":")
	hour, _ := strconv.Atoi(arr[0])
	min, _ := strconv.Atoi(arr[1])
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, loc)
	if !now.Before(next) {
		next = next.Add(24 * time.Hour)
	}
	return next.Sub(now)
}

// STU64 string -> uint64
func STU64(str string) uint64 {
	n, _ := strconv.ParseUint(str, 10, 64)
	return n
}

// STI64 string -> int64
func STI64(str string) int64 {
	n, _ := strconv.ParseInt(str, 10, 64)
	return n
}

// LTTS local time -> string
func LTTS(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

// LSTT local string -> time
func LSTT(str string) int64 {
	loc, _ := time.LoadLocation("Local")

	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", str, loc)
	if err == nil {
		return theTime.Unix()
	} else {
		return 0
	}
}

// TTS time -> string
func TTS(t int64) string {
	local, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		return time.Unix(t, 0).Format("2006-01-02 15:04:05")
	}
	return time.Unix(t, 0).In(local).Format("2006-01-02 15:04:05")
}

// STT string -> time
func STT(str string) int64 {
	loc, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		loc, _ = time.LoadLocation("Local")
	}

	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", str, loc)
	if err == nil {
		return theTime.Unix()
	} else {
		return 0
	}
}

// If 三目操作模拟函数
func If(x bool, a interface{}, b interface{}) interface{} {
	if x {
		return a
	}

	return b
}

// SubUint32Array ...
func SubUint32Array(arr1, arr2 []uint32) []uint32 {
	var arr []uint32
	for i := 0; i < len(arr1); i++ {
		if !InUint32Array(arr2, arr1[i]) {
			arr = append(arr, arr1[i])
		}
	}

	return arr
}

// InUint16Array n是否存在指定数组中
func InUint16Array(array []uint16, n uint16) bool {
	for i := 0; i < len(array); i++ {
		if n == array[i] {
			return true
		}
	}

	return false
}

// InUint32Array n是否存在指定数组中
func InUint32Array(array []uint32, n uint32) bool {
	for i := 0; i < len(array); i++ {
		if n == array[i] {
			return true
		}
	}

	return false
}

// InUint64Array n是否存在指定数组中
func InUint64Array(array []uint64, n uint64) bool {
	for i := 0; i < len(array); i++ {
		if n == array[i] {
			return true
		}
	}

	return false
}

// GetWeekRange 获取指定周日期返回，i表示查询第几周，i=0表示查询本周
func GetWeekRange(i int) (string, string) {
	now := time.Now()
	week := int(now.Weekday())
	if week == 0 {
		week = 7
	}
	begin := now.Add(time.Hour * (-24) * time.Duration(week-1))
	end := now.Add(time.Hour * (24) * time.Duration(7-week))

	begin = begin.Add(time.Hour * (24) * 7 * time.Duration(i))
	end = end.Add(time.Hour * (24) * 7 * time.Duration(i))

	return begin.Format("2006-01-02"), end.Format("2006-01-02")
}

func CopyInt64Array(src []int64) []int64 {
	if len(src) == 0 {
		return []int64{}
	}
	dst := make([]int64, len(src))
	copy(dst, src)

	return dst
}

func CopyUInt64Array(src []uint64) []uint64 {
	if len(src) == 0 {
		return []uint64{}
	}
	dst := make([]uint64, len(src))
	copy(dst, src)

	return dst
}

func CopyUInt32Array(src []uint32) []uint32 {
	if len(src) == 0 {
		return []uint32{}
	}
	dst := make([]uint32, len(src))
	copy(dst, src)

	return dst
}

func EqualUInt32Array(arr1 []uint32, arr2 []uint32) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := 0; i < len(arr1); i++ {
		ok := false
		for j := 0; j < len(arr2); j++ {
			if arr1[i] == arr2[j] {
				ok = true
				break
			}
		}

		if !ok {
			return false
		}
	}

	return true
}

// CopyInt32Array 复制一个数组
func CopyInt32Array(src []int32) []int32 {
	if len(src) == 0 {
		return []int32{}
	}
	dst := make([]int32, len(src))
	copy(dst, src)

	return dst
}

// EqualInt32Array 两个数组是否相等
func EqualInt32Array(arr1 []int32, arr2 []int32) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := 0; i < len(arr1); i++ {
		ok := false
		for j := 0; j < len(arr2); j++ {
			if arr1[i] == arr2[j] {
				ok = true
				break
			}
		}

		if !ok {
			return false
		}
	}

	return true
}

// IsSameDay 两个时间戳是否同一天
func IsSameDay(t1, t2 uint64) bool {
	return time.Unix(int64(t1), 0).Format("2006-01-02") == time.Unix(int64(t2), 0).Format("2006-01-02")
}

// UnquieInsert ...
func UnquieInsert(arr *[]string, v ...string) bool {
	count := 0
	for i := 0; i < len(v); i++ {
		find := false
		for j := 0; j < len(*arr); j++ {
			if (*arr)[j] == v[i] {
				find = true
				break
			}
		}

		if !find {
			*arr = append(*arr, v[i])
			count++
		}
	}

	return count == 0
}

// UnquieInsertUInt64 ...
func UnquieInsertUInt64(arr *[]uint64, v ...uint64) {
	for i := 0; i < len(v); i++ {
		find := false
		for j := 0; j < len(*arr); j++ {
			if (*arr)[j] == v[i] {
				find = true
				break
			}
		}

		if !find {
			*arr = append(*arr, v[i])
		}
	}
}

// UnquieInsertUInt32 ...
func UnquieInsertUInt32(arr *[]uint32, v ...uint32) {
	for i := 0; i < len(v); i++ {
		find := false
		for j := 0; j < len(*arr); j++ {
			if (*arr)[j] == v[i] {
				find = true
				break
			}
		}

		if !find {
			*arr = append(*arr, v[i])
		}
	}
}

// UnquieInsertInt32 ...
func UnquieInsertInt32(arr *[]int32, v ...int32) {
	for i := 0; i < len(v); i++ {
		find := false
		for j := 0; j < len(*arr); j++ {
			if (*arr)[j] == v[i] {
				find = true
				break
			}
		}

		if !find {
			*arr = append(*arr, v[i])
		}
	}
}

// RemoveUInt64Array ...
func RemoveUInt64Array(arr *[]uint64, v ...uint64) {
	for i := 0; i < len(v); i++ {
		index := -1
		for j := 0; j < len(*arr); j++ {
			if (*arr)[j] == v[i] {
				index = j
				break
			}
		}

		if index >= 0 {
			*arr = append((*arr)[:index], (*arr)[index+1:]...)
		}
	}
}

// RemBit ...
func RemBit(mask, i uint32) uint32 {
	return mask ^ (1 << i)
}

// SetBit ...
func SetBit(mask, i uint32) uint32 {
	return (1 << i) | mask
}

// GetBit ...
func GetBit(mask, i uint32) bool {
	return mask&(1<<i) != 0
}

// In 判断元素是否在数组中
func In(arr interface{}, v interface{}) (bool, error) {
	sVal := reflect.ValueOf(arr)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == v {
				return true, nil
			}
		}

		return false, nil
	}

	return false, fmt.Errorf("unsupport type")
}

// Uint32ArrayRepeat ...
func Uint32ArrayRepeat(arr1, arr2 []uint32) bool {
	for i := 0; i < len(arr1); i++ {
		for j := 0; j < len(arr2); j++ {
			if arr1[i] == arr2[j] {
				return true
			}
		}
	}

	return false
}

// Uint64ArratToString ...
func Uint64ArratToString(arr []uint64, sep string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(arr); i++ {
		if i > 0 {
			buffer.WriteString(sep)
		}

		buffer.WriteString(strconv.FormatUint(arr[i], 10))
	}

	return buffer.String()
}

// Int64ArratToString ...
func Int64ArratToString(arr []int64, sep string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(arr); i++ {
		if i > 0 {
			buffer.WriteString(sep)
		}

		buffer.WriteString(strconv.FormatInt(arr[i], 10))
	}

	return buffer.String()
}

// 2019-01-01 00:00:00 时间戳
const (
	ChinaLocTime   = 1546272000
	VietnamLocTime = 1546275600
	// HoChiMinhLocKey 越南时区
	HoChiMinhLocKey = "Asia/Ho_Chi_Minh"
)

// GetVietnamDay 获取越南从2019-01-01 00:00:00到ts经历的天数
func GetVietnamDay(ts uint64) uint64 {
	return (ts - VietnamLocTime) / (24 * 3600)
}

// GetVietnamDate 获取越南日期格式
func GetVietnamDate(ts uint64) string {
	HoChiMinhLoc, _ := time.LoadLocation(HoChiMinhLocKey)
	return time.Unix(int64(ts), 0).In(HoChiMinhLoc).Format("2006-01-02")
}

// GetVietnamDateTime 获取越南日期格式
func GetVietnamDateTime(ts uint64) string {
	HoChiMinhLoc, _ := time.LoadLocation(HoChiMinhLocKey)
	return time.Unix(int64(ts), 0).In(HoChiMinhLoc).Format("2006-01-02 15:04:05")
}

// GetVietnamBeginAndEndDateTime 获取指定日期越南开始结束日期的本地时间
func GetVietnamBeginAndEndDateTime(ts uint64) (string, string) {
	HoChiMinhLoc, _ := time.LoadLocation(HoChiMinhLocKey)

	now := time.Now()
	beginDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, HoChiMinhLoc)
	beginDayStr := beginDay.In(time.Local).Format("2006-01-02 15:04:05")
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, HoChiMinhLoc)
	endDateStr := endDate.In(time.Local).Format("2006-01-02 15:04:05")

	return beginDayStr, endDateStr
}

// JoinUint64Array ...
func JoinUint64Array(arr []uint64, sep string) string {
	var b strings.Builder
	for i, u := range arr {
		if i != 0 {
			b.WriteString(sep)
		}
		b.WriteString(fmt.Sprintf("%v", u))
	}

	return b.String()
}

// GetVersionFromStr ...
func GetVersionFromStr(str string) []int {
	var ret []int
	arr := strings.Split(str, ".")
	for _, s := range arr {
		i, _ := strconv.Atoi(s)
		ret = append(ret, i)
	}

	return ret
}

// CompareVersion 比较v1,v2版本,.1:v1>v2 0:v1=v2 -1:v1<v2
func CompareVersion(v1, v2 string) int {
	a1 := GetVersionFromStr(v1)
	a2 := GetVersionFromStr(v2)
	for i := 0; i < 3; i++ {
		if a1[i] > a2[i] {
			return 1
		} else if a1[i] < a2[i] {
			return -1
		}
	}
	return 0
}

func Upper(s string, n int) string {
	var buff strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if n > 0 {
			if b >= 'a' && b <= 'z' {
				b -= 'a' - 'A'
			}
			n--
		}
		buff.WriteByte(b)
	}
	return buff.String()
}

// ProtectCall 错误保护调用
func ProtectCall(f func(), failFunc func()) {
	fail := false
	defer func() {
		if err := recover(); err != nil {
			Error("错误信息", "err", err)
			Error("错误调用堆栈信息", "stack", string(debug.Stack()))
			fail = true
		}

		if fail && failFunc != nil {
			failFunc()
		}
	}()

	f()
}
