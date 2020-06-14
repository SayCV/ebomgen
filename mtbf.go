/*
 * =====
 * SPDX-License-Identifier { (GPL-2.0+ OR MIT) {
 *
 * !!! THIS IS NOT GUARANTEED TO WORK !!!
 *
 * Copyright (c) 2018-2020, SayCV
 * =====
 */

 package ebomgen

 import (
	"encoding/csv"
	"io"
	//"net/url"
	"os"
	"path"
	"path/filepath"
	//"regexp"
	"strconv"
	"strings"

	"github.com/saycv/ebomgen/pkg/configuration"
	"github.com/saycv/ebomgen/pkg/types"
	//"github.com/saycv/ebomgen/pkg/utils"
	"github.com/saycv/ebomgen/pkg/reliability"

	log "github.com/sirupsen/logrus"
 )
 
func CalcMtbfBasedPCPMain(config configuration.Configuration) error {
	outputFilenameAppend := ""

	filenameWithSuffix := path.Base(config.InputFile)
	fileSuffix := path.Ext(filenameWithSuffix)
	prjname := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	log.Infof("Project Name %s", prjname)

	outputFilename := filepath.ToSlash(filepath.Join(config.OutputFile, outputFilenameAppend))
	log.Infof("CSV Output File %s", outputFilename)

	csvFile, err := os.Open(config.InputFile)
	if err != nil {
		return err
	}
	reader := csv.NewReader(csvFile)
	var frParts []reliability.EBOMFrPart
	var bomParts []types.EBOMItem
	var cpart types.EBOMItem
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if line[0] == "Item" && line[1] == "References" {
			continue
		}
		if line[5] == "DNP" || line[5] == "TestPoint" {
			continue
		}
		//cpart.Item = line[0]
		cpart.References = strings.Split(line[1], ",")
		cpart.Quantity, _ = strconv.Atoi(line[2])
		cpart.Value = line[3]
		cpart.Footprint = line[4]
		cpart.Desc = line[5]
		//cpart.FrUnit = ""
		cpart.Attributes = map[string]string{
			"FrUnit":   "",
			"Description": line[5],
		}

		bomParts = append(bomParts, cpart)

		frpart := reliability.NewFrPart(cpart,
			reliability.WithFrType("RES-Film-Carbon"),
			reliability.WithFrProcess(""),
			reliability.WithClsEnv("GB"),
			reliability.WithClsQuality("C1"),
			reliability.WithOperatingTemp("GB"),
			reliability.WithCurrentStress("0.5"),
			reliability.WithVoltageStress("0.5"),
			reliability.WithPowerStress("0.5"),
		)
		frParts = append(frParts, *frpart)
	}

	for k, ipart := range frParts {
		if strings.HasPrefix(ipart.Desc, "Capacitor") {
			ipart.FrUnit, _ = ipart.FrCalcCap()
		} else if strings.HasPrefix(ipart.Desc, "Resistor") {
			ipart.FrUnit, _ = ipart.FrCalcRes()
		} else if strings.HasPrefix(ipart.Desc, "Inductor") {
			ipart.FrUnit, _ = ipart.FrCalcInd()
		} else if strings.HasPrefix(ipart.Desc, "Fuse") {
			ipart.FrUnit, _ = ipart.FrCalcRes()
		} else if strings.HasPrefix(ipart.Desc, "LED") {
			ipart.FrUnit, _ = ipart.FrCalcDiodeBjt()
		} else if strings.HasPrefix(ipart.Desc, "Diode") {
			ipart.FrUnit, _ = ipart.FrCalcDiodeBjt()
		} else if strings.HasPrefix(ipart.Desc, "Transistor") {
			ipart.FrUnit, _ = ipart.FrCalcDiodeBjt()
		} else if strings.HasPrefix(ipart.Desc, "FET") {
			ipart.FrUnit, _ = ipart.FrCalcDiodeBjt()
		} else if strings.HasPrefix(ipart.Desc, "Crystal") {
			ipart.FrUnit, _ = ipart.FrCalcXtal()
		} else if strings.HasPrefix(ipart.Desc, "Oscillator") {
			ipart.FrUnit, _ = ipart.FrCalcXtal()
		} else if strings.HasPrefix(ipart.Desc, "ConnRJ") {
			ipart.FrUnit, _ = ipart.FrCalcConn()
		} else if strings.HasPrefix(ipart.Desc, "ConnUSB") {
			ipart.FrUnit, _ = ipart.FrCalcConn()
		} else if strings.HasPrefix(ipart.Desc, "Connector") {
			ipart.FrUnit, _ = ipart.FrCalcConn()
		} else if strings.HasPrefix(ipart.Desc, "Switch") {
			ipart.FrUnit, _ = ipart.FrCalcSwitch()
		} else if strings.HasPrefix(ipart.Desc, "XFRM") {
			ipart.FrUnit, _ = ipart.FrCalcInd()
		} else if strings.HasPrefix(ipart.Desc, "IC") {
			ipart.FrUnit, _ = ipart.FrCalcIc()
		}
		bomParts[k].Attributes["FrUnit"] = ipart.FrUnit
	}

	BOM, err := types.NewBOM(bomParts, config)
	if err != nil {
		return err
	}

	output, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer output.Close()
	BOM.WriteCSV(output)
	log.Infof("Created BOM in '%s'", outputFilename)

	return nil
}

func CalcMtbfBasedPCP(config configuration.Configuration) error {
 
	return CalcMtbfBasedPCPMain(config)
}
