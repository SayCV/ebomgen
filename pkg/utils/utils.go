package utils

import (
	//"io"
	//"fmt"
	//"math"
	"regexp"
	"unicode"

	// "net/url"
	// "sort"
	"strconv"
	"strings"
	//"github.com/saycv/ebomgen/pkg/types"
)

func DeleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)      //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}

// GetFValFromEVal convert float from string
func GetFValFromEVal(evalue string) float64 {
	var _value string
	var baseval string
	var fvalue float64
	var multipliers = map[string]float64{
		"M": 1e6,
		"K": 1e3,
		"k": 1e3,
		"m": 1e-3,
		"U": 1e-6,
		"u": 1e-6,
		"N": 1e-9,
		"n": 1e-9,
		"P": 1e-12,
		"p": 1e-12,
	}
	multipliersKeys := make([]string, 0, len(multipliers))
	for k := range multipliers {
		multipliersKeys = append(multipliersKeys, k)
	}
	_value = strings.Replace(evalue, " ", "", -1)
	if len(_value) == 0 {
		fvalue = -1.0
		//return -1.0
	} else if unicode.IsDigit([]rune(evalue)[0]) {
		re := regexp.MustCompile("\\d*\\.?\\d+")
		baseval = string(re.FindAll([]byte(evalue), -1)[0])
		//log.Debugf("Check [%s] in [%s]", string([]rune(_value)[0+len(baseval)]), strings.Join(multipliersKeys," "))
		if len(baseval) == len(_value) { // no multiplier
			fvalue, _ = strconv.ParseFloat(baseval, 64)
		} else if strings.Contains(strings.Join(multipliersKeys, " "), string([]rune(_value)[0+len(baseval)])) { // multiplier existss
			fvalue, _ = strconv.ParseFloat(baseval, 64)
			fvalue = fvalue * multipliers[string((evalue)[0+len(baseval)])]
		} else {
			fvalue, _ = strconv.ParseFloat(baseval, 64)
		}
	} else if strings.HasPrefix(strings.ToUpper(evalue), "CRY-") || strings.HasPrefix(strings.ToUpper(evalue), "OSC-") {
		re := regexp.MustCompile("\\d*\\.?\\d+")
		baseval = string(re.FindAll([]byte(evalue), -1)[0])
		if len(baseval) == len(_value) { // no multiplier
			fvalue, _ = strconv.ParseFloat(baseval, 64)
		} else if strings.Contains(strings.Join(multipliersKeys, " "), string([]rune(_value)[0+len(baseval)])) { // multiplier existss
			fvalue, _ = strconv.ParseFloat(baseval, 64)
			fvalue = fvalue * multipliers[string((evalue)[4+len(baseval)])]
		} else {
			fvalue, _ = strconv.ParseFloat(baseval, 64)
		}
	} else {
		fvalue = -1.0
	}
	return fvalue
}
