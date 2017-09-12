package cmd

import (
	flag "github.com/spf13/pflag"

	"github.com/jetstack/tarmak/pkg/tarmak"
)

func terraformPFlags(fs *flag.FlagSet) {
	fs.StringSlice(
		tarmak.FlagTerraformStacks,
		[]string{},
		"Comma separated list of Terraform stacks to execute. "+
			"One or more of: state, network, tools, vault, kubernetes.",
	)
}
