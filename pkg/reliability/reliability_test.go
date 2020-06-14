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

func TestCase0(t *testing.T) {

	queryTempStr := "55"
	queryTemp, err := strconv.Atoi(queryTempStr)
	t.Log(err)
	t.Log(queryTemp)
}

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
		WithFrProcess(""),
		WithClsEnv("GF1"),
		WithClsQuality("C1"),
		WithOperatingTemp("GF1"),
		WithCurrentStress("0.5"),
		WithVoltageStress("0.5"),
		WithPowerStress("0.5"),
	)
	t.Log(frpart)

	results, err := frpart.GetFailureRateBaseImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ := strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results, err = frpart.GetFactorEnvImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ = strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results, err = frpart.GetFactorQualityImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ = strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results, err = frpart.GetFactorTemperatureImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ = strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results, err = frpart.GetFactorStressImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ = strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results, err = frpart.GetFactorChImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ = strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results, err = frpart.GetFactorProcessImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ = strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results, err = frpart.GetFactorApplicationImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ = strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results, err = frpart.GetFactorC1Imported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ = strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results, err = frpart.GetFactorC2Imported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ = strconv.ParseFloat(results, 64)
	t.Log(fvalue)

}

func TestCase3(t *testing.T) {
	part := types.EBOMItem{
		/*Quantity=*/ 0,
		/*References=*/ []string{""},
		/*Value=*/ "10k",
		/*FValue=*/ 0.0,
		/*Library=*/ "",
		/*Footprint=*/ "",
		/*Desc=*/ "",
		/*Attributes*/ map[string]string{},
		/*Group*/ []string{""},
		/*PartSpecs*/ types.EBOMWebPart{},
	}
	frpart := NewFrPart(part,
		WithFrType("RES-Film-Carbon"),
		WithFrProcess(""),
		WithClsEnv("GB"),
		WithClsQuality("C1"),
		WithOperatingTemp("GB"),
		WithCurrentStress("0.5"),
		WithVoltageStress("0.5"),
		WithPowerStress("0.5"),
	)
	t.Log(frpart)

	results, err := frpart.GetFailureRateBaseImported()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ := strconv.ParseFloat(results, 64)
	t.Log(fvalue)

}

func TestResFr(t *testing.T) {
	part := types.EBOMItem{
		/*Quantity=*/ 1,
		/*References=*/ []string{"U1"},
		/*Value=*/ "CPU8086",
		/*FValue=*/ 0.0,
		/*Library=*/ "",
		/*Footprint=*/ "BGA472",
		/*Desc=*/ "IC",
		/*Attributes*/ map[string]string{},
		/*Group*/ []string{""},
		/*PartSpecs*/ types.EBOMWebPart{},
	}
	frpart := NewFrPart(part,
		WithFrType("MPU-MOS"),
		WithFrProcess(""),
		WithClsEnv("GB"),
		WithClsQuality("C1"),
		WithOperatingTemp("GB"),
		WithCurrentStress("0.5"),
		WithVoltageStress("0.5"),
		WithPowerStress("0.5"),
	)
	t.Log(frpart)

	results, err := frpart.FrCalcRes()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results = strings.Replace(results, " ", "", -1)
	fvalue, _ := strconv.ParseFloat(results, 64)
	t.Log(fvalue)

	results2, err := frpart.FrCalcIc()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	results2 = strings.Replace(results2, " ", "", -1)
	fvalue2, _ := strconv.ParseFloat(results2, 64)
	t.Log(fvalue2)

}
