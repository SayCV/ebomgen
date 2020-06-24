package utils

import (
	//"math"
	//"unicode"
	"regexp"
	// "net/url"
	// "sort"
	//"strconv"
	"strings"

	"github.com/saycv/ebomgen/pkg/types"
	//log "github.com/sirupsen/logrus"
)

func mustRegexpMatch(pattern string, b []byte) bool {
	i, err := regexp.Match(pattern, b)
	if err != nil {
		panic(err)
	}
	return i
}

// NamerulesProcess classify the part
func NamerulesProcess(part types.EBOMItem, propvalue string, propfootprint string, propclass map[string]string) bool {
	_capREF := strings.ToUpper(part.References[0])
	_capVAL := strings.ToUpper(propvalue)
	rCharacter := regexp.MustCompile("[a-zA-Z]")

	// Step1 - process References common rules
	if strings.HasPrefix(_capREF, "C") && len(rCharacter.FindAllStringSubmatch(_capREF, -1)) < 4 {
		propclass["part"] = "Capacitor"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "R") && len(rCharacter.FindAllStringSubmatch(_capREF, -1)) < 4 {
		propclass["part"] = "Resistor"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "L") && len(rCharacter.FindAllStringSubmatch(_capREF, -1)) < 4 {
		propclass["part"] = "Inductor"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "F") && len(rCharacter.FindAllStringSubmatch(_capREF, -1)) < 4 {
		propclass["part"] = "Fuse"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "D") && len(rCharacter.FindAllStringSubmatch(_capREF, -1)) < 4 {
		propclass["part"] = "Diode"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "Q") && len(rCharacter.FindAllStringSubmatch(_capREF, -1)) < 4 {
		propclass["part"] = "Transistor"
		propclass["group"] = "Transistor"
	} else if strings.HasPrefix(_capREF, "X") {
		propclass["part"] = "Crystal"
		propclass["group"] = "Clock"
		if !strings.Contains(_capVAL, "MHZ") &&
			!strings.Contains(_capVAL, "M") &&
			!strings.Contains(_capVAL, "KHZ") &&
			!strings.Contains(_capVAL, "PPM") &&
			!strings.Contains(_capVAL, "XTAL") &&
			!strings.Contains(_capVAL, "OSC") {
			propclass["part"] = "Connector"
			propclass["group"] = "Connector"
		}
	} else if strings.HasPrefix(_capREF, "Y") {
		propclass["part"] = "Crystal"
		propclass["group"] = "Clock"
		if !strings.Contains(_capVAL, "MHZ") &&
			!strings.Contains(_capVAL, "M") &&
			!strings.Contains(_capVAL, "KHZ") &&
			!strings.Contains(_capVAL, "PPM") &&
			!strings.Contains(_capVAL, "XTAL") &&
			!strings.Contains(_capVAL, "OSC") {
			propclass["part"] = "Connector"
			propclass["group"] = "Connector"
		}
	} else if strings.HasPrefix(_capREF, "J") {
		propclass["part"] = "Connector"
		propclass["group"] = "Connector"
	} else if strings.HasPrefix(_capREF, "U") {
		propclass["part"] = "IC"
		propclass["group"] = "IC"
	} else if strings.HasPrefix(_capREF, "T") {
		propclass["part"] = "XFRM"
		propclass["group"] = "XFRM"
	} else if strings.HasPrefix(_capREF, "P") {
		propclass["part"] = "Connector"
		propclass["group"] = "Connector"
	} else if strings.HasPrefix(_capREF, "S") || strings.HasPrefix(_capREF, "K") {
		propclass["part"] = "Switch"
		propclass["group"] = "Mech"
	}

	// Step2 - process biliteral References
	if strings.HasPrefix(_capREF, "FB") {
		propclass["part"] = "FerritBead"
	} else if strings.HasPrefix(_capREF, "RV") {
		propclass["part"] = "TVS"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "RP") || strings.HasPrefix(_capREF, "RN") {
		propclass["part"] = "ResistorArray"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "CP") || strings.HasPrefix(_capREF, "CN") {
		propclass["part"] = "CapacitorArray"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "CM") {
		propclass["part"] = "Inductor"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "ED") {
		propclass["part"] = "DiodeESD"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "IC") {
		propclass["part"] = "IC"
		propclass["group"] = "IC"
	} else if strings.HasPrefix(_capREF, "JP") {
		propclass["part"] = "Connector"
	} else if strings.HasPrefix(_capREF, "CN") {
		propclass["part"] = "Connector"
	} else if strings.HasPrefix(_capREF, "BT") || strings.HasPrefix(_capREF, "BAT") {
		propclass["part"] = "BATT"
	} else if strings.HasPrefix(_capREF, "TP") || strings.Contains(_capVAL, "TEST") {
		propclass["part"] = "TestPoint"
	} else if strings.HasPrefix(_capREF, "TC") {
		propclass["part"] = "CapacitorTan"
		propclass["group"] = "Passive"
	}

	// Step3 - process trigram References
	if strings.HasPrefix(_capREF, "ANT") {
		//
	} else if strings.HasPrefix(_capREF, "ESD") {
		propclass["part"] = "DiodeESD"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "TVS") {
		propclass["part"] = "TVS"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "LED") {
		propclass["part"] = "LED"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "FET") {
		propclass["part"] = "FET"
		propclass["group"] = "Transistor"
	} else if strings.HasPrefix(_capREF, "REG") {
		propclass["part"] = "Reg"
		propclass["group"] = "IC"
	} else if strings.HasPrefix(_capREF, "COM") {
		propclass["part"] = "Connector"
		propclass["group"] = "Connector"
	} else if strings.HasPrefix(_capREF, "CON") {
		propclass["part"] = "Connector"
		propclass["group"] = "Connector"
	} else if strings.HasPrefix(_capREF, "OSC") {
		propclass["part"] = "Crystal"
		propclass["group"] = "Clock"
	} else if strings.HasPrefix(_capREF, "USB") {
		propclass["part"] = "ConnUSB"
		propclass["group"] = "Connector"
	} else if strings.HasPrefix(_capREF, "SIM") {
		propclass["part"] = "Connector"
		propclass["group"] = "Connector"
	}

	if strings.HasPrefix(_capREF, "MARK") {
		propclass["part"] = "MOUNT"
		propclass["group"] = "MOUNT"
	}

	// Step4 - process special footprint
	if strings.Contains(strings.ToUpper(propfootprint), "TP") {
		propclass["part"] = "TestPoint"
		propclass["group"] = "MOUNT"
	} else if strings.Contains(strings.ToUpper(propfootprint), "XFRM") {
		propclass["part"] = "XFRM"
		propclass["group"] = "XFRM"
	} else if strings.Contains(strings.ToUpper(propfootprint), "NET") &&
		strings.Contains(strings.ToUpper(propfootprint), "SHORT") {
		propclass["part"] = "NET_SHORT"
		propclass["group"] = "MOUNT"
	} else if strings.Contains(strings.ToUpper(propfootprint), "FERRIT") ||
		strings.Contains(strings.ToUpper(propfootprint), "BEAD") {
		propclass["part"] = "FerritBead"
		propclass["group"] = "Passive"
	}

	// Step5 - process special value
	if strings.Contains(_capVAL, "DNP") || strings.Contains(_capVAL, "DNI") ||
		(strings.Contains(_capVAL, "NP") && (!strings.Contains(_capVAL, "PNP") && !strings.Contains(_capVAL, "NPN") && !mustRegexpMatch(".*NPO", []byte(_capVAL)))) ||
		(strings.Contains(_capVAL, "NC") && propclass["group"] == "Passive") ||
		(strings.Contains(_capVAL, "NF") && !mustRegexpMatch("[0-9\\.]+NF", []byte(_capVAL)) && !mustRegexpMatch("[A-Z\\.]+NF[A-Z\\.]", []byte(_capVAL))) {
		propclass["part"] = "DNP"
		propclass["group"] = "MOUNT"
	} else if strings.Contains(_capVAL, "PAD") {
		propclass["part"] = "MOUNT"
		propclass["group"] = "MOUNT"
	} else if strings.Contains(_capVAL, "HOLE") ||
		(propclass["part"] == "unkownPart" && strings.HasPrefix(_capREF, "P")) ||
		(propclass["part"] == "unkownPart" && strings.HasPrefix(_capREF, "U")) ||
		strings.Contains(strings.Replace(_capVAL, " ", "", -1), "GOLDFINGERS") || strings.Contains(strings.Replace(_capVAL, " ", "", -1), "BOARDFINGERS") {
		propclass["part"] = "MOUNT"
		propclass["group"] = "MOUNT"
	} else if strings.Contains(_capVAL, "TP") && (!strings.Contains(_capVAL, "TPS") && !strings.Contains(_capVAL, "TPD")) {
		propclass["part"] = "TestPoint"
		propclass["group"] = "MOUNT"
	} else if strings.HasPrefix(_capREF, "D") &&
		(strings.Contains(_capVAL, "LED") || strings.Contains(_capVAL, "GRN") || strings.Contains(_capVAL, "RED") || strings.Contains(_capVAL, "YLW") || 
		strings.Contains(_capVAL, "GREEN") || strings.Contains(_capVAL, "YELLOW") || strings.Contains(propfootprint, "LED")) {
		propclass["part"] = "LED"
		propclass["group"] = "Passive"
	}

	if propclass["part"] == "Resistor" && (strings.Contains(_capVAL, "2%") || strings.Contains(_capVAL, "1%") ||
		strings.Contains(_capVAL, "0.5%") || strings.Contains(_capVAL, "0.1%")) {
		propclass["part"] = "ResistorHR"
	}
	if propclass["part"] == "Capacitor" && (IsCapTanFp(strings.ToUpper(propfootprint))) {
		propclass["part"] = "CapacitorTan"
	} else if propclass["part"] == "Capacitor" && IsCapAecFp(strings.ToUpper(propfootprint), _capVAL) {
		propclass["part"] = "CapacitorAec"
	}
	if propclass["part"] == "Crystal" && strings.Contains(_capVAL, "OSC") {
		propclass["part"] = "Oscillator"
	}

	if strings.HasPrefix(_capREF, "X") || strings.HasPrefix(_capREF, "Y") {
		if propclass["part"] == "Connector" {
			if strings.Contains(strings.ToUpper(propfootprint), "2016") || strings.Contains(strings.ToUpper(propfootprint), "2520") || strings.Contains(strings.ToUpper(propfootprint), "3225") ||
				strings.Contains(strings.ToUpper(propfootprint), "4025") || strings.Contains(strings.ToUpper(propfootprint), "5032") ||
				strings.Contains(strings.ToUpper(propfootprint), "6035") || strings.Contains(strings.ToUpper(propfootprint), "7050") ||
				strings.Contains(strings.ToUpper(propfootprint), "8045") {
				propclass["part"] = "Crystal"
			} else if strings.Contains(_capVAL, "32.768") {
				propclass["part"] = "Crystal"
			}
		}
	}
	if propclass["part"] == "Crystal" && (strings.Contains(_capVAL, "1.8V") || strings.Contains(_capVAL, "2.5V") || strings.Contains(_capVAL, "3.3V") || strings.Contains(_capVAL, "5V")) {
		propclass["part"] = "Oscillator"
	}

	if propclass["part"] == "Connector" {
		if mustRegexpMatch("RJ[0-9][0-9]", []byte(_capVAL)) {
			propclass["part"] = "ConnRJ"
		} else if strings.Contains(_capVAL, "USB") {
			propclass["part"] = "ConnUSB"
		}
	}

	if propclass["part"] == "unkownPart" {
		if strings.HasPrefix(_capREF, "S") &&
			strings.Contains(_capVAL, "SW") {
			propclass["part"] = "Switch"
			propclass["group"] = "Mech"
		} else if strings.HasPrefix(_capREF, "P") {
			propclass["part"] = "Connector"
			propclass["group"] = "Connector"
		}
	}

	return true
}
