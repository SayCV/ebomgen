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

	for _, ipart := range bomParts {
		value := ipart.Value
		fp := cpart.Footprint
		querympn := value
		if strings.HasPrefix(ipart.Attributes["Description"], "Capacitor") {
			querympn = strings.Join([]string{value, fp}, " ")
		}

		log.Infof(querympn)
		webpart, _ := FetchPriceFromDigikey(querympn)
		//log.Println(webpart)
		ipart.Attributes["UnitPrice"] = webpart.UnitPrice.Value
		log.Println(ipart)
	}

	return nil
}

func FetchPriceFromDigikey(query string) (types.EBOMWebPart, error) {
	hc := webecd.NewDigikeyClient()
	result, err := hc.QueryWDCall(query)
	if err != nil {
		log.Infof("Error with query call: " + err.Error())
	}
	return result, err
}
