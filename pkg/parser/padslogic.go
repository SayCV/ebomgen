package parser

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/prometheus/common/log"
	"github.com/saycv/ebomgen/pkg/types"
	"github.com/saycv/ebomgen/pkg/utils"
)

// parse ascii text file to retrieve parts
func parseTextParts(filename string) ([]types.EBOMItem, error) {
	// Load a PADS Logic file and extract the part information

	var partsList []types.EBOMItem
	var partName string
	var partDesc string
	var part types.EBOMItem
	var byteVal []byte
	var strVal string
	var err error

	byteAll, _ := ioutil.ReadFile(filename)
	reader := bufio.NewReader(bytes.NewBuffer(byteAll))
	newspacelines := 0
	partName = ""

	log.Debugf("Starting process: %v", filename)
	for {
		byteVal, _, err = reader.ReadLine()
		strVal = strings.TrimSpace(string(byteVal))
		if err == nil && strVal == "*PART*       ITEMS" {
			byteVal, _, err = reader.ReadLine()
			byteVal, _, err = reader.ReadLine()
			// byteVal, _, err = reader.ReadLine()
			newspacelines = 0
			for {
				if newspacelines > 1 {
					break
				}
				byteVal, _, err = reader.ReadLine()
				strVal = strings.TrimSpace(string(byteVal))
				if strVal == "\n" || strVal == "" {
					newspacelines = newspacelines + 1
					continue
				}
				// f_ = sch_f.readline().strip()
				if partName == "" && newspacelines == 1 {
					kv := strings.Split(strVal, " ")
					partName = strings.Split(kv[0], "-")[0]
					partDesc = kv[1]
					part.References = []string{partName}
					part.Desc = partDesc
					continue
				}
				for {
					byteVal, _, err = reader.ReadLine()
					strVal = strings.TrimSpace(string(byteVal))
					if strVal == "\n" || strVal == "" {
						break
					}
					if strings.HasPrefix(strings.ToUpper(strVal), "\"PCB DECAL\"") {
						part.Footprint = strings.Replace(strings.ToUpper(strVal), "\"PCB DECAL\"", "", -1)
					}
					if strings.HasPrefix(strings.ToUpper(strVal), "\"VALUE\"") {
						//part.Value = strings.Replace(strings.ToUpper(strVal), "\"VALUE\"", "", -1)
						part.Value = strings.TrimSpace(strVal[7:])
					}
					if strings.HasPrefix(strings.ToUpper(strVal), "\"CVPART NUMBER\"") {
						//part.Attributes["Manufacturer Part Number"] = strings.Replace(strings.ToUpper(strVal), "\"PART NUMBER\"","", -1)
					}
				}
				partsList = append(partsList, part)
				log.Infof("partName - %v", partName)
				partName = ""
			}
		} else if err == nil && strVal == "*BUSSES*" {
			log.Infof("*BUSSES*")
			continue
		} else if err == nil && strings.HasPrefix(strings.ToUpper(strVal), "*END*") {
			log.Infof("*END*")
			break
		}
	}
	//print(part in partsList)
	return partsList, err
}

// ExtractPADSLogicComponents Load a PADS Logic file and extract the part information
func ExtractPADSLogicComponents(filename string) bool {
	var propclass = map[string]string{
		"Description": "unkownDesc",
		"part":        "unkownPart",
		"group":       "unkownGroup",
	}
	//var components []types.EBOMItem

	parts, _ := parseTextParts(filename)

	propfootprint := "unkownFp"
	propvalue := "unkownval"

	for _, part := range parts {
		//primitive = Primitive(chips[part.description()])
		propfootprint = "unkownFp"
		propvalue = "unkownVal"
		propclass["part"] = "unkownPart"
		propclass["group"] = "unkownGroup"

		// for key, prop in enumerate(primitive.name().properties().items()) {
		//     if prop[0] == "JEDEC_TYPE"{
		//         propfootprint = prop[1]
		//     //} else if prop[0] == "CLASS"{
		//     // propclass[prop[0]] = prop[1]
		//     // ignore "CLASS" attribute
		//     } else if prop[0] == "VALUE"{
		//         propvalue = prop[1]
		//         #propvalue = str(prop[1]).lower().replace("pf", "pF").replace("nf", "nF").replace("uf", "uF")
		//     } else if prop[0] == "PART_NUMBER"{
		//         propclass["Manufacturer Part Number"] = prop[1]
		//     } else if prop[0] == "PART_NAME"{
		//         propclass["Description"] = prop[1]
		//     }
		// }

		utils.NamerulesProcess(part, propvalue, propfootprint, propclass)

		// Build BOMpart and append to list
		propclass["Description"] = propclass["part"]

		log.Infof("Add Part - %v", part)
		log.Infof("  propclass - %v", propclass)
	}
	return true
}
