package main

import (
	"log"

	"github.com/jkkkls/hjing/cmds/hjing/internal/project"
	"github.com/jkkkls/hjing/internal/utils"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "Hjing",
	Short:   "Hjing: An simple toolkit for Go microservices.",
	Long:    `Hjing: An simple toolkit for Go microservices.`,
	Version: version,
}

func init() {
	cmd.AddCommand(project.CmdNew)
	// cmd.AddCommand(proto.CmdProto)
	// cmd.AddCommand(upgrade.CmdUpgrade)
	// cmd.AddCommand(change.CmdChange)
	// cmd.AddCommand(run.CmdRun)
}

func main() {
	utils.FixRange(0, 0, 0)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
