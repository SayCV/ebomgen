package reliability

import (
	//"io"
	"reflect"

	//"bytes"
	//"fmt"
	//"net/url"
	//"sort"
	"strconv"
	"strings"

	//"github.com/saycv/ebomgen/pkg/configuration"
	"github.com/pkg/errors"
	"github.com/saycv/ebomgen/pkg/types"
	"github.com/saycv/ebomgen/pkg/utils"

	log "github.com/sirupsen/logrus"
	//yaml "gopkg.in/yaml.v2"
)

// EBOMFrPart for reliability
type EBOMFrPart struct {
	Quantity      int
	References    []string
	Value         string
	FValue        float64
	Footprint     string
	Desc          string
	FrType        string
	FrProcess     string
	FrFp          string
	FrSealed      bool
	FrSMD         bool
	FrPins        int
	ClsEnv        string
	ClsQuality    string
	OperatingTemp string
	CurrentStress string
	VoltageStress string
	PowerStress   string

	FrUnit string
}

type Setting func(frpart *EBOMFrPart)

func WithFrType(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.FrType = value
	}
}

func WithFrProcess(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.FrProcess = value
	}
}

func WithFrFp(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.FrFp = value
	}
}

func WithFrSealed(value bool) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.FrSealed = value
	}
}

func WithFrSMD(value bool) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.FrSMD = value
	}
}

func WithFrPins(value int) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.FrPins = value
	}
}

func WithClsEnv(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.ClsEnv = value
	}
}

func WithClsQuality(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.ClsQuality = value
	}
}

func WithOperatingTemp(value string) Setting {
	return func(frpart *EBOMFrPart) {
		frpart.OperatingTemp = value
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
	new.Quantity = part.Quantity
	new.References = part.References
	new.Value = part.Value
	new.FValue = part.FValue
	new.Footprint = part.Footprint
	new.Desc = part.Desc

	new.FrType = "RES-Film-Carbon"
	new.FrType = ""
	new.FrProcess = ""
	new.FrFp = ""
	new.FrSealed = true
	new.FrSMD = true
	new.FrPins = 2
	new.ClsEnv = "GB"
	new.ClsQuality = ""
	new.CurrentStress = "0.5"
	new.VoltageStress = "0.5"
	new.PowerStress = "0.5"

	new.FrUnit = ""

	new.FrPins, _ = utils.GetPinsFromFp(new.Desc, new.Footprint)

	for _, set := range settings {
		set(new)
	}

	return new
}

func (b *EBOMFrPart) GetFailureRateActiveContactImported(contactnbr int) (string, error) {

	if contactnbr <= 2 {
		return "0.0025", nil
	} else if contactnbr <= 3 {
		return "0.00375", nil
	} else if contactnbr <= 4 {
		return "0.005", nil
	} else if contactnbr <= 6 {
		return "0.0075", nil
	} else if contactnbr <= 8 {
		return "0.01", nil
	} else if contactnbr <= 9 {
		return "0.0113", nil
	} else if contactnbr <= 10 {
		return "0.0125", nil
	} else if contactnbr <= 12 {
		return "0.015", nil
	} else if contactnbr <= 18 {
		return "0.0225", nil
	} else {
		baseval := 0.00125
		val := baseval * float64(contactnbr)
		strval := strconv.FormatFloat(val, 'f', -1, 64)
		return strval, nil
	}
}

func (b *EBOMFrPart) GetFailureRateBaseImported() (string, error) {
	partType := b.FrType
	//partProcess := b.FrProcess
	tableData := FailureRateBaseImported
	reqValue, ok := tableData[partType].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", partType, reflect.TypeOf(tableData))
	}
	log.Info(partType)
	log.Info(reqValue)

	uriRefcnts := 0
	for {
		if reqValue[0] == "URI" {
			reqValue, ok = tableData[reqValue[1]].([]string)
			if !ok {
				return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
			}
			if uriRefcnts > 1 {
				log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
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

func (b *EBOMFrPart) GetFactorEnvImported() (string, error) {
	partType := b.FrType
	partProcess := b.FrProcess
	tableData := FactorEnvImported
	queryEnvIndexStr, ok := ClsEnv[b.ClsEnv].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.ClsEnv, reflect.TypeOf(ClsEnv))
	}
	queryEnvIndex, _ := strconv.Atoi(queryEnvIndexStr[0])

	if partType == "AIC" || partType == "FIFO" || partType == "DRAM" {
		partType = partType
	} else if partProcess == "Bipolar" || partProcess == "MOS" {
		partType = partType + "-" + partProcess
	}

	if strings.HasPrefix(partType, "CAP") {
		partType = "CAP"
	}

	reqValue, ok := tableData[b.FrType].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.FrType, reflect.TypeOf(tableData))
	}

	uriRefcnts := 0
	for {
		if reqValue[0] == "URI" {
			reqValue, ok = tableData[reqValue[1]].([]string)
			if !ok {
				return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
			}
			if uriRefcnts > 1 {
				log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
			}
			uriRefcnts = uriRefcnts + 1
		} else {
			break
		}
	}
	if reqValue[0] == "EXP" {
		reqValue = []string{"0"}
	}

	return reqValue[queryEnvIndex], nil
}

