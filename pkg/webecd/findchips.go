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
	findchipsHome        = "https://www.findchips.com"
	findchipsParamSearch = "https://www.findchips.com/parametric"
)

type FindchipsClient struct {
	RemoteHost string
	client     *http.Client
	infoCache  map[string]interface{}
}

func NewFindchipsClient() *FindchipsClient {
	hc := &FindchipsClient{
		RemoteHost: findchipsHome}
	hc.client = &http.Client{}
	hc.infoCache = make(map[string]interface{})
	return hc
}

func (hc *FindchipsClient) queryCallDetail(suburl string, partSpecs types.EBOMWebPart) (types.EBOMWebPart, error) {
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
		return partSpecs, errors.Errorf(findchipsHome + " queryCallDetail error: " + resp.Status)
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
	content := doc.Find(".dash-section")

	found1st := false
	content.Find(".part-details-list-item").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		band := s.Find("small").Text()
		title := s.Find("p").Text()
		title = strings.Replace(title, "\n", "", -1)
		title = utils.DeleteExtraSpace(title)
		title = strings.TrimSpace(title)
		//log.Printf("Found %d: %s - %s", i, band, title)
		if strings.HasPrefix(band, "Rohs") {
			partSpecs.RoHS = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Part Life Cycle") {
			partSpecs.Lifecycle = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Package Code") {
			partSpecs.PackageCase = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Moisture Sensitivity") {
			partSpecs.MoistureSensitive = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Peak Reflow") {
			partSpecs.ReflowTemperaturePeak = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Operating Temperature-Min") {
			partSpecs.OperatingTemperatureMin = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Operating Temperature-Max") {
			partSpecs.OperatingTemperatureMax = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Supply Voltage-Min") {
			partSpecs.SupplyVoltageMin = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Supply Voltage-Max") {
			partSpecs.SupplyVoltageMax = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Supply Voltage-Nom") {
			partSpecs.SupplyVoltageNom = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Supply Current-Min") {
			partSpecs.SupplyCurrentMin = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Supply Current-Max") {
			partSpecs.SupplyCurrentMax = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Supply Current-Nom") {
			partSpecs.SupplyCurrentNom = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Power Dissipation-Min") {
			partSpecs.PowerDissipationMin = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Power Dissipation-Max") {
			partSpecs.PowerDissipationMax = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Power Dissipation-Nom") {
			partSpecs.PowerDissipationNom = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Length") {
			partSpecs.UnitLength = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Width") {
			partSpecs.UnitWidth = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Seated Height") {
			partSpecs.UnitHeight = types.PartParameter{title, types.ParamFromFindchips}
		} else if strings.HasPrefix(band, "Weight") {
			partSpecs.UnitWeight = types.PartParameter{title, types.ParamFromFindchips}
		}
		return !found1st
	})
	return partSpecs, nil
}

func (hc *FindchipsClient) queryCall(mpn string) (types.EBOMWebPart, error) {
	var partSpecs types.EBOMWebPart
	var detaillink string
	paramString := mpn
	method := "search?term="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/parametric/" + method + paramStringUnescaped)
	resp, err := hc.client.Get(hc.RemoteHost + "/parametric/" + method + paramString)
	if err != nil {
		return partSpecs, err
	}
	if resp.StatusCode != 200 {
		return partSpecs, errors.Errorf(findchipsHome + " queryCall error: " + resp.Status)
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
	content := doc.Find(".parametric-content")

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
				detaillink,_ = s.Find("a").Attr("href")
				_val := s.Find("a").Text()
				_val = strings.Replace(_val, "\n", "", -1)
				_val = utils.DeleteExtraSpace(_val)
				_val = strings.TrimSpace(_val)
				log.Printf(_val)
				log.Printf(detaillink)
				partSpecs.MPN = types.PartParameter{_val, types.ParamFromFindchips}

				_val = s.Find("span").Text()
				partSpecs.MFR = types.PartParameter{_val, types.ParamFromFindchips}
				log.Printf(_val)

            case 3:
				datasheetlink,_ := s.Find("a").Attr("href")
				log.Printf("Found datasheet: %s", datasheetlink)
				partSpecs.Datasheet = types.PartParameter{datasheetlink, types.ParamFromFindchips}

            case 4:
				_val := s.Find("a").Text()
				_val = strings.Replace(_val, "\n", "", -1)
				_val = utils.DeleteExtraSpace(_val)
				_val = strings.TrimSpace(_val)
				log.Printf("Found price: %s", _val)
				partSpecs.UnitPrice = types.PartParameter{_val, types.ParamFromFindchips}
            default:
			}
			found1st = true
		})
		return !found1st
	})

	partSpecs, err = hc.queryCallDetail(detaillink, partSpecs)

	return partSpecs, nil
}
