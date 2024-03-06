package utils

/**
作者:guangbo
模块：日志模块
说明：
创建时间：2015-10-30
**/

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// log 异步 多pipe
const (
	// TraceLevel Trace级别
	TraceLevel = iota
	// DebugLevel Debug级别
	DebugLevel
	// WarnLevel Warn级别
	WarnLevel
	// ErrorLevel Error级别
	ErrorLevel
	// InfoLevel Info级别
	InfoLevel
	// FatalLevel Fatal级别
	FatalLevel
)

type GosLogContent struct {
	Text  string
	Now   time.Time
	File  string
	Line  int
	Level string
	Param []interface{}
}

type GosLog struct {
	Name    string
	Level   int
	LogFile *os.File
	Time    string
	Logs    chan *GosLogContent

	plugin []func(map[string]any)
}

// LogColor 颜色标签
type LogColor struct {
	LevelLeft  string
	LevelRight string
}

var (
	logPath         string
	logOutputScreen bool
	logScreenCache  chan *GosLogContent
	logInit         bool = false
	defaultLog      *GosLog
)

var logColor map[string]*LogColor = map[string]*LogColor{
	"TRAC": {LevelLeft: "", LevelRight: ""},
	"DEBU": {LevelLeft: "\033[35m", LevelRight: "\033[0m"},
	"WARN": {LevelLeft: "\033[34m", LevelRight: "\033[0m"},
	"ERRO": {LevelLeft: "\033[31m", LevelRight: "\033[0m"},
	"INFO": {LevelLeft: "\033[37m", LevelRight: "\033[0m"},
	"FATA": {LevelLeft: "\033[31m", LevelRight: "\033[0m"},
}

var logText = []string{"TRAC", "DEBU", "WARN", "ERRO", "INFO", "FATA"}

func init() {
	logScreenCache = make(chan *GosLogContent, 1024)
}

type LogParam func(*GosLog)

func WithPlugin(f func(map[string]any)) LogParam {
	return func(l *GosLog) {
		l.plugin = append(l.plugin, f)
	}
}

// InitLogger 日志模块初始化函数,程序启动时调用
func GosLogInit(name, path string, screen bool, level int, params ...LogParam) {
	// 只初始化一次
	if logInit {
		return
	}

	logOutputScreen = screen
	logPath = path
	if logPath == "" {
		logPath = "./log/"
	}
	logInit = true

	if !PathExists(logPath) {
		CreateDir(logPath)
	}

	defaultLog = GosLogNewRouter(name)
	defaultLog.Level = level

	for _, v := range params {
		v(defaultLog)
	}

	if !logOutputScreen {
		return
	}

	go func() {
		for {
			select {
			case info := <-logScreenCache:
				if info == nil {
					return
				}
				fmt.Printf(info.textFormat())
			}
		}
	}()
}

func SetLogLevel(level int) {
	defaultLog.Level = level
}

func getLogTime() string {
	return time.Now().Format("2006-01-02")
}

func GosLogNewRouter(name string) *GosLog {
	if !logInit {
		return nil
	}

	timeStr := getLogTime()
	file := logPath + "/" + name + ".log." + timeStr
	logFile, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
	if err != nil {
		return nil
	}
	l := &GosLog{
		Name:    name,
		Level:   TraceLevel,
		LogFile: logFile,
		Time:    timeStr,
		Logs:    make(chan *GosLogContent),
	}

	go func(gl *GosLog) {
		for {
			select {
			case info := <-l.Logs:
				if info == nil {
					return
				}
				gl.checkTime()
				data, m := info.jsonFormat()
				gl.LogFile.Write(data)
				gl.LogFile.WriteString("\n")

				Go(func() {
					for _, v := range gl.plugin {
						v(m)
					}
				})

				// if gl.ZincDomain != "" {
				// 	Submit(func() {
				// 		err := WriteZinc(gl.ZincDomain, gl.ZincUsername, gl.ZincPwd, gl.ZincIndex, data2)
				// 		if err != nil {
				// 			log.Println("WriteZinc", err)
				// 		}
				// 	})
				// }
			}
		}
	}(l)

	return l
}

func (info *GosLogContent) textFormat() string {
	short := info.File
	for i := len(info.File) - 1; i > 0; i-- {
		if info.File[i] == '/' {
			short = info.File[i+1:]
			break
		}
	}
	param := ""
	l := len(info.Param) / 2
	for i := 0; i < l; i++ {
		if i > 0 {
			param += " "
		}
		param += fmt.Sprintf("%v%v%v=%+v",
			logColor[info.Level].LevelLeft, info.Param[i*2], logColor[info.Level].LevelRight,
			info.Param[i*2+1])
	}
	year, month, day := info.Now.Date()
	hour, min, sec := info.Now.Clock()

	level := logColor[info.Level].LevelLeft + info.Level + logColor[info.Level].LevelRight
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d %s %s:%d %s [%s]\n", //
		year, month, day, hour, min, sec, info.Now.Nanosecond()/1000, level, short, info.Line, info.Text, param)
}

