package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/tarmak"
	"github.com/jetstack/tarmak/pkg/tarmak/provider"
)

var providerInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init providers",
	Run: func(cmd *cobra.Command, args []string) {
		t := tarmak.New(cmd)
		provider.Init(t)
	},
}

func init() {
	providerCmd.AddCommand(providerInitCmd)
}
