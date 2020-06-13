package reliability

import (
	//"io"
	//"strconv"

	//"bytes"
	//"fmt"
	//"net/url"
	//"sort"
	//"strconv"
	//"strings"

	//"github.com/saycv/ebomgen/pkg/configuration"
	"github.com/pkg/errors"
	"github.com/saycv/ebomgen/pkg/types"

	log "github.com/sirupsen/logrus"
	//yaml "gopkg.in/yaml.v2"
)

// EBOMFrPart for reliability
type EBOMFrPart struct {
	Value         string
	FValue        float64
	Footprint     string
	Desc          string
	FrType        string
	ClsEnv        string
	CurrentStress string
	VoltageStress string
	PowerStress   string
}

type Setting func(frpart *EBOMFrPart)

func WithFrType(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.FrType = value
	}
}

func WithClsEnv(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.ClsEnv = value
	}
}

func WithCurrentStress(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.CurrentStress = value
	}
}

func WithVoltageStress(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.VoltageStress = value
	}
}

func WithPowerStress(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.PowerStress = value
	}
}

// NewFrPart init
func NewFrPart(part types.EBOMItem, settings ...Setting) *EBOMFrPart {
	new := &EBOMFrPart{}
	new.Value = part.Value
	new.FValue = part.FValue
	new.Footprint = part.Footprint
	new.Desc = part.Desc
	new.ClsEnv = "GB"
	new.CurrentStress = "0.5"
	new.VoltageStress = "0.5"
	new.PowerStress = "0.5"

	new.FrType = "RES-Film-Carbon"

	for _, set := range settings {
		set(new)
	}

	return new
}

func (b *EBOMFrPart) GetFailureRateBaseImported() (string, error) {
	tableData := FailureRateBaseImported
	reqValue, err := tableData[b.FrType].([]string)
	if !err {
		return "", errors.Errorf("%s not found in FailureRateBaseImported", b.FrType)
	}
	uriRefcnts := 0
	log.Info(b.FrType)
	log.Info(reqValue)
	for {
		if reqValue[0] == "URI" {
			reqValue, err = tableData[reqValue[1]].([]string)
			if !err {
				return "", errors.Errorf("%s not found in FailureRateBaseImported", reqValue[1])
			}
			if uriRefcnts > 1 {
				log.Errorf("ERR: URI Ref Counts too many!!!")
			}
			uriRefcnts = uriRefcnts + 1
		} else {
			break
		}
	}
	if reqValue[0] == "EXP" {
		reqValue = []string{"0"}
	}
	return reqValue[0], nil
}

func (b *EBOMFrPart) getFactorQualityImported() error {

	return nil
}
