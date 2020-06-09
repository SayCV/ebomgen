package webecd

import (
	//"fmt"
	//"regexp"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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
	digikeyHome        = "https://www.digikey.com"
	digikeyParamSearch = "https://www.digikey.com/Search"
)

type DigikeyClient struct {
	RemoteHost string
	client     *http.Client
	infoCache  map[string]interface{}
}

func NewDigikeyClient() *DigikeyClient {
	hc := &DigikeyClient{
		RemoteHost: digikeyHome}
	hc.client = &http.Client{}
	hc.infoCache = make(map[string]interface{})
	return hc
}

func (hc *DigikeyClient) queryCallDetail(suburl string, partSpecs types.EBOMWebPart) (types.EBOMWebPart, error) {
	//var partSpecs types.EBOMWebPart
	paramString := suburl
	method := ""

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "" + method + paramStringUnescaped)
	ua := utils.GetUaHeaders()
	log.Infof("Headers: " + ua)
	//resp, err := hc.client.Get(hc.RemoteHost + "" + method + paramString)
	request, err := http.NewRequest("GET", hc.RemoteHost+""+method+paramString, nil)
	request.Header.Add("Cookie", "name=anny")
	request.Header.Add("User-Agent", ua)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return partSpecs, err
	}
	resp, err := hc.client.Do(request)
	if err != nil {
		return partSpecs, err
	}
	if resp.StatusCode != 200 {
		return partSpecs, errors.Errorf(digikeyHome + " queryCall error: " + resp.Status)
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
			partSpecs.RoHS = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Part Life Cycle") {
			partSpecs.Lifecycle = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Package / Case") {
			partSpecs.PackageCase = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Moisture Sensitivity") {
			partSpecs.MoistureSensitive = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Peak Reflow") {
			partSpecs.ReflowTemperaturePeak = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Minimum Operating Temperature") {
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
		return !found1st
	})
	return partSpecs, nil
}

func (hc *DigikeyClient) queryCall(mpn string) (types.EBOMWebPart, error) {
	var partSpecs types.EBOMWebPart
	var detaillink string
	paramString := mpn
	method := "en?keywords="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/products/" + method + paramStringUnescaped)
	ua := utils.GetUaHeaders()
	log.Infof("Headers: " + ua)
	request, err := http.NewRequest("GET", hc.RemoteHost+"/products/"+method+paramString, nil)
	request.Header.Add("Cookie", "name=anny")
	request.Header.Add("User-Agent", ua)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return partSpecs, err
	}
	resp, err := hc.client.Do(request)
	if err != nil {
		return partSpecs, err
	}
	if resp.StatusCode != 200 {
		return partSpecs, errors.Errorf(digikeyHome + " queryCall error: " + resp.Status)
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
	content := doc.Find("#productTable")

	found1st := false
	content.Find("tbody tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		log.Printf("Found %d: %s - %s", i, band, title)

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
				partSpecs.MPN = types.PartParameter{_val, types.ParamFromDigikey}
			default:
			}
			found1st = true
		})
		return !found1st
	})

	if !found1st {
		return partSpecs, errors.Errorf(digikeyHome + " not found vaildate data: " + resp.Status)
	}
	//partSpecs, err = hc.queryCallDetail(detaillink, partSpecs)

	return partSpecs, nil
}