func (b *EBOMFrPart) GetFactorQualityImported() (string, error) {
	partType := b.FrType
	partProcess := b.FrProcess
	tableData := FactorQualityImported
	queryEnvIndexStr, ok := ClsQuality[b.ClsQuality].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.ClsQuality, reflect.TypeOf(ClsQuality))
	}
	queryEnvIndex, _ := strconv.Atoi(queryEnvIndexStr[0])

	if partType == "AIC" || partType == "FIFO" || partType == "DRAM" {
		partType = partType
	} else if partProcess == "Bipolar" || partProcess == "MOS" {
		partType = partType + "-" + partProcess
	}

	if strings.HasPrefix(partType, "CAP") {
		partType = "CAP"
	}

	reqValue, ok := tableData[b.FrType].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.FrType, reflect.TypeOf(tableData))
	}

	uriRefcnts := 0
	for {
		if reqValue[0] == "URI" {
			reqValue, ok = tableData[reqValue[1]].([]string)
			if !ok {
				return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
			}
			if uriRefcnts > 1 {
				log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
			}
			uriRefcnts = uriRefcnts + 1
		} else {
			break
		}
	}
	if reqValue[0] == "EXP" {
		reqValue = []string{"0"}
	}

	return reqValue[queryEnvIndex], nil
}

