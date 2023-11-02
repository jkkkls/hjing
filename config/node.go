package config

const (
	NodeRegisterKey = "nodeRegister"
)

type FunctionInfo struct {
	Name string `json:"name,omitempty"`
	Cmd  uint16 `json:"cmd,omitempty"`
}

type ServiceInfo struct {
	Name string          `json:"name,omitempty"`
	Func []*FunctionInfo `json:"func,omitempty"`
}

// 节点注册函数
type NodeInfo struct {
	Id      uint64         `json:"id,omitempty"`
	Name    string         `json:"name,omitempty"`    //
	Type    string         `json:"type,omitempty"`    //
	Address string         `json:"address,omitempty"` //ip:port
	Service []*ServiceInfo `json:"service,omitempty"` //服务列表
	Region  uint32         `json:"region,omitempty"`
	Set     string         `json:"set,omitempty"`
}
