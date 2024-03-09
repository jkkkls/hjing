package main

import (
	"github.com/jkkkls/hjing/rpc"

	"{{projectName}}/services/monitor"
	// end import
)

func main() {
	rpc.NewApp("{{appName}}.yaml").
		WithRegister(monitor.NewMonitorService()).
		// end register
		Run()
}
