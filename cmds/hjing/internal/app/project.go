package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jkkkls/hjing/layout"
	"github.com/jkkkls/hjing/utils"
	"github.com/spf13/cobra"
)

// CmdAddApp represents the new command.
var CmdNew = &cobra.Command{
	Use:   "new [<projectName> | <domainName>]",
	Short: "Create a template",
	Long:  "Create a project using the repository template. Example: hjing new [<projectName> | <domainName>]",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires 2 args, example: hjing add-app <[<projectName> | <domainName>]")
		}

		if !isValidDiomainName(args[0]) {
			return fmt.Errorf("projectName is invalid, name: %v", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		arr := strings.Split(projectName, "/")
		dir := arr[len(arr)-1]
		if utils.PathExists(dir) {
			cmd.Usage()
			color.Red("project dir[%v] has been exitst\n", dir)
			return
		}

		upDir := strings.ToUpper(dir[:1]) + dir[1:]

		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}

		err = layout.CopyDir("project", dir, "{{projectName}}", projectName, "{{upProjectName}}", upDir)
		if err != nil {
			cmd.Usage()
			color.Red(err.Error())
			return
		}

		str, err := utils.ExecCmd(dir, "go", "mod", "init", projectName)
		if err != nil {
			cmd.Usage()
			color.Red(str, err.Error())
			return
		}

		color.Green("create project[%v] success\n", projectName)
	},
}
