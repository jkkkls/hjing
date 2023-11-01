package main

import (
	"github.com/jkkkls/hjing/cmds/hjing/internal/app"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "Hjing",
	Short:   "Hjing: An simple toolkit for Go microservices.",
	Long:    `Hjing: An simple toolkit for Go microservices.`,
	Version: version,
}

func init() {
	cmd.AddCommand(app.CmdNew)
	cmd.AddCommand(app.CmdAddApp)
	cmd.AddCommand(app.CmdAddSrv)
	cmd.AddCommand(app.CmdAddItf)
	cmd.AddCommand(app.CmdAddDb)
}

func main() {
	cmd.Execute()
}
