package app

import (
	"log"
	"os"
	"regexp"

	"github.com/jkkkls/hjing/layout"
	"github.com/spf13/cobra"
)

// CmdAddApp represents the new command.
var CmdAddApp = &cobra.Command{
	Use:   "add-app",
	Short: "Create a template",
	Long:  "Create a project using the repository template. Example: hjing add-app <AppName>",
	Run:   run,
}

// isValidAppName 检查格式，只允许字母开头，大小字母和数字组成
func isValidAppName(appName string) bool {
	ok, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]+$", appName)
	return ok
}

func run(cmd *cobra.Command, args []string) {
	appName := args[0]
	if !isValidAppName(appName) {
		log.Fatal("app name is invalid")
	}
	// upAppName = strings.ToUpper(appName[:1]) + appName[1:]

	err := os.MkdirAll("apps/"+appName, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = layout.CopyFile("app/app.main.go.tpl", "apps/"+appName+"/main.go", "{{appName}}", appName)
	if err != nil {
		log.Fatal(err)
	}
}
