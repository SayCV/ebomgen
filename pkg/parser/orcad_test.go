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

func TestExtractOrCADComponents(t *testing.T) {
	exfile := "../../test/orcad/SCH/allegro-ex5"

	t.Logf("Testing OrCAD netlist Parser")
	if _, e := ExtractOrCADComponents(exfile); e != nil {
		t.Errorf("ExtractOrCADComponents error: %v", e)
	}
	t.Logf("ExtractOrCADComponents test done.")

}
