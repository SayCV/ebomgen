package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"path/filepath"

	//"github.com/saycv/ebomgen"
	logsupport "github.com/saycv/ebomgen/pkg/log"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewRootCmd returns the root command
func NewRootCmd() *cobra.Command {

	var noHeaderFooter bool
	var outputName string
	var logLevel string
	var css string
	var attributes []string

	rootCmd := &cobra.Command{
		Use:   "ebomgen [flags] FILE",
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
			if len(args) == 0 {
				return helpCommand.RunE(cmd, args)
			}
			//attrs := parseAttributes(attributes)
			for _, sourcePath := range args {
				out, close := getOut(cmd, sourcePath, outputName)
				if out != nil {
					defer close()
					path, _ := filepath.Abs(sourcePath)
					log.Debugf("Starting to process file %v", path)
				}
			}
			return nil
		},
	}
	rootCmd.SilenceUsage = true
	flags := rootCmd.Flags()
	flags.BoolVarP(&noHeaderFooter, "no-header-footer", "s", false, "do not render header/footer (default: false)")
	flags.StringVarP(&outputName, "out-file", "o", "", "output file (default: based on path of input file); use - to output to STDOUT")
	flags.StringVar(&logLevel, "log", "warning", "log level to set [debug|info|warning|error|fatal|panic]")
	flags.StringVar(&css, "css", "", "the path to the CSS file to link to the document")
	flags.StringArrayVarP(&attributes, "attribute", "a", []string{}, "a document attribute to set in the form of name, name!, or name=value pair")
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
	if outputName == "-" {
		// outfile is STDOUT
		return cmd.OutOrStdout(), defaultCloseFunc()
	} else if outputName != "" {
		// outfile is specified in the command line
		outfile, e := os.Create(outputName)
		if e != nil {
			log.Warnf("Cannot create output file - %v, skipping", outputName)
		}
		return outfile, newCloseFileFunc(outfile)
	} else if sourcePath != "" {
		// outfile is based on sourcePath
		path, _ := filepath.Abs(sourcePath)
		outname := strings.TrimSuffix(path, filepath.Ext(path)) + ".html"
		outfile, err := os.Create(outname)
		if err != nil {
			log.Warnf("Cannot create output file - %v, skipping", outname)
			return nil, nil
		}
		return outfile, newCloseFileFunc(outfile)
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