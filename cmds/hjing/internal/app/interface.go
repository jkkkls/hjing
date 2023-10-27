package app

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// CmdAddApp represents the new command.
var CmdAddItf = &cobra.Command{
	Use:   "add-itf",
	Short: "Add a interface",
	Long:  "Add a interface using the repository template. Example: hjing add-svc <ServiceName> <interfaceName>",
	Run: func(cmd *cobra.Command, args []string) {
		svcName := args[0]
		itfName := args[1]
		if !isValidAppName(svcName) || !isValidAppName(itfName) {
			log.Fatal("name is invalid")
		}
		//服务文件是否存在
		if _, err := os.Stat("services/" + svcName + "/service.go"); err != nil {
			log.Fatalf("service[%v] is not exists", svcName)
		}
		// upAppName = strings.ToUpper(appName[:1]) + appName[1:]

		// err := os.MkdirAll("apps/"+appName, os.ModePerm)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// err = layout.CopyFile("app/app.main.go.tpl", "apps/"+appName+"/main.go", "{{appName}}", appName)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}}
