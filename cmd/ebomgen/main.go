package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := NewRootCmd()
	bomcostCmd := NewBomcostCmd()
	bommtbfCmd := NewBommtbfCmd()
	bommergeCmd := NewBommergeCmd()
	versionCmd := NewVersionCmd()
	rootCmd.AddCommand(bomcostCmd)
	rootCmd.AddCommand(bommtbfCmd)
	rootCmd.AddCommand(bommergeCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.SetHelpCommand(helpCommand)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var helpCommand = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about the command",
	RunE: func(c *cobra.Command, args []string) error {
		cmd, args, e := c.Root().Find(args)
		if cmd == nil || e != nil || len(args) > 0 {
			return errors.Errorf("unknown help topic: %v", strings.Join(args, " "))
		}
		helpFunc := cmd.HelpFunc()
		helpFunc(cmd, args)
		return nil
	},
}
