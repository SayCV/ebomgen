/*
 * =====
 * SPDX-License-Identifier: (GPL-2.0+ OR MIT):
 *
 * !!! THIS IS NOT GUARANTEED TO WORK !!!
 *
 * Copyright (c) 2018-2020, SayCV
 * =====
 */

package ebomgen

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/saycv/ebomgen/pkg/configuration"
	"github.com/saycv/ebomgen/pkg/parser"

	log "github.com/sirupsen/logrus"
)

var (
	// BuildCommit lastest build commit (set by Makefile)
	BuildCommit = ""
	// BuildTag if the `BuildCommit` matches a tag
	BuildTag = ""
	// BuildTime set by build script (set by Makefile)
	BuildTime = ""
)

// ExtractComponents converts the content of the given filename into an BOM document.
// The conversion result is written in the given writer `output`, whereas the document metadata (title, etc.) (or an error if a problem occurred) is returned
// as the result of the function call.
func ExtractComponents(config configuration.Configuration) error {
	log.Infof("ExtractComponents!!!")

	if strings.ToUpper(config.EDATool) == "PADSLOGIC" {
		_, err := parser.ExtractPADSLogicComponents(config.InputFile)
		return err
	} else {
		return errors.Errorf("unknown tools: %v", config.EDATool)
	}
}
