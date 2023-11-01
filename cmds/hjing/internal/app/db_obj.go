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

// CmdAddDb represents the new command.
var CmdAddDb = &cobra.Command{
	Use:   "add-db <dbObjName>",
	Short: "Add a db model",
	Long:  "Add a db model using the repository template. Example: hjing add-db <dbObjName>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires 2 args, example: hjing add-db <dbObjName>")
		}

		if !isValidAppName(args[0]) {
			return fmt.Errorf("name is invalid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		objName := args[0]
		upObjName := strings.ToUpper(objName[:1]) + objName[1:]
		//服务文件是否存在
		objFileName := "db/models.go"

		objBuff, err := os.ReadFile(objFileName)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}
		if bytes.Contains(objBuff, []byte("&pb."+upObjName+"{}")) {
			cmd.Usage()
			color.Red("error: obj[%v] is exists", upObjName)
			return
		}

		// 从go.mod读取domain
		domain, err := getDomainFromGoMod()
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}

		//检查proto文件
		protoName := fmt.Sprintf("pb/model_%v.proto", objName)
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
		newContent := fmt.Sprintf(`message %v{
	uint64 Id = 1; //[primarykey]
}`, upObjName)

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
		svcContent := fmt.Sprintf(`&pb.%v{},
	%v`, upObjName, string(modelRegMask))
		objBuff = bytes.ReplaceAll(objBuff, modelRegMask, []byte(svcContent))
		os.WriteFile(objFileName, objBuff, 0644)

		newBuff, err := utils.ExecCmd("", "goimports", objFileName)
		if err != nil {
			color.Red(err.Error())
			return
		}
		os.WriteFile(objFileName, []byte(newBuff), 0644)
		// utils.ExecCmd("", "go", "fmt", svcFileName)

		color.Green("create db model[%v] success", objName)
	},
}
