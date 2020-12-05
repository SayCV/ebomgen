package webecd

import (
	//"flag"
	//"fmt"
	//"strconv"
	//"strings"
	//"regexp"
	//"strconv"

	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func TestDigikeyQueryCall(t *testing.T) {
	hc := NewDigikeyClient()
	result, err := hc.QueryWDCall("DS1338Z-33+")
	if err != nil {
		t.Errorf("Error with query call: " + err.Error())
	}
	log.Println("TestQueryCall Done.")
	log.Println(result)
	hc.Close()
}

func TestDigikeyJson(t *testing.T) {
	data, err := ioutil.ReadFile("dgkdata.json")
	checkError(err)
	prodPriceMap := &NgDgkData{}
	err = json.Unmarshal(data, &prodPriceMap)
	checkError(err)
	PageProps := prodPriceMap.Props.PageProps
	Pricing := PageProps.Envelope.Data.PriceQuantity.Pricing
	PricingTiers := Pricing[0].PricingTiers

	valPrice := ""
	lastPrice := ""
	for _, pricing := range PricingTiers {
		qty, err := strconv.Atoi(strings.ReplaceAll(pricing.BreakQty, ",", ""))
		checkError(err)
		if qty <= 1000 {
			valPrice = pricing.UnitPrice
		} else if valPrice == "" {
			valPrice = pricing.UnitPrice
		} else {
			break
		}
		lastPrice = pricing.UnitPrice
		log.Println(lastPrice)
	}
	log.Println(valPrice)
}

func TestMisc(t *testing.T) {
	data := `<script id="__NEXT_DATA__" type="application/json">{test}</script>`
	//re, _ := regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	re, _ := regexp.Compile(`<script.*?>(.*)</script>`)
	data = re.ReplaceAllString(data, "$1")

	log.Println(data)
}
