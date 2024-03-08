package monitor

import (
	"sync"
	"{{projectName}}/pb"

	"github.com/jkkkls/hjing/rpc"
	"google.golang.org/protobuf/proto"
)

var (
	Processes sync.Map
	Online    sync.Map
)

// MonitorMgrService 进程监控管理服务
type MonitorMgrService struct{}

func (service *MonitorMgrService) GetName() string                       { return "MonitorMgr" }
func (service *MonitorMgrService) NodeConn(name string)                  {}
func (service *MonitorMgrService) NodeClose(name string)                 {}
func (service *MonitorMgrService) OnEvent(eventName string, args ...any) {}

func NewMonitorMgrService() *MonitorMgrService {
	return &MonitorMgrService{}
}

// Exit 退出处理
func (service *MonitorMgrService) Exit() {}

// Run 服务启动函数
func (service *MonitorMgrService) Run() error {
	return nil
}

func (service *MonitorMgrService) UpdateMonitor(conn *rpc.Context, req *pb.UpdateMonitorReq, rsp proto.Message) (uint16, error) {
	Processes.Store(req.Process.Name, req.Process)
	return 0, nil
}

func GetProcess() []*pb.Process {
	var processes []*pb.Process
	Processes.Range(func(key, value interface{}) bool {
		processes = append(processes, value.(*pb.Process))
		return true
	})

	return processes
}
