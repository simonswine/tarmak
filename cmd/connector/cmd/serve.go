// Copyright Jetstack Ltd. See LICENSE for details.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/connector"
)

var RootCmd = &cobra.Command{
	Use:   "connector",
	Short: "tarmak connector to facilitate tarmak and terraform communications",
}

func Execute() {
	flag.CommandLine.Parse([]string{})

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	serverCmd := connector.NewCommandStartConnector()
	serverCmd.Use = "serve"
	serverCmd.Flags().AddGoFlagSet(flag.CommandLine)
	RootCmd.AddCommand(serverCmd)
}
