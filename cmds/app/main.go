package main

import (
	"github.com/jkkkls/hjing/utils"
)

func main() {
	utils.NewApp("app.yaml").WithRegister(func(app *utils.App) error {
		// rpc.RegisterService(serviceName, service.servicename)
		//end register
		return nil
	}).Run()
}
