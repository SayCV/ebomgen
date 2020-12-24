package utils

import (
	//"io"
	//"fmt"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sync"
	"time"
	"unicode"

	// "net/url"
	// "sort"
	"strconv"
	"strings"

	//"github.com/saycv/ebomgen/pkg/types"

	"github.com/fedesog/webdriver"
	log "github.com/sirupsen/logrus"
)

var (
	initOnce  sync.Once
	uaHeaders = []string{
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.153 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:30.0) Gecko/20100101 Firefox/30.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/537.75.14",
		"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Win64; x64; Trident/6.0)",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; it; rv:1.8.1.11) Gecko/20071127 Firefox/2.0.0.11",
		"Opera/9.25 (Windows NT 5.1; U; en)",
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
		"Mozilla/5.0 (compatible; Konqueror/3.5; Linux) KHTML/3.5.5 (like Gecko) (Kubuntu)",
		"Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.8.0.12) Gecko/20070731 Ubuntu/dapper-security Firefox/1.5.0.12",
		"Lynx/2.8.5rel.1 libwww-FM/2.14 SSL-MM/1.4.1 GNUTLS/1.2.9",
		"Mozilla/5.0 (X11; Linux i686) AppleWebKit/535.7 (KHTML, like Gecko) Ubuntu/11.04 Chromium/16.0.912.77 Chrome/16.0.912.77 Safari/535.7",
		"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:10.0) Gecko/20100101 Firefox/10.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.30 Safari/537.36",
	}
)

func uaHeadersInit() {
	rand.Seed(time.Now().Unix())
}

func GetUaHeaders() string {
	initOnce.Do(uaHeadersInit)
	return uaHeaders[rand.Intn(len(uaHeaders))]
}

func InitChromeBrowser() (*webdriver.ChromeDriver, *webdriver.Session) {
	cd := os.Getenv("ChromeDriver")
	log.Println("ChromeDriver: " + cd)
	chromeDriver := webdriver.NewChromeDriver(cd)
	err := chromeDriver.Start()
	if err != nil {
		log.Fatal(err)
	}
	desired := webdriver.Capabilities{"Platform": "Linux"}
	required := webdriver.Capabilities{}
	var args []string
	// args = []string{"--headless", "--no-sandbox"}
	args = []string{"--disable-notifications"}

	desired["chromeOptions"] = webdriver.Capabilities{"args": args}

	session, err := chromeDriver.NewSession(desired, required)
	if err != nil {
		log.Fatal(err)
	}

	return chromeDriver, session
}

func DeleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)       //替换tab为空格
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

func GetPinsFromFp(desc string, fp string) (int, error) {
	pins := 1
	fp = strings.ToUpper(fp)

	if strings.HasPrefix(desc, "Capacitor") {
		pins = 2
		if strings.HasPrefix(desc, "CapacitorArray") {
			pins = 8
		}
	} else if strings.HasPrefix(desc, "Resistor") {
		pins = 2
		if strings.HasPrefix(desc, "ResistorArray") {
			pins = 8
		}
	} else if strings.HasPrefix(desc, "Inductor") {
		pins = 2
	} else if strings.HasPrefix(desc, "Diode") {
		pins = 2
	} else if strings.HasPrefix(desc, "LED") {
		pins = 2
	} else if strings.HasPrefix(desc, "Transistor") {
		pins = 3
	} else if strings.HasPrefix(desc, "Crystal") {
		pins = 4
	} else if strings.HasPrefix(desc, "Oscillator") {
		pins = 4
	} else if strings.HasPrefix(desc, "IC") || strings.HasPrefix(desc, "Reg") || strings.HasPrefix(desc, "XFRM") {
		_val := 0
		reVal := regexp.MustCompile("\\d*\\d+")
		findRet := reVal.FindAll([]byte(fp), -1)
		if findRet != nil {
			_val, _ = strconv.Atoi(string(findRet[0]))
		}
		if strings.Contains(fp, "SOT23-6") {
			_val = 6
		} else if strings.Contains(fp, "SOT23-5") || strings.Contains(fp, "SC70-5") || strings.Contains(fp, "TO263-5") {
			_val = 5
		} else if strings.Contains(fp, "SC70") ||
			strings.Contains(fp, "SOT23") ||
			strings.Contains(fp, "SOT89") ||
			strings.Contains(fp, "SOT323") {
			_val = 3
		}
		pins = _val
	} else if strings.HasPrefix(desc, "Connector") || strings.HasPrefix(desc, "Switch") {
		reVal, _ := regexp.Compile(`([0-9]+)[X]([0-9]+)`)
		value := reVal.FindAll([]byte(fp), -1)
		if value == nil {
			reVal := regexp.MustCompile("\\d*\\d+")
			findRet := reVal.FindAll([]byte(fp), -1)
			if findRet != nil {
				pins, _ = strconv.Atoi(string(findRet[0]))
			}
		} else {
			valuelist := strings.Split(string(value[0]), "X")
			row, _ := strconv.Atoi(valuelist[0])
			col, _ := strconv.Atoi(valuelist[1])
			pins = row * col
		}
	} else if strings.HasPrefix(desc, "ConnRJ") {
		pins = 12
	} else if strings.HasPrefix(desc, "ConnUSB") {
		pins = 4
	}

	return pins, nil
}

func IsCapTanFp(fp string) bool {
	if strings.Contains(fp, "3216") || strings.Contains(fp, "3528") || strings.Contains(fp, "6032") || strings.Contains(fp, "7343") {
		return true
	}
	return false
}

func IsCapAecFp(fp string, value string) bool {
	fvalue := GetFValFromEVal(value)
	if fvalue <= 22*1e-06 {
		return false
	}
	reVal, _ := regexp.Compile(`([0-9]+)[X]([0-9]+)`)
	fpvalue := reVal.FindAll([]byte(fp), -1)
	if fpvalue != nil {
		return true
	}
	return false
}

func ElapsedTime() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("Time elapsed = %v\n", tc)
	}
}
