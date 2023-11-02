package rpc

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/jkkkls/hjing/config"
	"github.com/jkkkls/hjing/utils"
)

var (
	GitTag    = "2000.01.01.release"
	PcName    = "anyone"
	BuildTime = "2000-01-01T00:00:00+0800"
	GitSHA    = "xxx"
	Area      = ""
	file      = flag.String("config", "", "配置名")
)

type App struct {
	plugins    []AppParam
	services   []AppParam
	configName string
	area       string
}

type AppParam func(*App) error

func NewApp(configName string) *App {
	return &App{configName: configName}
}

func (app *App) WithPlugin(f func(app *App) error) *App {
	app.plugins = append(app.plugins, f)
	return app
}

func (app *App) WithRegister(f func(app *App) error) *App {
	app.services = append(app.services, f)
	return app
}

func (app *App) Run() {
	//编译传进来的参数
	fmt.Println("Version      : " + GitTag)
	fmt.Println("Git SHA      : " + GitSHA)
	fmt.Println("Build PC     : " + PcName)
	fmt.Println("Build Time   : " + BuildTime)
	fmt.Println("Build Area   : " + Area)
	flag.Parse()

	err := utils.SetUlimit()
	if err != nil {
		fmt.Println(err)
	}

	var configName string
	if *file != "" {
		configName = *file
	} else {
		configName = app.configName
	}

	nodeConfig, err := config.LoadConf(configName)
	if err != nil {
		fmt.Println(err)
		return
	}

	//打开pprof
	utils.RunMonitor()

	// 初始化日志
	utils.GosLogInit(nodeConfig.App.Name, nodeConfig.Log.Dir, nodeConfig.Log.Screen, nodeConfig.Log.Level)

	//pid文件
	os.WriteFile("./"+nodeConfig.App.Name+".pid", []byte(strconv.Itoa(os.Getpid())), 0666)

	//
	for _, v := range app.plugins {
		err = v(app)
		if err != nil {
			fmt.Println("初始化失败", err.Error())
			return
		}
	}

	//连接etcd
	iRegister, err := NewEtcdRegister(nodeConfig.Etcds...)
	if err != nil {
		fmt.Println("初始化ETCD失败", err.Error())
		return
	}
	iRegister.WatchNode(func(key string, info *config.NodeInfo) {
		if info == nil {
			DisconnectNode(key)
		} else {
			ConnectNewNode(key)
		}
	})

	for _, v := range app.services {
		err = v(app)
		if err != nil {
			fmt.Println("初始化服务失败", err.Error())
			return
		}
	}

	//rpc节点服务
	// api.RegisterRpcCallBack()
	InitNode(iRegister, &NodeConfig{
		Id:       nodeConfig.App.Id,
		Nodename: nodeConfig.App.Name,
		Nodetype: nodeConfig.App.Type,
		Host:     nodeConfig.App.Host,
		Port:     nodeConfig.App.Port,
		// Cmds:     pb.RequestCmd_value,
		// HttpPort: nodeConfig.App.HttpPort,
	})
}
