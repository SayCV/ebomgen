package webecd

import (
	//"fmt"
	//"regexp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

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
	resp, err := hc.client.Get(hc.RemoteHost + "" + method + paramString)
	if err != nil {
		return partSpecs, err
	}
	if resp.StatusCode != 200 {
		return partSpecs, errors.Errorf(mouserHome + " queryCallDetail error: " + resp.Status)
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
	content.Find("#FindSimilarPnl").Find("div-table-body div-table-row").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("a").Text()
		title = strings.Replace(title, "\n", "", -1)
		title = utils.DeleteExtraSpace(title)
		title = strings.TrimSpace(title)
		log.Printf("Found %d: %s - %s", i, band, title)

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
	reqest, err := http.NewRequest("GET", hc.RemoteHost + "/Search/" + method + paramString, nil)
	reqest.Header.Add("Cookie", "name=anny")
    reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.30 Safari/537.36")
    reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return partSpecs, err
	}
	resp, err := hc.client.Do(reqest)
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
			_val := strings.Replace(s.Text(), "\n", "", -1)
			_val = utils.DeleteExtraSpace(_val)
			log.Printf("Found %d-%d: %s", i, j, _val)
			switch j {
            case 2:
				detaillink,_ = s.Find("a").Attr("href")
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
	//partSpecs, err = hc.queryCallDetail(detaillink, partSpecs)

	return partSpecs, nil
}
