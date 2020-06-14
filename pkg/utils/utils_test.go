// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	//"flag"
	//"fmt"
	"math"
	"strconv"
	//"strings"
	"testing"

	"github.com/saycv/ebomgen/pkg/types"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestGetFValFromEVal(t *testing.T) {
	val := "2.2uF"
	if e := GetFValFromEVal(val); !almostEqual(e, 0.0000022) {
		t.Errorf("utils.GetFValFromEVal error: got %f; want 0.0000022", e)
	}
	val = "10k"
	if e := GetFValFromEVal(val); !almostEqual(e, 10000) {
		t.Errorf("utils.GetFValFromEVal error: got %f; want 10000", e)
	}
	val = "10M"
	if e := GetFValFromEVal(val); !almostEqual(e, 10000000) {
		t.Errorf("utils.GetFValFromEVal error: got %f; want 10000000", e)
	}
}

// https://stackoverflow.com/questions/23205419/how-do-you-print-in-a-go-test-using-the-testing-package
// vscode: "go.testFlags": ["-v"]
func TestNamerulesProcess(t *testing.T) {
	var part types.EBOMItem
	part.References = []string{"c1"}
	part.Value = "10pF"
	part.Library = "10pF CAP"
	part.Footprint = "C0603"
	var propclass1 = map[string]string{
		"Description": "CAP",
		"part":        "unkownPart",
		"group":       "unkownGroup",
	}
	if e := NamerulesProcess(part, part.Value, part.Footprint, propclass1); !e {
		t.Errorf("utils.NamerulesProcess error: %v", e)
	}
	t.Logf("utils.NamerulesProcess results: %v", propclass1)

	part.References = []string{"r1"}
	part.Value = "10M"
	part.Library = "10M RES"
	part.Footprint = "R0603"
	var propclass2 = map[string]string{
		"Description": "RES",
		"part":        "unkownPart",
		"group":       "unkownGroup",
	}
	if e := NamerulesProcess(part, part.Value, part.Footprint, propclass2); !e {
		t.Errorf("utils.NamerulesProcess error: %v", e)
	}
	t.Logf("utils.NamerulesProcess results: %v", propclass2)

}

func TestCase0(t *testing.T) {

	queryTempStr := "GF1"
	queryTemp, err := strconv.Atoi(queryTempStr)
	t.Log(err)
	t.Log(queryTemp)

	__spFloat := 0.6
	result := int(__spFloat*10)-1
	t.Log(result)
}
