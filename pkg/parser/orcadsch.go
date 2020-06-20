
package parser

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"

	//"strconv"
	"strings"

	"github.com/saycv/ebomgen/pkg/types"
	"github.com/saycv/ebomgen/pkg/utils"

	//cfg "github.com/larspensjo/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func stripQuotes(name string) string {
	ind1 := 0
	ind2 := 0
	if (len(name) > 1) {
        ind1 = strings.Index(name, "'")
		ind2 = strings.Index(name[ind1 + 1:], "'")
	}
    if ind1 < 0 || ind2 < 0 {
		return name
	}
    return name[ind1 + 1:ind2 + ind1 + 1]
}

func parseField(field string) ([]string, error) {
    ind := strings.Index(field, "=")
    if (len(field) > 1) && ind > 0 {
        fieldkey := field[0:ind]
        fieldvalue := stripQuotes(field[ind + 1:])
		return []string{fieldkey, fieldvalue}, nil
	}
	//log.Errorf("No validate fileld")
	return []string{""}, errors.Errorf("No validate fileld")
}

func parsePin(p string) (string, error) {
    pn := p
    pinnumber := strings.Split(pn, ",")
    for _, pin := range pinnumber {
        if pin != "0" {
			return pin, nil
		}
	}
    log.Errorf("Unable to parse pin")
    return "", errors.Errorf("Unable to parse pin")
}

// OrcadSplit split spaces unless between quotes
func OrcadSplit(s string) []string {
	indQuote1 := strings.Index(s, "'")
	indSpace := strings.Index(s, " ")
    if indSpace < 0 {
		return []string{s}
	}
    if indQuote1 < 0 {
		return strings.Split(s, " ")
	}
	indQuote2 := strings.Index(string(s[indQuote1+1:]), "'")
	return []string{s[:indQuote1], s[indQuote1+1:indQuote2 + indQuote1 + 1]}
}

func ParseXprt(filename string) (map[string]types.EBOMItem, error) {
	partsList := make(map[string]types.EBOMItem)
	//var partName string
	//var partDesc string
	var part types.EBOMItem
	var byteVal []byte
	var strVal string
	var err error

	byteAll, _ := ioutil.ReadFile(filename)
	reader := bufio.NewReader(bytes.NewBuffer(byteAll))
	
	ignoreReadLine := false
	for {
		if !ignoreReadLine {
			byteVal, _, err = reader.ReadLine()
			strVal = strings.TrimSpace(string(byteVal))
		}
		ignoreReadLine = false
        if err == nil && strVal == "PART_NAME" {
            byteVal, _, err = reader.ReadLine()
			strVal = strings.TrimSpace(string(byteVal))
			vallist := OrcadSplit(strVal)
            part.References = []string{vallist[0]}
			part.Desc = vallist[1]
			part.Attributes = make(map[string]string)
            for {
                byteVal, _, err = reader.ReadLine()
				strVal = strings.TrimSpace(string(byteVal))
                if err == nil && (strVal == "PART_NAME" || strVal == "END.") {
					ignoreReadLine = true
					break
                } else {
					propLine, err := parseField(strVal)
					if err == nil {
						part.Attributes[propLine[0]] = propLine[1]
					}
				}
			}
            partsList[part.References[0]] = part
		} else if err == nil && strVal == "END." {
			break
		}
	}
    return partsList, nil
}

func ParseXnet(filename string) (map[string]types.EBOMItem, error) {
	partsList := make(map[string]types.EBOMItem)
	//var partName string
	//var partDesc string
	var part types.EBOMItem
	var byteVal []byte
	var strVal string
	var err error

	byteAll, _ := ioutil.ReadFile(filename)
	reader := bufio.NewReader(bytes.NewBuffer(byteAll))
	
	ignoreReadLine := false
	for {
		if !ignoreReadLine {
			byteVal, _, err = reader.ReadLine()
			strVal = strings.TrimSpace(string(byteVal))
		}
		ignoreReadLine = false
        if err == nil && strVal == "PART_NAME" {
            byteVal, _, err = reader.ReadLine()
			strVal = strings.TrimSpace(string(byteVal))
			vallist := OrcadSplit(strVal)
            part.References = []string{vallist[0]}
			part.Desc = vallist[1]
			part.Attributes = make(map[string]string)
            for {
                byteVal, _, err = reader.ReadLine()
				strVal = strings.TrimSpace(string(byteVal))
                if err == nil && (strVal == "PART_NAME" || strVal == "END.") {
					ignoreReadLine = true
					break
                } else {
					propLine, err := parseField(strVal)
					if err == nil {
						part.Attributes[propLine[0]] = propLine[1]
					}
				}
			}
            partsList[part.References[0]] = part
		} else if err == nil && strVal == "END." {
			break
		}
	}
    return partsList, nil
}

