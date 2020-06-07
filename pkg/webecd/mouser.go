package webecd

import (
	//"fmt"
	//"regexp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/saycv/ebomgen/pkg/types"
	"github.com/saycv/ebomgen/pkg/utils"

	"github.com/PuerkitoBio/goquery"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	mouserHome        = "https://www.mouser.com"
	mouserParamSearch = "https://www.mouser.com/Search"
)

type MouserClient struct {
	RemoteHost string
	client     *http.Client
	infoCache  map[string]interface{}
}

func NewMouserClient() *MouserClient {
	hc := &MouserClient{
		RemoteHost: mouserHome}
	hc.client = &http.Client{}
	hc.infoCache = make(map[string]interface{})
	return hc
}

func (hc *MouserClient) queryCallDetail(suburl string, partSpecs types.EBOMWebPart) (types.EBOMWebPart, error) {
	//var partSpecs types.EBOMWebPart
	paramString := suburl
	method := ""

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "" + method + paramStringUnescaped)
	//resp, err := hc.client.Get(hc.RemoteHost + "" + method + paramString)
	request, err := http.NewRequest("GET", hc.RemoteHost+""+method+paramString, nil)
	request.Header.Add("Cookie", "name=anny")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return partSpecs, err
	}
	resp, err := hc.client.Do(request)
	if err != nil {
		return partSpecs, err
	}
	if resp.StatusCode != 200 {
		return partSpecs, errors.Errorf(mouserHome + " queryCall error: " + resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return partSpecs, err
	}

	utfBody := strings.NewReader(string(body))

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal(err)
	}

	// Find the key items
	content := doc.Find("#pdpMainContentDiv")

	found1st := false
	content.Find(".div-table-body .div-table-row").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		band := s.Find(".col-xs-4").Find("label").Text()
		band = strings.Replace(band, "\n", "", -1)
		band = strings.TrimSpace(band)
		title := s.Find(".col-xs-5").Text()
		title = strings.Replace(title, "\n", "", -1)
		title = strings.TrimSpace(title)
		//log.Printf("Found %d: %s - %s", i, band, title)

		if strings.HasPrefix(band, "RoHS") {
			partSpecs.RoHS = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Part Life Cycle") {
			partSpecs.Lifecycle = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Package / Case") {
			partSpecs.PackageCase = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Moisture Sensitivity") {
			partSpecs.MoistureSensitive = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Peak Reflow") {
			partSpecs.ReflowTemperaturePeak = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Minimum Operating Temperature") {
			partSpecs.OperatingTemperatureMin = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Maximum Operating Temperature") {
			partSpecs.OperatingTemperatureMax = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Supply Voltage - Min") {
			partSpecs.SupplyVoltageMin = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Supply Voltage-Max") {
			partSpecs.SupplyVoltageMax = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Supply Voltage-Nom") {
			partSpecs.SupplyVoltageNom = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Supply Current-Min") {
			partSpecs.SupplyCurrentMin = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Supply Current-Max") {
			partSpecs.SupplyCurrentMax = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Operating Supply Current") {
			partSpecs.SupplyCurrentNom = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Power Dissipation-Min") {
			partSpecs.PowerDissipationMin = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Power Dissipation-Max") {
			partSpecs.PowerDissipationMax = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Pd - Power Dissipation") {
			partSpecs.PowerDissipationNom = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Length") {
			partSpecs.UnitLength = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Width") {
			partSpecs.UnitWidth = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Height") {
			partSpecs.UnitHeight = types.PartParameter{title, types.ParamFromMouser}
		} else if strings.HasPrefix(band, "Unit Weight") {
			partSpecs.UnitWeight = types.PartParameter{title, types.ParamFromMouser}
		}
		return !found1st
	})
	return partSpecs, nil
}

