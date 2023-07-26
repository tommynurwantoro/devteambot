package cmd

import (
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use: "khugabot",
}

func init() {
	command.AddCommand(GetCommand())
}

func Execute() {
	if err := command.Execute(); err != nil {
		panic(err)
	}
}
