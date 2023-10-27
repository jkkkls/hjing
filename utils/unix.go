package utils

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/seveye/goms/util/sys"
)

var threadMap sync.Map

type ReqEventCall struct {
	Uid    uint64
	CallId uint64
}

// SetCallId : 通过线程id保存事件id，但一旦用异步时将失效（go）
func SetCallId(uid, callId uint64) {
	runtime.LockOSThread()
	tid := sys.GetThreadId()
	threadMap.Store(tid, &ReqEventCall{Uid: uid, CallId: callId})
}

func UnLockCallId() {
	runtime.UnlockOSThread()
}

// GetCallId : 通过线程id获取事件id，但一旦用异步时将失效（go）
func GetCallId() uint64 {
	tid := sys.GetThreadId()
	if id, ok := threadMap.Load(tid); ok {
		return id.(*ReqEventCall).CallId
	}
	return 0
}

// GetReqEventCall : return uid and callId
func GetReqEventCall() (uint64, uint64) {
	tid := sys.GetThreadId()
	if data, ok := threadMap.Load(tid); ok {
		if v, ok := data.(*ReqEventCall); ok {
			return v.Uid, v.CallId
		}
	}
	return 0, 0
}

// GetRunPath : 获取程序运行的文件路径，skip 可跳过最下条数
func GetRunPath(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		path := strings.Split(file, "go/src")
		if len(path) == 2 {
			return fmt.Sprintf("%s:%d", path[1], line)
		}
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "???"
}
