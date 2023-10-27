package main

import (
	"log"

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
	// cmd.AddCommand(proto.CmdProto)
	// cmd.AddCommand(upgrade.CmdUpgrade)
	// cmd.AddCommand(change.CmdChange)
	// cmd.AddCommand(run.CmdRun)
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
