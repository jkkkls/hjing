package project

import "github.com/spf13/cobra"

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a service template",
	Long:  "Create a service project using the repository template. Example: hjing new helloworld",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {

}
