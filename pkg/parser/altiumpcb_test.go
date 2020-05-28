// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import ( //"flag"
	//"fmt"
	//"strconv"
	//"strings"
	//"regexp"
	//"strconv"
	"testing"
)

// https://stackoverflow.com/questions/23205419/how-do-you-print-in-a-go-test-using-the-testing-package
// vscode: "go.testFlags": ["-v"]
func TestExtractAltiumPCBComponents(t *testing.T) {
	exfile := "../../test/altium/PCB/ex4.PcbDoc"

	t.Logf("Testing Altium PcbDoc Ascii Text Parser")
	if _, e := ExtractAltiumPCBComponents(exfile); e != nil {
		t.Errorf("ExtractAltiumPCBComponents error: %v", e)
	}
	t.Logf("ExtractAltiumPCBComponents test done.")

}