func (b *EBOMFrPart) GetFactorTemperatureImported() (string, error) {
	partType := b.FrType
	partProcess := b.FrProcess
	tableData := FactorTemperatureImported
	//queryEnvIndexStr, ok := ClsEnv[b.ClsEnv].([]string)
	//if !ok {
	//	return "", errors.Errorf("%s not found in %v", b.ClsEnv, reflect.TypeOf(ClsEnv))
	//}
	//queryEnvIndex, _ := strconv.Atoi(queryEnvIndexStr[0])

	partTempType := "RES-TEMP-CLASS"
	if strings.HasPrefix(partType, "AIC") || strings.HasPrefix(partType, "DIC") ||
		strings.HasPrefix(partType, "RAM") || strings.HasPrefix(partType, "ROM") ||
		strings.HasPrefix(partType, "FLASH") || strings.HasPrefix(partType, "MPU") ||
		strings.HasPrefix(partType, "DRAM") || strings.HasPrefix(partType, "GaAsMMIC") {
		partTempType = "IC-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "CAP") {
		partTempType = "CAP-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "IND") {
		partTempType = "IND-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "Diode") {
		partTempType = "DIODE-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "LED") {
		partTempType = "LED-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "NPN") || strings.HasPrefix(partType, "PNP") {
		partTempType = "BJT-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "XTAL") || strings.HasPrefix(partType, "OSC") {
		partTempType = "XTAL-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "OptoElectronicDevices") {
		partTempType = "OptoElectronicDevices-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "XFMR") {
		partTempType = "XFMR-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "Relay") {
		partTempType = "Relay-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "Switch") {
		partTempType = "Switch-TEMP-CLASS"
	} else if strings.HasPrefix(partType, "CONN") {
		partTempType = "CONN-TEMP-CLASS"
	}

	reqTempValue, ok := tableData[partTempType].([]int)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.FrType, reflect.TypeOf(tableData))
	}

	//uriRefcnts := 0
	//for {
	//	if reqValue[0] == "URI" {
	//		reqValue, ok = tableData[reqValue[1]].([]string)
	//		if !ok {
	//			return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
	//		}
	//		if uriRefcnts > 1 {
	//			log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
	//		}
	//		uriRefcnts = uriRefcnts + 1
	//	} else {
	//		break
	//	}
	//}
	//if reqValue[0] == "EXP" {
	//	reqValue = []string{"0"}
	//}

	//log.Info(partTempType)
	//log.Info(reqTempValue)

	// OperatingTemp support "GB" or "30"
	queryTempStr := b.OperatingTemp
	queryTemp, err := strconv.Atoi(queryTempStr)
	if err != nil {
		queryEnvIndexStr, ok := ClsEnv[queryTempStr].([]string)
		if !ok {
			return "", errors.Errorf("%s not found in %v", queryTempStr, reflect.TypeOf(ClsEnv))
		}
		queryTemp, err = strconv.Atoi(queryEnvIndexStr[3])
	}

	queryTempIndex := 0
	for k, temp := range reqTempValue {
		if queryTemp <= temp {
			queryTempIndex = k
			break
		}
	}
	log.Info("queryTempIndex: ", queryTempIndex)

	if partType == "AIC" || partType == "FIFO" || partType == "DRAM" {
		partType = partType
	} else if partProcess == "Bipolar" || partProcess == "MOS" {
		partType = partType + "-" + partProcess
	}

	if strings.HasPrefix(partType, "CAP") {
		partType = "CAP"
	}

	reqValue, ok := tableData[b.FrType].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.FrType, reflect.TypeOf(tableData))
	}

	uriRefcnts := 0
	for {
		if reqValue[0] == "URI" {
			reqValue, ok = tableData[reqValue[1]].([]string)
			if !ok {
				return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
			}
			if uriRefcnts > 1 {
				log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
			}
			uriRefcnts = uriRefcnts + 1
		} else {
			break
		}
	}
	if reqValue[0] == "EXP" {
		reqValue = []string{"0"}
	}

	return reqValue[queryTempIndex], nil
}

func (b *EBOMFrPart) GetFactorStressImported() (string, error) {
	partType := b.FrType
	//partProcess := b.FrProcess
	tableData := FactorStressImported
	//queryEnvIndexStr, ok := ClsEnv[b.ClsEnv].([]string)
	//if !ok {
	//	return "", errors.Errorf("%s not found in %v", b.ClsEnv, reflect.TypeOf(ClsEnv))
	//}
	//queryEnvIndex, _ := strconv.Atoi(queryEnvIndexStr[0])

	reqValue, ok := tableData[partType].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.FrType, reflect.TypeOf(tableData))
	}

	spreqValue := make(map[string]interface{})
	special_part := ""

	uriRefcnts := 0
	foundValue := false
	for {
		if reqValue[0] == "URI" && !foundValue {
			special_part = reqValue[1]
			if special_part == "BJT" {
				spreqValue, ok = tableData[reqValue[1]].(map[string]interface{})
				if !ok {
					return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
				}
				foundValue = true
			} else if special_part == "Diode" {
				spreqValue, ok = tableData[reqValue[1]].(map[string]interface{})
				if !ok {
					return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
				}
				foundValue = true
			} else {
				reqValue, ok = tableData[reqValue[1]].([]string)
				if !ok {
					return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
				}
			}
			if uriRefcnts > 1 {
				log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
			}
			uriRefcnts = uriRefcnts + 1
		} else {
			break
		}
	}
	if reqValue[0] == "EXP" {
		reqValue = []string{"0"}
	}

	result := ""
	if special_part == "BJT" {
		__svStr := b.VoltageStress
		//__siStr := b.CurrentStress
		__spStr := b.PowerStress
		__spFloat, _ := strconv.ParseFloat(__spStr, 64)
		reqValue, ok = spreqValue[__svStr].([]string)
		if !ok {
			return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
		}
		result = reqValue[int(__spFloat*10)-1]
	} else if special_part == "Diode" {
		//__svStr := b.VoltageStress
		__siStr := b.CurrentStress
		__spStr := b.PowerStress
		__spFloat, _ := strconv.ParseFloat(__spStr, 64)
		reqValue, ok = spreqValue[__siStr].([]string)
		if !ok {
			return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
		}
		result = reqValue[int(__spFloat*10)-1]
	} else {
		__svStr := b.VoltageStress
		__svFloat, _ := strconv.ParseFloat(__svStr, 64)
		result = reqValue[int(__svFloat*10)-1]
	}
	return result, nil
}

