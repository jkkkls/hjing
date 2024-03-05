// Code generated by hjing. DO NOT EDIT.

package main

import (
	"github.com/jkkkls/hjing/config"
	"github.com/jkkkls/hjing/rpc"
	"github.com/jkkkls/hjing/utils"

	"{{projectName}}/services/monitor"
	//end import
)

func main() {
	rpc.NewApp("{{appName}}.yaml").WithRegister(func(app *rpc.App) error {
		//
		rpc.RegisterService("Monitor", &monitor.MonitorService{
			Name:      config.ConfInstance.App.Name,
			GitSHA:    rpc.GitSHA,
			PcName:    rpc.PcName,
			BuildTime: rpc.BuildTime,
			GitTag:    rpc.GitTag,
			Time:      utils.Now(),
		})

		//end register
		return nil
	}).Run()
}