func (hc *MouserClient) queryCall(mpn string) (types.EBOMWebPart, error) {
	var partSpecs types.EBOMWebPart
	var detaillink string
	paramString := mpn
	method := "Refine?Keyword="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/Search/" + method + paramStringUnescaped)
	//resp, err := hc.client.Get(hc.RemoteHost + "/Search/" + method + paramString)
	request, err := http.NewRequest("GET", hc.RemoteHost+"/Search/"+method+paramString, nil)
	request.Header.Add("Cookie", "name=anny")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return partSpecs, err
	}
	resp, err := hc.client.Do(request)
	if err != nil {
		return partSpecs, err
	}
	if resp.StatusCode != 200 {
		return partSpecs, errors.Errorf(mouserHome + " queryCall error: " + resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return partSpecs, err
	}
	//log.Infof(string(body))
	utfBody := strings.NewReader(string(body))

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal(err)
	}

	// Find the key items
	content := doc.Find(".search-table-wrapper")

	found1st := false
	content.Find("tbody tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		//band := s.Find("a").Text()
		//title := s.Find("i").Text()
		//log.Printf("Found %d: %s - %s", i, band, title)

		s.Children().Each(func(j int, s *goquery.Selection) {
			//_val := strings.Replace(s.Text(), "\n", "", -1)
			//_val = utils.DeleteExtraSpace(_val)
			//log.Printf("Found %d-%d: %s", i, j, _val)
			switch j {
			case 2:
				detaillink, _ = s.Find("a").Attr("href")
				_val := s.Find("a").Text()
				_val = strings.Replace(_val, "\n", "", -1)
				_val = utils.DeleteExtraSpace(_val)
				_val = strings.TrimSpace(_val)
				log.Printf(_val)
				log.Printf(detaillink)
				partSpecs.MPN = types.PartParameter{_val, types.ParamFromMouser}
			default:
			}
			found1st = true
		})
		return !found1st
	})

	if !found1st {
		return partSpecs, errors.Errorf(mouserHome + " not found vaildate data: " + resp.Status)
	}
	partSpecs, err = hc.queryCallDetail(detaillink, partSpecs)

	return partSpecs, nil
}

