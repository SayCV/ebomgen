// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import ( //"flag"
	//"fmt"
	//"strconv"
	//"strings"
	"regexp"
	"strconv"
	"testing"
)

// https://stackoverflow.com/questions/23205419/how-do-you-print-in-a-go-test-using-the-testing-package
// vscode: "go.testFlags": ["-v"]
func TestExtractPADSLogicComponents(t *testing.T) {
	exfile := "../../test/pads/SCH/ex1.txt"

	t.Logf("Testing PADS Logic Ascii Text Parser")
	if _, e := ExtractPADSLogicComponents(exfile); e != nil {
		t.Errorf("ExtractPADSLogicComponents error: %v", e)
	}
	t.Logf("ExtractPADSLogicComponents test done.")

}

func TestCase1(t *testing.T) {
	fp := "gp234"
	var err error
	_val := 0

	t.Logf("Testing TestCase1")
	reVal := regexp.MustCompile("\\d*\\d+")
	findRet := reVal.FindAll([]byte(fp), -1)
	if findRet != nil {
		_val, err = strconv.Atoi(string(findRet[0]))
	}

	if err != nil {
		t.Errorf("TestCase1 error: %v", err)
	}
	t.Logf("TestCase1 test done: %v.", _val)

}
