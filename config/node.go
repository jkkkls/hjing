package config

const (
	NodeRegisterKey = "nodeRegister"
)

// 节点注册函数
type NodeInfo struct {
	Id       uint64   `json:"id,omitempty" yaml:"id"`           //
	Name     string   `json:"name,omitempty" yaml:"name"`       //
	Type     string   `json:"type,omitempty" yaml:"type"`       //
	Address  string   `json:"address,omitempty" yaml:"address"` //ip:port
	Services []string `json:"service,omitempty" yaml:"service"` //服务列表
	Region   uint32   `json:"region,omitempty" yaml:"region"`
	Set      string   `json:"set,omitempty" yaml:"set"`
}
