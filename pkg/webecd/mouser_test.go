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

func TestMouserQueryCall(t *testing.T) {
	hc := NewMouserClient()
	result, err := hc.queryCall("MAX706TESA")
	if err != nil {
		t.Errorf("Error with query call: " + err.Error())
	}
	log.Println("TestQueryCall Done.")
	log.Println(result)
}

func TestMisc(t *testing.T) {

}