func (hc *MouserClient) queryWDCall(mpn string) (types.EBOMWebPart, error) {
	var partSpecs types.EBOMWebPart
	var detaillink webdriver.WebElement
	//var cookie webdriver.Cookie
	paramString := mpn
	method := "Refine?Keyword="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/Search/" + method + paramStringUnescaped)

	chromeDriver, session := utils.InitChromeBrowser()
	//err := session.SetCookie(cookie)
	//if err != nil {
	//	return partSpecs, err
	//}
	wh := session.GetCurrentWindowHandle()
	size, err := wh.GetSize()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("width: %d, height: %d", size.Width, size.Height)
	//err = wh.MaximizeWindow()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	size2, err := wh.GetSize()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("width: %d, height: %d", size2.Width, size2.Height)

	err = session.Url(hc.RemoteHost + "/Search/" + method + paramString)
	if err != nil {
		return partSpecs, err
	}
	//time.Sleep(10 * time.Second)

	we, err := session.FindElement(webdriver.ClassName, "search-table-wrapper")
	if err != nil {
		return partSpecs, err
	}
	tbody, err := we.FindElement(webdriver.TagName, "tbody")
	if err != nil {
		return partSpecs, err
	}
	//log.Println(we)
	trs, err := tbody.FindElements(webdriver.TagName, "tr")
	if err != nil {
		return partSpecs, err
	}
	found1st := false
	for k, trv := range trs {
		log.Println(k, trv)
		tds, err := trv.FindElements(webdriver.TagName, "td")
		if err != nil {
			return partSpecs, err
		}
		for j, tdv := range tds {
			switch j {
			case 2:
				//detaillink = tdv
				_val, err := tdv.FindElements(webdriver.TagName, "a") //[0].GetAttribute("href")
				if err != nil {
					return partSpecs, err
				}
				detaillink = _val[0]
				href, err := detaillink.GetAttribute("href")
				if err != nil {
					return partSpecs, err
				}
				log.Printf(href)
				//partSpecs.MPN = types.PartParameter{href, types.ParamFromDigikey}
			default:
			}
		}
		found1st = true
		break
	}

	if !found1st {
		return partSpecs, errors.Errorf(digikeyHome+" not found vaildate data: %v", err)
	}

	//partSpecs, err = hc.queryWDCallDetail(detaillink, partSpecs)

	//err = detaillink.SendKeys("\n")
	//session.ExecuteScript("arguments[0].click();", []interface{}{detaillink})
	if err != nil {
		return partSpecs, err
	}
	err = detaillink.Click()
	if err != nil {
		return partSpecs, err
	}
	time.Sleep(2 * time.Second)
	// #pdp_content
	//   #product-photo
	//   product-details-procurement
	//   product-details-overview
	//     product-overview-photo-spacer
	//     product-details-documents-media product-details-section
	//     product-details-product-attributes product-details-section
	//     product-details-environmental-export product-details-section
	prodAttr, err := session.FindElement(webdriver.ID, "product-attribute-table")
	if err != nil {
		return partSpecs, err
	}
	prodAttrTrs, err := prodAttr.FindElements(webdriver.TagName, "tr")
	if err != nil {
		return partSpecs, err
	}
	for k, trv := range prodAttrTrs {
		log.Println(k, trv)
		if err != nil {
			continue
		}
		if k == 0 {
			continue
		}
		_val, err := trv.FindElement(webdriver.TagName, "th")
		if err != nil {
			continue
		}
		band, err := _val.Text()
		if err != nil {
			continue
		}
		//band = strings.Replace(band, "\n", "", -1)
		//band = strings.TrimSpace(band)

		tds, err := trv.FindElements(webdriver.TagName, "td")
		if err != nil {
			continue
		}
		title, err := tds[0].Text()
		if err != nil {
			continue
		}
		if strings.HasPrefix(band, "RoHS") {
			partSpecs.RoHS = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Part Status") {
			partSpecs.Lifecycle = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Package / Case") {
			partSpecs.PackageCase = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Moisture Sensitivity") {
			partSpecs.MoistureSensitive = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Peak Reflow") {
			partSpecs.ReflowTemperaturePeak = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Operating Temperature") {
			partSpecs.OperatingTemperatureMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Maximum Operating Temperature") {
			partSpecs.OperatingTemperatureMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage - Min") {
			partSpecs.SupplyVoltageMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage-Max") {
			partSpecs.SupplyVoltageMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage-Nom") {
			partSpecs.SupplyVoltageNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Current-Min") {
			partSpecs.SupplyCurrentMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Current-Max") {
			partSpecs.SupplyCurrentMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Operating Supply Current") {
			partSpecs.SupplyCurrentNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Power Dissipation-Min") {
			partSpecs.PowerDissipationMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Power Dissipation-Max") {
			partSpecs.PowerDissipationMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Pd - Power Dissipation") {
			partSpecs.PowerDissipationNom = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Length") {
			partSpecs.UnitLength = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Width") {
			partSpecs.UnitWidth = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Height") {
			partSpecs.UnitHeight = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Unit Weight") {
			partSpecs.UnitWeight = types.PartParameter{title, types.ParamFromDigikey}
		}
	}
	prodPrice, err := session.FindElement(webdriver.ClassName, "product-dollars")
	if err != nil {
		return partSpecs, err
	}
	prodPriceTrs, err := prodPrice.FindElements(webdriver.TagName, "tr")
	if err != nil {
		return partSpecs, err
	}
	prodPriceBand, err := prodPriceTrs[0].Text()
	if err != nil {
		return partSpecs, err
	}
	log.Printf(prodPriceBand)
	prodPriceTitle, err := prodPriceTrs[1].Text()
	if err != nil {
		return partSpecs, err
	}
	log.Printf(prodPriceTitle)
	_val := strings.Split(prodPriceTitle, " ")
	partSpecs.UnitPrice = types.PartParameter{_val[1], types.ParamFromDigikey}

	session.Delete()
	chromeDriver.Stop()

	return partSpecs, nil
}
