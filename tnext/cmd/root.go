package cmd

import (
	"github.com/spf13/cobra"
	"github.com/trento-project/trento/tnext/cmd/agent"
)

func NewTNextCmd() *cobra.Command {
	tNext := &cobra.Command{
		Use:   "tNext",
		Short: "Trento NEXT",
	}

	tNext.AddCommand(agent.NewAgentCmd())

	return tNext
}
