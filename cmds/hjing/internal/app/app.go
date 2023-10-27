package app

import (
	"log"

	"github.com/spf13/cobra"
)

// CmdAddApp represents the new command.
var CmdAddApp = &cobra.Command{
	Use:   "add-app",
	Short: "Create a template",
	Long:  "Create a project using the repository template. Example: hjing add-app game",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	log.Println("add-app called", args)
}
