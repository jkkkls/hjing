package main

import (
	"{{projectName}}/services/monitor"
	"{{projectName}}/services/web_backend"

	"github.com/jkkkls/hjing/rpc"
	// end import
)

func main() {
	rpc.NewApp("admin.yaml").WithRegister(func(app *rpc.App) error {
		//
		rpc.RegisterService("MonitorMgr", &monitor.MonitorMgrService{})
		// end register

		web_backend.RunWebServices()
		return nil
	}).Run()
}
