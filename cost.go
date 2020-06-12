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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/saycv/ebomgen/pkg/configuration"
	"github.com/saycv/ebomgen/pkg/types"
	"github.com/saycv/ebomgen/pkg/utils"
	"github.com/saycv/ebomgen/pkg/webecd"

	log "github.com/sirupsen/logrus"
)

func FetchPriceFromWebecd(config configuration.Configuration) error {
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
		cpart.Attributes = map[string]string{
			"UnitPrice":   "",
			"Description": line[5],
		}

		bomParts = append(bomParts, cpart)
	}

	hcDigikey := webecd.NewDigikeyClient()
	hcSzlcsc := webecd.NewSzlcscClient()
	for _, ipart := range bomParts {
		if strings.HasPrefix(ipart.Attributes["Description"], "DNP") ||
			strings.HasPrefix(ipart.Attributes["Description"], "TestPoint") {
			continue
		}
		value := ipart.Value
		fp := ipart.Footprint
		regVal, err := regexp.Compile("[^a-zA-Z0-9%\\.]+")
		value = regVal.ReplaceAllString(value, " ")
		if strings.HasPrefix(ipart.Attributes["Description"], "Conn") {
			regVal2, err := regexp.Compile(`([0-9]+)[dDpP]([0-9]+)`)
			value = regVal2.ReplaceAllString(value, "${1}.${2}")
			regVal3, err := regexp.Compile(`.HDR`)
			processedString = regVal3.ReplaceAllString(value, " header ")
		}

		reg, err := regexp.Compile("[^0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		digitfp := fp
		if strings.HasPrefix(ipart.Attributes["Description"], "CapacitorArray") || strings.HasPrefix(ipart.Attributes["Description"], "ResistorArray") {
			digitfp = fp[3:]
		}
		log.Println(digitfp)
		digitfp = reg.ReplaceAllString(digitfp, "")
		log.Println(digitfp)

		querympn := value
		if strings.HasPrefix(ipart.Attributes["Description"], "Capacitor") {
			fvalue := strconv.FormatFloat(utils.GetFValFromEVal(value), 'E', -1, 64)
			log.Println(fvalue)
			valPref := ""
			if strings.HasPrefix(ipart.Attributes["Description"], "CapacitorArray") {
				valPref = "Capacitor Array"
			}
			querympn = strings.Join([]string{valPref, value, digitfp}, " ")
			if fvalue == "-1E+00" {
				querympn = strings.Join([]string{valPref, "0.1uF", digitfp}, " ")
			}
		}  else if strings.HasPrefix(ipart.Attributes["Description"], "Resistor") {
			fvalue := strconv.FormatFloat(utils.GetFValFromEVal(value), 'E', -1, 64)
			log.Println(fvalue)
			valPref := ""
			if strings.HasPrefix(ipart.Attributes["Description"], "ResistorArray") {
				valPref = "Resistor Array"
			}
			querympn = strings.Join([]string{valPref, value, digitfp}, " ")
			if fvalue == "-1E+00" {
				querympn = strings.Join([]string{valPref, "22R", digitfp}, " ")
			}
		} else if strings.HasPrefix(ipart.Attributes["Description"], "IC") {
			if strings.Contains(value, " ") {
				_val := strings.Split(value, " ")
				querympn = _val[0]
			}
		} else if strings.HasPrefix(ipart.Attributes["Description"], "LED") {
			querympn = strings.Join([]string{value, "LED"}, " ")
		}

		log.Infof(querympn)
		webpart, err := FetchPriceFromDigikey(hcDigikey, url.QueryEscape(querympn))
		//log.Println(webpart)
		if webpart.UnitPrice.Value == "" {
			log.Infof("Try get from 2nd websource")
			webpart, err = FetchPriceFromSzlcsc(hcSzlcsc, url.QueryEscape(querympn))
		}
		ipart.Attributes["UnitPrice"] = webpart.UnitPrice.Value
		log.Println(ipart)
	}
	hcDigikey.Close()

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

func FetchPriceFromDigikey(hc *webecd.DigikeyClient, query string) (types.EBOMWebPart, error) {
	//hc := webecd.NewDigikeyClient()
	result, err := hc.QueryWDCall(query)
	if err != nil {
		log.Infof("Error with query call: " + err.Error())
	}
	return result, err
}

func FetchPriceFromSzlcsc(hc *webecd.SzlcscClient, query string) (types.EBOMWebPart, error) {
	//hc := webecd.NewSzlcscClient()
	result, err := hc.QueryCall(query)
	if err != nil {
		log.Infof("Error with query call: " + err.Error())
	}
	return result, err
}
