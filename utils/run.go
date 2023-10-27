package utils

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/jkkkls/hjing/config"
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

	err := SetUlimit()
	if err != nil {
		fmt.Println(err)
	}

	var configName string
	if *file != "" {
		configName = *file
	} else {
		configName = app.configName
	}

	nodeConfig, err := config.LoadConf("global.yaml", configName)
	if err != nil {
		fmt.Println(err)
		return
	}

	//打开pprof
	RunMonitor()

	// 初始化日志
	GosLogInit(nodeConfig.App.Name, nodeConfig.Log.Dir, nodeConfig.Log.Screen, nodeConfig.Log.Level)

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

	//连接watch
	// client, err := watch.NewWatchClient(nodeConfig.Watch.Host)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// client.RegisterCallback(watch_config.NodeRegisterKey+":", rpc.WatchNodeRegister)
	// if nodeConfig.Node.Set != "" {
	// 	client.RegisterCallback(fmt.Sprintf("%v:%v", watch_config.NodeRegisterKey, nodeConfig.Node.Set), rpc.WatchNodeRegister)
	// }
	// client.Start()

	for _, v := range app.services {
		err = v(app)
		if err != nil {
			fmt.Println("初始化服务失败", err.Error())
			return
		}
	}

	//rpc节点服务
	// api.RegisterRpcCallBack()
	// rpcConf := config.GetRpcNode()
	// rpc.InitNode(&rpc.NodeConfig{
	// 	// Client:   client,
	// 	Id:       rpcConf.Id,
	// 	Nodename: nodeConfig.App.Name,
	// 	Nodetype: rpcConf.Type,
	// 	// Set:      rpcConf.Set,
	// 	Host:   rpcConf.Ip,
	// 	Port:   rpcConf.Port,
	// 	Region: 0,
	// 	// Cmds:     pb.RequestCmd_value,
	// 	// HttpPort: nodeConfig.App.HttpPort,
	// })
}
