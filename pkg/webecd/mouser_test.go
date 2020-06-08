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
	"time"

	"github.com/saycv/ebomgen/pkg/utils"
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

func TestMouserQueryWDCall(t *testing.T) {
	hc := NewMouserClient()
	result, err := hc.queryWDCall("MAX706TESA")
	if err != nil {
		t.Errorf("Error with query call: " + err.Error())
	}
	log.Println("TestQueryCall Done.")
	log.Println(result)
}

func TestMisc1(t *testing.T) {
	result := utils.GetUaHeaders()
	log.Println(result)
}

func TestMisc2(t *testing.T) {
	chromeDriver, session := utils.InitChromeBrowser()
	err := session.Url("http://bing.com")
	if err != nil {
		log.Println(err)
	}
	time.Sleep(10 * time.Second)
	session.Delete()
	chromeDriver.Stop()
}
