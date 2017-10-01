package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/wing"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "rung wing daemon",
	Run: func(cmd *cobra.Command, args []string) {
		w := wing.New(flags)
		w.Must(w.Start())
	},
}

func init() {
	initMeshFlags(agentCmd.PersistentFlags(), flags)
	RootCmd.AddCommand(agentCmd)
}
