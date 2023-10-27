package app

import (
	"github.com/jkkkls/hjing/layout"
	"github.com/jkkkls/hjing/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

// CmdAddApp represents the new command.
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a template",
	Long:  "Create a project using the repository template. Example: hjing new <AppName>",
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		arr := strings.Split(projectName, "/")
		dir := arr[len(arr)-1]
		if !isValidAppName(dir) {
			log.Fatal("app name is invalid")
		}

		if utils.PathExists(dir) {
			log.Fatalf("project dir[%v] has been exitst", dir)
		}

		upDir := strings.ToUpper(dir[:1]) + dir[1:]

		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		err = layout.CopyDir("project", dir, "{{projectName}}", projectName, "upProjectName", upDir)
		if err != nil {
			log.Fatal(err)
		}

		str, err := utils.ExecCmd(dir, "go", "mod", "init", projectName)
		if err != nil {
			log.Fatal(str, err)
		}
	}}
