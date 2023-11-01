package app

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/jkkkls/hjing/layout"
	"github.com/spf13/cobra"
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

		err := os.MkdirAll("apps/"+appName, os.ModePerm)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
		}

		err = layout.CopyFile("app/app.main.go.tpl", "apps/"+appName+"/main.go", "{{appName}}", appName)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
		}

		color.Green("create app[%v] success", appName)
	},
}

// isValidAppName 检查格式，只允许字母开头，大小字母和数字组成
func isValidAppName(appName string) bool {
	ok, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]+$", appName)
	return ok
}
