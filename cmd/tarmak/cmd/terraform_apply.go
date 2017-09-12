package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/tarmak"
	"github.com/jetstack/tarmak/pkg/tarmak/utils"
)

var terraformApplyCmd = &cobra.Command{
	Use:     "terraform-apply",
	Aliases: []string{"t-a"},
	Short:   "This applies the set of stacks in the current context",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := utils.GetContext()
		defer cancel()
		t := tarmak.New(cmd)
		t.Must(t.CmdTerraformApply(args, ctx))
	},
}

func init() {
	terraformPFlags(terraformApplyCmd.PersistentFlags())
	RootCmd.AddCommand(terraformApplyCmd)
}
