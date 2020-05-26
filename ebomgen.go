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

	"github.com/saycv/ebomgen/pkg/utils"

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

var partgroups []types.EBOMGroup

func AddGroup3(groupdeflist []types.EBOMGroup, reference string, pt string, gt string) []types.EBOMGroup {
	groupdeflist = append(groupdeflist,
		types.EBOMGroup{reference, pt, gt, 0, "?"})
	return groupdeflist
}

func AddGroup(groupdeflist []types.EBOMGroup, reference string, pt string, gt string, unt string) []types.EBOMGroup {
	groupdeflist = append(groupdeflist,
		types.EBOMGroup{reference, pt, gt, 0, unt})
	return groupdeflist
}

func SetPrecedence(groupdeflist []types.EBOMGroup) []types.EBOMGroup {
	var ret []types.EBOMGroup

	GRPID_FACTOR := 1000
	PASSIVE_GRPID := 0 * GRPID_FACTOR
	DISCRETE_GRPID := 1 * GRPID_FACTOR
	IC_GRPID := 2 * GRPID_FACTOR
	CONNECTOR_GRPID := 3 * GRPID_FACTOR
	UNK_GRPID := 4 * GRPID_FACTOR
	MECH_GRPID := 5 * GRPID_FACTOR
	MOUNT_GRPID := 6 * GRPID_FACTOR

	PARTID_FACTOR := 10
	PARTID_NBRS := 10
	BASE_PARTID := PASSIVE_GRPID

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	CAPGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//CAPSM_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//CAPSM3_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	CAPSM8_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//CAPSMPOL_PARTID := BASE_PARTID + 3*PARTID_FACTOR
	//CAPTM_PARTID := BASE_PARTID + 4*PARTID_FACTOR
	//CAPTM3_PARTID := BASE_PARTID + 5*PARTID_FACTOR
	//CAPTM8_PARTID := BASE_PARTID + 6*PARTID_FACTOR
	//CAPTMPOL_PARTID := BASE_PARTID + 7*PARTID_FACTOR

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	RESGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//RESSM_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//RESSM3_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	RESSM8_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//RESSMPOL_PARTID := BASE_PARTID + 3*PARTID_FACTOR
	//RESTM_PARTID := BASE_PARTID + 4*PARTID_FACTOR
	//RESTM3_PARTID := BASE_PARTID + 5*PARTID_FACTOR
	//RESTM8_PARTID := BASE_PARTID + 6*PARTID_FACTOR
	//RESTMPOL_PARTID := BASE_PARTID + 7*PARTID_FACTOR

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	INDGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//INDSM_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//INDSM3_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	//INDSM8_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//INDXFRM_PARTID := BASE_PARTID + 2*PARTID_FACTOR + 1
	INDFBSM_PARTID := BASE_PARTID + 3*PARTID_FACTOR
	INDFUSESM_PARTID := BASE_PARTID + 4*PARTID_FACTOR
	//INDTM_PARTID := BASE_PARTID + 5*PARTID_FACTOR
	//INDTM3_PARTID := BASE_PARTID + 6*PARTID_FACTOR
	//INDTM8_PARTID := BASE_PARTID + 7*PARTID_FACTOR
	//INDXFRM_PARTID = BASE_PARTID + 7*PARTID_FACTOR + 1
	//INDFBTM_PARTID := BASE_PARTID + 8*PARTID_FACTOR
	//INDFUSETM_PARTID := BASE_PARTID + 9*PARTID_FACTOR

	PARTID_FACTOR = 10
	PARTID_NBRS = 10
	BASE_PARTID = DISCRETE_GRPID

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	DIODEGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//DIODESM_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//DIODESM3_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	//DIODESM8_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//DIODESMPOL_PARTID := BASE_PARTID + 3*PARTID_FACTOR
	//DIODETM_PARTID := BASE_PARTID + 4*PARTID_FACTOR
	//DIODETM3_PARTID := BASE_PARTID + 5*PARTID_FACTOR
	//DIODETM8_PARTID := BASE_PARTID + 6*PARTID_FACTOR
	//DIODETMPOL_PARTID := BASE_PARTID + 7*PARTID_FACTOR

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	//ZENERGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//ZENERSM_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//ZENERSM3_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	//ZENERTM_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//ZENERTM3_PARTID := BASE_PARTID + 3*PARTID_FACTOR

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	//TVSGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//TVSSM_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//TVSSM3_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	//TVSTM_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//TVSTM3_PARTID := BASE_PARTID + 3*PARTID_FACTOR

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	LEDGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//LEDSM_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//LEDSM3_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	//LEDTM_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//LEDTM3_PARTID := BASE_PARTID + 3*PARTID_FACTOR

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	TRANSGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//TRANSSM3_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//MOSSM3_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	//TRANSTM3_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//MOSTM3_PARTID := BASE_PARTID + 3*PARTID_FACTOR

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	CLKGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//CLKSM2_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	CLKSM4_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	//CLKSM6_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//CLKTM2_PARTID := BASE_PARTID + 3*PARTID_FACTOR
	//CLKTM4_PARTID := BASE_PARTID + 4*PARTID_FACTOR
	//CLKTM6_PARTID := BASE_PARTID + 5*PARTID_FACTOR

	PARTID_FACTOR = 10
	PARTID_NBRS = 10
	BASE_PARTID = IC_GRPID

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	ICGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//MCU_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//CPU_PARTID := BASE_PARTID + 0*PARTID_FACTOR + 1
	//DSP_PARTID := BASE_PARTID + 0*PARTID_FACTOR + 2
	REG_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	//TXRX_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//COMP_PARTID := BASE_PARTID + 3*PARTID_FACTOR
	//OPAMP_PARTID := BASE_PARTID + 4*PARTID_FACTOR
	//IC_PARTID := BASE_PARTID + 5*PARTID_FACTOR
	XFRM_PARTID := BASE_PARTID + 6*PARTID_FACTOR

	PARTID_FACTOR = 10
	PARTID_NBRS = 10
	BASE_PARTID = CONNECTOR_GRPID

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	CONGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//CONSM_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	CONTM_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	//CONSW_PARTID := BASE_PARTID + 3*PARTID_FACTOR
	//CONSD_PARTID := BASE_PARTID + 4*PARTID_FACTOR
	CONUSB_PARTID := BASE_PARTID + 5*PARTID_FACTOR
	//CONDB_PARTID := BASE_PARTID + 6*PARTID_FACTOR
	//CONLEMO_PARTID := BASE_PARTID + 7*PARTID_FACTOR

	PARTID_FACTOR = 10
	PARTID_NBRS = 10
	BASE_PARTID = UNK_GRPID

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	SWITCH_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	BATT_PARTID := BASE_PARTID + 1*PARTID_FACTOR

	UNKGEN_PARTID := BASE_PARTID + 9*PARTID_FACTOR

	PARTID_FACTOR = 10
	PARTID_NBRS = 10
	BASE_PARTID = MECH_GRPID

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	//MECHGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//MECHSM_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//MECHTM_PARTID := BASE_PARTID + 1*PARTID_FACTOR

	PARTID_FACTOR = 10
	PARTID_NBRS = 10
	BASE_PARTID = MOUNT_GRPID

	BASE_PARTID = BASE_PARTID + PARTID_NBRS*PARTID_FACTOR
	MNTGEN_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//MNTSM_PARTID := BASE_PARTID + 0*PARTID_FACTOR
	//MNTTM_PARTID := BASE_PARTID + 1*PARTID_FACTOR
	MNTTP_PARTID := BASE_PARTID + 2*PARTID_FACTOR
	MNTNETSHORT_PARTID := BASE_PARTID + 2*PARTID_FACTOR + 1
	MNTDNP_PARTID := BASE_PARTID + 3*PARTID_FACTOR

	DFLT_PARTID := BASE_PARTID + PARTID_NBRS*PARTID_FACTOR

	for _, groupdef := range groupdeflist {
		groupdef.Precedence = DFLT_PARTID

		if groupdef.GroupType == "IC" {
			groupdef.Precedence = ICGEN_PARTID
		} else if groupdef.GroupType == "DNP" {
			groupdef.Precedence = MNTDNP_PARTID
		}

		if groupdef.PartType == "Capacitor" {
			groupdef.Precedence = CAPGEN_PARTID
		} else if groupdef.PartType == "CapacitorArray" {
			groupdef.Precedence = CAPSM8_PARTID
		} else if groupdef.PartType == "Resistor" {
			groupdef.Precedence = RESGEN_PARTID
		} else if groupdef.PartType == "ResistorArray" {
			groupdef.Precedence = RESSM8_PARTID
		} else if groupdef.PartType == "Inductor" {
			groupdef.Precedence = INDGEN_PARTID
		} else if groupdef.PartType == "FerritBead" {
			groupdef.Precedence = INDFBSM_PARTID
		} else if groupdef.PartType == "Fuse" {
			groupdef.Precedence = INDFUSESM_PARTID
		} else if groupdef.PartType == "Diode" {
			groupdef.Precedence = DIODEGEN_PARTID
		} else if groupdef.PartType == "LED" {
			groupdef.Precedence = LEDGEN_PARTID
		} else if groupdef.PartType == "Transistor" {
			groupdef.Precedence = TRANSGEN_PARTID
		} else if groupdef.PartType == "FET" {
			groupdef.Precedence = TRANSGEN_PARTID
		} else if groupdef.PartType == "Crystal" {
			groupdef.Precedence = CLKGEN_PARTID
		} else if groupdef.PartType == "Oscillator" {
			groupdef.Precedence = CLKSM4_PARTID
		} else if groupdef.PartType == "XFRM" {
			groupdef.Precedence = XFRM_PARTID
		} else if groupdef.PartType == "Switch" {
			groupdef.Precedence = SWITCH_PARTID
		} else if groupdef.PartType == "BATT" {
			groupdef.Precedence = BATT_PARTID
		} else if groupdef.PartType == "Connector" {
			groupdef.Precedence = CONGEN_PARTID
		} else if groupdef.PartType == "ConnRJ" {
			groupdef.Precedence = CONTM_PARTID
		} else if groupdef.PartType == "ConnUSB" {
			groupdef.Precedence = CONUSB_PARTID
		} else if groupdef.PartType == "unkownPart" {
			groupdef.Precedence = UNKGEN_PARTID
		} else if groupdef.PartType == "TestPoint" {
			groupdef.Precedence = MNTTP_PARTID
		} else if groupdef.PartType == "NET_SHORT" {
			groupdef.Precedence = MNTNETSHORT_PARTID
		} else if groupdef.PartType == "MOUNT" {
			groupdef.Precedence = MNTGEN_PARTID
		} else if groupdef.PartType == "Reg" {
			groupdef.Precedence = REG_PARTID
		} else if groupdef.PartType == "Other" {
			groupdef.Precedence = DFLT_PARTID
		}

		ret = append(ret, groupdef)
	}
	return ret
}