func (info *GosLogContent) jsonFormat() ([]byte, map[string]any) {
	short := info.File
	for i := len(info.File) - 1; i > 0; i-- {
		if info.File[i] == '/' {
			short = info.File[i+1:]
			break
		}
	}
	year, month, day := info.Now.Date()
	hour, min, sec := info.Now.Clock()

	m := make(map[string]any)
	m["time"] = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d", year, month, day, hour, min, sec, info.Now.Nanosecond()/1000)
	m["line"] = fmt.Sprintf("%s:%d", short, info.Line)
	m["level"] = fmt.Sprintf("%s", info.Level)
	m["msg"] = info.Text

	l := len(info.Param) / 2
	for i := 0; i < l; i++ {
		s, ok := info.Param[i*2].(string)
		if ok {
			m[s] = info.Param[i*2+1]
		}
	}

	buff, err := json.Marshal(m)
	if err != nil {
		fmt.Println("log error", err, m["line"])
	}

	return buff, m
}

func (l *GosLog) SetLevel(level int) {
	if level > FatalLevel || level < TraceLevel {
		l.Level = TraceLevel
	} else {
		l.Level = level
	}
}

func (l *GosLog) Free() {
	l.Logs <- nil
}

func (l *GosLog) checkTime() {
	temp := getLogTime()
	if temp == l.Time {
		return
	}
	l.Time = temp
	l.LogFile.Close()

	file := logPath + "/" + l.Name + ".log." + l.Time
	l.LogFile, _ = os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
}

// Log 不是很重要类型日志
func (l *GosLog) Log(level int, skip int, text string, v ...interface{}) {
	if l.Level > level {
		return
	}

	info := &GosLogContent{
		Text:  text,
		Now:   time.Now(),
		Level: logText[level],
		Param: v,
	}
	var ok bool
	_, info.File, info.Line, ok = runtime.Caller(skip)
	if !ok {
		info.File = "???"
		info.Line = 0
	}

	l.Logs <- info

	if logOutputScreen {
		logScreenCache <- info
	}
}

// Trace 不是很重要类型日志
func (l *GosLog) Trace(text string, v ...interface{}) {
	l.Log(TraceLevel, 3, text, v...)
}

// TraceSkip 不是很重要类型日志
func (l *GosLog) TraceSkip(skip int, text string, v ...interface{}) {
	l.Log(TraceLevel, skip, text, v...)
}

// Debug 调试类型日志
func (l *GosLog) Debug(text string, v ...interface{}) {
	l.Log(DebugLevel, 3, text, v...)
}

// DebugSkip 调试类型日志
func (l *GosLog) DebugSkip(skip int, text string, v ...interface{}) {
	l.Log(DebugLevel, skip, text, v...)
}

// Warn 警告类型日志
func (l *GosLog) Warn(text string, v ...interface{}) {
	l.Log(WarnLevel, 3, text, v...)
}

// WarnSkip 警告类型日志
func (l *GosLog) WarnSkip(skip int, text string, v ...interface{}) {
	l.Log(WarnLevel, skip, text, v...)
}

// Error 错误类型日志
func (l *GosLog) Error(text string, v ...interface{}) {
	l.Log(ErrorLevel, 3, text, v...)
}

// ErrorSkip 错误类型日志
func (l *GosLog) ErrorSkip(skip int, text string, v ...interface{}) {
	l.Log(ErrorLevel, skip, text, v...)
}

// Info 程序信息类型日志
func (l *GosLog) Info(text string, v ...interface{}) {
	l.Log(InfoLevel, 3, text, v...)
}

// InfoSkip 程序信息类型日志
func (l *GosLog) InfoSkip(skip int, text string, v ...interface{}) {
	l.Log(InfoLevel, skip, text, v...)
}

// Fatal 致命错误类型日志
func (l *GosLog) Fatal(text string, v ...interface{}) {
	l.Log(FatalLevel, 3, text, v...)
}

// FatalSkip 致命错误类型日志
func (l *GosLog) FatalSkip(skip int, text string, v ...interface{}) {
	l.Log(FatalLevel, skip, text, v...)
}

func Trace(text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.Trace(text, v...)
}

func Debug(text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.Debug(text, v...)
}

func Info(text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.Info(text, v...)
}

func Warn(text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.Warn(text, v...)
}

func Error(text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.Error(text, v...)
}

func Fatal(text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.Fatal(text, v...)
}

func TraceSkip(skip int, text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.TraceSkip(skip, text, v...)
}

func DebugSkip(skip int, text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.DebugSkip(skip, text, v...)
}

func InfoSkip(skip int, text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.InfoSkip(skip, text, v...)
}

func WarnSkip(skip int, text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.WarnSkip(skip, text, v...)
}

func ErrorSkip(skip int, text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.ErrorSkip(skip, text, v...)
}

func FatalSkip(skip int, text string, v ...interface{}) {
	if defaultLog == nil {
		log.Println(text, v)
		return
	}
	defaultLog.FatalSkip(skip, text, v...)
}
