// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import ( //"flag"
	//"fmt"
	"log"
	//"strconv"
	"strings"
	"testing"
)

func TestOrcadCase0(t *testing.T) {
	exfile := "../../test/orcad/SCH/allegro-ex5"

	t.Logf("Testing Orcad Parser")
	results, e := ExtractOrcadSchComponents(exfile); 
	if e != nil {
		t.Errorf("ExtractOrcadSchComponents error: %v", e)
	}
	t.Log(results)
	t.Logf("ExtractOrcadSchComponents test done.")
}

func TestOrcadCase1(t *testing.T) {
	name := "primitive 'RES_RESC1608X55N_DISCRETE_154 K';"
	ind1 := strings.Index(name, "'")
	log.Println(len(name),ind1)
	name2 := string(name[ind1+1:])
	log.Println(len(name2),name2)
	ind2 := strings.Index(name2, "'")
	log.Println(len(name2),ind2)
	log.Println("ind1, ind2 = ", ind1, ", ",ind2)

	log.Println("newStr = ", name[ind1 + 1:ind2 + ind1 + 1])

	log.Println(stripQuotes(name))

	log.Println(OrcadSplit(name)[0], ", ", OrcadSplit(name)[1])
}
