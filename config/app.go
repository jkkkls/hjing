package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// 基础配置
type App struct {
	Name      string `yaml:"name"`
	Desc      string `yaml:"desc"`
	PprofAddr string `yaml:"pprofAddr"` //"0.0.0.0:9061"
	Pro       bool   `yaml:"pro"`

	Version   string
	GitSHA    string
	PcName    string
	BuildTime string
}

type Log struct {
	Level  int    `yaml:"level"`
	Dir    string `yaml:"dir"`
	Screen bool   `yaml:"screen"`
}

// RpcNode配置
type RpcNode struct {
	Id   uint64 `yaml:"id"`
	Name string `yaml:"name"`
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`

	Type   string `yaml:"type"`
	Weight int    `yaml:"weight"`
}

// gin配置
type Gin struct {
	Port int `yaml:"port"`
}

// fast配置
type Fast struct {
	Port int `yaml:"port"`
}

type RpcNodes struct {
	Nodes map[string]RpcNode `yaml:"nodes"`
}

type NodeConf struct {
	App      App         `yaml:"app"`
	Gin      Gin         `yaml:"gin"`
	Fast     Fast        `yaml:"fast"`
	Log      Log         `yaml:"log"`
	Custom   map[any]any `yaml:"custom"`
	RpcNodes RpcNodes    `yaml:"-"`
}

var ConfInstance *NodeConf

func UpdateApp(f func(c *NodeConf)) {
	f(ConfInstance)
}

func LoadGlobalConf(global string) error {
	ConfInstance = &NodeConf{}
	//
	buff, err := os.ReadFile(global)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buff, &ConfInstance.RpcNodes)
	if err != nil {
		return err
	}

	return nil
}

func loadConf(global, app string) (*NodeConf, error) {
	tmpIns := ConfInstance
	//
	temp := &NodeConf{}
	//
	buff, err := os.ReadFile(app)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buff, temp)
	if err != nil {
		return nil, err
	}
	//
	buff, err = os.ReadFile(global)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buff, &temp.RpcNodes)
	if err != nil {
		return nil, err
	}

	if tmpIns != nil {
		temp.App.Version = tmpIns.App.Version
		temp.App.PcName = tmpIns.App.PcName
		temp.App.BuildTime = tmpIns.App.BuildTime
		temp.App.GitSHA = tmpIns.App.GitSHA
	}

	ConfInstance = temp
	return ConfInstance, nil
}

func LoadConf(global, app string) (*NodeConf, error) {

	return loadConf(global, app)
}

type C struct {
	a any
}

func (c *C) Get(field any) *C {
	if k, ok := field.(string); ok {
		c.a = c.a.(map[any]any)[k]
	} else {
		i := field.(int)
		c.a = c.a.([]any)[i]
	}

	return c
}

// Get 获取字符串, field只能是string int，int表示数组位置
func GetString(fields ...any) string {
	c := &C{a: ConfInstance.Custom}
	for _, field := range fields {
		c = c.Get(field)
	}
	v, ok := c.a.(string)
	if !ok {
		return ""
	}

	return v
}

// GetSize 获取数组大小, field只能是string int，int表示数组位置
func GetSize(fields ...any) int {
	c := &C{a: ConfInstance.Custom}
	for _, field := range fields {
		c = c.Get(field)
	}
	arr, ok := c.a.([]any)
	if !ok {
		return 0
	}
	return len(arr)
}

// Get 获取数值, field只能是string int，int表示数组位置
func GetInt(fields ...any) int {
	c := &C{a: ConfInstance.Custom}
	for _, field := range fields {
		c = c.Get(field)
	}

	v, ok := c.a.(int)
	if !ok {
		return 0
	}

	return v
}

// GetRpcNode 返回节点grpc信息
func GetRpcNode() *RpcNode {
	v, ok := ConfInstance.RpcNodes.Nodes[ConfInstance.App.Name]
	if !ok {
		return nil
	}
	return &v
}
