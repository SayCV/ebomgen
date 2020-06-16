package webecd

import (
	//"flag"
	//"fmt"
	//"strconv"
	//"strings"
	//"regexp"
	//"strconv"
	"net/url"
	
	"log"
	"testing"
)

func TestSzlcscQueryCall(t *testing.T) {
	hc := NewSzlcscClient()
	result, err := hc.QueryCall(url.QueryEscape("MAX706TESA"))
	if err != nil {
		t.Errorf("Error with query call: " + err.Error())
	}
	log.Println("TestQueryCall Done.")
	log.Println(result)
}