func ParseChip(filename string) (map[string]types.EBOMItem, error) {
	partsList := make(map[string]types.EBOMItem)
	//var partName string
	//var partDesc string
	var part types.EBOMItem
	var byteVal []byte
	var strVal string
	var err error

	byteAll, _ := ioutil.ReadFile(filename)
	reader := bufio.NewReader(bytes.NewBuffer(byteAll))
	
	ignoreReadLine := false
	for {
		if !ignoreReadLine {
			byteVal, _, err = reader.ReadLine()
			strVal = strings.TrimSpace(string(byteVal))
			//log.Infof(strVal)
		}
		ignoreReadLine = false
        if err == nil && strings.HasPrefix(strVal, "primitive") {
            //byteVal, _, err = reader.ReadLine()
			//strVal = strings.TrimSpace(string(byteVal))
			//log.Infof(strVal)
			vallist := OrcadSplit(strVal)
			if len(vallist) < 2 {
				continue
			}
            part.References = []string{vallist[0]}
			part.Desc = vallist[1]
			part.Attributes = make(map[string]string)
			part.Attributes["PIN_NUMBER"] = ""
			
            for {
                byteVal, _, err = reader.ReadLine()
				strVal = strings.TrimSpace(string(byteVal))
				//log.Infof(strVal)
                if err == nil && (strings.HasPrefix(strVal, "primitive") || strVal == "END.") {
					ignoreReadLine = true
					break
                } else {
					propLine, err := parseField(strVal)
					if err == nil {
						if propLine[0] == "PIN_NUMBER" {
							pin,_ := parsePin(propLine[1])
							pin = strings.Replace(pin, "(", "", -1)
							pin = strings.Replace(pin, ")", "", -1)
							part.Attributes[propLine[0]] = part.Attributes[propLine[0]] + pin + ","
						} else if propLine[0] == "PART_NAME" {
							part.Attributes[propLine[0]] = propLine[1]
						} else if propLine[0] != "PINUSE" {
							part.Attributes[propLine[0]] = propLine[1]
						}
					}
				}
			}
            partsList[part.Desc] = part
		} else if err == nil && strVal == "END." {
			break
		}
	}
    return partsList, nil
}

// ExtractOrcadSchComponents Load a sch file and extract the part information
func ExtractOrcadSchComponents(filename string) ([]types.EBOMItem, error) {
	var propclass = map[string]string{
		"Description": "unkownDesc",
		"part":        "unkownPart",
		"group":       "unkownGroup",
	}
	var components []types.EBOMItem

    chip := filename + "/pstchip.dat"
    xnet := filename + "/pstxnet.dat"
	xprt := filename + "/pstxprt.dat"
	
	input, err := os.Open(chip)
	if err != nil {
		return components, err
	}
	defer input.Close()
	
	input, err = os.Open(xnet)
	if err != nil {
		return components, err
	}
	defer input.Close()
	
	input, err = os.Open(xprt)
	if err != nil {
		return components, err
	}
	defer input.Close()

	parts, _ := ParseXprt(xprt)
	chips, _ := ParseChip(chip)

	propfootprint := "unkownFp"
	propvalue := "unkownval"

	for _, part := range parts {
		//primitive = Primitive(chips[part.description()])
		propfootprint = "unkownFp"
		propvalue = "unkownVal"
		propclass["part"] = "unkownPart"
		propclass["group"] = "unkownGroup"
		chipdesc := chips[part.Desc]

		part.Attributes["chip"] = part.Desc

		if chipdesc.Attributes["VALUE"] != "" {
			propvalue = chipdesc.Attributes["VALUE"]
			part.Value = propvalue
		}
		if chipdesc.Attributes["JEDEC_TYPE"] != "" {
			propfootprint = chipdesc.Attributes["JEDEC_TYPE"]
			part.Footprint = propfootprint
		}
		if chipdesc.Attributes["PART_NUMBER"] != "" {
			propclass["PART_NUMBER"] = chipdesc.Attributes["PART_NUMBER"]
			part.Attributes["PART_NUMBER"] = propclass["PART_NUMBER"]
		}
		if chipdesc.Attributes["PART_NAME"] != "" {
			propclass["Description"] = chipdesc.Attributes["PART_NAME"]
			part.Desc = propclass["Description"]
		}

		utils.NamerulesProcess(part, propvalue, propfootprint, propclass)

		// Build BOMpart and append to list
		propclass["Description"] = propclass["part"]

		part.Attributes["Description"] = propclass["Description"]
		part.Attributes["part"] = propclass["part"]
		part.Attributes["group"] = propclass["group"]
		if part.Quantity == 0 {
			part.Quantity = 1
		}

		components = append(components, part)

		//log.Infof("Add Part - %v", part)
		//log.Infof("  propclass - %v", propclass)
	}
	return components, err
}
