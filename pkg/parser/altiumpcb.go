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

	cfg "github.com/larspensjo/config"
	log "github.com/sirupsen/logrus"
)

// parse ascii text file to retrieve parts
func parsePCBDocParts(filename string) (map[string]types.EBOMItem, error) {
	// Load a PADS Logic file and extract the part information

	partsList := make(map[string]types.EBOMItem)
	var partName string
	//var partDesc string
	var part types.EBOMItem
	var byteVal []byte
	var strVal string
	var err error
	//var partDecalNbr int64

	c := cfg.NewDefault()
	if len(c.Sections()) != 1 {
		log.Infof("Sections failure: invalid length")
	}

	byteAll, _ := ioutil.ReadFile(filename)
	reader := bufio.NewReader(bytes.NewBuffer(byteAll))
	//newspacelines := 0
	partName = ""
	part.Value = ""
	part.Attributes = map[string]string{
		"Manufacturer Part Number": "",
		"Description":              "unkownDesc",
		"part":                     "unkownPart",
		"group":                    "unkownGroup",
	}

	//pkgList := make(map[string]types.EBOMItem)
	//var pkgDesc types.EBOMItem
	//pkgName := ""

	//log.Debugf("Starting process: %v", filename)
	for {
		byteVal, _, err = reader.ReadLine()
		strVal = strings.TrimSpace(string(byteVal))
		//log.Infof("ext - %s", strVal)
		if err == nil && strings.HasPrefix(strVal, "|RECORD=Component") {
			listVal := strings.Split(strVal, "|")
			//log.Infof(strings.Join(listVal[1:], ","))
			c.AddSection(listVal[2])
			for _, v := range listVal[1:] {
				opt := strings.Split(v, "=")
				c.AddOption(listVal[2], opt[0], opt[1])
			}

		} else if err == nil && strings.HasPrefix(strVal, "|RECORD=ParamItem") {
			listVal := strings.Split(strVal, "|")
			//log.Infof(strings.Join(listVal[1:], ","))
			_val := strings.Split(listVal[2], "=")
			index := strings.Split(_val[1], "#")
			for _, v := range listVal[1:] {
				opt := strings.Split(v, "=")
				c.AddOption("ID="+string(index[1]), opt[0], opt[1])
			}

		} else if err != nil {
			break
		}
	}

	for _, orderedSection := range c.Sections() {
		log.Infof(orderedSection)
		if orderedSection == "DEFAULT" {
			continue
		}
		partName, _ = c.String(orderedSection, "SOURCEDESIGNATOR")
		part.References = []string{partName}
		part.Value, _ = c.String(orderedSection, "VALUE")
		part.Footprint, _ = c.String(orderedSection, "PATTERN")
		part.Desc, _ = c.String(orderedSection, "FOOTPRINTDESCRIPTION")

		partsList[partName] = part
		//log.Infof("partName - %v", partName)
		partName = ""
		part.Value = ""
		part.Attributes = map[string]string{
			"Manufacturer Part Number": "",
			"Description":              "unkownDesc",
			"part":                     "unkownPart",
			"group":                    "unkownGroup",
		}
	}
	//print(part in partsList)
	log.Infof("Parse Done.")
	//c.WriteFile("config.cfg", 0644, "Components list file")
	return partsList, err
}

// ExtractAltiumPCBComponents Load a Altium PcbDoc ascii file and extract the part information
func ExtractAltiumPCBComponents(filename string) ([]types.EBOMItem, error) {
	var propclass = map[string]string{
		"Description": "unkownDesc",
		"part":        "unkownPart",
		"group":       "unkownGroup",
	}
	var components []types.EBOMItem

	input, err := os.Open(filename)
	if err != nil {
		return components, err
	}
	defer input.Close()

	parts, _ := parsePCBDocParts(filename)

	propfootprint := "unkownFp"
	propvalue := "unkownval"

	for _, part := range parts {
		//primitive = Primitive(chips[part.description()])
		propfootprint = "unkownFp"
		propvalue = "unkownVal"
		propclass["part"] = "unkownPart"
		propclass["group"] = "unkownGroup"

		if part.Value != "" {
			propvalue = part.Value
		}
		if part.Footprint != "" {
			propfootprint = part.Footprint
		}
		if part.Attributes["Manufacturer Part Number"] != "" {
			propclass["Manufacturer Part Number"] = part.Attributes["Manufacturer Part Number"]
		}

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

		part.Attributes = map[string]string{
			"Description": propclass["Description"],
			"part":        propclass["part"],
			"group":       propclass["group"],
		}
		if part.Quantity == 0 {
			part.Quantity = 1
		}

		components = append(components, part)

		//log.Infof("Add Part - %v", part)
		//log.Infof("  propclass - %v", propclass)
	}
	return components, err
}
