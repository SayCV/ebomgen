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
	// Step1 - process References common rules
	if strings.HasPrefix(_capREF, "C") {
		propclass["part"] = "Capacitor"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "R") {
		propclass["part"] = "Resistor"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "L") {
		propclass["part"] = "Inductor"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "F") {
		propclass["part"] = "Fuse"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "D") {
		propclass["part"] = "Diode"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "Q") {
		propclass["part"] = "Transistor"
		propclass["group"] = "Transistor"
	} else if strings.HasPrefix(_capREF, "X") {
		propclass["part"] = "Crystal"
		propclass["group"] = "Clock"
		if !strings.Contains(_capVAL, "MHZ") &&
			!strings.Contains(_capVAL, "M") &&
			!strings.Contains(_capVAL, "KHZ") &&
			!strings.Contains(_capVAL, "PPM") {
			propclass["part"] = "Connector"
			propclass["group"] = "Connector"
		}
	} else if strings.HasPrefix(_capREF, "Y") {
		propclass["part"] = "Crystal"
		propclass["group"] = "Clock"
		if !strings.Contains(_capVAL, "MHZ") &&
			!strings.Contains(_capVAL, "M") &&
			!strings.Contains(_capVAL, "KHZ") &&
			!strings.Contains(_capVAL, "PPM") {
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
	} else if strings.HasPrefix(_capREF, "RP") {
		propclass["part"] = "ResistorArray"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "CP") {
		propclass["part"] = "CapacitorArray"
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
	if strings.HasPrefix(_capREF, "LED") {
		propclass["part"] = "LED"
		propclass["group"] = "Passive"
	} else if strings.HasPrefix(_capREF, "FET") {
		propclass["part"] = "FET"
		propclass["group"] = "Transistor"
	} else if strings.HasPrefix(_capREF, "REG") {
		propclass["part"] = "Reg"
		propclass["group"] = "IC"
	} else if strings.HasPrefix(_capREF, "CON") {
		propclass["part"] = "Connector"
		propclass["group"] = "Connector"
	} else if strings.HasPrefix(_capREF, "OSC") {
		propclass["part"] = "Crystal"
		propclass["group"] = "Clock"
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
		strings.Contains(_capVAL, "NC") ||
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
		(strings.Contains(_capVAL, "LED") || strings.Contains(_capVAL, "GRN") || strings.Contains(_capVAL, "RED") || strings.Contains(_capVAL, "YLW")) {
		propclass["part"] = "LED"
		propclass["group"] = "Passive"
	}

	if propclass["part"] == "Capacitor" && ( strings.Contains(strings.ToUpper(propfootprint), "3216") ||
		strings.Contains(strings.ToUpper(propfootprint), "3528") ||
		strings.Contains(strings.ToUpper(propfootprint), "6032") ||
		strings.Contains(strings.ToUpper(propfootprint), "7343") ) {
		propclass["part"] = "CapacitorTan"
	}
	if propclass["part"] == "Crystal" && strings.Contains(_capVAL, "OSC") {
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
