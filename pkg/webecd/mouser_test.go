package webecd

import (
	//"flag"
	//"fmt"
	//"strconv"
	"strings"
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

	value := "RN4-0402"
	_vallist := strings.Split(value, "-")
	if len(_vallist) > 2 {
		value = strings.Join(_vallist[:2], "-")
	}
	log.Println(value)
}

func TestMisc2(t *testing.T) {
	example := "RN4-0402"
	example = example[3:]
	log.Println(example)
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(example, "")

	fmt.Printf("A string of %s becomes %s \n", example, processedString)
}

func TestMisc3(t *testing.T) {
	example := "1X6PIN-2.54mm-CIR-VERT"
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(example, " ")
	fmt.Printf("A string of %s becomes %s \n", example, processedString)

	reg2, err := regexp.Compile(`([0-9]+)[*X]([0-9]+)`)
	if err != nil {
		log.Fatal(err)
	}
	processedString = reg2.ReplaceAllString(example, "${1}.${2}")
	fmt.Printf("A string of %s becomes %s \n", example, processedString)

	reg3, err := regexp.Compile(`.HDR`)
	if err != nil {
		log.Fatal(err)
	}
	processedString = reg3.ReplaceAllString(processedString, " header ")
	fmt.Printf("A string of %s becomes %s \n", example, processedString)
}
