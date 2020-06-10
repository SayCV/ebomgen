package webecd

import (
	"fmt"
	//"regexp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/saycv/ebomgen/pkg/types"
	"github.com/saycv/ebomgen/pkg/utils"

	"github.com/PuerkitoBio/goquery"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	szlcscHome        = "https://so.szlcsc.com"
	szlcscParamSearch = "https://so.szlcsc.com/parametric"
)

type SzlcscClient struct {
	RemoteHost string
	client     *http.Client
	infoCache  map[string]interface{}
}

func NewSzlcscClient() *SzlcscClient {
	hc := &SzlcscClient{
		RemoteHost: szlcscHome}
	hc.client = &http.Client{}
	hc.infoCache = make(map[string]interface{})
	return hc
}

func (hc *SzlcscClient) QueryCallDetail(suburl string, partSpecs types.EBOMWebPart) (types.EBOMWebPart, error) {

	return partSpecs, nil
}

func (hc *SzlcscClient) QueryCall(mpn string) (types.EBOMWebPart, error) {
	var partSpecs types.EBOMWebPart
	//var detaillink string
	paramString := mpn
	method := "global.html?k="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/" + method + paramStringUnescaped)
	resp, err := hc.client.Get(hc.RemoteHost + "/" + method + paramString)
	if err != nil {
		return partSpecs, err
	}
	if resp.StatusCode != 200 {
		return partSpecs, errors.Errorf(szlcscHome + " queryCall error: " + resp.Status)
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
	content := doc.Find("#shop-list")

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
			case 1:
				s.Find("span").Each(func(j int, s *goquery.Selection) {
					_val := s.Text()
					//log.Printf(_val)
					switch j {
					case 6:
						partSpecs.MPN = types.PartParameter{_val, types.ParamFromSzlcsc}
					}
				})

			case 2:
				valPrice := ""
				lastPrice := ""
				s.Find("li").Each(func(j int, s *goquery.Selection) {
					if j != 0 {
						_val := strings.Replace(s.Text(), "\n", "", -1)
						_val = utils.DeleteExtraSpace(_val)
						_val = strings.TrimSpace(_val)
						//log.Printf(_val)
						_vallist := strings.Split(_val, " ")
						priceBreak, _ := strconv.Atoi(strings.Replace(_vallist[0], "+：", "", -1))
						_valPrice := strings.Replace(_vallist[1], "￥", "", -1)
						if priceBreak <= 1000 {
							valPrice = _valPrice
						} else if valPrice == "" {
							valPrice = _valPrice
						}
						lastPrice = _valPrice
					}
				})
				if valPrice == "" {
					valPrice = lastPrice
				}
				log.Println(valPrice)
				priceCny, _ := strconv.ParseFloat(valPrice, 64)
				priceUsd := priceCny / types.USD2CNY
				valPrice = fmt.Sprintf("%.5f", priceUsd)
				partSpecs.UnitPrice = types.PartParameter{valPrice, types.ParamFromSzlcsc}
			default:
			}
			found1st = true
		})
		return !found1st
	})

	//partSpecs, err = hc.queryCallDetail(detaillink, partSpecs)

	return partSpecs, nil
}
