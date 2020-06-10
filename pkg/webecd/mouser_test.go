package webecd

import (
	//"flag"
	//"fmt"
	//"strconv"
	//"strings"
	//"regexp"
	//"strconv"

	"fmt"
	"log"
	"regexp"
	"testing"

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
	example := "rn4-0402"
	log.Println(example[3:])
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(example, "")

	fmt.Printf("A string of %s becomes %s \n", example, processedString)
}

func TestMisc3(t *testing.T) {
	example := "rn4-0402"
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(example, " ")

	fmt.Printf("A string of %s becomes %s \n", example, processedString)
}
