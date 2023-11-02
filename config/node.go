package config

const (
	NodeRegisterKey = "nodeRegister"
)

type FunctionInfo struct {
	Name string `json:"name,omitempty"`
	Cmd  uint16 `json:"cmd,omitempty"`
}

// 节点注册函数
type NodeInfo struct {
	Id       uint64   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`    //
	Type     string   `json:"type,omitempty"`    //
	Address  string   `json:"address,omitempty"` //ip:port
	Services []string `json:"service,omitempty"` //服务列表
	Region   uint32   `json:"region,omitempty"`
	Set      string   `json:"set,omitempty"`
}
