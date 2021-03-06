package types

import (
	"io"
	"strconv"

	//"bytes"
	"fmt"
	//"net/url"
	//"sort"
	//"strconv"
	"strings"

	"github.com/prometheus/common/log"
	"github.com/saycv/ebomgen/pkg/configuration"
	//"github.com/pkg/errors"
	//"github.com/sirupsen/logrus"
	//log "github.com/sirupsen/logrus"
	//yaml "gopkg.in/yaml.v2"
)

// EBOMGroup part group
type EBOMGroup struct {
	Ref        string
	PartType   string
	GroupType  string
	Precedence int
	Unit       string
}

type ParamSourceType int32

const (
	ParamFromInit ParamSourceType = iota
	ParamFromOctopart
	ParamFromFindchips
	ParamFromDigikey
	ParamFromMouser
	ParamFromSzlcsc
	ParamFromCalc
)

type PartParameter struct {
	Value  string
	Source ParamSourceType
}

type EBOMWebPart struct {
	Category                PartParameter
	Subcategory             PartParameter
	MFR                     PartParameter
	MPN                     PartParameter
	Lifecycle               PartParameter
	RoHS                    PartParameter
	Datasheet               PartParameter
	PackageCase             PartParameter
	MoistureSensitive       PartParameter
	ReflowTemperaturePeak   PartParameter
	OperatingTemperatureMin PartParameter
	OperatingTemperatureMax PartParameter
	SupplyVoltageMin        PartParameter
	SupplyVoltageMax        PartParameter
	SupplyVoltageNom        PartParameter
	SupplyCurrentMin        PartParameter
	SupplyCurrentMax        PartParameter
	SupplyCurrentNom        PartParameter
	PowerDissipationMin     PartParameter
	PowerDissipationMax     PartParameter
	PowerDissipationNom     PartParameter

	UnitLength PartParameter
	UnitWidth  PartParameter
	UnitHeight PartParameter
	UnitWeight PartParameter

	UnitPrice PartParameter

	Attributes map[string]PartParameter
}

type EBOMItem struct {
	Quantity   int
	References []string
	Value      string
	FValue     float64
	Library    string
	Footprint  string
	Desc       string
	Attributes map[string]string
	Group      []string
	PartSpecs  EBOMWebPart
}

// EBOMSheet include all parts
type EBOMSheet struct {
	Headers       []string
	Items         []EBOMItem
	CustomHeaders []string
	Config        configuration.Configuration
}

func NewBOM(bomParts []EBOMItem, config configuration.Configuration) (*EBOMSheet, error) {
	res := &EBOMSheet{}
	res.Items = bomParts
	res.Config = config

	return res, nil
}

func (b *EBOMSheet) appendField(fieldName string) {
	for _, f := range b.Headers {
		if f == fieldName {
			return
		}
	}
	b.CustomHeaders = append(b.CustomHeaders, fieldName)
	b.Headers = append(b.Headers, fieldName)
}

func (b *EBOMSheet) generateHeaders() error {
	if b.Config.OnePartRows && strings.Contains(strings.ToUpper(b.Config.EDATool), "PCB") {
		b.Headers = []string{"Item", "References", "Quantity", "Value", "Footprint", "Description", "Rotation", "Layer"}
	} else if b.Config.Command == "bomcost" {
		b.Headers = []string{"Item", "References", "Quantity", "Value", "Footprint", "Description", "UnitPrice", "TotalPrice"}
	} else if b.Config.Command == "bommtbf" {
		b.Headers = []string{"Item", "References", "Quantity", "Value", "Footprint", "Description", "FrUnit", "FrTot"}
	} else {
		b.Headers = []string{"Item", "References", "Quantity", "Value", "Footprint", "Description"}
	}
	return nil
}

func (b *EBOMSheet) makeUniqueIdentifier(comp EBOMItem) string {
	ident := fmt.Sprintf("Value=%s_Footprint=%s", comp.Value, comp.Footprint)

	return ident
}

