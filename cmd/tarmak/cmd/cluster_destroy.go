package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/tarmak"
)

var clusterDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "This destroys the current cluster context",
	Run: func(cmd *cobra.Command, args []string) {
		t := tarmak.New(cmd)
		t.Must(t.CmdTerraformDestroy(args))
	},
}

func init() {
	terraformPFlags(clusterDestroyCmd.PersistentFlags())
	clusterDestroyCmd.PersistentFlags().Bool(tarmak.FlagForceDestroyStateStack, false, "destroy the state stack as well, possibly dangerous")
	clusterCmd.AddCommand(clusterDestroyCmd)
}
