# 长链接网关开发说明

# 实例，新增gate app，修改main.go

``` go
func main() {
	rpc.NewApp("data.yaml").WithRegister(func(app *rpc.App) error {
		//
		rpc.RegisterService("Monitor", &monitor.MonitorService{
			Name:      config.ConfInstance.App.Name,
			GitSHA:    rpc.GitSHA,
			PcName:    rpc.PcName,
			BuildTime: rpc.BuildTime,
			GitTag:    rpc.GitTag,
			Time:      utils.Now(),
		})

		// end register
		return nil
	}).WithPlugin(func(app *rpc.App) error {
		utils.Go(func() {
			gater := net.NewGater(net.TcpParam(true, &net.ConnCoder{}))
			err := gater.RunTCPServer(config.GetInt("tcpPort"))
			if err != nil {
				panic(err)
			}
		})
		return nil
	}).WithCmds(map[int32]string{
		1: "Data.Get",
		2: "Data.Set",
	}).Run()
}
```