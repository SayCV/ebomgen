// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reliability

import (
	//"flag"
	//"fmt"
	//"strconv"
	//"strings"
	//"regexp"
	//"strconv"
	"strconv"
	"strings"
	"testing"

	"github.com/saycv/ebomgen/pkg/types"
)

func TestCase1(t *testing.T) {

	t.Logf("%d", RpmPartsCount)
	t.Logf("%d", RpmPartStress)
	t.Logf("%d", RpmPseudoStress)
	t.Logf("%d", RpmUndef)

	t.Log(len(ClsEnv))
	t.Log(ClsEnv)

	t.Log(len(FactorEnvImported))
	t.Log(FactorEnvImported)

	t.Log(len(FactorQualityImported))
	t.Log(FactorQualityImported)

	t.Log(len(FactorTemperatureImported))
	t.Log(FactorTemperatureImported)

	t.Log(len(FactorStressImported))
	t.Log(FactorStressImported)
}

func TestCase2(t *testing.T) {
	part := types.EBOMItem{
		/*Quantity=*/ 0,
		/*References=*/ []string{""},
		/*Value=*/ "",
		/*FValue=*/ 0.0,
		/*Library=*/ "",
		/*Footprint=*/ "",
		/*Desc=*/ "",
		/*Attributes*/ map[string]string{},
		/*Group*/ []string{""},
		/*PartSpecs*/ types.EBOMWebPart{},
	}
	frpart := NewFrPart(part,
		WithFrType("NPN-Si-LP"),
		WithClsEnv("GF1"),
		WithCurrentStress("0.8"),
		WithVoltageStress("0.8"),
		WithPowerStress("0.8"),
	)
	results, err := frpart.GetFailureRateBaseImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ := strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	t.Log(frpart)
}