func (hc *DigikeyClient) QueryWDCall(mpn string) (types.EBOMWebPart, error) {
	var partSpecs types.EBOMWebPart
	var detaillink webdriver.WebElement
	//var cookie webdriver.Cookie
	reDigit := regexp.MustCompile("\\d*\\.?\\d+")

	paramString := mpn
	method := "en?keywords="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/products/" + method + paramStringUnescaped)

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
	err = wh.MaximizeWindow()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	size2, err := wh.GetSize()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("width: %d, height: %d", size2.Width, size2.Height)

	err = session.Url(hc.RemoteHost + "/products/" + method + paramString)
	if err != nil {
		return partSpecs, err
	}
	//time.Sleep(10 * time.Second)
	_, err = session.FindElement(webdriver.ID, "noResults")
	if err == nil {
		return partSpecs, errors.Errorf(digikeyHome + " noResults")
	}

	qpLinkList, err := session.FindElement(webdriver.ID, "qpLinkList") // when search keywords is not met
	if err == nil {
		trs, err := qpLinkList.FindElements(webdriver.TagName, "tr")
		if err != nil {
			return partSpecs, err
		}
		for k, trv := range trs {
			log.Println(k, trv)
			tds, err := trv.FindElements(webdriver.TagName, "td")
			if err != nil {
				return partSpecs, err
			}
			for j, tdv := range tds {
				switch j {
				case 0:
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
			break
		}
		_, err = session.GetAlertText()
		if err == nil {
			err = session.DismissAlert()
		}
		err = detaillink.Click()
		if err != nil {
			return partSpecs, err
		}
		time.Sleep(2 * time.Second)
	} else {
		productIndexList, err := session.FindElement(webdriver.ID, "productIndexList")
		if err == nil {
			type productIndexItem struct {
				detaillink webdriver.WebElement
				href       string
				items      int
			}
			var prodItem []productIndexItem
			lis, err := productIndexList.FindElements(webdriver.TagName, "li")
			if err != nil {
				return partSpecs, err
			}
			for k, liv := range lis {
				log.Println(k, liv)
				_val, err := liv.FindElement(webdriver.TagName, "a")
				if err != nil {
					return partSpecs, err
				}
				detaillink = _val
				href, err := _val.GetAttribute("href")
				if err != nil {
					return partSpecs, err
				}
				items, err := liv.Text()
				baseval := string(reDigit.FindAll([]byte(items), -1)[0])
				itemsInt, err := strconv.Atoi(baseval)
				log.Println(itemsInt)
				_prod := productIndexItem{detaillink, href, itemsInt}
				prodItem = append(prodItem, _prod)
			}
			index := 0
			lastItems := prodItem[0].items
			for k, prod := range prodItem {
				items := prod.items
				if lastItems < items {
					lastItems = items
					index = k
				}
			}
			detaillink = prodItem[index].detaillink
			_, err = session.GetAlertText()
			if err == nil {
				err = session.DismissAlert()
			}
			err = detaillink.Click()
			if err != nil {
				return partSpecs, err
			}
			time.Sleep(2 * time.Second)
		}
	}

	clickcnts := 0
	for {
		we, err := session.FindElement(webdriver.ID, "lnkPart") // productTable
		if err != nil {
			return partSpecs, err
		}
		//log.Println(we)
		trs, err := we.FindElements(webdriver.TagName, "tr")
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

		session.ExecuteScript("window.scrollBy(0, 200)", make([]interface{}, 0))
		err = detaillink.Click()
		clickcnts = clickcnts + 1

		if err != nil {
			log.Infof("click failed: %d try again!!!", clickcnts)
			if err != nil && clickcnts == 5 {
				return partSpecs, err
			}
			//err = session.SendKeysOnActiveElement(string(keyboard.KeyTab))
			//err = session.SendKeysOnActiveElement(string(keyboard.KeyEsc))
			err = session.Refresh()
			if err != nil {
				log.Infof("SendKeysOnActiveElement failed")
			}
			time.Sleep(1 * time.Second)
		} else {
			break
		}
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
		//log.Println(k, trv)
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
			_val := strings.Split(title, "~")
			baseval := string(reDigit.FindAll([]byte(_val[0]), -1)[0])
			partSpecs.OperatingTemperatureMin = types.PartParameter{baseval, types.ParamFromDigikey}
			baseval = string(reDigit.FindAll([]byte(_val[1]), -1)[0])
			partSpecs.OperatingTemperatureMax = types.PartParameter{baseval, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Maximum Operating Temperature") {
			partSpecs.OperatingTemperatureMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage - Min") {
			partSpecs.SupplyVoltageMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Voltage-Max") {
			partSpecs.SupplyVoltageMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Voltage - Supply") {
			_val := strings.Split(title, "~")
			baseval := string(reDigit.FindAll([]byte(_val[0]), -1)[0])
			partSpecs.SupplyVoltageMin = types.PartParameter{baseval, types.ParamFromDigikey}
			baseval = string(reDigit.FindAll([]byte(_val[1]), -1)[0])
			partSpecs.SupplyVoltageMax = types.PartParameter{baseval, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Current-Min") {
			partSpecs.SupplyCurrentMin = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Supply Current-Max") {
			partSpecs.SupplyCurrentMax = types.PartParameter{title, types.ParamFromDigikey}
		} else if strings.HasPrefix(band, "Current - Supply") {
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
