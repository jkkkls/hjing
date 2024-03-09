package rpc

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

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
	plugins      []AppParam
	services     []GxService
	configName   string
	logParams    []utils.LogParam
	rpcNodeParam []GxNodeParam
	iRegister    IRegister

	// 长链接模块路由转发表，key为cmd id，value为模块名(Player.Login或者Player_Login)
	// cmd.proto中定义的cmd格式为
	// enum RequestCmd {
	// 	RequestBegin = 0x0000;

	// 	Game_Heartbeat           = 0x0001; //心跳
	// 	Game_Login               = 0x0002; //登陆
	// 	Game_TestCmd             = 0x0003; //测试命令
	cmds map[int32]string
}

type AppParam func(*App) error

func NewApp(configName string) *App {
	if configName == "" {
		configName = *file
	}
	_, err := config.LoadConf(configName)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return &App{configName: configName}
}

// WithLogParam 向App结构体中的logParams切片追加日志参数，并返回修改后的App指针
func (app *App) WithLogParam(params ...utils.LogParam) *App {
	app.logParams = append(app.logParams, params...)
	return app
}

// WithGxNodeParam 初始化rpc节点参数
func (app *App) WithGxNodeParam(params ...GxNodeParam) *App {
	app.rpcNodeParam = append(app.rpcNodeParam, params...)
	return app
}

// WithPlugin 新增初始化插件，并返回修改后的App指针
// 可以通过WithPlugin初始化全局属性，例如数据库链接等
func (app *App) WithPlugin(f func(app *App) error) *App {
	app.plugins = append(app.plugins, f)
	return app
}

// WithRegister 函数接收一个函数作为参数，该函数接收一个App指针作为参数，返回一个error。
// WithRegister函数将传入的函数添加到App结构体的services切片中，并返回App指针。
func (app *App) WithRegister(service GxService) *App {
	app.services = append(app.services, service)
	return app
}

// WithIRegister 自定义服务发现模块
func (app *App) WithIRegister(iRegister IRegister) *App {
	app.iRegister = iRegister
	return app
}

// WithCmds 函数用于设置长链接模块路由转发表
func (app *App) WithCmds(cmds map[int32]string) *App {
	app.cmds = cmds
	return app
}

func (app *App) Run() {
	// 编译传进来的参数
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

	nodeConfig := config.ConfInstance

	// 打开pprof
	utils.RunMonitor()

	// 初始化日志
	utils.GosLogInit(nodeConfig.App.Name, nodeConfig.Log.Dir, nodeConfig.Log.Screen, nodeConfig.Log.Level, app.logParams...)

	// pid文件
	os.WriteFile("./"+nodeConfig.App.Name+".pid", []byte(strconv.Itoa(os.Getpid())), 0o666)

	app.iRegister, err = NewEtcdRegister(nodeConfig.Etcds...)
	if err != nil {
		fmt.Println(err)
		return
	}

	node, err := InitNode(app.iRegister, &NodeConfig{
		Id:       nodeConfig.App.Id,
		Nodename: nodeConfig.App.Name,
		Nodetype: nodeConfig.App.Type,
		Host:     nodeConfig.App.Host,
		Port:     nodeConfig.App.Port,
		Cmds:     app.cmds,
		HttpPort: nodeConfig.App.HttpPort,
	}, app.rpcNodeParam...)
	if err != nil {
		fmt.Println("初始化rpc node失败", err.Error())
		return
	}

	// 初始化插件
	for _, v := range app.plugins {
		err = v(app)
		if err != nil {
			fmt.Println("初始化失败", err.Error())
			return
		}
	}

	// 初始化服务注册模块
	if app.iRegister == nil {
		app.iRegister, err = NewEtcdRegister(nodeConfig.Etcds...)
		if err != nil {
			fmt.Println("初始化ETCD失败", err.Error())
			return
		}
	}
	app.iRegister.WatchNode(func(key string, info *config.NodeInfo) {
		arr := strings.Split(key, ":")
		name := arr[len(arr)-1]
		if info == nil {
			node.DisconnectNode(name)
		} else {
			node.ConnectNewNode(name)
		}
	})

	// 初始化服务
	for _, v := range app.services {
		node.RegisterService(v)
	}

	ctx := context.Background()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// rpc节点服务
	node.Start()

	select {
	case <-signalChan:
	case <-ctx.Done():
	}

	node.Exit()
}
