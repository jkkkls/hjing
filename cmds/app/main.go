package main

import (
	"github.com/jkkkls/hjing/utils"
)

func main() {
	utils.NewApp("app.yaml").WithRegister(func(app *utils.App) error {
		//end register
		return nil
	}).Run()
}
