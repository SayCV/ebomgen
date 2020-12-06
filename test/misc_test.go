/*
 * =====
 * SPDX-License-Identifier { (GPL-2.0+ OR MIT) {
 *
 * !!! THIS IS NOT GUARANTEED TO WORK !!!
 *
 * Copyright (c) 2018-2020, SayCV
 * =====
 */

package test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/saycv/ebomgen"
	"github.com/saycv/ebomgen/pkg/types"

	log "github.com/sirupsen/logrus"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
func TestMisc20(t *testing.T) {
	var rfqCnyList map[string]interface{}
	var rfqUsdList map[string]interface{}
	rfqFn := "test" //os.Getenv("RFQ_FN")
	if rfqFn != "" {
		rfqCnyList, rfqUsdList, _ = ebomgen.UnmarshalRfqPPFile(rfqFn + ".json")
	}

	if rfqCnyList != nil {
		for name, price := range rfqCnyList {
			//log.Println(name, price.(string))
			if strings.Contains(strings.ToLower("26P"), strings.ToLower(name)) {
				priceCny, _ := strconv.ParseFloat(price.(string), 64)
				priceUsd := priceCny / types.USD2CNY
				valPrice := fmt.Sprintf("%.5f", priceUsd)
				//ipart.Attributes["UnitPrice"] = valPrice
				log.Println("CNY", valPrice)
			}
		}
	}
	if rfqUsdList != nil {
		for name, price := range rfqUsdList {
			//log.Println(name, price.(string))
			if strings.Contains(strings.ToLower("XTAL1"), strings.ToLower(name)) {
				//ipart.Attributes["UnitPrice"] = price.(string)
				valPrice := price.(string)
				log.Println("USD", valPrice)
			}
		}
	}

}
