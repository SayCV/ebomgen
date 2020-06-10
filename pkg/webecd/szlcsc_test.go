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

func TestSzlcsQueryCall(t *testing.T) {
	hc := NewSzlcsClient()
	result, err := hc.queryCall("MAX706TESA")
	if err != nil {
		t.Errorf("Error with query call: " + err.Error())
	}
	log.Println("TestQueryCall Done.")
	log.Println(result)
}
