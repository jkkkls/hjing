package app

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jkkkls/hjing/layout"
	"github.com/jkkkls/hjing/utils"
	"github.com/spf13/cobra"
)

// CmdAddItf represents the new command.
var CmdAddItf = &cobra.Command{
	Use:   "add-itf <serviceName> <interfaceName>",
	Short: "Add a interface",
	Long:  "Add a interface using the repository template. Example: hjing add-svc <serviceName> <interfaceName>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("requires 2 args, example: hjing add-svc <serviceName> <interfaceName>")
		}

		if !isValidAppName(args[0]) || !isValidAppName(args[1]) {
			return fmt.Errorf("name is invalid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		svcName := args[0]
		itfName := args[1]
		upSvcName := strings.ToUpper(svcName[:1]) + svcName[1:]
		upItfName := strings.ToUpper(itfName[:1]) + itfName[1:]
		//服务文件是否存在
		svcFileName := "services/" + svcName + "/service.go"
		svcBuff, err := os.ReadFile(svcFileName)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}
		if bytes.Contains(svcBuff, []byte(upItfName+"(context *rpc.Context")) {
			cmd.Usage()
			log.Fatal("error: interface is exists")
		}

		// 从go.mod读取domain
		domain, err := getDomainFromGoMod()
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}

		//检查proto文件
		protoName := fmt.Sprintf("pb/%v.proto", svcName)
		buff, err := os.ReadFile(protoName)
		if err != nil {
			err = layout.CopyFile("app/proto.tpl", protoName, "{{domain}}", domain)
			if err != nil {
				cmd.Usage()
				color.Red(err.Error())
				return
			}
			buff, err = os.ReadFile(protoName)
			if err != nil {
				cmd.Usage()
				color.Red(err.Error())
				return
			}
		}

		//生成协议
		newContent := fmt.Sprintf(`message %vReq {}
message %vRsp {}
			`, upItfName, upItfName)

		buff = append(buff, []byte(newContent)...)
		err = os.WriteFile(protoName, buff, os.ModePerm)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}

		//生成go文件
		msg, err := utils.ExecCmd("", "protoc", "--gogu_out", "./", "--gogu_opt", "paths=source_relative", protoName)
		if err != nil {
			color.Red(msg, err)
			return
		}

		//生成接口
		svcContent := fmt.Sprintf(`func (service *%vService) %v(context *rpc.Context, req *pb.%vReq, rsp *pb.%vRsp) (ret uint16, err error) {
	return
	}`, upSvcName, upItfName, upItfName, upItfName)
		svcBuff = append(svcBuff, []byte(svcContent)...)
		os.WriteFile(svcFileName, svcBuff, 0644)

		newBuff, err := utils.ExecCmd("", "goimports", svcFileName)
		if err != nil {
			color.Red(err.Error())
			return
		}
		os.WriteFile(svcFileName, []byte(newBuff), 0644)
		// utils.ExecCmd("", "go", "fmt", svcFileName)

		color.Green("create interface[%v] for %v success", itfName, svcName)
	},
}
