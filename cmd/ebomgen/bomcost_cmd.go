package main

import (
	"fmt"

	"github.com/saycv/ebomgen"
	"github.com/spf13/cobra"
)

// NewBomcostCmd returns the root command
func NewBomcostCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bomcost",
		Short: "Fetch the bomcost info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "build time: %s\n", ebomgen.BuildTime)
		},
	}
}
