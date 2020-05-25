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
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/saycv/ebomgen/pkg/configuration"
	"github.com/saycv/ebomgen/pkg/parser"
	"github.com/saycv/ebomgen/pkg/types"

	log "github.com/sirupsen/logrus"
)

var (
	// BuildCommit lastest build commit (set by Makefile)
	BuildCommit = ""
	// BuildTag if the `BuildCommit` matches a tag
	BuildTag = ""
	// BuildTime set by build script (set by Makefile)
	BuildTime = ""
)

// ExtractComponents converts the content of the given filename into an BOM document.
// The conversion result is written in the given writer `output`, whereas the document metadata (title, etc.) (or an error if a problem occurred) is returned
// as the result of the function call.
func ExtractComponents(config configuration.Configuration) error {
	var bomParts []types.EBOMItem
	var err error

	filenameWithSuffix := path.Base(config.InputFile)
	fileSuffix := path.Ext(filenameWithSuffix)
	prjname := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	log.Infof("Project Name %s", prjname)

	if strings.ToUpper(config.EDATool) == "PADSLOGIC" {
		bomParts, err = parser.ExtractPADSLogicComponents(config.InputFile)
	} else {
		err = errors.Errorf("unknown tools %v", config.EDATool)
	}
	numberofparts := len(bomParts)
	log.Infof("numberofparts %d", numberofparts)

	combinedBOMparts, _ := combineBOMparts(bomParts)
	outputFilename := filepath.ToSlash(filepath.Join(config.OutputFile, prjname+"_BOM.csv"))
	log.Infof("CSV Output File %s", outputFilename)

	for k, ipart := range combinedBOMparts {
		combinedBOMparts[k].References = sortComponentRef(ipart)
	}

	BOM, err := types.NewBOM(combinedBOMparts)
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

func combineBOMparts(bomparts []types.EBOMItem) ([]types.EBOMItem, error) {
	// This function combines identical parts, updates the count
	// and returns a list of combined parts

	var combinedparts []types.EBOMItem
	var foundFlag bool

	// Check if part is in combined parts
	// if so append update the qty for the part in combined parts and update refs
	// if not append it to combined parts
	// Not the most elegant but it works
	for _, ipart := range bomparts {
		foundFlag = false
		for i, cpart := range combinedparts {
			if cpart.Value == ipart.Value &&
				cpart.Footprint == ipart.Footprint {
				// Match found, update combined parts list
				foundFlag = true
				combinedparts[i].References = append(cpart.References, ipart.References[0])
				combinedparts[i].Quantity++
			}
		}
		if !foundFlag {
			// part not already in combined parts list
			// add it
			combinedparts = append(combinedparts, ipart)
		}
	}
	return combinedparts, nil
}

type Items struct {
	data  interface{}
	field string
}

func (items *Items) Len() int {
	if reflect.ValueOf(items.data).Kind() != reflect.Slice {
		return -1
	}
	return reflect.ValueOf(items.data).Len()
}

func (items *Items) Less(i, j int) bool {
	a := reflect.ValueOf(items.data).Index(i)
	b := reflect.ValueOf(items.data).Index(j)
	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}
	if b.Kind() == reflect.Ptr {
		b = b.Elem()
	}
	reRef := regexp.MustCompile("[A-Z\\.]+")
	reVal := regexp.MustCompile("\\d*\\d+")
	// log.Infof("a -- ", a)
	// log.Infof("a.FieldByName -- ", a.FieldByName(items.field))
	// log.Infof("b -- ", b)
	// log.Infof("b.FieldByName -- ", b.FieldByName(items.field))

	refVal1 := string(reRef.FindAll([]byte(strings.ToUpper(a.FieldByName(items.field).String())), -1)[0])
	numVal1, _ := strconv.Atoi(string(reVal.FindAll([]byte(a.FieldByName(items.field).String()), -1)[0]))
	refVal2 := string(reRef.FindAll([]byte(strings.ToUpper(b.FieldByName(items.field).String())), -1)[0])
	numVal2, _ := strconv.Atoi(string(reVal.FindAll([]byte(b.FieldByName(items.field).String()), -1)[0]))

	return refVal1 < refVal2 || numVal1 < numVal2
}

func (items *Items) Swap(i, j int) {
	reflect.Swapper(items.data)(i, j)
}

func SortItems(i interface{}, str string) {
	if reflect.ValueOf(i).Kind() != reflect.Slice {
		return
	}
	a := &Items{
		data:  i,
		field: str,
	}
	sort.Sort(a)
}

func sortComponentRef(self types.EBOMItem) []string {
	type sortItem struct {
		References string
	}
	var sortItems []sortItem
	var ret []string

	for _, v := range self.References {
		sortItems = append(sortItems, sortItem{v})
	}

	SortItems(sortItems, "References")
	for _, v := range sortItems {
		//log.Infof("after sorted ref -- ", v)
		ret = append(ret, v.References)
	}
	return ret
}
