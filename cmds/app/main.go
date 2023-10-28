package main

import (
	"github.com/jkkkls/hjing/rpc"
)

func main() {
	rpc.NewApp("app.yaml").WithRegister(func(app *rpc.App) error {
		// rpc.RegisterService(serviceName, service.servicename)
		//end register
		return nil
	}).Run()
}
