package app

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jkkkls/hjing/layout"
	"github.com/jkkkls/hjing/utils"
	"github.com/spf13/cobra"
)

// CmdAddItf represents the new command.
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
		upSvcName := strings.ToUpper(svcName[:1]) + svcName[1:]
		upItfName := strings.ToUpper(itfName[:1]) + itfName[1:]
		//服务文件是否存在
		svcFileName := "services/" + svcName + "/service.go"
		svcBuff, err := os.ReadFile(svcFileName)
		if err != nil {
			log.Fatal(err)
		}
		if bytes.Contains(svcBuff, []byte(upItfName+"(context *rpc.Context")) {
			log.Fatal("error: interface is exists")
		}

		// 从go.mod读取domain
		domain, err := getDomainFromGoMod()
		if err != nil {
			log.Fatal(err)
		}

		//检查proto文件
		protoName := fmt.Sprintf("pb/%v.proto", svcName)
		buff, err := os.ReadFile(protoName)
		if err != nil {
			err = layout.CopyFile("app/proto.tpl", protoName, "{{domain}}", domain)
			if err != nil {
				log.Fatal(err)
			}
			buff, err = os.ReadFile(protoName)
			if err != nil {
				log.Fatal(err)
			}
		}

		//生成协议
		newContent := fmt.Sprintf(`message %vReq {}
message %vRsp {}
			`, upItfName, upItfName)

		buff = append(buff, []byte(newContent)...)
		err = os.WriteFile(protoName, buff, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		//生成go文件
		msg, err := utils.ExecCmd("", "protoc", "--gogu_out", "./", "--gogu_opt", "paths=source_relative", protoName)
		if err != nil {
			log.Fatal(msg, err)
		}

		//生成接口
		svcContent := fmt.Sprintf(`func (service *%vService) %v(context *rpc.Context, req *pb.%vReq, rsp *pb.%vRsp) (ret uint16, err error) {
	return			
	}`, upSvcName, upItfName, upItfName, upItfName)
		svcBuff = append(svcBuff, []byte(svcContent)...)
		os.WriteFile(svcFileName, svcBuff, 0644)

		newBuff, err := utils.ExecCmd("", "goimports", svcFileName)
		if err != nil {
			log.Fatal(err)
		}
		os.WriteFile(svcFileName, []byte(newBuff), 0644)
		// utils.ExecCmd("", "go", "fmt", svcFileName)
	}}
