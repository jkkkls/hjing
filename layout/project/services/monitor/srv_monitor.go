package monitor

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
	"{{projectName}}/pb"

	"github.com/jkkkls/hjing/rpc"
	"github.com/shirou/gopsutil/v3/process"
)

var NodeName string

// MonitorService 进程监控服务
type MonitorService struct {
	Name      string
	GitTag    string
	PcName    string
	BuildTime string
	GitSHA    string
	Time      uint64
}

func (service *MonitorService) NodeConn(name string)                  {}
func (service *MonitorService) NodeClose(name string)                 {}
func (service *MonitorService) OnEvent(eventName string, args ...any) {}

// Exit 退出处理
func (service *MonitorService) Exit() {}

// Run 服务启动函数
func (service *MonitorService) Run() error {
	NodeName = service.Name
	go func() {
		for {
			time.Sleep(10 * time.Second)
			process, err := GetProcessInfo()
			if err != nil {
				continue
			}
			process.Name = service.Name
			process.Time = uint64(time.Now().Unix())
			process.GitTag = service.GitTag
			process.PcName = service.PcName
			process.BuildTime = service.BuildTime
			process.RunTime = service.Time
			if len(service.GitSHA) > 8 {
				process.GitSHA = service.GitSHA[:8]
			}

			rpc.Call(rpc.EmptyContext(), "MonitorMgr.UpdateMonitor", &pb.UpdateMonitorReq{Process: process}, nil)
		}
	}()
	return nil
}

// GetProcessInfo 获取进程信息
func GetProcessInfo() (*pb.Process, error) {
	info := &pb.Process{Pid: int32(os.Getpid())}
	p, err := process.NewProcess(info.Pid)
	if err != nil {
		return nil, err
	}
	if v, err := p.Parent(); err == nil {
		info.ParentId = v.Pid
	}
	if v, err := p.MemoryInfo(); err == nil {
		info.Memory = v.RSS / 1024
	}
	if v, err := p.MemoryPercent(); err == nil {
		info.MemoryPercent = fmt.Sprintf("%.3f%%", v)
	}
	if v, err := p.CPUPercent(); err == nil {
		info.CPUPercent = fmt.Sprintf("%.3f%%", v)
	}
	if v, err := p.Username(); err == nil {
		info.Username = v
	}
	if v, err := p.Cmdline(); err == nil {
		info.Cmd = v
	}

	info.NumGoroutine = int32(runtime.NumGoroutine())
	info.OSthreads = int32(pprof.Lookup("threadcreate").Count())
	info.GOMAXPROCS = int32(runtime.GOMAXPROCS(0))
	info.CPUNum = int32(runtime.NumCPU())

	return info, nil
}
