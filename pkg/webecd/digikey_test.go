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

func TestDigikeyQueryCall(t *testing.T) {
	hc := NewDigikeyClient()
	result, err := hc.queryWDCall("stm32f407")
	if err != nil {
		t.Errorf("Error with query call: " + err.Error())
	}
	log.Println("TestQueryCall Done.")
	log.Println(result)
}
