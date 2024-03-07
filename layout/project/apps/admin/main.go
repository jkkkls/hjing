package main

import (
	"github.com/jkkkls/test_app/services/monitor"
	"github.com/jkkkls/test_app/services/web_backend"

	"github.com/jkkkls/hjing/rpc"
	// end import
)

func main() {
	rpc.NewApp("admin.yaml").
		WithRegister(monitor.NewMonitorMgrService()).
		WithPlugin(func(app *rpc.App) error {
			web_backend.RunWebServices()
			return nil
		}).Run()
}
