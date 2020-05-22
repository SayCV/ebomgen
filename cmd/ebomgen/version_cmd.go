package main

import (
	"fmt"

	"github.com/saycv/ebomgen"
	"github.com/spf13/cobra"
)

// NewVersionCmd returns the root command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version and build info",
		Run: func(cmd *cobra.Command, args []string) {
			if ebomgen.BuildTag != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "version:    %s\n", ebomgen.BuildTag)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "commit:     %s\n", ebomgen.BuildCommit)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "build time: %s\n", ebomgen.BuildTime)
		},
	}
}
