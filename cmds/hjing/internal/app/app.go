package app

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/jkkkls/hjing/layout"
	"github.com/spf13/cobra"
)

var (
	appPbMask    = []byte("#pb tag")
	appBuildMask = []byte(`#build tag
build:`)
)

// CmdAddApp represents the new command.
var CmdAddApp = &cobra.Command{
	Use:   "add-app <appName>",
	Short: "Create a template",
	Long:  "Create a app using the repository template. Example: hjing add-app <appName>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires 2 args, example: hjing add-app <<appName>")
		}

		if !isValidAppName(args[0]) {
			return fmt.Errorf("name is invalid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]

		mkBuff, err := os.ReadFile("Makefile")
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}
		// 从go.mod读取domain
		domain, err := getDomainFromGoMod()
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}

		err = os.MkdirAll("apps/"+appName, os.ModePerm)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
		}

		err = layout.CopyFile("app/app.main.go.tpl", "apps/"+appName+"/main.go", "{{appName}}", appName, "{{projectName}}", domain)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
		}

		err = layout.CopyFile("app/app.yaml.tpl", "apps/"+appName+"/"+appName+".yaml", "{{appName}}", appName)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
		}

		mask := `{{appName}}:
	@$(GOBUILD) -o build/{{appName}} {{domain}}/apps/{{appName}}
	@cp apps/{{appName}}/{{appName}}.yaml ./build/
	@echo "编译{{appName}}完成"
` + string(appPbMask)
		mask = strings.ReplaceAll(mask, "{{appName}}", appName)
		mask = strings.ReplaceAll(mask, "{{domain}}", domain)
		//
		mkBuff = bytes.ReplaceAll(mkBuff, []byte(".PHONY: "), []byte(".PHONY: "+appName+" "))
		mkBuff = bytes.ReplaceAll(mkBuff, appBuildMask, []byte(string(appBuildMask)+" "+appName))
		mkBuff = bytes.ReplaceAll(mkBuff, appPbMask, []byte(mask))
		os.WriteFile("Makefile", mkBuff, 0o644)

		color.Green("create app[%v] success", appName)
	},
}

// isValidAppName 检查格式，只允许字母开头，大小字母和数字组成
func isValidAppName(appName string) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9]+$`, appName)
	return ok
}

func isValidDiomainName(appName string) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9/.-_]+$`, appName)
	return ok
}
