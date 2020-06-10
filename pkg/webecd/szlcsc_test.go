package webecd

import (
	//"flag"
	//"fmt"
	//"strconv"
	//"strings"
	//"regexp"
	//"strconv"
	"log"
	"testing"
)

func TestSzlcscQueryCall(t *testing.T) {
	hc := NewSzlcscClient()
	result, err := hc.QueryCall("MAX706TESA")
	if err != nil {
		t.Errorf("Error with query call: " + err.Error())
	}
	log.Println("TestQueryCall Done.")
	log.Println(result)
}
