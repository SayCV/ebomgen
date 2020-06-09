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
	"os"
	"path"
	"path/filepath"
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

	hc := webecd.NewDigikeyClient()
	for _, ipart := range bomParts {
		value := ipart.Value
		fp := ipart.Footprint
		value = strings.Replace(value, "/", " ", -1)
		value = strings.Replace(value, "-", " ", -1)
		querympn := value
		if strings.HasPrefix(ipart.Attributes["Description"], "Capacitor") {
			fvalue := strconv.FormatFloat(utils.GetFValFromEVal(value), 'E', -1, 64)
			log.Println(fvalue)
			querympn = strings.Join([]string{value, fp}, " ")
			if fvalue == "-1E+00" {
				querympn = strings.Join([]string{"0.1uF", fp}, " ")
			}
		}

		log.Infof(querympn)
		webpart, _ := FetchPriceFromDigikey(hc, querympn)
		//log.Println(webpart)
		ipart.Attributes["UnitPrice"] = webpart.UnitPrice.Value
		log.Println(ipart)
	}
	hc.Close()

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