func createGroupsList(partgroups []types.EBOMGroup) []types.EBOMGroup {

	// Passive Parts
	partgroups = AddGroup(partgroups, "C", "Capacitor", "Passive", "f")
	partgroups = AddGroup(partgroups, "CP", "CapacitorArray", "Passive", "f")
	partgroups = AddGroup(partgroups, "R", "Resistor", "Passive", "R")
	partgroups = AddGroup(partgroups, "R", "Resistor", "Passive", "K")
	partgroups = AddGroup(partgroups, "R", "Resistor", "Passive", "M")
	partgroups = AddGroup(partgroups, "R", "Resistor", "Passive", "ohm")
	partgroups = AddGroup(partgroups, "RP", "ResistorArray", "Passive", "R")
	partgroups = AddGroup(partgroups, "RP", "ResistorArray", "Passive", "K")
	partgroups = AddGroup(partgroups, "RP", "ResistorArray", "Passive", "M")
	partgroups = AddGroup(partgroups, "RP", "ResistorArray", "Passive", "ohm")
	partgroups = AddGroup(partgroups, "R", "Potentiometer", "Passive", "ohm")
	partgroups = AddGroup3(partgroups, "L", "Inductor", "Passive")
	partgroups = AddGroup3(partgroups, "FB", "FerritBead", "Passive")
	partgroups = AddGroup3(partgroups, "D", "Diode", "Passive")
	partgroups = AddGroup3(partgroups, "TVS", "TVS", "Passive")
	partgroups = AddGroup3(partgroups, "LED", "LED", "Passive")
	partgroups = AddGroup3(partgroups, "F", "Fuse", "Passive")
	// Transistors
	partgroups = AddGroup3(partgroups, "U", "FET", "Transistor")
	partgroups = AddGroup3(partgroups, "U", "Transistor", "General")
	partgroups = AddGroup3(partgroups, "U", "NPN", "Transistor")
	partgroups = AddGroup3(partgroups, "U", "PNP", "Transistor")
	// MCUs
	partgroups = AddGroup3(partgroups, "U", "atmega", "mcu")
	partgroups = AddGroup3(partgroups, "U", "pic", "mcu")
	// ICs
	partgroups = AddGroup3(partgroups, "U", "Device", "IC")
	partgroups = AddGroup3(partgroups, "IC", "Device", "IC")
	partgroups = AddGroup3(partgroups, "REG", "Reg", "IC")
	partgroups = AddGroup3(partgroups, "U", "Reg", "IC")
	partgroups = AddGroup3(partgroups, "U", "Transceiver", "IC")
	partgroups = AddGroup3(partgroups, "U", "Comparator", "IC")
	partgroups = AddGroup3(partgroups, "U", "ESD", "IC")
	// Mechanical
	partgroups = AddGroup3(partgroups, "SW", "Switch", "Mech")
	partgroups = AddGroup3(partgroups, "K", "Switch", "Mech")
	partgroups = AddGroup3(partgroups, "S", "Switch", "Mech")
	// Crystals
	partgroups = AddGroup3(partgroups, "Q", "Crystal", "Clock")
	partgroups = AddGroup3(partgroups, "X", "Crystal", "Clock")
	partgroups = AddGroup(partgroups, "X", "Crystal", "Clock", "HZ")
	partgroups = AddGroup3(partgroups, "Y", "Crystal", "Clock")
	partgroups = AddGroup3(partgroups, "OSC", "Oscillator", "Clock")

	// Mechanical
	partgroups = AddGroup3(partgroups, "J", "Connector", "Connector")
	partgroups = AddGroup3(partgroups, "J", "ConnRJ", "Connector")
	partgroups = AddGroup3(partgroups, "J", "ConnUSB", "Connector")
	partgroups = AddGroup3(partgroups, "JP", "Connector", "Connector")
	partgroups = AddGroup3(partgroups, "JP", "ConnRJ", "Connector")
	partgroups = AddGroup3(partgroups, "JP", "ConnUSB", "Connector")

	partgroups = AddGroup3(partgroups, "BT", "BATT", "BATT")

	partgroups = AddGroup3(partgroups, "T", "XFRM", "XFRM")

	partgroups = AddGroup3(partgroups, "*", "unkownPart", "unkownGroup")

	partgroups = AddGroup3(partgroups, "TP", "TestPoint", "MOUNT")
	partgroups = AddGroup3(partgroups, "*", "NET_SHORT", "MOUNT")

	partgroups = AddGroup(partgroups, "H", "MOUNT", "MOUNT", "HOLE")
	partgroups = AddGroup(partgroups, "TH", "MOUNT", "MOUNT", "HOLE")
	partgroups = AddGroup(partgroups, "HOLE", "MOUNT", "MOUNT", "HOLE")
	partgroups = AddGroup(partgroups, "P", "MOUNT", "MOUNT", "PAD")
	partgroups = AddGroup(partgroups, "*", "*", "DNP", "NF")

	// Set Precedence
	partgroups = SetPrecedence(partgroups)
	return partgroups
}