func (b *EBOMFrPart) GetFactorChImported() (string, error) {
	chip := b.FrSMD
	partType := b.FrType
	//partProcess := b.FrProcess
	tableData := FactorChImported
	//queryEnvIndexStr, ok := ClsEnv[b.ClsEnv].([]string)
	//if !ok {
	//	return "", errors.Errorf("%s not found in %v", b.ClsEnv, reflect.TypeOf(ClsEnv))
	//}
	//queryEnvIndex, _ := strconv.Atoi(queryEnvIndexStr[0])

	reqValue, ok := tableData[partType].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.FrType, reflect.TypeOf(tableData))
	}

	uriRefcnts := 0
	for {
		if reqValue[0] == "URI" {
			reqValue, ok = tableData[reqValue[1]].([]string)
			if !ok {
				return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
			}
			if uriRefcnts > 1 {
				log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
			}
			uriRefcnts = uriRefcnts + 1
		} else {
			break
		}
	}
	if reqValue[0] == "EXP" {
		reqValue = []string{"0"}
	}

	if chip {
		return reqValue[0], nil
	} else {
		return reqValue[1], nil
	}
}

func (b *EBOMFrPart) GetFactorProcessImported() (string, error) {
	pins := b.FrPins
	partType := b.FrType
	//partProcess := b.FrProcess
	tableData := FactorProcessImported
	//queryEnvIndexStr, ok := ClsEnv[b.ClsEnv].([]string)
	//if !ok {
	//	return "", errors.Errorf("%s not found in %v", b.ClsEnv, reflect.TypeOf(ClsEnv))
	//}
	//queryEnvIndex, _ := strconv.Atoi(queryEnvIndexStr[0])

	reqValue, ok := tableData[partType].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.FrType, reflect.TypeOf(tableData))
	}

	uriRefcnts := 0
	for {
		if reqValue[0] == "URI" {
			reqValue, ok = tableData[reqValue[1]].([]string)
			if !ok {
				return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
			}
			if uriRefcnts > 1 {
				log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
			}
			uriRefcnts = uriRefcnts + 1
		} else {
			break
		}
	}
	if reqValue[0] == "EXP" {
		reqValue = []string{"0"}
	}

	if partType == "Connector" {
		if pins <= 1 {
			pins = 1
		}
		if pins <= 20 {
			return reqValue[pins-1], nil
		} else if pins > 100 {
			return reqValue[len(reqValue)-1], nil
		} else {
			index := (pins - 20) / 5
			return reqValue[19+index], nil
		}
	} else {
		return reqValue[0], nil
	}
}

