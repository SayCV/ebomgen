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

// import (
//	 "fmt"
//	 "io"
//	 "os"
//	 "strings"
// 
//	 "path/filepath"
// 
//	 log "github.com/sirupsen/logrus"
//	 "github.com/spf13/cobra"
// )
 
var (
	// BuildCommit lastest build commit (set by Makefile)
	BuildCommit = ""
	// BuildTag if the `BuildCommit` matches a tag
	BuildTag = ""
	// BuildTime set by build script (set by Makefile)
	BuildTime = ""
)

