package utils

// Copyright 2017 guangbo. All rights reserved.

// 常用接口

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/constraints"
)

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
func If[T any](x bool, a T, b T) T {
	if x {
		return a
	}

	return b
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

// IsSameDay 两个时间戳是否同一天
func IsSameDay(t1, t2 uint64) bool {
	return time.Unix(int64(t1), 0).Format("2006-01-02") == time.Unix(int64(t2), 0).Format("2006-01-02")
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

// Uint64ArratToString ...
func ArrayToString[T any](arr []T, sep string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(arr); i++ {
		if i > 0 {
			buffer.WriteString(sep)
		}

		buffer.WriteString(fmt.Sprint(arr[i]))
	}

	return buffer.String()
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

// FixRange 矫正数值, 必须大于等于下限，小于等于上限
func FixRange[T constraints.Integer | constraints.Float](arr ...T) T {
	var empty T
	if len(arr) == 0 {
		return empty
	} else if len(arr) == 1 {
		return arr[0]
	} else if len(arr) >= 2 && arr[0] < arr[1] {
		return arr[1]
	} else if len(arr) >= 3 && arr[0] > arr[2] {
		return arr[2]
	}

	return arr[0]
}

func ExecCmd(dir, cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	command.Dir = dir
	// 给标准输入以及标准错误初始化一个buffer，每条命令的输出位置可能是不一样的，
	// 比如有的命令会将输出放到stdout，有的放到stderr
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}

	err := command.Run()
	if err != nil {
		// 打印程序中的错误以及命令行标准错误中的输出
		return command.Stderr.(*bytes.Buffer).String(), err
	}
	// 打印命令行的标准输出
	return command.Stdout.(*bytes.Buffer).String(), nil
}

// Pointer 返回指针
func Pointer[T any](t T, returnNull bool) *T {
	if returnNull {
		return nil
	}
	return &t
}

// SwapHander 返回json格式的handler
func SwapHandler(f func(c *gin.Context) (ok bool, msg string)) func(c *gin.Context) {
	return func(c *gin.Context) {
		ok, msg := f(c)
		if msg != "" {
			log.Println("--------------->", c.Request.URL.String(), ok, msg)
			c.JSON(http.StatusOK, gin.H{
				"success":      ok,
				"errorMessage": msg,
			})
		}
	}
}

func GetName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// GenArray 将一个类型的数组转换为另一个类型的数组
func GenArray[T1 any, T2 any](l []T2, f func(T2) T1) []T1 {
	var ret []T1
	for _, v := range l {
		ret = append(ret, f(v))
	}
	return ret
}

// Array 将一个类型的数组转换为一个map
func GenMap[K comparable, V any, T2 any](l []T2, f func(T2) (K, V)) map[K]V {
	ret := make(map[K]V)
	for _, v := range l {
		k, v := f(v)
		ret[k] = v
	}
	return ret
}

// Download 下载
func Download(url string) ([]byte, error) {
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode() == http.StatusNotFound {
		return nil, nil
	}
	return rsp.Body(), nil
}