func componentGroupInit() {
	partgroups = createGroupsList(partgroups)

	log.Infof("Init partgroups.")

}

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
		if err != nil {
			return err
		}
	} else if strings.ToUpper(config.EDATool) == "PADSPCB" {
		bomParts, err = parser.ExtractPADSPCBComponents(config.InputFile)
		if err != nil {
			return err
		}
	} else {
		err = errors.Errorf("unknown tools %v", config.EDATool)
		return err
	}
	numberofparts := len(bomParts)
	log.Infof("numberofparts %d", numberofparts)

	combinedBOMparts, _ := combineBOMparts(bomParts)
	outputFilename := filepath.ToSlash(filepath.Join(config.OutputFile, prjname+"_BOM.csv"))
	log.Infof("CSV Output File %s", outputFilename)

	componentGroupInit()

	for k, ipart := range combinedBOMparts {
		combinedBOMparts[k].FValue = utils.GetFValFromEVal(combinedBOMparts[k].Value)
		combinedBOMparts[k].References = sortComponentRef(ipart)
		combinedBOMparts[k].SetComponentGroup(partgroups, true)
	}

	combinedBOMparts = sortComponentList(combinedBOMparts)

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
			if ((!strings.Contains("M", strings.ToUpper(cpart.Value)) && strings.ToUpper(cpart.Value) == strings.ToUpper(ipart.Value)) ||
				(strings.Contains("M", strings.ToUpper(cpart.Value)) && cpart.Value == ipart.Value)) &&
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

