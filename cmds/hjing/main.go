package main

import (
	"fmt"

	"github.com/jkkkls/hjing/cmds/hjing/internal/app"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "hjing",
	Short:   fmt.Sprintf("hjing[%v]: An simple toolkit for Go microservices.", version),
	Long:    fmt.Sprintf("hjing[%v]: An simple toolkit for Go microservices...", version),
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
