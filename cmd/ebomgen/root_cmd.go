package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"path/filepath"

	"github.com/saycv/ebomgen"
	"github.com/saycv/ebomgen/pkg/configuration"
	logsupport "github.com/saycv/ebomgen/pkg/log"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewRootCmd returns the root command
func NewRootCmd() *cobra.Command {

	var verbose bool
	var writeCSV bool
	var writeXLSX bool
	var input string
	var output string
	var edaTool string
	var logLevel string

	rootCmd := &cobra.Command{
		Use:   "ebomgen -i infile -o outfile -t [eagle|kicad|orcad|padslogic] [-w] [-v]",
		Short: `ebomgen is a tool to auto generate bom from EDA design file, it support Orcad, Altium or Mentor Graphics`,
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
			out, close := getOut(cmd, input, output)
			if out != nil {
				defer close()
				path, _ := filepath.Abs(input)
				path = filepath.ToSlash(path)

				config := configuration.NewConfiguration(
					configuration.WithInputFile(input),
					configuration.WithOutputFile(output),
					configuration.WithEDATool(edaTool))

				err := ebomgen.ExtractComponents(config)
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
	flags.BoolVarP(&writeCSV, "writeCSV", "w", true, "Write BOM to CSV file")
	flags.BoolVarP(&writeXLSX, "writeXLSX", "x", true, "Write BOM to XLSX file")
	flags.StringVarP(&input, "input", "i", "../../test/padslogic/SCH/ex1.txt", "The path to the input schematic or netlist file")
	flags.StringVarP(&output, "output", "o", "../../test/padslogic/BOM/", "The path for the output file")
	flags.StringVarP(&edaTool, "edaTool", "t", "padslogic", "Define what EDA tool created the input file")
	flags.StringVarP(&logLevel, "log", "l", "debug", "log level to set [debug|info|warning|error|fatal|panic]")
	return rootCmd
}

type closeFunc func() error

func defaultCloseFunc() closeFunc {
	return func() error { return nil }
}

func newCloseFileFunc(c io.Closer) closeFunc {
	return func() error {
		return c.Close()
	}
}

func getOut(cmd *cobra.Command, sourcePath, outputName string) (io.Writer, closeFunc) {
	// Create the mod cache so we can rename it later, even if we don't need it.
	if err := os.MkdirAll(outputName, 0755); err != nil {
		return cmd.OutOrStdout(), defaultCloseFunc()
	}
	return cmd.OutOrStdout(), defaultCloseFunc()
}

// converts the `name`, `!name` and `name=value` into a map
func parseAttributes(attributes []string) map[string]string {
	result := make(map[string]string, len(attributes))
	for _, attr := range attributes {
		data := strings.Split(attr, "=")
		if len(data) > 1 {
			result[data[0]] = data[1]
		} else {
			result[data[0]] = ""
		}
	}
	return result
}