func (b *EBOMFrPart) GetFactorApplicationImported() (string, error) {
	partType := b.FrType
	//partProcess := b.FrProcess
	tableData := FactorApplicationImported
	//queryEnvIndexStr, ok := ClsEnv[b.ClsEnv].([]string)
	//if !ok {
	//	return "", errors.Errorf("%s not found in %v", b.ClsEnv, reflect.TypeOf(ClsEnv))
	//}
	//queryEnvIndex, _ := strconv.Atoi(queryEnvIndexStr[0])

	reqValue, ok := tableData[partType].([]string)
	if !ok {
		return "", errors.Errorf("%s not found in %v", b.FrType, reflect.TypeOf(tableData))
	}

	uriRefcnts := 0
	for {
		if reqValue[0] == "URI" {
			reqValue, ok = tableData[reqValue[1]].([]string)
			if !ok {
				return "", errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
			}
			if uriRefcnts > 1 {
				log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
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

func (b *EBOMFrPart) GetFactorC1Imported() (string, error) {
	partType := b.FrType
	partProcess := b.FrProcess
	tableData := FactorC1Imported
	//queryEnvIndexStr, ok := ClsEnv[b.ClsEnv].([]string)
	//if !ok {
	//	return "", errors.Errorf("%s not found in %v", b.ClsEnv, reflect.TypeOf(ClsEnv))
	//}
	//queryEnvIndex, _ := strconv.Atoi(queryEnvIndexStr[0])

	if partType == "AIC" || partType == "FIFO" || partType == "DRAM" {
		partType = partType
	} else if partProcess == "Bipolar" || partProcess == "MOS" {
		partType = partType + "-" + partProcess
	}

	if strings.HasPrefix(partType, "CAP") {
		partType = "CAP"
	}

	//log.Info(partType)
	reqValueMap := make(map[string]interface{})
	reqValue, ok := tableData[partType].([]string)
	if ok {
		uriRefcnts := 0
		for {
			if reqValue[0] == "URI" {
				partType = reqValue[1]
				_reqValue, ok := tableData[partType].([]string)
				if ok {
					if uriRefcnts > 1 {
						log.Errorf("ERR: URI Ref Counts too many!!!: %s, %s", reqValue[0], reqValue[1])
					}
					uriRefcnts = uriRefcnts + 1
					reqValue = _reqValue
					continue
				}
				break
			} else {
				break
			}
		}
		if reqValue[0] == "EXP" {
			reqValue = []string{"0"}
		}
	}
	//log.Info(partType)
	reqValueMap, ok = tableData[partType].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("%s not found in %v", reqValue[1], reflect.TypeOf(tableData))
	}

	// transistors on chip
	toc, _ := utils.GetTocOfIc(b.Value, b.Desc, b.Footprint)

	// 37 lines
	//totlines := 36
	query_c1_index := 0
	for k, _ := range reqValueMap {
		if reqValueMap[k] == "-     " {
			continue
		}
		tempInt, _ := strconv.Atoi(k)
		//log.Info(tempInt)
		if toc >= tempInt && query_c1_index < tempInt {
			query_c1_index = tempInt
			//log.Info("update")
			continue
		}
	}
	query_c1_indexStr := strconv.Itoa(query_c1_index)

	//log.Info(toc)
	//log.Info(query_c1_index)
	return reqValueMap[query_c1_indexStr].([]string)[0], nil
}

func (b *EBOMFrPart) GetFactorC2Imported() (string, error) {
	fp_map := map[string]int{
		// "sealed" ： DIP, FP, GEN, LCC, Metal "|||||" "unsealed" ： DIP, FP, GEN, LCC
		"sealed-DIP": 0, "sealed-FP": 1, "sealed-GEN": 2, "sealed-LCC": 3, "sealed-Metal": 4,
		"UNDEFINE":     5,
		"unsealed-DIP": 6, "unsealed-FP": 7, "unsealed-GEN": 8, "unsealed-LCC": 9,
	}
	pins := b.FrPins
	sealed := b.FrSealed
	partType := b.FrType
	partProcess := b.FrProcess
	tableData := FactorC2Imported
	//queryEnvIndexStr, ok := ClsEnv[b.ClsEnv].([]string)
	//if !ok {
	//	return "", errors.Errorf("%s not found in %v", b.ClsEnv, reflect.TypeOf(ClsEnv))
	//}
	//queryEnvIndex, _ := strconv.Atoi(queryEnvIndexStr[0])

	if partType == "AIC" || partType == "FIFO" || partType == "DRAM" {
		partType = partType
	} else if partProcess == "Bipolar" || partProcess == "MOS" {
		partType = partType + "-" + partProcess
	}

	if strings.HasPrefix(partType, "CAP") {
		partType = "CAP"
	}

	partFp := "GEN"
	fp_map_key := ""
	if sealed {
		fp_map_key = "sealed" + "-" + partFp
	} else {
		fp_map_key = "unsealed" + "-" + partFp
	}

	//log.Infof("C2/pins: %d", pins)
	pinsAvailabe := pins
	for k, _ := range tableData {
		valInt, _ := strconv.Atoi(k)
		if pins >= valInt && pinsAvailabe < valInt {
			pinsAvailabe = valInt
			continue
		}
	}
	pinsAvailabeStr := strconv.Itoa(pinsAvailabe)
	//log.Infof("C2/pinsAvailabe: %s", pinsAvailabeStr)

	reqValue, ok := tableData[pinsAvailabeStr].([]string)
	if !ok {
		return "", errors.Errorf("ERR1 : %d not found in %v", pinsAvailabe, reflect.TypeOf(tableData))
	}

	return reqValue[fp_map[fp_map_key]], nil
}
