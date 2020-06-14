package utils

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// "net/url"
// "sort"
//"github.com/saycv/ebomgen/pkg/types"

var (
	DEFAULT_TRANS_COUNT = map[string]interface{}{
		"GEN-50":   50.0,
		"GEN-100":  100.0,
		"GEN-500":  500.0,
		"GEN-1000": 1000.0,

		"GEN1-4004":     2300.0,
		"GEN1-8008":     3500.0,
		"GEN2-8080":     6000.0,
		"GEN2-8085":     9000.0,
		"GEN3-8086":     0.029 * 1e+06,
		"GEN3-80286":    0.134 * 1e+06,
		"GEN4-80386":    0.275 * 1e+06,
		"GEN4-80486":    1.25 * 1e+06,
		"GEN4-80486BL2": 1.4 * 1e+06,
		"GEN4-80486BL3": 1.4 * 1e+06,
		"GEN5-Pentium1": 3.1 * 1e+06,
		"GEN5-K6MMX":    4.3 * 1e+06,
		"GEN5-Pentium2": 7.5 * 1e+06,
		"GEN5-Pentium3": 28 * 1e+06,
		"GEN5-K7Athlon": 22 * 1e+06,
		"GEN5-Pentium4": 55 * 1e+06,

		"PL-Virtex":       70 * 1e+06,
		"PL-Virtex-E":     200 * 1e+06,
		"PL-Virtex-2":     350 * 1e+06,
		"PL-Virtex-2-pro": 430 * 1e+06,
		"PL-Virtex-4":     1000 * 1e+06,
		"PL-Virtex-5":     1100 * 1e+06,

		"PL-Stratix-4": 2500 * 1e+06,
		"PL-Stratix-5": 3800 * 1e+06,
		"PL-Arria-10":  5300 * 1e+06,
		"PL-Artix-7":   6800 * 1e+06,

		"PL-Virtex-7":   6800 * 1e+06,
		"PL-Stratix-10": 17000 * 1e+06,
		"PL-Virtex-us":  20000 * 1e+06,
		"PL-Everest":    50000 * 1e+06,

		"XC7Z015": []string{"URI", "PL-Artix-7"},

		"SRAM-1GB": 8 * 0.6 * 1e+09,
		"DDR3-1GB": []string{"URI", "SRAM-1GB"},

		"FLASH-1MB":   8 * 0.6 * 1e+06,
		"FLASH-128MB": 128 * 8 * 0.6 * 1e+06,
		"FLASH-1GB":   8 * 0.6 * 1e+09,
	}
)

func GetTocOfIc(value string, desc string, fp string) (int, error) {
	partType := "DIC"
	partICType := "GEN-50"
	tableData := DEFAULT_TRANS_COUNT

	fp = strings.ToUpper(fp)

	pins, _ := GetPinsFromFp(desc, fp)

	log.Infof("processing %s, %s, %d", value, fp, pins)

	if true {
		if pins >= 1000 {
			partICType = "PL-Everest"
		} else if pins >= 800 {
			partICType = "PL-Virtex-us"
		} else if pins >= 450 {
			partICType = "PL-Artix-7"
		} else if pins >= 200 {
			partICType = "PL-Stratix-4"
		} else if pins >= 100 {
			partICType = "PL-Virtex"
		} else if pins >= 50 {
			partICType = "GEN4-80486"
		} else if pins >= 45 {
			partICType = "GEN4-80386"
		} else if pins >= 40 {
			partICType = "GEN3-80286"
		} else if pins >= 35 {
			partICType = "GEN3-8086"
		} else if pins >= 30 {
			partICType = "GEN2-8085"
		} else if pins >= 25 {
			partICType = "GEN2-8080"
		} else if pins >= 20 {
			partICType = "GEN1-4004"
		}

	}

	if strings.Contains(fp, "QFP") {
		if pins >= 1000 {
			partICType = "PL-Everest"
		} else if pins >= 800 {
			partICType = "PL-Virtex-us"
		} else if pins >= 450 {
			partICType = "PL-Artix-7"
		} else if pins >= 200 {
			partICType = "PL-Stratix-4"
		} else if pins >= 100 {
			partICType = "PL-Virtex"
		} else if pins >= 50 {
			partICType = "GEN4-80486"
		}
		partType = "DIC"
	} else if strings.Contains(fp, "BGA") {
		if pins >= 1000 {
			partICType = "PL-Everest"
		} else if pins >= 800 {
			partICType = "PL-Virtex-us"
		} else if pins >= 450 {
			partICType = "PL-Artix-7"
		} else if pins >= 200 {
			partICType = "PL-Stratix-4"
		} else if pins >= 100 {
			partICType = "PL-Virtex"
		} else if pins >= 50 {
			partICType = "GEN4-80486"
		}
		partType = "MPU"
	}
	if strings.Contains(value, "256M16") {
		partICType = "DDR3-1GB"
		partType = "DRAM"
	} else if strings.Contains(value, "128MB") && strings.Contains(fp, "TSOP48") {
		partICType = "FLASH-128MB"
		partType = "ROM"
	}

	reqValueFloat := float64(0.0)
	reqValueCls, ok := tableData[partICType]
	if !ok {
		return 0, errors.Errorf("%s not found in %v", partType, reflect.TypeOf(tableData))
	}

	reqValue, ok := reqValueCls.([]string)
	if !ok {
		reqValueFloat = reqValueCls.(float64)
	} else {

		uriRefcnts := 0
		for {
			if reqValue[0] == "URI" {
				reqValueCls, ok = tableData[reqValue[1]]
				if !ok {
					return 0, errors.Errorf("%s not found in %v, %v", reqValue[1], reflect.ValueOf(tableData), reflect.TypeOf(tableData))
				}
				reqValue, ok := reqValueCls.([]string)
				if !ok {
					reqValueFloat = reqValueCls.(float64)
					break
				}

				if uriRefcnts > 1 {
					log.Errorf("ERR { URI Ref Counts too many!!! { %s, %s", reqValue[0], reqValue[1])
				}
				uriRefcnts = uriRefcnts + 1
			} else {
				break
			}
		}
		if reqValue[0] == "EXP" {
			reqValue = []string{"0"}
			reqValueFloat = 0.0
		}
	}

	toc := reqValueFloat
	log.Info(toc)

	return int(toc), nil
}
