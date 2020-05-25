// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import ( //"flag"
	//"fmt"
	//"strconv"
	//"strings"
	"testing"
)

// https://stackoverflow.com/questions/23205419/how-do-you-print-in-a-go-test-using-the-testing-package
// vscode: "go.testFlags": ["-v"]
func TestExtractPADSPCBComponents(t *testing.T) {
	exfile := "../../test/pads/PCB/ex2.asc"

	t.Logf("Testing PADS PCB Ascii Text Parser")
	if _, e := ExtractPADSPCBComponents(exfile); e != nil {
		t.Errorf("ExtractPADSPCBComponents error: %v", e)
	}
	t.Logf("ExtractPADSPCBComponents test done.")

}
