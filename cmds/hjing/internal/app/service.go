package app

import (
	"github.com/spf13/cobra"
)

// CmdAddSrv represents the new command.
var CmdAddSrv = &cobra.Command{
	Use:   "add-app",
	Short: "Create a template",
	Long:  "Create a services using the repository template. Example: hjing add-app <AppName>",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
