package webecd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	urlHome        = "https://www.findchips.com"
	urlParamSearch = "https://www.findchips.com/parametric"
)

type FindchipsClient struct {
	RemoteHost string
	client     *http.Client
	infoCache  map[string]interface{}
}

func NewFindchipsClient() *FindchipsClient {
	hc := &FindchipsClient{
		RemoteHost: urlHome}
	hc.client = &http.Client{}
	hc.infoCache = make(map[string]interface{})
	return hc
}

// https://www.findchips.com/parametric/search?term=max706
func (hc *FindchipsClient) queryCall(mpn string) (map[string]interface{}, error) {
	paramString := mpn
	method := "search?term="

	paramStringUnescaped, _ := url.QueryUnescape(paramString)
	log.Infof("Fetching: " + hc.RemoteHost + "/parametric/" + method + paramStringUnescaped)
	resp, err := hc.client.Get(hc.RemoteHost + "/parametric/" + method + paramString)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf(urlHome + " queryCall error: " + resp.Status)
	}
	result := make(map[string]interface{})
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	utfBody := strings.NewReader(string(body))

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal(err)
	}

	// Find the key items
	content := doc.Find(".parametric-content")
	content.Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})

	return result, nil
}
