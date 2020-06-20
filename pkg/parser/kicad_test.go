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

func TestExtractKicadComponents(t *testing.T) {
	exfile := "../../test/kicad/SCH/ex9.sch"

	t.Logf("Testing Kicad netlist Parser")
	if _, e := ExtractKicadSchComponents(exfile); e != nil {
		t.Errorf("ExtractKicadComponents error: %v", e)
	}
	t.Logf("ExtractKicadComponents test done.")

}
