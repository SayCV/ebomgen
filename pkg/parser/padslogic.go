package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/saycv/ebomgen/pkg/types"
	"github.com/saycv/ebomgen/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type SchConfig struct {
	FpMap []any
}

// parse ascii text file to retrieve parts
func parseTextParts(filename string) (map[string]types.EBOMItem, error) {
	// Load a PADS Logic file and extract the part information

	partsList := make(map[string]types.EBOMItem)
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
	part.Value = ""
	part.Attributes = map[string]string{
		"Manufacturer Part Number": "",
		"Description":              "unkownDesc",
		"part":                     "unkownPart",
		"group":                    "unkownGroup",
	}

	//log.Debugf("Starting process: %v", filename)
	for {
		byteVal, _, err = reader.ReadLine()
		strVal = strings.TrimSpace(string(byteVal))
		//log.Infof("ext - %s", strVal)
		if err == nil && strVal == "*PART*       ITEMS" {
			byteVal, _, err = reader.ReadLine()
			byteVal, _, err = reader.ReadLine()
			// byteVal, _, err = reader.ReadLine()
			newspacelines = 0
			for {
				if err != nil || newspacelines > 1 {
					break
				}
				byteVal, _, err = reader.ReadLine()
				strVal = strings.TrimSpace(string(byteVal))
				if err != nil || strVal == "\n" || strVal == "" {
					newspacelines = newspacelines + 1
					continue
				}
				// f_ = sch_f.readline().strip()
				if partName == "" && newspacelines == 1 {
					kv := strings.Fields(strVal)
					partName = strings.Split(kv[0], "-")[0]
					partDesc = kv[1]
					part.References = []string{partName}
					part.Desc = partDesc
					//if partName == "P1" {
					//	log.Infof("Add Part - %s， %s", part.Value, part.Desc)
					//}
					continue
				}
				for {
					byteVal, _, err = reader.ReadLine()
					strVal = strings.TrimSpace(string(byteVal))
					//log.Infof("line - %s", strVal)
					if err != nil || strVal == "\n" || strVal == "" {
						break
					}
					if strings.HasPrefix(strings.ToUpper(strVal), "\"PCB DECAL\"") {
						part.Footprint = strings.Replace(strings.ToUpper(strVal), "\"PCB DECAL\"", "", -1)
						part.Footprint = strings.TrimSpace(part.Footprint)
					}
					if strings.HasPrefix(strings.ToUpper(strVal), "\"VALUE\"") {
						//part.Value = strings.Replace(strings.ToUpper(strVal), "\"VALUE\"", "", -1)
						part.Value = strings.TrimSpace(strVal[7:])
					}
					if strings.HasPrefix(strings.ToUpper(strVal), "\"PART NUMBER\"") {
						part.Attributes["Manufacturer Part Number"] = strings.Replace(strings.ToUpper(strVal), "\"PART NUMBER\"", "", -1)
					}
				}

				if part.Value == "" {
					if part.Desc != "" {
						part.Value = part.Desc
					} else if part.Attributes["Manufacturer Part Number"] != "" {
						part.Value = part.Attributes["Manufacturer Part Number"]
					}
				}
				//partsList = append(partsList, part)
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
		} else if err == nil && strVal == "*BUSSES*" {
			//log.Infof("*BUSSES*")
			continue
		} else if err == nil && strings.HasPrefix(strings.ToUpper(strVal), "*END*") {
			//log.Infof("*END*")
			break
		} else if err != nil {
			break
		}
	}
	//print(part in partsList)
	log.Infof("Parse Done.")
	return partsList, err
}

// ExtractPADSLogicComponents Load a PADS Logic file and extract the part information
func ExtractPADSLogicComponents(filename string) ([]types.EBOMItem, error) {
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

	dirname := filepath.Dir(filename)
	basename := filepath.Base(filename)
	suffix := filepath.Ext(filename)
	filenameWithoutSuffix := strings.TrimSuffix(basename, suffix)
	fnSlice := strings.Split(filenameWithoutSuffix, "-")

	configfile_toml := path.Join(dirname, strings.Join(fnSlice[:1], "-")+"-metadata.toml")
	if _, err := os.Stat(configfile_toml); err != nil {
		configfile_toml = path.Join(dirname, fnSlice[0]+"-metadata.toml")
	}
	if _, err := os.Stat(configfile_toml); err != nil {
		configfile_toml = path.Join(dirname, "metadata.toml")
	}

	var fpmaps []interface{}
	if tree, err := toml.LoadFile(configfile_toml); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Try to prase metadata from %s", configfile_toml)
		fpmaps = tree.Get("FpMap").([]interface{})
	}

	parts, _ := parseTextParts(filename)

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
			// fmt.Printf("FpMap(): %s\n", fpmaps)
			for _, entry := range fpmaps {
				// 转换为数组类型
				fpMapEntry := entry.([]interface{})
				if len(fpMapEntry) == 2 {
					// 获取包名和类型
					packageOld := fpMapEntry[0].(string)
					packageNew := fpMapEntry[1].(string)
					//fmt.Printf("packageOld: %s, packageNew: %s\n", packageOld, packageNew)

					if part.Footprint == packageOld {
						log.Infof("preprocess %v / %v / %v", part.Footprint, part.Value, part.References)
						propfootprint = packageNew
						part.Footprint = propfootprint
					}
				}
			}

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
