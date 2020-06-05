package webecd

import (
	"io/ioutil"
	"net/http"
	"net/url"

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

	log.Print(string(body))

	return result, nil
}
