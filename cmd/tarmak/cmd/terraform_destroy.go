package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/tarmak"
	"github.com/jetstack/tarmak/pkg/tarmak/utils"
)

// tfdestroyCmd represents the tfdestroy command
var terraformDestroyCmd = &cobra.Command{
	Use:     "terraform-destroy",
	Aliases: []string{"t-d"},
	Short:   "This applies the set of stacks in the current context",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := utils.GetContext()
		defer cancel()
		t := tarmak.New(cmd)
		t.Must(t.CmdTerraformDestroy(args, ctx))
	},
}

func init() {
	terraformPFlags(terraformDestroyCmd.PersistentFlags())
	terraformDestroyCmd.PersistentFlags().Bool(tarmak.FlagForceDestroyStateStack, false, "destroy the state stack as well, possibly dangerous")
	RootCmd.AddCommand(terraformDestroyCmd)
}
