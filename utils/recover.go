// Copyright 2017 guangbo. All rights reserved.
package utils

import (
	"runtime/debug"

	"github.com/google/gops/agent"
)

type NginxInfo struct {
	Address       string
	FileDir       string
	PlayerIconDir string
	ClubIconDir   string
	ErrorDir      string
}

// Recover ...
func Recover() {
	if err := recover(); err != nil {
		Error("错误信息", "err", err)
		Error("错误堆栈", "stack", string(debug.Stack()))
	}
}

// RunMonitor 运行节点监控逻辑
func RunMonitor() {
	if err := agent.Listen(agent.Options{}); err != nil {
		Error("RunMonitor error", "err", err)
	}
}

func CheckError(info string, err error) {
	if err == nil {
		return
	}

	Warn(info, "error", err)
}
