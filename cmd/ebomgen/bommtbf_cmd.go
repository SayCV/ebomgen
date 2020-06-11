package main

import (
	"fmt"

	"path/filepath"

	"github.com/saycv/ebomgen"
	"github.com/saycv/ebomgen/pkg/configuration"
	logsupport "github.com/saycv/ebomgen/pkg/log"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewBommtbfCmd returns the bommtbf command
func NewBommtbfCmd() *cobra.Command {

	var verbose bool
	var input string
	var output string
	var edaTool string
	var logLevel string

	rootCmd := &cobra.Command{
		Use:   "bommtbf -i infile -o outfile",
		Short: `Fetch the bommtbf info`,
		Args:  cobra.ArbitraryArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			lvl, err := log.ParseLevel(logLevel)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "unable to parse log level '%v'", logLevel)
				return err
			}
			logsupport.Setup()
			log.SetLevel(lvl)
			log.SetOutput(cmd.OutOrStdout())
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			//if len(args) == 0 {
			//	return helpCommand.RunE(cmd, args)
			//}
			//attrs := parseAttributes(attributes)

			//for _, sourcePath := range args {
			//}
			out, close := getOut(cmd, input, filepath.Dir(output))
			if out != nil {
				defer close()
				path, _ := filepath.Abs(input)
				path = filepath.ToSlash(path)

				config := configuration.NewConfiguration(
					configuration.WithInputFile(input),
					configuration.WithOutputFile(output),
					configuration.WithEDATool(edaTool),
					configuration.WithCommand("bommtbf"))

				err := ebomgen.CalcMtbfBasedPCP(config)
				if err != nil {
					return err
				}
				log.Infof("finished!!!")
			}
			return nil
		},
	}
	rootCmd.SilenceUsage = true
	flags := rootCmd.Flags()
	flags.BoolVarP(&verbose, "verbose", "v", true, "verbose")
	flags.StringVarP(&input, "input", "i", "test/pads/BOM/ex1_BOM.csv", "The path to the input schematic or netlist file")
	flags.StringVarP(&output, "output", "o", "test/pads/BOM/ex1_BOM.bommtbf.csv", "The path for the output file")
	flags.StringVarP(&logLevel, "log", "l", "debug", "log level to set [debug|info|warning|error|fatal|panic]")
	return rootCmd
}
