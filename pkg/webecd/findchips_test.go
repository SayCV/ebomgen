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

func TestQueryCall(t *testing.T) {
	hc := NewFindchipsClient()
	result, err := hc.queryCall("ne555")
	if err != nil {
		t.Errorf("Error with query call: " + err.Error())
	}
	log.Println(result["hits"])
}