func (b *EBOMSheet) writeItem(w io.Writer, k int, i EBOMItem) error {
	res := make([]string, 0, len(i.Attributes)+5)
	res = append(res, fmt.Sprintf("%d", k+1))
	res = append(res, fmt.Sprintf(`"%s"`, strings.Join(i.References, ",")))
	res = append(res, fmt.Sprintf("%d", i.Quantity))
	res = append(res, `"`+i.Value+`"`)
	res = append(res, `"`+i.Footprint+`"`)
	res = append(res, `"`+i.Attributes["Description"]+`"`)
	if false {
		res = append(res, `"`+i.Group[0]+`"`)
		res = append(res, `"`+i.Group[1]+`"`)
	}

	if b.Config.Command == "bomcost" {
		res = append(res, `"`+i.Attributes["UnitPrice"]+`"`)
		str := fmt.Sprintf("=C%d*G%d", k+2, k+2)
		res = append(res, `"`+str+`"`)
	} else if b.Config.Command == "bommtbf" {
		res = append(res, `"`+i.Attributes["FrUnit"]+`"`)
		str := fmt.Sprintf("=C%d*G%d", k+2, k+2)
		res = append(res, `"`+str+`"`)
	}

	rotate, ok := i.Attributes["rotate"]
	if ok && rotate != "" && b.Config.OnePartRows && strings.Contains(strings.ToUpper(b.Config.EDATool), "PCB") {
		res = append(res, `"`+rotate+`"`)
	}
	layer, ok := i.Attributes["layer"]
	if ok && layer != "" && b.Config.OnePartRows && strings.Contains(strings.ToUpper(b.Config.EDATool), "PCB") {
		res = append(res, `"`+layer+`"`)
	}

	_, err := fmt.Fprintln(w, strings.Join(res, ","))

	return err
}

