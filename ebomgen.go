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
	"github.com/saycv/ebomgen/pkg/configuration"
	"github.com/saycv/ebomgen/pkg/parser"
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
func ExtractComponents(config configuration.Configuration) (bool, error) {
	return parser.ExtractPADSLogicComponents(config.InputFile), nil
}