type RefItems struct {
	data  interface{}
	field string
}

func (items *RefItems) Len() int {
	if reflect.ValueOf(items.data).Kind() != reflect.Slice {
		return -1
	}
	return reflect.ValueOf(items.data).Len()
}

func (items *RefItems) Less(i, j int) bool {
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

func (items *RefItems) Swap(i, j int) {
	reflect.Swapper(items.data)(i, j)
}

func sortRefItems(i interface{}, str string) {
	if reflect.ValueOf(i).Kind() != reflect.Slice {
		return
	}
	a := &RefItems{
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

	sortRefItems(sortItems, "References")
	for _, v := range sortItems {
		//log.Infof("after sorted ref -- ", v)
		ret = append(ret, v.References)
	}
	return ret
}

type ComponentItems struct {
	data  interface{}
	field string
}

func (items *ComponentItems) Len() int {
	if reflect.ValueOf(items.data).Kind() != reflect.Slice {
		return -1
	}
	return reflect.ValueOf(items.data).Len()
}

func (items *ComponentItems) Less(i, j int) bool {
	a := reflect.ValueOf(items.data).Index(i)
	b := reflect.ValueOf(items.data).Index(j)
	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}
	if b.Kind() == reflect.Ptr {
		b = b.Elem()
	}
	_va := a.FieldByName(items.field).Interface()
	_vb := b.FieldByName(items.field).Interface()
	objA := _va.(types.EBOMItem)
	objB := _vb.(types.EBOMItem)
	log.Infof("a -- ", _va, _va.(types.EBOMItem))
	log.Infof("b -- ", _vb, _vb.(types.EBOMItem))

	groupA := objA.Group[1]
	groupB := objB.Group[1]
	log.Infof("a -- ", objA, groupA)
	log.Infof("b -- ", objB, groupB)

	fvalueA := objA.FValue
	fvalueB := objB.FValue

	return groupA < groupB || fvalueA < fvalueB
}

func (items *ComponentItems) Swap(i, j int) {
	reflect.Swapper(items.data)(i, j)
}

func sortComponentItems(i interface{}, str string) {
	if reflect.ValueOf(i).Kind() != reflect.Slice {
		return
	}
	a := &ComponentItems{
		data:  i,
		field: str,
	}
	sort.Sort(a)
}

func sortComponentList(self []types.EBOMItem) []types.EBOMItem {
	type sortItem struct {
		EBOMItem types.EBOMItem
	}
	var sortItems []sortItem
	var ret []types.EBOMItem

	for _, v := range self {
		sortItems = append(sortItems, sortItem{v})
	}

	sortComponentItems(sortItems, "EBOMItem")
	for _, v := range sortItems {
		//log.Infof("after sorted ref -- ", v)
		ret = append(ret, v.EBOMItem)
	}
	return ret
}