// WriteCSV saveas csv file
func (b *EBOMSheet) WriteCSV(w io.Writer) error {

	b.generateHeaders()
	_, err := fmt.Fprintln(w, strings.Join(b.Headers, ","))
	if err != nil {
		return err
	}

	indexCnt := 0
	for k, i := range b.Items {
		indexCnt = k
		err = b.writeItem(w, k, i)
		if err != nil {
			return err
		}
	}

	if b.Config.Command == "bomcost" {
		totalBomPrice := fmt.Sprintf("=sum(H%d:H%d)", 2, indexCnt+2)

		totalBomPriceItem := []string{"", "", "", "", "", "", "Total BOMCOST", "", "USD"}
		totalBomPriceItem[7] = totalBomPrice
		_, err := fmt.Fprintln(w, strings.Join(totalBomPriceItem, ","))
		if err != nil {
			return err
		}
		totalBomPriceItem = []string{"", "", "", "", "", "", "Total BOMCOST", "", "CNY"}
		usd2rmb := USD2CNY
		totalBomPriceItem[7] = fmt.Sprintf("=sum(H%d*%f)", indexCnt+3, usd2rmb)
		_, err = fmt.Fprintln(w, strings.Join(totalBomPriceItem, ","))
		if err != nil {
			return err
		}
	} else if b.Config.Command == "bommtbf" {
		totalBomPrice := fmt.Sprintf("=sum(H%d:H%d)", 2, indexCnt+2)

		totalBomPriceItem := []string{"", "", "", "", "", "", "Total BOMFR", "", ""}
		totalBomPriceItem[7] = totalBomPrice
		_, err := fmt.Fprintln(w, strings.Join(totalBomPriceItem, ","))
		if err != nil {
			return err
		}
		totalBomPriceItem = []string{"", "", "", "", "", "", "Total BOMMTBF", "", "h"}
		totalBomPriceItem[7] = fmt.Sprintf("=10^6/(H%d)", indexCnt+3)
		_, err = fmt.Fprintln(w, strings.Join(totalBomPriceItem, ","))
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *EBOMItem) SetComponentGroup(partgroups []EBOMGroup, verbose_p bool) {
	current_confidence := 0
	highest_confidence := 0
	//ThreshHoldMatch := []int
	var MostLikelyMatches []EBOMGroup
	confidence_threshold := 2
	//MostLikelyMatch := [2]string{"Other", "1000"}
	//self := b.Group

	// Print if verbose is true
	if verbose_p {
		log.Infof("-")
	}

	for index, definition := range partgroups {
		// Reset per definiton comparison variables
		current_confidence = 0
		// Build Confidence in a match
		if strings.Contains(strings.ToUpper(self.References[0]), strings.ToUpper(definition.Ref)) {
			current_confidence += 1
			if verbose_p {
				log.Infof("{%s} matched in {%s} at {%d}",
					definition.Ref, self.References[0], index)
			}
		}
		if strings.Contains(strings.ToUpper(self.Library), strings.ToUpper(definition.PartType)) ||
			strings.Contains(strings.ToUpper(self.Footprint), strings.ToUpper(definition.PartType)) {
			current_confidence += 1
			if verbose_p {
				log.Infof("{%s} matched in {%s} or {%s} at {%d}",
					definition.PartType, self.Library, self.Footprint,
					index)
			}
		}
		if strings.Contains(strings.ToUpper(self.Library), strings.ToUpper(definition.GroupType)) ||
			strings.Contains(strings.ToUpper(self.Footprint), strings.ToUpper(definition.GroupType)) {
			current_confidence += 1
			if verbose_p {
				log.Infof("{%s} matched in {%s} or {%s} at {%s}",
					definition.GroupType, self.Library, self.Footprint,
					index)
			}
		}
		if strings.Contains(strings.ToUpper(self.Value), strings.ToUpper(definition.Unit)) {
			current_confidence += 1
			if verbose_p {
				log.Infof("{%s} mached in {%s}", definition.Unit,
					self.Value)
			}
		}

		if true && (self.Attributes["part"] == "CapacitorTan") {
			if definition.PartType == "CapacitorTan" {
				current_confidence += 100
				if verbose_p {
					log.Infof("{%s} mached in {%s}",
						definition.Unit, self.Value)
				}
			} else {
				current_confidence = 0
			}
		}
		if true && (self.Attributes["part"] == "IC" ||
			self.Attributes["part"] == "Reg") {
			if definition.GroupType == "IC" {
				current_confidence += 100
				if verbose_p {
					log.Infof("{%s} mached in {%s}",
						definition.Unit, self.Value)
				}
			} else {
				current_confidence = 0
			}
		}
		if true && self.Attributes["part"] == "DNP" {
			if definition.GroupType == "DNP" {
				current_confidence += 100
				if verbose_p {
					log.Infof("{%s} mached in {%s}",
						definition.Unit, self.Value)
				}
			} else {
				current_confidence = 0
			}
		}
		if true && self.Attributes["part"] == "TestPoint" {
			if definition.PartType == "TestPoint" {
				current_confidence += 100
				if verbose_p {
					log.Infof("{%s} mached in {%s}",
						definition.Unit, self.Value)
				}
			} else {
				current_confidence = 0
			}
		}
		if true && strings.Contains(self.Attributes["part"], "Connector") {
			if strings.Contains(definition.PartType, "Connector") {
				current_confidence += 100
				if verbose_p {
					log.Infof("{%s} mached in {%s}",
						definition.Unit, self.Value)
				}
			} else {
				current_confidence = 0
			}
		}
		if true && strings.Contains(self.Attributes["part"], "unkownPart") {
			if strings.Contains(definition.PartType, "unkownPart") {
				current_confidence += 100
				if verbose_p {
					log.Infof("{%s} mached in {%s}",
						definition.Unit, self.Value)
				}
			} else {
				current_confidence = 0
			}
		}
		if true && self.Attributes["part"] == definition.PartType {
			current_confidence += 100
			if verbose_p {
				log.Infof("{%s} mached in {%s}", definition.Unit,
					self.Value)
			}
		}
		if current_confidence > highest_confidence {
			highest_confidence = current_confidence
			self.Group = []string{definition.GroupType + " " + definition.PartType, strconv.Itoa(definition.Precedence)}

		}
		if current_confidence >= confidence_threshold {
			MostLikelyMatches = append(MostLikelyMatches, definition)
			if verbose_p {
				log.Infof("Threshold Met")
			}
		}
	}

	if verbose_p {
		log.Infof("Part:")
		log.Infof("    Ref:{%s}", self.References[0])
		log.Infof("    Value:{%s}", self.Value)
		log.Infof("    Footprint:{%s}", self.Footprint)
		log.Infof("    Library: {%s}", self.Library)
		log.Infof("    Highest Confidence: {%d}", highest_confidence)
		log.Infof("Most Likely Groups:")
		if len(MostLikelyMatches) == 0 {
			log.Infof("No Match")
		} else {
			for _, match := range MostLikelyMatches {
				log.Infof(match.Ref, match.PartType, match.GroupType)
			}
		}
	}
}
