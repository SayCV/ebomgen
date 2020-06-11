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

func TestExtractEagleComponents(t *testing.T) {
	exfile := "../../test/eagle/SCH/ex7.sch"

	t.Logf("Testing Eagle netlist Parser")
	if _, e := ExtractEagleComponents(exfile); e != nil {
		t.Errorf("ExtractEagleComponents error: %v", e)
	}
	t.Logf("ExtractEagleComponents test done.")

}
