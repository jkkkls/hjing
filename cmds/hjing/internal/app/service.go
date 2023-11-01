package app

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jkkkls/hjing/layout"
	"github.com/jkkkls/hjing/utils"
	"github.com/spf13/cobra"
)

var (
	svcRegMask    = []byte("//end register")
	svcImportMask = []byte("//end import")
	modelRegMask  = []byte("//end models")
)

func getDomainFromGoMod() (string, error) {
	f, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	buf.ReadFrom(f)
	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module") {
			return strings.Split(line, " ")[1], nil
		}
	}
	return "", fmt.Errorf("go.mod format error")
}

// CmdAddSrv represents the new command.
var CmdAddSrv = &cobra.Command{
	Use:   "add-svc <appName> <serviceName>",
	Short: "Create a template",
	Long:  "Create a services using the repository template. Example: hjing add-svc <appName> <serviceName>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("requires 2 args, example: hjing add-svc <appName> <serviceName>")
		}

		if !isValidAppName(args[0]) || !isValidAppName(args[1]) {
			return fmt.Errorf("name is invalid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		svcName := args[1]
		upSvcName := strings.ToUpper(svcName[:1]) + svcName[1:]

		mainFile := fmt.Sprintf("./apps/%v/main.go", appName)
		buff, err := os.ReadFile(mainFile)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}
		if bytes.Contains(buff, []byte(upSvcName)) {
			cmd.Usage()
			color.Red("service[%v] is exists", upSvcName)
			return
		}
		if !bytes.Contains(buff, svcRegMask) || !bytes.Contains(buff, svcImportMask) {
			cmd.Usage()
			color.Red("main.json format error")
			return
		}

		// 从go.mod读取domain
		domain, err := getDomainFromGoMod()
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}

		err = os.MkdirAll("services/"+svcName, os.ModePerm)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}

		//替换引用
		newContent := fmt.Sprintf(`"%v/services/%v"
		%v`,
			domain, svcName, string(svcImportMask))
		buff = bytes.ReplaceAll(buff, svcImportMask, []byte(newContent))

		//注册服务
		newContent = fmt.Sprintf(`rpc.RegisterService("%v", &%v.%vService{})
		%v`,
			upSvcName, svcName, upSvcName, string(svcRegMask))
		buff = bytes.ReplaceAll(buff, svcRegMask, []byte(newContent))
		os.WriteFile(mainFile, buff, 0644)

		utils.ExecCmd("", "go", "fmt", mainFile)

		//生成服务模版
		err = layout.CopyFile("app/service.go.tpl", "services/"+svcName+"/service.go", "{{lowServiceName}}", svcName, "{{serviceName}}", upSvcName)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}
		color.Green("create service[%v] success", svcName)
	},
}
