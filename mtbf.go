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
	"github.com/saycv/ebomgen/pkg/utils"
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

		frType := ""
		frProcess := ""
		clsQuality := config.FrClsQuality
		ambientTemp := config.FrClsEnv
		operatingTemp := config.FrOpsEnv
		currentStress := config.FrDegrade
		voltageStress := config.FrDegrade
		powerStress := config.FrDegrade

		if strings.HasPrefix(cpart.Desc, "Capacitor") {
			frType = "CAP-Ceramic-1"
			if strings.HasPrefix(cpart.Desc, "CapacitorTan") {
				frType = "CAP-TAN"
			}
		} else if strings.HasPrefix(cpart.Desc, "Resistor") {
			frType = "RES-Film-Carbon"
		} else if strings.HasPrefix(cpart.Desc, "Inductor") {
			frType = "IND"
		} else if strings.HasPrefix(cpart.Desc, "Fuse") {
			frType = "RES-Wirewound-Power"
		} else if strings.HasPrefix(cpart.Desc, "LED") {
			frType = "LED"
		} else if strings.HasPrefix(cpart.Desc, "Diode") {
			frType = "Diode-Ge-LP"
			if strings.HasPrefix(cpart.Footprint, "POWER") {
				frType = "Diode-Ge-HP"
			}
		} else if strings.HasPrefix(cpart.Desc, "Transistor") {
			frType = "NPN-Si-LP"
			if strings.HasPrefix(cpart.Footprint, "POWER") {
				frType = "NPN-Si-HP"
			}
		} else if strings.HasPrefix(cpart.Desc, "FET") {
			frType = "FET-Si-Switch"
			if strings.HasPrefix(cpart.Footprint, "POWER") {
				frType = "FET-Si-Amp"
			}
		} else if strings.HasPrefix(cpart.Desc, "Crystal") {
			frType = "XTAL"
		} else if strings.HasPrefix(cpart.Desc, "Oscillator") {
			frType = "OSC"
		} else if strings.HasPrefix(cpart.Desc, "ConnRJ") {
			frType = "CONN-PCB"
		} else if strings.HasPrefix(cpart.Desc, "ConnUSB") {
			frType = "CONN-PCB"
		} else if strings.HasPrefix(cpart.Desc, "Connector") {
			frType = "CONN-PCB"
		} else if strings.HasPrefix(cpart.Desc, "Switch") {
			frType = "Switch"
		} else if strings.HasPrefix(cpart.Desc, "XFRM") {
			frType = "XFMR-LF"
		} else if strings.HasPrefix(cpart.Desc, "IC") {
			frType = "DIC-MOS"
			pins, _ := utils.GetPinsFromFp(cpart.Desc, cpart.Footprint)
			if pins>=100 {
				frType = "MPU-MOS"
			}
			if (pins==96 || pins==178) && strings.HasPrefix(cpart.Footprint, "BGA") {
				frType = "DRAM"
			}
			if (pins==48 && strings.HasPrefix(cpart.Footprint, "SO")) || (pins==169 && strings.HasPrefix(cpart.Footprint, "BGA")) {
				frType = "FLASH-MOS"
			}
		}

		frpart := reliability.NewFrPart(cpart,
			reliability.WithFrType(frType),
			reliability.WithFrProcess(frProcess),
			reliability.WithClsEnv(ambientTemp),
			reliability.WithClsQuality(clsQuality),
			reliability.WithOperatingTemp(operatingTemp),
			reliability.WithCurrentStress(currentStress),
			reliability.WithVoltageStress(voltageStress),
			reliability.WithPowerStress(powerStress),
		)
		frParts = append(frParts, *frpart)
	}

	for k, cpart := range frParts {
		log.Info(cpart)
		if strings.HasPrefix(cpart.Desc, "Capacitor") {
			cpart.FrUnit, _ = cpart.FrCalcCap()
		} else if strings.HasPrefix(cpart.Desc, "Resistor") {
			cpart.FrUnit, _ = cpart.FrCalcRes()
		} else if strings.HasPrefix(cpart.Desc, "Inductor") {
			cpart.FrUnit, _ = cpart.FrCalcInd()
		} else if strings.HasPrefix(cpart.Desc, "Fuse") {
			cpart.FrUnit, _ = cpart.FrCalcRes()
		} else if strings.HasPrefix(cpart.Desc, "LED") {
			cpart.FrUnit, _ = cpart.FrCalcOptoElectronicDevices()
		} else if strings.HasPrefix(cpart.Desc, "Diode") {
			cpart.FrUnit, _ = cpart.FrCalcDiodeBjt()
		} else if strings.HasPrefix(cpart.Desc, "Transistor") {
			cpart.FrUnit, _ = cpart.FrCalcDiodeBjt()
		} else if strings.HasPrefix(cpart.Desc, "FET") {
			cpart.FrUnit, _ = cpart.FrCalcDiodeBjt()
		} else if strings.HasPrefix(cpart.Desc, "Crystal") {
			cpart.FrUnit, _ = cpart.FrCalcXtal()
		} else if strings.HasPrefix(cpart.Desc, "Oscillator") {
			cpart.FrUnit, _ = cpart.FrCalcXtal()
		} else if strings.HasPrefix(cpart.Desc, "ConnRJ") {
			cpart.FrUnit, _ = cpart.FrCalcConn()
		} else if strings.HasPrefix(cpart.Desc, "ConnUSB") {
			cpart.FrUnit, _ = cpart.FrCalcConn()
		} else if strings.HasPrefix(cpart.Desc, "Connector") {
			cpart.FrUnit, _ = cpart.FrCalcConn()
		} else if strings.HasPrefix(cpart.Desc, "Switch") {
			cpart.FrUnit, _ = cpart.FrCalcSwitch()
		} else if strings.HasPrefix(cpart.Desc, "XFRM") {
			cpart.FrUnit, _ = cpart.FrCalcInd()
		} else if strings.HasPrefix(cpart.Desc, "IC") {
			cpart.FrUnit, _ = cpart.FrCalcIc()
		}
		bomParts[k].Attributes["FrUnit"] = cpart.FrUnit
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
